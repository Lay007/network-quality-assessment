package bench

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"server_sfp_sla/internal/measurement"
	"time"
)

type Report struct {
	Version   int       `json:"schema_version"`
	ID        string    `json:"bench_session_id"`
	Status    string    `json:"status"`
	Start     time.Time `json:"started_at"`
	End       time.Time `json:"ended_at"`
	Agents    []Info    `json:"agents"`
	Sessions  []Session `json:"sessions"`
	Attempted int       `json:"packets_attempted"`
	Received  int       `json:"packets_received"`
	Lost      int       `json:"packets_lost"`
	Errors    []string  `json:"errors"`
	Note      string    `json:"note"`
}

type Client struct {
	http  *http.Client
	token string
}

func NewClient(token string) (*Client, error) {
	if len(token) < 16 || len(token) > 4096 {
		return nil, fmt.Errorf("token must have 16..4096 bytes")
	}
	return &Client{token: token, http: &http.Client{Timeout: 2 * time.Second, Transport: &http.Transport{Proxy: nil, DialContext: (&net.Dialer{Timeout: time.Second}).DialContext, MaxConnsPerHost: 2, MaxIdleConnsPerHost: 2, IdleConnTimeout: 5 * time.Second}, CheckRedirect: func(*http.Request, []*http.Request) error { return fmt.Errorf("redirects disabled") }}}, nil
}
func (c *Client) Close() { c.http.CloseIdleConnections() }

func endpoint(base string) error {
	u, err := url.Parse(base)
	if err != nil {
		return err
	}
	if u.Scheme != "http" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Path != "" {
		return fmt.Errorf("plain loopback HTTP endpoint required")
	}
	addr, err := measurement.LoopbackAddress(u.Host)
	if err != nil || addr.Port() == 0 {
		return fmt.Errorf("numeric loopback control endpoint required")
	}
	return nil
}

func (c *Client) request(ctx context.Context, base, method, path string, input, output any) error {
	if err := endpoint(base); err != nil {
		return err
	}
	var body io.Reader
	if input != nil {
		data, err := json.Marshal(input)
		if err != nil {
			return err
		}
		body = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, base+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	response, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("control request failed: %w", err)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(io.LimitReader(response.Body, 8*1024*1024+1))
	if err != nil {
		return err
	}
	if len(data) > 8*1024*1024 {
		return fmt.Errorf("control response exceeds limit")
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("control status %d", response.StatusCode)
	}
	if output != nil {
		if err := json.Unmarshal(data, output); err != nil {
			return fmt.Errorf("invalid control response: %w", err)
		}
	}
	return nil
}

// Run coordinates two agents; it does not synchronize their clocks or start times.
// It returns a report on failure so partial evidence can still be exported.
func (c *Client) Run(parent context.Context, bases []string, config measurement.Config) (Report, error) {
	r := Report{Version: 1, Start: time.Now().UTC(), Status: "running", Errors: []string{}, Note: "Loopback software RTT or explicit synthetic model; no synchronized launch, one-way or hardware timing qualification."}
	if len(bases) != 2 || bases[0] == bases[1] {
		return r, fmt.Errorf("two distinct agent endpoints required")
	}
	for _, base := range bases {
		if err := endpoint(base); err != nil {
			return r, err
		}
	}
	id, err := NewID()
	if err != nil {
		return r, err
	}
	r.ID = id
	ctx, cancel := context.WithTimeout(parent, 30*time.Second)
	defer cancel()
	finish := func(err error) (Report, error) {
		r.End = time.Now().UTC()
		r.Status = "completed"
		if err != nil {
			r.Status = "failed"
			if errors.Is(err, context.Canceled) {
				r.Status = "canceled"
			}
			if errors.Is(err, context.DeadlineExceeded) {
				r.Status = "timeout"
			}
			r.Errors = append(r.Errors, err.Error())
		}
		for _, s := range r.Sessions {
			if s.Result != nil {
				r.Attempted += s.Result.Attempted
				r.Received += s.Result.Received
				r.Lost += s.Result.Lost
			}
		}
		return r, err
	}
	for _, base := range bases {
		var info Info
		if err := c.request(ctx, base, "GET", "/v1/info", nil, &info); err != nil {
			return finish(err)
		}
		if !identifier.MatchString(info.ID) {
			return finish(fmt.Errorf("invalid agent identity"))
		}
		if addr, err := measurement.LoopbackAddress(info.ProbeAddress); err != nil || addr.Port() == 0 {
			return finish(fmt.Errorf("invalid probe endpoint"))
		}
		r.Agents = append(r.Agents, info)
	}
	if r.Agents[0].ID == r.Agents[1].ID {
		return finish(fmt.Errorf("agent identities must differ"))
	}
	// Record client-generated control IDs before POST so cleanup can cancel even if
	// the response is lost after the agent accepted the start request.
	defer func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		for i, s := range r.Sessions {
			if s.State == "running" || s.State == "stopping" {
				_ = c.request(cleanup, bases[i], "DELETE", "/v1/sessions/"+s.ID, nil, nil)
			}
		}
	}()
	for i, base := range bases {
		cfg := config
		if cfg.Mode == "udp" {
			cfg.Target = r.Agents[1-i].ProbeAddress
		}
		if err := cfg.Validate(); err != nil {
			return finish(err)
		}
		id, err := NewID()
		if err != nil {
			return finish(err)
		}
		r.Sessions = append(r.Sessions, Session{ID: id, AgentID: r.Agents[i].ID, State: "running"})
		var started Session
		if err := c.request(ctx, base, "POST", "/v1/sessions", StartRequest{ID: id, Config: cfg}, &started); err != nil {
			return finish(err)
		}
		if started.ID != id || started.AgentID != r.Agents[i].ID {
			return finish(fmt.Errorf("session identity mismatch"))
		}
	}
	var failures error
	for {
		terminal := 0
		for i, s := range r.Sessions {
			if s.State != "running" && s.State != "stopping" {
				terminal++
				continue
			}
			var current Session
			if err := c.request(ctx, bases[i], "GET", "/v1/sessions/"+s.ID, nil, &current); err != nil {
				return finish(err)
			}
			if current.ID != s.ID || current.AgentID != r.Agents[i].ID {
				return finish(fmt.Errorf("session identity mismatch"))
			}
			r.Sessions[i] = current
			if current.State != "running" && current.State != "stopping" {
				terminal++
				if current.Result == nil || current.State != "completed" {
					failures = errors.Join(failures, fmt.Errorf("agent %s session ended with %s", current.AgentID, current.State))
				}
			}
		}
		if terminal == 2 {
			return finish(failures)
		}
		timer := time.NewTimer(20 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return finish(ctx.Err())
		case <-timer.C:
		}
	}
}

func Export(base string, r Report) (string, error) {
	if !controlID.MatchString(r.ID) {
		return "", fmt.Errorf("invalid bench ID")
	}
	if err := os.MkdirAll(base, 0755); err != nil {
		return "", err
	}
	dest := filepath.Join(base, r.ID)
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		return "", fmt.Errorf("bench destination exists or is inaccessible")
	}
	tmp, err := os.MkdirTemp(base, ".bench-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmp)
	for _, session := range r.Sessions {
		if session.Result != nil && session.Result.ID != "" {
			result := *session.Result
			result.Samples = session.Samples
			if _, err := measurement.Export(filepath.Join(tmp, "agents"), result); err != nil {
				return "", err
			}
		}
	}
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(tmp, "aggregate.json"), append(data, '\n'), 0600); err != nil {
		return "", err
	}
	if err := os.Rename(tmp, dest); err != nil {
		return "", err
	}
	return dest, nil
}

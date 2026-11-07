package bench

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"server_sfp_sla/internal/measurement"
	"strings"
	"testing"
	"time"
)

const testToken = "nqa-test-token-for-loopback-only"

func testAgent(t *testing.T, id string) (*Agent, *httptest.Server) {
	t.Helper()
	a, err := NewAgent(context.Background(), id, "127.0.0.1:0", testToken)
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(a.Handler())
	t.Cleanup(func() { server.Close(); a.Close() })
	return a, server
}
func testClient(t *testing.T) *Client {
	t.Helper()
	c, err := NewClient(testToken)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(c.Close)
	return c
}
func udpConfig() measurement.Config {
	c := measurement.DefaultConfig()
	c.Mode = "udp"
	c.Interface = "loopback"
	c.Workers = 1
	c.Packets = 5
	c.Interval = 10 * time.Millisecond
	c.ProbeTimeout = 200 * time.Millisecond
	c.LossEvery = 0
	return c
}
func terminal(t *testing.T, c *Client, base, id string) Session {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	for {
		var s Session
		if err := c.request(ctx, base, "GET", "/v1/sessions/"+id, nil, &s); err != nil {
			t.Fatal(err)
		}
		if s.State != "running" && s.State != "stopping" {
			return s
		}
		select {
		case <-ctx.Done():
			t.Fatal("session did not terminate")
		case <-time.After(5 * time.Millisecond):
		}
	}
}

func TestTwoAgentUDPIntegration(t *testing.T) {
	_, first := testAgent(t, "first")
	_, second := testAgent(t, "second")
	client := testClient(t)
	report, err := client.Run(context.Background(), []string{first.URL, second.URL}, udpConfig())
	if err != nil {
		t.Fatal(err)
	}
	if report.Status != "completed" || report.Attempted != 10 || report.Received != 10 || report.Lost != 0 {
		t.Fatalf("unexpected report status=%s attempted=%d received=%d lost=%d", report.Status, report.Attempted, report.Received, report.Lost)
	}
	for _, session := range report.Sessions {
		r := session.Result
		if r.Provenance.Kind != "software" || r.Provenance.MeasurementKind != "round-trip" || r.RTTMean == nil || r.Sent != 5 || r.Unfinished != 0 {
			t.Fatal("incorrect timing/accounting provenance")
		}
		for _, s := range session.Samples {
			if s.Delay != nil || s.Jitter != nil || s.RTT == nil || *s.RTT < 0 {
				t.Fatal("software RTT mislabeled as one-way")
			}
		}
	}
	dest, err := Export(t.TempDir(), report)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dest, "aggregate.json"))
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte(testToken)) {
		t.Fatal("export leaked token")
	}
	for _, session := range report.Sessions {
		dir := filepath.Join(dest, "agents", session.Result.ID)
		if _, err := os.Stat(filepath.Join(dir, "rtt.csv")); err != nil {
			t.Fatal(err)
		}
		if _, err := os.Stat(filepath.Join(dir, "trace.csv")); !os.IsNotExist(err) {
			t.Fatal("RTT exported as SLA trace")
		}
	}
	if _, err := Export(filepath.Dir(dest), report); err == nil {
		t.Fatal("overwrote bench")
	}
}

func TestTwoAgentSyntheticIntegration(t *testing.T) {
	_, first := testAgent(t, "first")
	_, second := testAgent(t, "second")
	c := measurement.DefaultConfig()
	c.Packets = 40
	report, err := testClient(t).Run(context.Background(), []string{first.URL, second.URL}, c)
	if err != nil {
		t.Fatal(err)
	}
	if report.Lost != 4 || report.Received != 76 {
		t.Fatal("model accounting changed")
	}
	if report.Sessions[0].Result.Provenance.Kind != "synthetic" {
		t.Fatal("model mislabeled")
	}
}

func TestSessionCancellationAndConcurrencyBound(t *testing.T) {
	a, server := testAgent(t, "cancel")
	client := testClient(t)
	c := udpConfig()
	c.Target = a.Info().ProbeAddress
	c.Packets = 1000
	ids := []string{}
	for range 2 {
		id, err := NewID()
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id)
		if err := client.request(context.Background(), server.URL, "POST", "/v1/sessions", StartRequest{ID: id, Config: c}, nil); err != nil {
			t.Fatal(err)
		}
	}
	id, _ := NewID()
	err := client.request(context.Background(), server.URL, "POST", "/v1/sessions", StartRequest{ID: id, Config: c}, nil)
	if err == nil || !strings.Contains(err.Error(), "429") {
		t.Fatalf("unbounded sessions: %v", err)
	}
	for _, id := range ids {
		if err := client.request(context.Background(), server.URL, "DELETE", "/v1/sessions/"+id, nil, nil); err != nil {
			t.Fatal(err)
		}
		s := terminal(t, client, server.URL, id)
		if s.State != "canceled" || len(s.Result.Errors) != 1 {
			t.Fatal("cancel did not reach engine")
		}
	}
}

func TestControllerConnectionLossCancelsAcceptedSession(t *testing.T) {
	a, first := testAgent(t, "first")
	_, second := testAgent(t, "second")
	client := testClient(t)
	started := make(chan string, 1)
	wrapped := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			var body StartRequest
			// Copy the accepted control ID without changing the handler request body.
			raw := new(bytes.Buffer)
			_, _ = raw.ReadFrom(r.Body)
			r.Body = io.NopCloser(bytes.NewReader(raw.Bytes()))
			_ = json.Unmarshal(raw.Bytes(), &body)
			a.Handler().ServeHTTP(w, r)
			started <- body.ID
			return
		}
		a.Handler().ServeHTTP(w, r)
	}))
	defer wrapped.Close()
	done := make(chan error, 1)
	go func() {
		c := udpConfig()
		c.Packets = 100
		_, err := client.Run(context.Background(), []string{wrapped.URL, second.URL}, c)
		done <- err
	}()
	var id string
	select {
	case id = <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("controller did not start")
	}
	second.Close()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("connection loss not reported")
		}
	case <-time.After(4 * time.Second):
		t.Fatal("controller hung")
	}
	// Query the first agent through its original endpoint after controller cleanup.
	s := terminal(t, client, first.URL, id)
	if s.State != "canceled" {
		t.Fatalf("orphaned session: %s", s.State)
	}
}

func TestAgentRetentionBound(t *testing.T) {
	a, server := testAgent(t, "retention")
	client := testClient(t)
	c := measurement.DefaultConfig()
	c.Packets = 1
	for range 20 {
		id, _ := NewID()
		if err := client.request(context.Background(), server.URL, "POST", "/v1/sessions", StartRequest{ID: id, Config: c}, nil); err != nil {
			t.Fatal(err)
		}
		terminal(t, client, server.URL, id)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.sessions) != 16 {
		t.Fatalf("retained %d sessions", len(a.sessions))
	}
}

func TestAuthenticationAndEndpointValidation(t *testing.T) {
	_, server := testAgent(t, "auth")
	response, err := http.Get(server.URL + "/v1/info")
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != 401 {
		t.Fatal("unauthenticated access")
	}
	for _, address := range []string{"http://192.0.2.1:80", "http://localhost:80", "https://127.0.0.1:80", "http://user:pass@127.0.0.1:80", "http://127.0.0.1:80/?q=x", "http://127.0.0.1:80/"} {
		if err := endpoint(address); err == nil {
			t.Fatalf("unsafe endpoint accepted: %s", address)
		}
	}
	if _, err := NewAgent(context.Background(), "bad id", "127.0.0.1:0", testToken); err == nil {
		t.Fatal("invalid identity accepted")
	}
	if _, err := NewAgent(context.Background(), "id", "0.0.0.0:0", testToken); err == nil {
		t.Fatal("non-loopback bind accepted")
	}
	if _, err := NewClient("short"); err == nil {
		t.Fatal("short token accepted")
	}
}

func TestServeCancellationJoinsWorkers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	a, err := NewAgent(ctx, "serve", "127.0.0.1:0", testToken)
	if err != nil {
		t.Fatal(err)
	}
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	done := make(chan error, 1)
	go func() { done <- Serve(ctx, listener, a) }()
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("server shutdown hung")
	}
}
func TestControllerTimeoutCancelsBothAgents(t *testing.T) {
	a, first := testAgent(t, "timeout-a")
	b, second := testAgent(t, "timeout-b")
	client := testClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	c := udpConfig()
	c.Packets = 1000
	report, err := client.Run(ctx, []string{first.URL, second.URL}, c)
	if err == nil || report.Status != "timeout" {
		t.Fatalf("missing controller timeout: status=%s err=%v", report.Status, err)
	}
	for i, agent := range []*Agent{a, b} {
		agent.mu.Lock()
		ids := []string{}
		for id := range agent.sessions {
			ids = append(ids, id)
		}
		agent.mu.Unlock()
		for _, id := range ids {
			s := terminal(t, client, []string{first.URL, second.URL}[i], id)
			if s.State != "canceled" {
				t.Fatalf("session survived controller timeout: %s", s.State)
			}
		}
	}
}

func TestLostAcceptanceResponseStillCancelsSession(t *testing.T) {
	a, first := testAgent(t, "accepted-a")
	_, second := testAgent(t, "accepted-b")
	client := testClient(t)
	accepted := make(chan string, 1)
	wrapped := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			raw, _ := io.ReadAll(r.Body)
			var start StartRequest
			_ = json.Unmarshal(raw, &start)
			r.Body = io.NopCloser(bytes.NewReader(raw))
			accepted <- start.ID
			recorder := httptest.NewRecorder()
			a.Handler().ServeHTTP(recorder, r)
			conn, _, err := w.(http.Hijacker).Hijack()
			if err == nil {
				conn.Close()
			}
			return
		}
		a.Handler().ServeHTTP(w, r)
	}))
	defer wrapped.Close()
	c := udpConfig()
	c.Packets = 1000
	if _, err := client.Run(context.Background(), []string{wrapped.URL, second.URL}, c); err == nil {
		t.Fatal("lost response not reported")
	}
	var acceptedID string
	select {
	case acceptedID = <-accepted:
	default:
		t.Fatal("start request did not arrive")
	}
	if s := terminal(t, client, first.URL, acceptedID); s.State != "canceled" {
		t.Fatalf("accepted session orphaned: %s", s.State)
	}
}

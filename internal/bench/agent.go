package bench

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"runtime"
	"server_sfp_sla/internal/measurement"
	"sync"
	"time"
)

var identifier = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`)
var controlID = regexp.MustCompile(`^[a-f0-9]{32}$`)

func NewID() (string, error) {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes[:]), nil
}

type Info struct {
	ID           string                 `json:"agent_id"`
	ProbeAddress string                 `json:"probe_address"`
	Platform     string                 `json:"platform"`
	Software     string                 `json:"software_version"`
	Provenance   measurement.Provenance `json:"timestamp_provenance"`
}
type StartRequest struct {
	ID     string             `json:"control_session_id"`
	Config measurement.Config `json:"config"`
}
type Session struct {
	ID      string               `json:"control_session_id"`
	AgentID string               `json:"agent_id"`
	State   string               `json:"state"`
	Result  *measurement.Result  `json:"result,omitempty"`
	Samples []measurement.Sample `json:"samples,omitempty"`
}
type entry struct {
	session Session
	cancel  context.CancelFunc
	order   uint64
}

type Agent struct {
	info      Info
	token     string
	ctx       context.Context
	cancel    context.CancelFunc
	udp       *net.UDPConn
	mu        sync.Mutex
	sessions  map[string]*entry
	order     uint64
	closed    bool
	wg        sync.WaitGroup
	closeOnce sync.Once
}

// NewAgent binds only loopback and owns both probe and measurement workers.
func NewAgent(parent context.Context, id, probeAddress, token string) (*Agent, error) {
	if !identifier.MatchString(id) || len(token) < 16 || len(token) > 4096 {
		return nil, fmt.Errorf("agent ID must be 1..64 safe characters and token 16..4096 bytes")
	}
	addr, err := measurement.LoopbackAddress(probeAddress)
	if err != nil {
		return nil, err
	}
	udp, err := net.ListenUDP("udp", net.UDPAddrFromAddrPort(addr))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(parent)
	a := &Agent{info: Info{ID: id, ProbeAddress: udp.LocalAddr().String(), Platform: runtime.GOOS + "/" + runtime.GOARCH, Software: measurement.SoftwareVersion(), Provenance: measurement.SoftwareProvenance()}, token: token, ctx: ctx, cancel: cancel, udp: udp, sessions: make(map[string]*entry)}
	a.wg.Go(a.echo)
	return a, nil
}

func (a *Agent) Close() {
	a.closeOnce.Do(func() {
		a.mu.Lock()
		a.closed = true
		a.cancel()
		a.udp.Close()
		a.mu.Unlock()
		a.wg.Wait()
	})
}

func (a *Agent) echo() {
	defer a.udp.Close()
	stop := context.AfterFunc(a.ctx, func() { a.udp.Close() })
	defer stop()
	var packet [65]byte
	for {
		n, peer, err := a.udp.ReadFromUDP(packet[:])
		if err != nil {
			return
		}
		if !peer.IP.IsLoopback() || !measurement.ValidProbe(packet[:n], a.token) {
			continue
		}
		if err := a.udp.SetWriteDeadline(time.Now().Add(50 * time.Millisecond)); err != nil {
			return
		}
		if _, err := a.udp.WriteToUDP(packet[:n], peer); err != nil && a.ctx.Err() != nil {
			return
		}
	}
}

func (a *Agent) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/info", func(w http.ResponseWriter, r *http.Request) { respond(w, http.StatusOK, a.info) })
	mux.HandleFunc("POST /v1/sessions", a.start)
	mux.HandleFunc("GET /v1/sessions/{id}", a.get)
	mux.HandleFunc("DELETE /v1/sessions/{id}", a.stop)
	expected := sha256.Sum256([]byte("Bearer " + a.token))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		supplied := sha256.Sum256([]byte(r.Header.Get("Authorization")))
		if subtle.ConstantTimeCompare(expected[:], supplied[:]) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		mux.ServeHTTP(w, r)
	})
}

func respond(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func (a *Agent) start(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 65536)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	var req StartRequest
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid request JSON", 400)
		return
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		http.Error(w, "one JSON object required", 400)
		return
	}
	if !controlID.MatchString(req.ID) {
		http.Error(w, "invalid control session ID", 400)
		return
	}
	if err := req.Config.Validate(); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.closed || a.ctx.Err() != nil {
		http.Error(w, "agent stopping", 503)
		return
	}
	if _, exists := a.sessions[req.ID]; exists {
		http.Error(w, "session already exists", 409)
		return
	}
	running := 0
	for _, e := range a.sessions {
		if e.session.State == "running" || e.session.State == "stopping" {
			running++
		}
	}
	if running >= 2 {
		http.Error(w, "agent busy", 429)
		return
	}
	// Keep at most sixteen completed/running records; evict the oldest terminal one.
	if len(a.sessions) >= 16 {
		oldest := ""
		var order uint64
		for id, e := range a.sessions {
			if e.session.State != "running" && e.session.State != "stopping" && (oldest == "" || e.order < order) {
				oldest = id
				order = e.order
			}
		}
		if oldest == "" {
			http.Error(w, "session retention full", 429)
			return
		}
		delete(a.sessions, oldest)
	}
	ctx, cancel := context.WithCancel(a.ctx)
	a.order++
	e := &entry{session: Session{ID: req.ID, AgentID: a.info.ID, State: "running"}, cancel: cancel, order: a.order}
	a.sessions[req.ID] = e
	a.wg.Go(func() {
		defer cancel()
		var result measurement.Result
		var runErr error
		if req.Config.Mode == "udp" {
			result, runErr = measurement.RunUDP(ctx, req.Config, a.token)
		} else {
			result, runErr = measurement.Run(ctx, req.Config)
		}
		a.mu.Lock()
		defer a.mu.Unlock()
		if runErr != nil && result.Status == "" {
			result.Status = "failed"
			result.Errors = []measurement.SessionError{{Code: "engine", Message: runErr.Error()}}
		}
		state := result.Status
		if state == "" {
			state = "failed"
		}
		e.session.State = state
		e.session.Samples = result.Samples
		result.Samples = nil
		e.session.Result = &result
	})
	respond(w, http.StatusAccepted, e.session)
}

func (a *Agent) get(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	e, ok := a.sessions[r.PathValue("id")]
	if !ok {
		http.Error(w, "unknown session", 404)
		return
	}
	respond(w, 200, e.session)
}
func (a *Agent) stop(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	e, ok := a.sessions[r.PathValue("id")]
	if !ok {
		http.Error(w, "unknown session", 404)
		return
	}
	if e.session.State == "running" {
		e.session.State = "stopping"
		e.cancel()
	}
	respond(w, 200, e.session)
}

// Serve shuts down HTTP connections and joins agent workers when canceled.
func Serve(ctx context.Context, listener net.Listener, a *Agent) error {
	defer a.Close()
	server := &http.Server{Handler: a.Handler(), ReadHeaderTimeout: 2 * time.Second, ReadTimeout: 3 * time.Second, WriteTimeout: 3 * time.Second, IdleTimeout: 5 * time.Second, MaxHeaderBytes: 8192}
	done := make(chan error, 1)
	go func() { done <- server.Serve(listener) }()
	select {
	case err := <-done:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	case <-ctx.Done():
		cleanup, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if err := server.Shutdown(cleanup); err != nil {
			server.Close()
		}
		err := <-done
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
}

func (a *Agent) Info() Info { return a.info }

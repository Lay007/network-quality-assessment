package measurement

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

const ConfigVersion = 1

// Config bounds memory, CPU concurrency and session duration before work starts.
type Config struct {
	Version   int           `json:"version"`
	Mode      string        `json:"mode"`
	Interface string        `json:"interface"`
	Packets   int           `json:"packets"`
	Workers   int           `json:"workers"`
	Interval  time.Duration `json:"interval_ns"`
	Timeout   time.Duration `json:"timeout_ns"`
	Seed      uint64        `json:"seed"`
	LossEvery int           `json:"loss_every"`
	Delay     time.Duration `json:"delay_ns"`
	Variation time.Duration `json:"variation_ns"`
}

func DefaultConfig() Config {
	return Config{Version: 1, Mode: "synthetic", Interface: "synthetic", Packets: 1000, Workers: 2, Interval: time.Millisecond, Timeout: 10 * time.Second, Seed: 1, LossEvery: 20, Delay: 300 * time.Microsecond, Variation: 10 * time.Microsecond}
}

func (c Config) Validate() error {
	if c.Version != ConfigVersion {
		return errors.New("unsupported configuration version")
	}
	if c.Mode != "synthetic" || c.Interface != "synthetic" {
		return errors.New("this engine requires synthetic mode and interface")
	}
	if c.Packets < 1 || c.Packets > 10000 || c.Workers < 1 || c.Workers > 16 {
		return errors.New("packets must be 1..10000; workers must be 1..16")
	}
	if c.Interval < time.Microsecond || c.Interval > time.Second || c.Timeout < time.Millisecond || c.Timeout > time.Minute {
		return errors.New("interval must be 1us..1s; timeout must be 1ms..1m")
	}
	if c.LossEvery < 0 || c.Delay < 0 || c.Delay > time.Second || c.Variation < 0 || c.Variation > c.Delay {
		return errors.New("invalid loss or delay model")
	}
	return nil
}

type Provenance struct {
	Kind            string `json:"kind"`
	Source          string `json:"source"`
	ClockDomain     string `json:"clock_domain"`
	Synchronization string `json:"synchronization"`
	MeasurementKind string `json:"measurement_kind"`
	Note            string `json:"note"`
}

// Sample uses integer nanoseconds; nil receive/delay fields represent loss.
type Sample struct {
	PacketID int    `json:"packet_id"`
	TX       int64  `json:"tx_timestamp_ns"`
	RX       *int64 `json:"rx_timestamp_ns"`
	Delay    *int64 `json:"one_way_delay_ns"`
	Jitter   *int64 `json:"jitter_ns"`
	Lost     bool   `json:"lost"`
}

type SessionError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Result struct {
	SchemaVersion int            `json:"schema_version"`
	ID            string         `json:"session_id"`
	Software      string         `json:"software_version"`
	Platform      string         `json:"platform"`
	Config        Config         `json:"config"`
	ConfigHash    string         `json:"config_sha256"`
	Provenance    Provenance     `json:"timestamp_provenance"`
	Status        string         `json:"status"`
	Start         time.Time      `json:"started_at"`
	End           time.Time      `json:"ended_at"`
	Attempted     int            `json:"packets_attempted"`
	Received      int            `json:"packets_received"`
	Lost          int            `json:"packets_lost"`
	LossPercent   float64        `json:"loss_percent"`
	Errors        []SessionError `json:"errors"`
	Samples       []Sample       `json:"-"`
}

func SoftwareVersion() string {
	version := "development"
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				version = setting.Value
			}
			if setting.Key == "vcs.modified" && setting.Value == "true" {
				version += "+dirty"
			}
		}
	}
	return version + " (" + runtime.Version() + ")"
}

func newResult(c Config) (Result, error) {
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return Result{}, fmt.Errorf("session identifier: %w", err)
	}
	encoded, err := json.Marshal(c)
	if err != nil {
		return Result{}, err
	}
	hash := sha256.Sum256(encoded)
	return Result{SchemaVersion: 1, ID: hex.EncodeToString(id[:]), Software: SoftwareVersion(), Platform: runtime.GOOS + "/" + runtime.GOARCH, Config: c, ConfigHash: hex.EncodeToString(hash[:]), Status: "running", Start: time.Now().UTC(), Errors: []SessionError{}, Samples: make([]Sample, c.Packets)}, nil
}

// Run models a trace without sleeping or sending packets. Context controls actual
// computation time; modeled trace duration is independent of execution duration.
func Run(ctx context.Context, c Config) (Result, error) {
	if err := c.Validate(); err != nil {
		return Result{}, fmt.Errorf("configuration: %w", err)
	}
	r, err := newResult(c)
	if err != nil {
		return r, err
	}
	r.Provenance = Provenance{Kind: "synthetic", Source: "deterministic-integer-model-v1", ClockDomain: "modeled-relative-nanoseconds", Synchronization: "not-applicable", MeasurementKind: "synthetic-one-way", Note: "Simulated timestamps; no packets, host clock accuracy or hardware evidence."}
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()
	jobs := make(chan int, c.Workers*2)
	ready := make([]bool, c.Packets)
	var wg sync.WaitGroup
	for range c.Workers {
		wg.Go(func() {
			for i := range jobs {
				if ctx.Err() != nil {
					return
				}
				tx := int64(i) * int64(c.Interval)
				sample := Sample{PacketID: i, TX: tx, Lost: c.LossEvery > 0 && (i+1)%c.LossEvery == 0}
				if !sample.Lost {
					// A per-packet integer mixer makes results independent of worker scheduling.
					x := uint64(i) + c.Seed + 0x9e3779b97f4a7c15
					x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
					x = (x ^ (x >> 27)) * 0x94d049bb133111eb
					x ^= x >> 31
					delay := int64(c.Delay) + int64(x%uint64(2*c.Variation+1)) - int64(c.Variation)
					rx := tx + delay
					sample.RX = &rx
					sample.Delay = &delay
				}
				r.Samples[i] = sample
				ready[i] = true
			}
		})
	}
enqueue:
	for i := range c.Packets {
		select {
		case <-ctx.Done():
			break enqueue
		case jobs <- i:
		}
	}
	close(jobs)
	wg.Wait()
	// Emit only the contiguous completed prefix on cancellation. Unattempted packets
	// are not counted as network loss and ordering stays compatible with the trace.
	n := 0
	for n < len(ready) && ready[n] {
		n++
	}
	r.Samples = r.Samples[:n]
	var previous *int64
	for i := range r.Samples {
		sample := &r.Samples[i]
		r.Attempted++
		if sample.Lost {
			r.Lost++
			continue
		}
		r.Received++
		jitter := int64(0)
		if previous != nil {
			jitter = *sample.Delay - *previous
		}
		sample.Jitter = &jitter
		previous = sample.Delay
	}
	if r.Attempted > 0 {
		r.LossPercent = 100 * float64(r.Lost) / float64(r.Attempted)
	}
	r.Status = "completed"
	r.End = time.Now().UTC()
	if err := ctx.Err(); err != nil {
		code := "canceled"
		if errors.Is(err, context.DeadlineExceeded) {
			code = "timeout"
		}
		r.Status = code
		r.Errors = append(r.Errors, SessionError{Code: code, Message: err.Error()})
		return r, err
	}
	return r, nil
}

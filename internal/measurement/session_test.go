package measurement

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestSyntheticRepeatabilityAndLoss(t *testing.T) {
	c := DefaultConfig()
	c.Packets = 100
	first, err := Run(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	c.Workers = 7
	second, err := Run(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	if first.ID == second.ID || !reflect.DeepEqual(first.Samples, second.Samples) {
		t.Fatal("session uniqueness or deterministic trace failed")
	}
	if first.Attempted != 100 || first.Lost != 5 || first.Received != 95 || first.LossPercent != 5 || first.Status != "completed" {
		t.Fatalf("invalid counters: %+v", first)
	}
	var previous *int64
	for i, s := range first.Samples {
		if s.PacketID != i || s.TX != int64(i)*int64(c.Interval) {
			t.Fatal("invalid ordering")
		}
		if s.Lost {
			if s.RX != nil || s.Delay != nil || s.Jitter != nil {
				t.Fatal("loss has timestamps")
			}
			continue
		}
		if *s.RX-s.TX != *s.Delay {
			t.Fatal("invalid delay")
		}
		want := int64(0)
		if previous != nil {
			want = *s.Delay - *previous
		}
		if *s.Jitter != want {
			t.Fatal("jitter must span received packets")
		}
		previous = s.Delay
	}
}

func TestConfigurationBounds(t *testing.T) {
	for _, modify := range []func(*Config){
		func(c *Config) { c.Version = 2 }, func(c *Config) { c.Mode = "hardware" }, func(c *Config) { c.Interface = "eth0" },
		func(c *Config) { c.Packets = 0 }, func(c *Config) { c.Packets = 10001 }, func(c *Config) { c.Workers = 0 }, func(c *Config) { c.Workers = 17 },
		func(c *Config) { c.Interval = 0 }, func(c *Config) { c.Timeout = 0 }, func(c *Config) { c.LossEvery = -1 }, func(c *Config) { c.Variation = c.Delay + 1 },
	} {
		c := DefaultConfig()
		modify(&c)
		if _, err := Run(context.Background(), c); err == nil {
			t.Fatal("invalid configuration accepted")
		}
	}
}

func TestCancellationAndTimeoutAreExportable(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	timed, stop := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))
	defer stop()
	for _, test := range []struct {
		ctx    context.Context
		status string
		err    error
	}{{ctx, "canceled", context.Canceled}, {timed, "timeout", context.DeadlineExceeded}} {
		r, err := Run(test.ctx, DefaultConfig())
		if !errors.Is(err, test.err) || r.Status != test.status || r.Attempted != 0 || r.Lost != 0 || len(r.Errors) != 1 {
			t.Fatalf("invalid cancellation: status=%s,err=%v", r.Status, err)
		}
		if _, err := Export(t.TempDir(), r); err != nil {
			t.Fatal(err)
		}
	}
}

func TestExportRoundTripAndNoOverwrite(t *testing.T) {
	c := DefaultConfig()
	c.Packets = 40
	r, err := Run(context.Background(), c)
	if err != nil {
		t.Fatal(err)
	}
	base := t.TempDir()
	dest, err := Export(base, r)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dest, "metadata.json"))
	if err != nil {
		t.Fatal(err)
	}
	var metadata Result
	if err := json.Unmarshal(data, &metadata); err != nil {
		t.Fatal(err)
	}
	if metadata.ID != r.ID || metadata.Provenance.Kind != "synthetic" || metadata.ConfigHash != r.ConfigHash {
		t.Fatal("metadata provenance changed")
	}
	f, err := os.Open(filepath.Join(dest, "trace.csv"))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 41 || rows[20][2] != "" || rows[20][5] != "1" || len(rows[0]) != 8 {
		t.Fatal("trace schema or missing timestamp changed")
	}
	if _, err := Export(base, r); err == nil {
		t.Fatal("overwrote session")
	}
	r.ID = "../invalid"
	if _, err := Export(base, r); err == nil {
		t.Fatal("accepted path as session ID")
	}
}

func BenchmarkSynthetic1000(b *testing.B) {
	c := DefaultConfig()
	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		if _, err := Run(context.Background(), c); err != nil {
			b.Fatal(err)
		}
	}
}

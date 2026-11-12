package main

import (
	"bytes"
	"context"
	"sync"
	"testing"
	"time"
)

type readyWriter struct {
	ready chan struct{}
	once  sync.Once
}

func (w *readyWriter) Write(p []byte) (int, error) {
	w.once.Do(func() { close(w.ready) })
	return len(p), nil
}

func TestServeCLIShutdown(t *testing.T) {
	t.Setenv("NQA_BENCH_TOKEN", "serve-cli-test-token-only")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := &readyWriter{ready: make(chan struct{})}
	done := make(chan error, 1)
	go func() {
		done <- run(ctx, []string{"serve", "--id", "cli", "--listen", "127.0.0.1:0", "--probe", "127.0.0.1:0"}, out)
	}()
	select {
	case <-out.ready:
	case err := <-done:
		t.Fatal(err)
	case <-time.After(2 * time.Second):
		t.Fatal("agent failed to start")
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("agent did not stop")
	}
}
func TestServeCLIRejectsUnsafeConfiguration(t *testing.T) {
	t.Setenv("NQA_BENCH_TOKEN", "serve-cli-test-token-only")
	for _, args := range [][]string{{"serve", "--listen", "0.0.0.0:80"}, {"serve", "extra"}, {"serve", "--id", "bad id"}, {"--mode", "udp", "--target", "192.0.2.1:9"}} {
		if err := run(context.Background(), args, &bytes.Buffer{}); err == nil {
			t.Fatal("unsafe agent arguments accepted")
		}
	}
}

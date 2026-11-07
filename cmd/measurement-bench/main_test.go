package main

import (
	"bytes"
	"context"
	"net/http/httptest"
	"os"
	"path/filepath"
	"server_sfp_sla/internal/bench"
	"strings"
	"testing"
)

func TestBenchCLIIntegration(t *testing.T) {
	const token = "cli-test-loopback-token-only"
	t.Setenv("NQA_BENCH_TOKEN", token)
	bases := []string{}
	for _, id := range []string{"cli-a", "cli-b"} {
		a, err := bench.NewAgent(context.Background(), id, "127.0.0.1:0", token)
		if err != nil {
			t.Fatal(err)
		}
		server := httptest.NewServer(a.Handler())
		t.Cleanup(func() { server.Close(); a.Close() })
		bases = append(bases, server.URL)
	}
	var out bytes.Buffer
	if err := run(context.Background(), []string{"--agents", strings.Join(bases, ","), "--packets", "3", "--output", t.TempDir()}, &out); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(strings.TrimSpace(out.String()), "aggregate.json")); err != nil {
		t.Fatal(err)
	}
}

func TestBenchCLIInvalidArguments(t *testing.T) {
	t.Setenv("NQA_BENCH_TOKEN", "")
	for _, args := range [][]string{{"--timeout", "0s"}, {"unexpected"}, {"--timeout", "31s"}, {}} {
		if err := run(context.Background(), args, &bytes.Buffer{}); err == nil {
			t.Fatal("invalid controller arguments accepted")
		}
	}
}

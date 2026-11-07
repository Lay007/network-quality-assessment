package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIProducesPackage(t *testing.T) {
	var out bytes.Buffer
	base := t.TempDir()
	if err := run(context.Background(), []string{"--packets", "10", "--output", base}, &out); err != nil {
		t.Fatal(err)
	}
	dest := strings.TrimSpace(out.String())
	for _, name := range []string{"metadata.json", "trace.csv"} {
		if _, err := os.Stat(filepath.Join(dest, name)); err != nil {
			t.Fatal(err)
		}
	}
}

func TestCLIRejectsBadArgumentsAndConfiguration(t *testing.T) {
	for _, args := range [][]string{{"--packets", "0"}, {"unexpected"}, {"--workers", "17"}} {
		if err := run(context.Background(), args, &bytes.Buffer{}); err == nil {
			t.Fatal("accepted invalid flags")
		}
	}
	file := filepath.Join(t.TempDir(), "config.json")
	for _, body := range []string{`{"unknown":1}`, `{} {}`, `{"version":99}`} {
		if err := os.WriteFile(file, []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
		if err := run(context.Background(), []string{"--config", file}, &bytes.Buffer{}); err == nil {
			t.Fatal("accepted invalid JSON configuration")
		}
	}
	if err := run(context.Background(), []string{"--config", file, "--packets", "10"}, &bytes.Buffer{}); err == nil {
		t.Fatal("ambiguous flag precedence")
	}
}

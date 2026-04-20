package main

import "testing"

func TestEnvOrDefault(t *testing.T) {
	t.Setenv("SFP_SLA_SAMPLE_ENV", "")
	if got := envOrDefault("SFP_SLA_SAMPLE_ENV", "fallback"); got != "fallback" {
		t.Fatalf("envOrDefault fallback mismatch: got %q", got)
	}

	t.Setenv("SFP_SLA_SAMPLE_ENV", "value")
	if got := envOrDefault("SFP_SLA_SAMPLE_ENV", "fallback"); got != "value" {
		t.Fatalf("envOrDefault value mismatch: got %q", got)
	}
}

func TestDBDSNLocalSocket(t *testing.T) {
	t.Setenv("SFP_SLA_DB_USER", "demo")
	t.Setenv("SFP_SLA_DB_PASSWORD", "secret")
	t.Setenv("SFP_SLA_DB_NAME", "metrics")
	t.Setenv("SFP_SLA_DB_ADDR", "")

	got := dbDSN()
	want := "demo:secret@/metrics"
	if got != want {
		t.Fatalf("dbDSN local mismatch: got %q want %q", got, want)
	}
}

func TestDBDSNTCP(t *testing.T) {
	t.Setenv("SFP_SLA_DB_USER", "demo")
	t.Setenv("SFP_SLA_DB_PASSWORD", "secret")
	t.Setenv("SFP_SLA_DB_NAME", "metrics")
	t.Setenv("SFP_SLA_DB_ADDR", "127.0.0.1:3306")

	got := dbDSN()
	want := "demo:secret@tcp(127.0.0.1:3306)/metrics"
	if got != want {
		t.Fatalf("dbDSN tcp mismatch: got %q want %q", got, want)
	}
}

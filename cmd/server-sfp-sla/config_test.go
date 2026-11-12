//go:build linux

package main

import (
	"github.com/go-sql-driver/mysql"
	"testing"
)

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
	want := "demo:secret@/metrics?readTimeout=5s&timeout=5s&writeTimeout=5s"
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
	want := "demo:secret@tcp(127.0.0.1:3306)/metrics?readTimeout=5s&timeout=5s&writeTimeout=5s"
	if got != want {
		t.Fatalf("dbDSN tcp mismatch: got %q want %q", got, want)
	}
}

func TestDBConfigurationEscapesName(t *testing.T) {
	t.Setenv("SFP_SLA_DB_NAME", "metrics/name")
	cfg, err := mysql.ParseDSN(dbDSN())
	if err != nil {
		t.Fatal(err)
	}
	if cfg.DBName != "metrics/name" || cfg.Timeout <= 0 || cfg.ReadTimeout <= 0 || cfg.WriteTimeout <= 0 {
		t.Fatalf("invalid DB configuration: name or timeouts")
	}
}

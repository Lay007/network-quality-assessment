//go:build linux

package main

import (
	"database/sql"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

const (
	defaultDBUser     = "sfp_user"
	defaultDBPassword = ""
	defaultDBName     = "server_sfp_sla"
)

var verboseLogs = os.Getenv("SFP_SLA_VERBOSE") == "1"

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func dbDSN() string {
	cfg := mysql.NewConfig()
	cfg.User = envOrDefault("SFP_SLA_DB_USER", defaultDBUser)
	cfg.Passwd = envOrDefault("SFP_SLA_DB_PASSWORD", defaultDBPassword)
	cfg.DBName = envOrDefault("SFP_SLA_DB_NAME", defaultDBName)
	if addr := os.Getenv("SFP_SLA_DB_ADDR"); addr != "" {
		cfg.Net, cfg.Addr = "tcp", addr
	}
	cfg.Timeout, cfg.ReadTimeout, cfg.WriteTimeout = 5*time.Second, 5*time.Second, 5*time.Second
	return cfg.FormatDSN()
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", dbDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(4)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}

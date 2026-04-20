package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
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
	user := envOrDefault("SFP_SLA_DB_USER", defaultDBUser)
	password := envOrDefault("SFP_SLA_DB_PASSWORD", defaultDBPassword)
	name := envOrDefault("SFP_SLA_DB_NAME", defaultDBName)
	addr := os.Getenv("SFP_SLA_DB_ADDR")

	if addr == "" {
		return fmt.Sprintf("%s:%s@/%s", user, password, name)
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, addr, name)
}

func openDB() (*sql.DB, error) {
	return sql.Open("mysql", dbDSN())
}

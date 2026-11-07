package measurement

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

var sessionID = regexp.MustCompile(`^[a-f0-9]{32}$`)

// Export publishes a complete package by renaming a private temporary directory.
// Existing session directories are never overwritten.
func Export(base string, r Result) (string, error) {
	if !sessionID.MatchString(r.ID) {
		return "", fmt.Errorf("invalid session identifier")
	}
	if err := os.MkdirAll(base, 0755); err != nil {
		return "", err
	}
	dest := filepath.Join(base, r.ID)
	if _, err := os.Stat(dest); !os.IsNotExist(err) {
		return "", fmt.Errorf("session destination already exists or is inaccessible")
	}
	tmp, err := os.MkdirTemp(base, ".session-")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmp)
	metadata, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(tmp, "metadata.json"), append(metadata, '\n'), 0600); err != nil {
		return "", err
	}
	f, err := os.OpenFile(filepath.Join(tmp, "trace.csv"), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if err != nil {
		return "", err
	}
	writeErr := writeTrace(f, r.Samples)
	closeErr := f.Close()
	if writeErr != nil {
		return "", writeErr
	}
	if closeErr != nil {
		return "", closeErr
	}
	if err := os.Rename(tmp, dest); err != nil {
		return "", err
	}
	return dest, nil
}

func integer(v *int64) string {
	if v == nil {
		return ""
	}
	return strconv.FormatInt(*v, 10)
}
func writeTrace(f *os.File, samples []Sample) error {
	w := csv.NewWriter(f)
	if err := w.Write([]string{"packet_id", "tx_timestamp_ns", "rx_timestamp_ns", "one_way_delay_ns", "jitter_ns", "lost", "burst_id", "scenario"}); err != nil {
		return err
	}
	for _, s := range samples {
		lost := "0"
		if s.Lost {
			lost = "1"
		}
		if err := w.Write([]string{strconv.Itoa(s.PacketID), strconv.FormatInt(s.TX, 10), integer(s.RX), integer(s.Delay), integer(s.Jitter), lost, "0", "synthetic-go-v1"}); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

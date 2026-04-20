package main

import (
	"bytes"
	"encoding/binary"
	"math"
	"testing"
)

func TestChecksumReturnsAllOnesForEmptyPayload(t *testing.T) {
	if got := checksum(nil); got != 0xFFFF {
		t.Fatalf("checksum(nil) = %#x, want 0xFFFF", got)
	}
}

func TestIPHeaderChecksumMatchesManualCalculation(t *testing.T) {
	h := iphdr{
		vhl:   0x45,
		tos:   0,
		iplen: 60,
		id:    0x1234,
		off:   0,
		ttl:   64,
		proto: 17,
		src:   [4]byte{10, 0, 0, 1},
		dst:   [4]byte{10, 0, 0, 2},
	}

	h.checksum()
	if h.csum == 0 {
		t.Fatal("computed header checksum must be non-zero")
	}

	hNoChecksum := h
	hNoChecksum.csum = 0
	var b bytes.Buffer
	if err := binary.Write(&b, binary.BigEndian, hNoChecksum); err != nil {
		t.Fatalf("binary.Write failed: %v", err)
	}

	want := checksum(b.Bytes())
	if h.csum != want {
		t.Fatalf("header checksum mismatch: got %#x want %#x", h.csum, want)
	}
}

func TestDelayToMicroseconds(t *testing.T) {
	got := delayToMicroseconds(1 << 32)
	if math.Abs(float64(got-1_000_000)) > 0.01 {
		t.Fatalf("delayToMicroseconds mismatch: got %f want 1000000", got)
	}
}

func TestJitterToMicroseconds(t *testing.T) {
	if got := jitterToMicroseconds(0); got != 0 {
		t.Fatalf("jitterToMicroseconds(0) = %f, want 0", got)
	}

	got := jitterToMicroseconds(float32(1 << 31))
	if math.Abs(float64(got-500_000)) > 0.01 {
		t.Fatalf("jitterToMicroseconds mismatch: got %f want 500000", got)
	}
}

func TestDelayAverageAndJitter(t *testing.T) {
	var st testSLA

	if got := st.getDelayAvg(100); got != 100 {
		t.Fatalf("first getDelayAvg mismatch: got %d want 100", got)
	}
	if got := st.getDelayAvg(300); got != 200 {
		t.Fatalf("second getDelayAvg mismatch: got %d want 200", got)
	}
	if got := st.getJitter(); got != 200 {
		t.Fatalf("getJitter mismatch: got %f want 200", got)
	}
}

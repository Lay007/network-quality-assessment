package measurement

import (
	"context"
	"net"
	"testing"
	"time"
)

func localUDPConfig(target string) Config {
	c := DefaultConfig()
	c.Mode = "udp"
	c.Target = target
	c.Interface = "loopback"
	c.Workers = 1
	c.Packets = 2
	c.Interval = 10 * time.Millisecond
	c.ProbeTimeout = 10 * time.Millisecond
	return c
}

func TestProbeAuthentication(t *testing.T) {
	token := "test-only-loopback-token"
	id := "0123456789abcdef0123456789abcdef"
	packet, err := ProbeMessage(id, 1, token)
	if err != nil {
		t.Fatal(err)
	}
	if !ValidProbe(packet, token) || ValidProbe(packet, "wrong") || ValidProbe(packet[:63], token) {
		t.Fatal("invalid authentication")
	}
	packet[24] ^= 1
	if ValidProbe(packet, token) {
		t.Fatal("accepted tampered packet")
	}
	if _, err := ProbeMessage("bad", 0, token); err == nil {
		t.Fatal("accepted invalid ID")
	}
}

func TestUDPLossIsObservationTimeout(t *testing.T) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := localUDPConfig(conn.LocalAddr().String())
	r, err := RunUDP(context.Background(), c, "test-only-loopback-token")
	if err != nil {
		t.Fatal(err)
	}
	if r.Attempted != 2 || r.Lost != 2 || r.Sent != 2 || r.Unfinished != 0 || r.RTTMean != nil {
		t.Fatal("invalid no-response accounting")
	}
	for _, s := range r.Samples {
		if !s.Lost || s.RTT != nil || s.Delay != nil || s.RX != nil {
			t.Fatal("loss has timing values")
		}
	}
}

func TestUDPCancellationDoesNotInventLoss(t *testing.T) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := localUDPConfig(conn.LocalAddr().String())
	c.ProbeTimeout = 200 * time.Millisecond
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	r, err := RunUDP(ctx, c, "test-only-loopback-token")
	if err == nil || r.Status != "timeout" || r.Lost != 0 || r.Attempted != 0 || r.Unfinished != r.Sent {
		t.Fatalf("cancellation counted as network loss: status=%s attempted=%d lost=%d", r.Status, r.Attempted, r.Lost)
	}
}

func TestUDPValidation(t *testing.T) {
	for _, modify := range []func(*Config){func(c *Config) { c.Target = "192.0.2.1:80" }, func(c *Config) { c.Target = "localhost:80" }, func(c *Config) { c.Workers = 2 }, func(c *Config) { c.Interval = time.Microsecond }, func(c *Config) { c.ProbeTimeout = time.Second }} {
		c := localUDPConfig("127.0.0.1:1")
		modify(&c)
		if _, err := RunUDP(context.Background(), c, "test-only-loopback-token"); err == nil {
			t.Fatal("accepted invalid UDP configuration")
		}
	}
}

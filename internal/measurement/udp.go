package measurement

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"time"
)

// TimestampEvent is a future capture boundary, not a device/register contract.
// A hardware provider must supply its own clock domain, capture point and evidence.
type TimestampEvent struct {
	PacketID     uint64
	Nanoseconds  int64
	ClockDomain  string
	CapturePoint string
}
type TimestampInput interface {
	Read(context.Context) (TimestampEvent, error)
	Provenance() Provenance
	Close() error
}

func LoopbackAddress(value string) (netip.AddrPort, error) {
	address, err := netip.ParseAddrPort(value)
	if err != nil || !address.Addr().IsLoopback() {
		return netip.AddrPort{}, fmt.Errorf("numeric loopback address and port required")
	}
	return address, nil
}

func validateUDP(c Config) error {
	target, err := LoopbackAddress(c.Target)
	if err != nil || target.Port() == 0 {
		return fmt.Errorf("UDP target must be a numeric loopback endpoint with nonzero port")
	}
	if c.Interface != "loopback" || c.Workers != 1 || c.Packets > 1000 || c.Interval < 10*time.Millisecond || c.ProbeTimeout < time.Millisecond || c.ProbeTimeout > 200*time.Millisecond {
		return fmt.Errorf("UDP requires loopback interface, one worker, <=1000 packets, interval >=10ms, probe timeout 1..200ms")
	}
	return nil
}

func SoftwareProvenance() Provenance {
	return Provenance{Kind: "software", Source: "Go time.Now and monotonic time.Since", ClockDomain: "local-host-monotonic-elapsed", Synchronization: "unsynchronized; remote clock unused", MeasurementKind: "round-trip", Note: "Application send-to-reply RTT includes scheduling, socket queues and echo processing. Wall timestamps are diagnostic only. No one-way or FPGA accuracy claim."}
}

// Probe messages are exactly 64 bytes. HMAC prevents the echo service from being
// an unauthenticated reflector; one outstanding request bounds receiver state.
func ProbeMessage(id string, sequence int, token string) ([]byte, error) {
	decoded, err := hex.DecodeString(id)
	if err != nil || len(decoded) != 16 {
		return nil, fmt.Errorf("invalid session ID")
	}
	packet := make([]byte, 64)
	copy(packet, "NQABEN01")
	copy(packet[8:24], decoded)
	binary.BigEndian.PutUint64(packet[24:32], uint64(sequence))
	if _, err := rand.Read(packet[32:48]); err != nil {
		return nil, err
	}
	mac := hmac.New(sha256.New, []byte(token))
	mac.Write(packet[:48])
	copy(packet[48:], mac.Sum(nil)[:16])
	return packet, nil
}

func ValidProbe(packet []byte, token string) bool {
	if len(packet) != 64 || string(packet[:8]) != "NQABEN01" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(token))
	mac.Write(packet[:48])
	return hmac.Equal(packet[48:], mac.Sum(nil)[:16])
}

// RunUDP uses the authenticated loopback bench wire protocol. The token stays
// outside configuration and exported metadata. A session has one UDP connection.
func RunUDP(ctx context.Context, c Config, token string) (Result, error) {
	if err := c.Validate(); err != nil {
		return Result{}, err
	}
	if c.Mode != "udp" || len(token) < 16 {
		return Result{}, fmt.Errorf("UDP mode and token of at least 16 bytes required")
	}
	return runUDPToken(ctx, c, token)
}

func runUDPToken(ctx context.Context, c Config, token string) (Result, error) {
	r, err := newResult(c)
	if err != nil {
		return r, err
	}
	r.Samples = nil
	r.Provenance = SoftwareProvenance()
	ctx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()
	var total int64
	finish := func(err error) (Result, error) {
		r.End = time.Now().UTC()
		r.Status = "completed"
		r.Unfinished = r.Sent - r.Attempted
		if r.Received > 0 {
			mean := float64(total) / float64(r.Received)
			r.RTTMean = &mean
		}
		if r.Attempted > 0 {
			r.LossPercent = 100 * float64(r.Lost) / float64(r.Attempted)
		}
		if err != nil {
			code := "network"
			if errors.Is(err, context.Canceled) {
				code = "canceled"
			}
			if errors.Is(err, context.DeadlineExceeded) {
				code = "timeout"
			}
			r.Status = code
			r.Errors = append(r.Errors, SessionError{Code: code, Message: err.Error()})
		}
		return r, err
	}
	if err := ctx.Err(); err != nil {
		return finish(err)
	}
	target, _ := net.ResolveUDPAddr("udp", c.Target)
	conn, err := net.DialUDP("udp", nil, target)
	if err != nil {
		return finish(err)
	}
	defer conn.Close()
	stop := context.AfterFunc(ctx, func() { conn.Close() })
	defer stop()
	var previous *int64
	for i := range c.Packets {
		if err := ctx.Err(); err != nil {
			return finish(err)
		}
		if i > 0 {
			timer := time.NewTimer(c.Interval)
			select {
			case <-ctx.Done():
				timer.Stop()
				return finish(ctx.Err())
			case <-timer.C:
			}
		}
		packet, err := ProbeMessage(r.ID, i, token)
		if err != nil {
			return finish(err)
		}
		tx := time.Now()
		deadline := tx.Add(c.ProbeTimeout)
		if end, ok := ctx.Deadline(); ok && end.Before(deadline) {
			deadline = end
		}
		if err := conn.SetDeadline(deadline); err != nil {
			return finish(err)
		}
		n, err := conn.Write(packet)
		if err != nil {
			return finish(err)
		}
		if n != len(packet) {
			return finish(fmt.Errorf("short UDP write"))
		}
		r.Sent++
		sample := Sample{PacketID: i, TX: tx.UnixNano()}
		var buffer [65]byte
		for {
			n, err := conn.Read(buffer[:])
			rx := time.Now()
			if ctx.Err() != nil {
				return finish(ctx.Err())
			}
			if err != nil {
				var netErr net.Error
				if errors.As(err, &netErr) && netErr.Timeout() {
					if end, ok := ctx.Deadline(); ok && !time.Now().Before(end) {
						return finish(context.DeadlineExceeded)
					}
					sample.Lost = true
					break
				}
				return finish(err)
			}
			if n != len(packet) || !bytes.Equal(packet, buffer[:n]) {
				continue
			}
			// time.Sub retains the monotonic component; wall-clock steps do not alter RTT.
			elapsed := rx.Sub(tx).Nanoseconds()
			wall := rx.UnixNano()
			sample.RX = &wall
			sample.RTT = &elapsed
			delta := int64(0)
			if previous != nil {
				delta = elapsed - *previous
			}
			sample.RTTDelta = &delta
			previous = sample.RTT
			total += elapsed
			if r.RTTMin == nil || elapsed < *r.RTTMin {
				v := elapsed
				r.RTTMin = &v
			}
			if r.RTTMax == nil || elapsed > *r.RTTMax {
				v := elapsed
				r.RTTMax = &v
			}
			break
		}
		r.Attempted++
		if sample.Lost {
			r.Lost++
		} else {
			r.Received++
		}
		r.Samples = append(r.Samples, sample)
	}
	if r.Received > 0 {
		mean := float64(total) / float64(r.Received)
		r.RTTMean = &mean
	}
	return finish(nil)
}

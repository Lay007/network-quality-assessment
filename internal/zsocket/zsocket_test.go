package zsocket

import (
	"syscall"
	"testing"

	"github.com/newtools/zsocket/nettypes"
)

func TestPacketOffsetIsAligned(t *testing.T) {
	if PacketOffset()%TPACKET_ALIGNMENT != 0 {
		t.Fatalf("PacketOffset must be aligned to %d, got %d", TPACKET_ALIGNMENT, PacketOffset())
	}
	if PacketOffset() <= 0 {
		t.Fatalf("PacketOffset must be positive, got %d", PacketOffset())
	}
}

func TestErrnoErr(t *testing.T) {
	if got := errnoErr(0); got != nil {
		t.Fatalf("errnoErr(0) = %v, want nil", got)
	}
	if got := errnoErr(syscall.EAGAIN); got == nil || got.Error() != "try again" {
		t.Fatalf("errnoErr(EAGAIN) mismatch: got %v", got)
	}
	if got := errnoErr(syscall.EINVAL); got == nil || got.Error() != "invalid argument" {
		t.Fatalf("errnoErr(EINVAL) mismatch: got %v", got)
	}
	if got := errnoErr(syscall.ENOENT); got == nil || got.Error() != "no such file or directory" {
		t.Fatalf("errnoErr(ENOENT) mismatch: got %v", got)
	}
	if got := errnoErr(syscall.EPERM); got != syscall.EPERM {
		t.Fatalf("errnoErr(EPERM) mismatch: got %v want %v", got, syscall.EPERM)
	}
}

func TestCopyFx(t *testing.T) {
	dst := make([]byte, 4)
	src := []byte{1, 2, 3, 4}

	if got := copyFx(dst, src, 4); got != 4 {
		t.Fatalf("copyFx length mismatch: got %d want 4", got)
	}
	for i := range dst {
		if dst[i] != src[i] {
			t.Fatalf("copyFx dst[%d] = %d, want %d", i, dst[i], src[i])
		}
	}
}

func TestCalculateLargestFrame(t *testing.T) {
	if got := calculateLargestFrame(uint(MINIMUM_FRAME_SIZE)); got != uint(MINIMUM_FRAME_SIZE) {
		t.Fatalf("calculateLargestFrame(MINIMUM_FRAME_SIZE) = %d, want %d", got, MINIMUM_FRAME_SIZE)
	}
	if got := calculateLargestFrame(3000); got != uint(MINIMUM_FRAME_SIZE) {
		t.Fatalf("calculateLargestFrame(3000) = %d, want %d", got, MINIMUM_FRAME_SIZE)
	}
	if got := calculateLargestFrame(5000); got != 4096 {
		t.Fatalf("calculateLargestFrame(5000) = %d, want 4096", got)
	}
}

func TestNewZSocketValidationErrorsBeforeSocketOpen(t *testing.T) {
	var ethType nettypes.EthType

	if _, err := NewZSocket(1, ENABLE_RX, uint(MINIMUM_FRAME_SIZE-1), 16, ethType); err == nil {
		t.Fatal("expected validation error for too small frame size")
	}
	if _, err := NewZSocket(1, ENABLE_RX, uint(MINIMUM_FRAME_SIZE+1), 16, ethType); err == nil {
		t.Fatal("expected validation error for non power-of-two frame size")
	}
	if _, err := NewZSocket(1, ENABLE_RX, uint(MINIMUM_FRAME_SIZE), 8, ethType); err == nil {
		t.Fatal("expected validation error for maxTotalFrames < 16")
	}
	if _, err := NewZSocket(1, ENABLE_RX, uint(MINIMUM_FRAME_SIZE), 17, ethType); err == nil {
		t.Fatal("expected validation error for maxTotalFrames not multiple of 8")
	}
}

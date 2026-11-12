//go:build linux

package zsocket

import (
	"context"
	"errors"
	"net"
	"sync"
	"syscall"
	"testing"
	"time"

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

func TestTXRejectsInvalidBuffers(t *testing.T) {
	zs := &ZSocket{txEnabled: true, txFrameSize: 64}
	for _, length := range []uint16{4, 65} {
		if _, err := zs.WriteToBuffer([]byte{1, 2, 3}, length); err == nil {
			t.Fatal("accepted out-of-bounds write")
		}
		if _, err := zs.CopyToBuffer([]byte{1, 2, 3}, length, copyFx); err == nil {
			t.Fatal("accepted out-of-bounds copy")
		}
	}
}

func TestCloseUnmapsAndIsIdempotent(t *testing.T) {
	fds := []int{0, 0}
	if err := syscall.Pipe(fds); err != nil {
		t.Fatal(err)
	}
	defer syscall.Close(fds[1])
	data, err := syscall.Mmap(-1, 0, 4096, syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_PRIVATE|syscall.MAP_ANON)
	if err != nil {
		t.Fatal(err)
	}
	zs := &ZSocket{socket: fds[0], raw: data}
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			if err := zs.Close(); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if zs.raw != nil {
		t.Fatal("mapping retained")
	}
	if _, err := zs.WriteToBuffer(nil, 0); !errors.Is(err, net.ErrClosed) {
		t.Fatalf("write after close: %v", err)
	}
	if err := syscall.Close(fds[0]); err != syscall.EBADF {
		t.Fatalf("descriptor not closed: %v", err)
	}
}

func TestListenerCancellationAndBounds(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	zs := &ZSocket{rxEnabled: true}
	if err := zs.ListenContext(ctx, func(*nettypes.Frame, uint16, uint16) {}); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	// A truncated kernel frame must be rejected before slicing it.
	rf := &ringFrame{raw: make([]byte, MINIMUM_FRAME_SIZE)}
	rf.setTpSnapLen(MINIMUM_FRAME_SIZE)
	rf.txSet()
	rf.raw[0] = 1
	zs.rxFrames = []*ringFrame{rf}
	zs.frameNum = 1
	// tp_mac is placed at the end of the mapped frame to force an invalid extent.
	rf.raw[_TP_MAC_START] = 0xff
	rf.raw[_TP_MAC_START+1] = 0x7f
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := zs.ListenContext(ctx, func(*nettypes.Frame, uint16, uint16) { t.Fatal("callback on invalid frame") }); err == nil {
		t.Fatal("accepted invalid extent")
	}
}

func TestNewSocketRejectsOptionsBeforePrivileges(t *testing.T) {
	for _, option := range []int{0, 8} {
		if _, err := NewZSocket(1, option, 2048, 16, nettypes.EthType{}); err == nil {
			t.Fatal("accepted invalid options")
		}
	}
}

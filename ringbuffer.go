package ringbuffer

import (
	"math"
	"syscall"
)

const (
	// 64 kb
	defaultCacheSize = 64 * 1024
)

// RingBuffer lock-free cache for a read-write goroutine.
type RingBuffer struct {
	p   []byte
	r   int
	w   int
	cap int
}

// NewRingBuffer .
func NewRingBuffer(cap int) *RingBuffer {
	if cap > math.MaxUint || cap <= 0 {
		cap = defaultCacheSize
	}

	if (cap & (cap - 1)) != 0 {
		cap = adjust(cap)
	}

	return &RingBuffer{
		p:   make([]byte, 0, cap),
		r:   0,
		w:   0,
		cap: cap,
	}
}

func adjust(n int) int {
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	return n + 1
}

// CopyFromFd .
func (r *RingBuffer) CopyFromFd(fd int) (int, error) {
	if r.full() {
		return 0, syscall.EAGAIN
	}

	if r.r > r.w {
		return syscall.Read(fd, r.p[r.w:r.r])
	}

	bs := [][]byte{
		r.p[r.w:],
		r.p[:r.r],
	}
	return Readv(fd, bs)
}

// Write .
func (r *RingBuffer) Write(p []byte) error {
	if r.full() {
		return syscall.EAGAIN
	}

	l := len(p)
	if l <= 0 {
		return nil
	}

	// return syscall.EAGAIN when the buffer capacity is insufficient.
	if r.writeableSize() < l {
		return syscall.EAGAIN
	}

	n := copy(p, r.p[r.w:])
	if n < l {
		copy(p[n:], r.p)
	}

	r.w += l
	return nil
}

// Read .
func (r *RingBuffer) Read(p []byte) (int, error) {
	if r.empty() {
		return 0, syscall.EAGAIN
	}

	l := len(p)
	if l <= 0 {
		return 0, nil
	}

	if r.r < r.w {
		return copy(p, r.p[r.r:r.w]), nil
	}

	readableSize := min(r.readableSize(), l)
	n := copy(p, r.p[r.r:])
	if n < l {
		copy(p[n:], r.p[:readableSize-n])
	}

	return readableSize, nil
}

// Bytes .
func (r *RingBuffer) Bytes() []byte {
	p := make([]byte, r.cap)
	i := copy(p, r.p[r.w:])
	copy(p[i:], r.p[:i])
	return p
}

// Len .
func (r *RingBuffer) Len() int {
	return r.readableSize()
}

// WriteString .
func (r *RingBuffer) WriteString(s string) error {
	return r.Write([]byte(s))
}

// IsEmpty .
func (r *RingBuffer) IsEmpty() bool {
	return r.readableSize() == 0
}

// Release .
func (r *RingBuffer) Release(n int) {
	if r.readableSize() < n {
		r.r = r.w
		return
	}

	r.r += n
}

// Clear .
func (r *RingBuffer) Clear() {
	r.r, r.w, r.cap, r.p = 0, 0, 0, nil
}

func (r *RingBuffer) full() bool {
	if r.r == (r.w+1)%r.cap {
		return true
	}

	return false
}

func (r *RingBuffer) empty() bool {
	return r.r == r.w
}

func (r *RingBuffer) writeableSize() int {
	if r.w > r.r {
		return r.cap + r.r - r.w
	}
	return r.r - r.w
}

func (r *RingBuffer) readableSize() int {
	return r.cap - r.writeableSize()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

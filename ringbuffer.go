package ringbuffer

// RingBuffer lock-free cache for a read-write goroutine.
type RingBuffer struct {
	data []byte
	r    uint64
	w    uint64
	size uint64
}

// NewRingBuffer .
func NewRingBuffer(size uint64) *RingBuffer {
	return &RingBuffer{
		data: make([]byte, 0, size),
		r:    0,
		w:    0,
		size: 0,
	}
}

// CopyFromFd .
func (r *RingBuffer) CopyFromFd(fd int) (int, error) {
	return 0, nil
}

func (r *RingBuffer) Write(bytes []byte) error {
	//TODO implement me
	panic("implement me")
}

func (r *RingBuffer) Read(bytes []byte) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RingBuffer) Bytes() []byte {
	//TODO implement me
	panic("implement me")
}

func (r *RingBuffer) Len() int {
	//TODO implement me
	panic("implement me")
}

func (r *RingBuffer) WriteString(s string) error {
	//TODO implement me
	panic("implement me")
}

func (r *RingBuffer) IsEmpty() bool {
	//TODO implement me
	panic("implement me")
}

func (r *RingBuffer) Release(n int) {
	//TODO implement me
	panic("implement me")
}

func (r *RingBuffer) Clear() {
	//TODO implement me
	panic("implement me")
}

package ringbuffer

// Buffer net in out p buffer
type Buffer interface {
	// CopyFromFd copy fd to buffer
	CopyFromFd(fd int) (int, error)
	// Write bytes to buffer, return syscall.EAGAIN when the buffer capacity is insufficient
	Write(bytes []byte) (int, error)
	// Read buffer p to  bytes
	Read(bytes []byte) (int, error)
	// Bytes will return buffer bytes
	Bytes() []byte
	// Len will return buffer readable length
	Len() int
	// WriteString string to buffer
	WriteString(s string) (int, error)
	// IsEmpty will return true if buffer len is zero
	IsEmpty() bool
	// Release will release length n buffer p
	Release(n int)
	// Clear will clear buffer
	Clear()
}

package ringbuffer

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"syscall"
	"testing"
)

func TestNewRingBuffer(t *testing.T) {
	type args struct {
		cap int
	}
	tests := []struct {
		name    string
		args    args
		wantCap int
	}{
		{
			name:    "normal",
			args:    args{cap: 1},
			wantCap: 1,
		},
		{
			name:    "normal cap != 2^n",
			args:    args{cap: 3},
			wantCap: 4,
		},
		{
			name:    "normal cap = maxCacheSize",
			args:    args{cap: maxCacheSize},
			wantCap: maxCacheSize,
		},
		{
			name:    "normal",
			args:    args{cap: 1},
			wantCap: 1,
		},
		{
			name:    "cap <= 0",
			args:    args{cap: -1},
			wantCap: defaultCacheSize,
		},
		{
			name:    "cap > maxCacheSize",
			args:    args{cap: maxCacheSize + 1},
			wantCap: maxCacheSize,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRingBuffer(tt.args.cap); !reflect.DeepEqual(got.Cap(), tt.wantCap) {
				t.Errorf("NewRingBuffer() cap= %v, want %v", got, tt.wantCap)
			}
		})
	}
}

func TestRingBuffer_RW(t *testing.T) {
	var (
		n   int
		err error
	)

	ringBuffer := NewRingBuffer(10)
	assert.Equal(t, true, ringBuffer.IsEmpty())

	n, err = ringBuffer.WriteString("helloworld")
	assert.Nil(t, err)
	assert.Equal(t, 10, n)

	assert.Equal(t, "helloworld", string(ringBuffer.Bytes()))
	assert.Equal(t, 10, ringBuffer.Len())
	assert.Equal(t, 16, ringBuffer.Cap())

	n, err = ringBuffer.WriteString("helloworld")
	assert.Nil(t, err)
	assert.Equal(t, 6, n)

	assert.Equal(t, "helloworldhellow", string(ringBuffer.Bytes()))
	assert.Equal(t, 16, ringBuffer.Len())

	n, err = ringBuffer.WriteString("test")
	assert.Equal(t, err, syscall.EAGAIN)
	assert.Equal(t, 0, n)

	assert.Equal(t, "helloworldhellow", string(ringBuffer.Bytes()))
	assert.Equal(t, 16, ringBuffer.Len())

	ringBuffer.Release(5)

	assert.Equal(t, "worldhellow", string(ringBuffer.Bytes()))
	assert.Equal(t, 11, ringBuffer.Len())

	n, err = ringBuffer.WriteString("123456")
	assert.Equal(t, 5, n)
	assert.Nil(t, err)

	assert.Equal(t, "worldhellow12345", string(ringBuffer.Bytes()))
	assert.Equal(t, 16, ringBuffer.Len())

	ringBuffer.Release(10)

	assert.Equal(t, "w12345", string(ringBuffer.Bytes()))
	assert.Equal(t, 6, ringBuffer.Len())

	n, err = ringBuffer.WriteString("789")
	assert.Equal(t, 3, n)
	assert.Nil(t, err)

	assert.Equal(t, "w12345789", string(ringBuffer.Bytes()))
	assert.Equal(t, 9, ringBuffer.Len())

	ringBuffer.Release(ringBuffer.Len() + 1)

	var emptyValue []byte
	assert.Equal(t, emptyValue, ringBuffer.Bytes())
	assert.Equal(t, 0, ringBuffer.Len())

	ringBuffer.Release(1)

	assert.Equal(t, 0, ringBuffer.Len())

	n, err = ringBuffer.Write(nil)
	assert.Equal(t, 0, n)
	assert.Nil(t, err)
}

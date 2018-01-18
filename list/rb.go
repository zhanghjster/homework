package list

import (
	"io"

	"github.com/pkg/errors"
)

var (
	ErrFull = errors.New("cache if full")
)

// Ring Buffer implemented by array, no overwriting
type RingBuffer struct {
	data     []byte
	start    int
	end      int
	length   int
	capacity int
}

func NewRingBuffer(capacity int) *RingBuffer {
	rb := &RingBuffer{
		data:     make([]byte, capacity),
		capacity: capacity,
	}
	return rb
}

func (b *RingBuffer) Read(p []byte) (n int, err error) {
	if b.length <= 0 {
		return n, io.EOF
	}

	for n < len(p) && b.length > 0 {
		var e int
		if b.end >= b.start {
			e = b.end
		} else {
			e = b.capacity
		}

		w := copy(p[n:], b.data[b.start:e])
		b.start = (b.start + w) % b.capacity

		n += w
		b.length -= w
	}

	return n, nil
}

func (b *RingBuffer) Write(p []byte) (n int, err error) {
	if b.length > b.capacity {
		return n, ErrFull
	}

	for n < len(p) && b.length < b.capacity {
		var e int
		if b.end >= b.start {
			e = b.capacity
		} else {
			e = b.start
		}

		w := copy(b.data[b.end:e], p[n:])
		b.end = (b.end + w) % b.capacity

		n += w
		b.length += w
	}

	return
}

func (b *RingBuffer) Capacity() int {
	return b.capacity
}

func (b *RingBuffer) Length() int {
	return b.length
}

func (b *RingBuffer) Start() int {
	return b.start
}

func (b *RingBuffer) End() int {
	return b.end
}

func (b *RingBuffer) Reset() {
	b.length, b.end, b.start = 0, 0, 0
}

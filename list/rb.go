package list

import (
	"github.com/pkg/errors"
	"io"
)

var (
	ErrOutOfRange = errors.New("out of range")
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
	if b.length < 1 {
		return n, io.EOF
	}

	for n < len(p) && b.length > 0 {
		var l int
		if b.end >= b.start {
			l = copy(p[n:], b.data[b.start:b.end])
		} else {
			l = copy(p[n:], b.data[b.start:b.capacity])
		}

		b.start = (b.start + l) % b.capacity

		n += l
		b.length -= l
	}

	return n, nil
}

func (b *RingBuffer) Write(p []byte) (n int, err error) {
	if b.length+len(p) > b.capacity {
		return n, ErrOutOfRange
	}

	for n < len(p) {
		l := copy(b.data[b.end:], p[n:])

		b.end = (b.end + l) % b.capacity

		n += l
		b.length += l
	}

	return
}

func (b *RingBuffer) Cap() int {
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

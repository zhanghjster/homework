package list

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestRingBuffer_ReadWrite(t *testing.T) {

	var src = []byte{'a','b','c','d','e'}
	b := NewRingBuffer(7)

	// write
	n, err := b.Write(src)
	assert.Nil(t, err)
	assert.Equal(t, len(src), n)
	assert.Equal(t, b.Start(), 0)
	assert.Equal(t, b.End(), len(src))

	// read
	var dst = make([]byte, 8)
	n, err = b.Read(dst)
	assert.Nil(t, err)
	assert.Equal(t, n, len(src))
	assert.Equal(t, b.Start(), len(src))
	assert.Equal(t, b.End(), len(src))

	// write
	_, err = b.Write([]byte{'f','g','h'})
	assert.Nil(t, err)

	// read
	_, err = b.Read(dst[n:])
	assert.Nil(t, err)
	assert.Equal(t, b.Start(), 1)
	assert.Equal(t, b.End(), 1)
	assert.Equal(t, int32(dst[len(dst)-1]), int32('h'))
}

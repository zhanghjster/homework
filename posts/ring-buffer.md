---
title: Ring Buffer
date: 2017-09-19 00:30:43
tags: 
	- ring buffer 
	- circular buffer
	- struct
categories:
	- struct and algrithm
---

### 简述

环形缓冲区(ring buffer), 又称圆形队列（circular queue), 循环缓冲区(cyclic buffer), 圆形缓冲区(circular buffer)。

它适合于实现明确大小的FIFO缓冲区，通常由一个数组实现，start和end两个索引来表示数据的开始和结束，length表示当前元素个数，capacity表示缓冲区容量。当有元素删除或插入时只需要移动start和end，其他元素不需要移动存储位置。

<!-- more -->

### 工作过程

缓冲区数据结构为一个capacity长度的数组，初始start、end、length为0

* 有新元素插入时，插入到end索引的位置，end后移，length增加，当end在超过数组的最大索引时，end为0，当length等于capacity时表示缓冲区满
* 读取元素时，从start开始读取，start增加，length减少，当start超过数组索引最大值时，start为0,当length为0是表示已经全部读出

### 代码

```
package ringBuffer

import (
	"github.com/pkg/errors"
	"io"
)

var (
	ErrFull = errors.New("cache is full")
)

// Ring Buffer implemented by array, no overwriting
type RingBuffer struct {
	data     []byte
	start    int
	end      int
	length   int
	capacity int
}

func New(capacity int) *RingBuffer {
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

```


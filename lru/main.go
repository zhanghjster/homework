package lru

import (
	"container/list"
)

// 原理：
// 	当cache满了时丢掉最久未被使用的单元，需要追踪每个单元最近被使用的时间
// 	使用双向链表，最新使用的在header，最久被使用的在tail
type Entry struct {
	key, value interface{}
}

type Cache struct {
	size    int
	dll     *list.List
	entries map[interface{}]*list.Element
}

func New(size int) *Cache {
	return &Cache{
		size:    size,
		entries: make(map[interface{}]*list.Element),
	}
}

func (c *Cache) Get(key interface{}) (interface{}, bool) {
	if e, ok := c.entries[key]; ok {
		node := e.Value.(*Entry)
		c.dll.MoveToFront(e)
		return node.value, ok
	}

	return nil, false
}

func (c *Cache) Set(key, value interface{}) {
	if _, ok := c.Get(key); ok {
		return
	}

	var e *list.Element
	if c.size <= c.dll.Len() && c.dll.Back() != nil {
		entry := c.dll.Back().Value.(*Entry)

		delete(c.entries, entry.key)

		entry.key = key
		entry.value = value

		c.dll.MoveToFront(e)
	} else {
		e = c.dll.PushFront(&Entry{key, value})
		c.entries[key] = e
	}
}

func (c *Cache) Peek(key interface{}) (interface{}, bool) {
	if e, ok := c.entries[key]; ok {
		return e.Value.(*Entry).value, ok
	}
	return nil, false
}

func (c *Cache) Pure() {
	for k := range c.entries {
		delete(c.entries, k)
	}
	c.dll.Init()
}

func (c *Cache) Contain(key interface{}) bool {
	_, ok := c.entries[key]
	return ok
}

func (c *Cache) Len() int {
	return c.dll.Len()
}

package lru

import (
	"container/list"
)

// 原理：
// 	当cache满了时丢掉最久未被使用的单元，需要追踪每个单元最近被使用的时间
// 	使用双向链表，最新使用的在header，最久被使用的在tail

// 非多线程安全
type Entry struct {
	key, value interface{}
}

type Cache struct {
	size    int
	marks   *list.List
	entries map[interface{}]*list.Element
}

func New(size int) *Cache {
	if size <= 0 {
		panic("cache size <= 0")
	}

	return &Cache{
		size:    size,
		entries: make(map[interface{}]*list.Element),
	}
}

// 获取内容，并将其移动到最近被使用
func (c *Cache) Get(key interface{}) (interface{}, bool) {
	if e, ok := c.entries[key]; ok {
		c.marks.MoveToFront(e)
		return e.Value.(*Entry).value, ok
	}

	return nil, false
}

func (c *Cache) Set(key, value interface{}) {
	// 已经存在就移动到header
	if _, ok := c.Get(key); ok {
		return
	}

	// 如果满了就移动到最久未使用的
	if c.IsFull() {
		e := c.marks.Back()
		c.marks.Remove(e)
		delete(c.entries, e.Value.(*Entry).key)
	}

	// 将新的元素放到header
	e := c.marks.PushFront(&Entry{key, value})

	c.entries[key] = e
}

// 返回缓存内容但不移动到最新
func (c *Cache) Peek(key interface{}) (interface{}, bool) {
	if e, ok := c.entries[key]; ok {
		return e.Value.(*Entry).value, ok
	}
	return nil, false
}

// 清空
func (c *Cache) Pure() {
	for k := range c.entries {
		delete(c.entries, k)
	}
	c.marks.Init()
}

// 检查key是否存在
func (c *Cache) Contain(key interface{}) bool {
	_, ok := c.entries[key]
	return ok
}

func (c *Cache) Len() int {
	return c.marks.Len()
}

func (c *Cache) IsFull() bool {
	return c.size <= c.Len()
}

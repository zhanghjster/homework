package list

import "sync"

type Item struct {
	list  *DoublyLinkList
	pre   *Item
	next  *Item
	Value interface{}
}

func (i *Item) Next() *Item {
	if i.next == &i.list.root {
		return nil
	}

	return i.next
}

type DoublyLinkList struct {
	mux    sync.RWMutex
	root   Item
	length int
}

func NewDoublyLinkList() *DoublyLinkList { return new(DoublyLinkList).Init() }

func (l *DoublyLinkList) Init() *DoublyLinkList {
	l.mux.Lock()
	defer l.mux.Unlock()

	l.root.list = l
	l.root.pre = &l.root
	l.root.next = &l.root
	l.length = 0
	return l
}

func (l *DoublyLinkList) PushFront(v interface{}) (item *Item) {
	return l.InsertAfter(v, &l.root)
}

func (l *DoublyLinkList) PushTail(v interface{}) (item *Item) {
	return l.InsertBefore(v, &l.root)
}

func (l *DoublyLinkList) InsertBefore(v interface{}, e *Item) (item *Item) {
	return l.InsertValue(v, e.pre)
}

func (l *DoublyLinkList) InsertAfter(v interface{}, e *Item) (item *Item) {
	return l.InsertValue(v, e)
}

func (l *DoublyLinkList) InsertValue(v interface{}, dst *Item) *Item {
	return l.Insert(&Item{Value: v}, dst)
}

func (l *DoublyLinkList) Insert(item, dst *Item) *Item {
	l.mux.Lock()
	defer l.mux.Unlock()

	l.length++
	item.list = l
	dst.next.pre = item
	item.next = dst.next
	item.pre = dst
	dst.next = item
	return item
}

func (l *DoublyLinkList) Delete(item *Item) {
	l.mux.Lock()
	defer l.mux.Unlock()

	l.length--
	if item.pre != nil {
		item.pre.next = item.next
	}

	if item.next != nil {
		item.next.pre = item.pre
	}

	item.pre = nil
	item.next = nil
	item.list = nil
}

func (l *DoublyLinkList) Move(src, dst *Item) {
	l.Delete(src)
	l.Insert(src, dst)
}

func (l *DoublyLinkList) MoveAfter(src, dst *Item) {
	l.Move(src, dst)
}

func (l *DoublyLinkList) MoveBefore(src, dst *Item) {
	l.Move(src, dst.pre)
}

func (l *DoublyLinkList) Front() *Item {
	l.mux.RLock()
	defer l.mux.RUnlock()

	if l.length == 0 {
		return nil
	}

	return l.root.next
}

func (l *DoublyLinkList) Tail() *Item {
	l.mux.RLock()
	defer l.mux.RUnlock()

	if l.length == 0 {
		return nil
	}

	return l.root.pre
}

func (l *DoublyLinkList) Len() int {
	l.mux.RLock()
	defer l.mux.RUnlock()

	return l.length
}

package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDoublyLinkList(t *testing.T) {
	list := NewDoublyLinkList()
	var value = 10
	item := list.PushFront(value)
	assert.Equal(t, value, item.Value.(int))
	front := list.Front()
	assert.ObjectsAreEqual(item, front)

	list.Init()

	item = list.PushFront(10)
	list.PushFront(11)
	list.Delete(item)
	front = list.Front()
	assert.Equal(t, 11, front.Value.(int))

	list.Init()
	var nums = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	for _, v := range nums {
		list.PushTail(v)
	}

	assert.Equal(t, len(nums), list.Len())

	var i = 0
	for item := list.Front(); item != nil; item = item.Next() {
		assert.Equal(t, nums[i], item.Value.(int))
		i++
	}
}

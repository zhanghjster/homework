package ds

import (
	"fmt"
	"testing"
)

func TestBTree_Insert(t *testing.T) {
	tree := NewBTree(3)
	//for _, k := range []int{10, 20, 5, 6, 12, 30, 7, 17} {
	for _, k := range []int{10, 20, 5, 6, 12} {
		fmt.Printf("insert %d\n", k)
		tree.Insert(k)
	}

	tree.Traverse()
}

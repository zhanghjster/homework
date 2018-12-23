package main

import (
	"github.com/zhanghjster/homework/ds"
)

func main() {
	tree := ds.NewBTree(3)
	//var n = []int{10, 20, 5, 6, 12, 30, 7, 17, 1, 60, 20, 35, 25, 40, 80, 100, 3, 4, 5, 2}
	var n = []int{1, 3, 7, 10, 11, 13, 14, 15, 18, 16, 19, 24, 25, 26, 21, 4, 5, 20, 22, 2, 17, 12, 6}
	for _, k := range n {
		tree.Insert(k)
	}

	tree.Delete(6)
	tree.Traverse()

	tree.Delete(13)
	tree.Traverse()

	tree.Delete(7)
	tree.Traverse()

	tree.Delete(4)
	tree.Traverse()
	tree.Delete(2)
	tree.Traverse()
	tree.Delete(16)
	tree.Traverse()
}

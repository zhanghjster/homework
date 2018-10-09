package main

import (
	"github.com/zhanghjster/homework/ds"
)

func main() {
	tree := ds.NewBTree(3)
	for _, k := range []int{10, 20, 5, 6, 12, 30, 7, 17, 1, 60, 20, 35, 25, 40, 80, 100, 3, 4, 5} {
		tree.Insert(k)
	}

	tree.Traverse()
}

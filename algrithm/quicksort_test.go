package algrithm

import (
	"testing"
	"fmt"
)

func TestQuickSort(t *testing.T) {
	var a = []int{9,4,1,5,2,8,7,3,5}
	QuickSort(a)
	for i, v := range a {
		fmt.Println(i, v)
	}
}
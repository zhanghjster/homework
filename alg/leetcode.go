package alg

import (
	"container/list"
	"fmt"
)

func TwoSum(a []int, t int) (int, int) {
	var m = make(map[int]int)
	for i, v := range a {
		if j, ok := m[v]; ok {
			return j, i
		} else {
			m[t-v] = i
		}
	}

	return -1, -1
}

func AddTwo(l1, l2 *list.List) *list.List {
	var l = list.New()

	var e1, e2 = l1.Back(), l2.Back()
	for d := 0; e1 != nil || e2 != nil; d /= 10 {
		var v1, v2 int
		if e1 != nil {
			v1 = e1.Value.(int)
			e1 = e1.Prev()
		}
		if e2 != nil {
			v2 = e2.Value.(int)
			e2 = e2.Prev()
		}

		d += v1 + v2
		l.PushFront(d % 10)
	}

	return l
}

func IntToList(v int) *list.List {
	var l = list.New()
	for ; v > 0; v /= 10 {
		l.PushFront(v % 10)
	}
	return l
}

func ListToInt(l *list.List) int {
	var v int
	for e := l.Front(); e != nil; e = e.Next() {
		v = 10*v + e.Value.(int)
	}
	return v
}

func PrintIntList(l *list.List) {
	if l == nil {
		return
	}

	for e := l.Front(); e != nil; e = e.Next() {
		print(e.Value.(int))
	}

	fmt.Println()
}

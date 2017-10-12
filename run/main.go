package main

import (
	"fmt"
	"unsafe"
	"container/list"
)

func main() {
	var a = intToList(1023)
	var b = intToList(2)

	var c = AddTwo(a, b)

	printIntList(c)

	// show a b
	printIntList(a)
	printIntList(b)
}

func AddTwo(l1, l2 *list.List) *list.List {
	var l = list.New()

	var e1, e2 = l1.Back(), l2.Back()
	for d := 0;e1 != nil || e2 !=nil; d /= 10 {
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
		l.PushFront(d%10)
	}

	return l
}

func intToList(v int) *list.List {
	var l = list.New()
	for ; v > 0; v /= 10{
		l.PushFront(v%10)
	}
	return l
}

func printIntList(l *list.List) {
	if l == nil { return }

	for e := l.Front(); e != nil; e = e.Next() {
		print(e.Value.(int))
	}

	fmt.Println()
}

func InspectSlice(slice []int) {

	// 数组的地址
	addr := unsafe.Pointer(&slice)
	fmt.Printf("%v, %v\n", addr, &slice[0])

	// len字段的地址
	lenAddr := uintptr(addr) + uintptr(8)
	capAddr := uintptr(addr) + uintptr(16)

	lenPtr := (*int)(unsafe.Pointer(lenAddr))
	capPtr := (*int)(unsafe.Pointer(capAddr))

	// a = (*uintptr)(addr)add转换成指向uintptr类型数据的指针
	// ptr = *(*uintptr)(addr))取a这个指针所指向的值, 为slice底层数据结构的第一个字段ptr的值
	// unsafe.Pointer(ptr)将ptr转换成指针
	arrPtr := unsafe.Pointer(*(*uintptr)(addr))

	fmt.Printf("Slice Addr[%p], Len Addr[0x%x] Cap Addr[0x%x]\n", addr, lenAddr, capAddr)
	fmt.Printf("Slice length[%d] Cap[%d]\n", *lenPtr, *capPtr)

	for i :=0; i < *lenPtr; i++ {
		p := unsafe.Pointer(uintptr(arrPtr) + uintptr(i) * uintptr(unsafe.Sizeof(int(0))))
		fmt.Printf("[%d] %p %d\n", i, (*int)(p), *(*int)(p))
	}
}
package main

import (
	"fmt"
	"unsafe"
	"container/list"
)

func main() {
	var a, b  = 234, 72

	var la, lb = list.New(), list.New()
	for ;a > 0; a /= 10 {
		la.PushFront(a%10)
	}
	for ;b >0; b /= 10 {
		lb.PushFront(b%10)
	}

	// show a b

	for ha := la.Front(); ha!=nil; ha = ha.Next() {
	}



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
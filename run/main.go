package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var i int64 = (1 << 9) | (1 << 17)
	fmt.Printf("%b\n", i)
	println(int8(i >> 8))

	var a uint8 = 10
	var b uint8 = 5
	fmt.Printf("%b,%b\n", a, b)
	fmt.Printf("%b\n", a^b)
	fmt.Printf("%b\n", a&^b)

	var aa= []int{1, 2, 3, 4, 5, 6}
	var bb= aa[1:5:5]
	fmt.Printf("%v\n", bb)
	p := unsafe.Pointer(uintptr(unsafe.Pointer(&bb[0])) + uintptr(4*unsafe.Sizeof(int(0))))
	*(*int)(p) = 0
	fmt.Printf("%v\n%v\n", aa, bb)

	var slice= []int{1, 2, 3, 4, 5, 6}

	InspectSlice(slice)

	//InspectSlice(slice[2:3])

	var newSlice = slice[2:4:4]

	InspectSlice(newSlice)

	//newSlice = append(newSlice, 1)

	newnewSlice := newSlice[:2:3]

	InspectSlice(newnewSlice)
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
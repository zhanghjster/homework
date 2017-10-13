package main

import (
	"container/list"
	"fmt"
	"reflect"
	"unsafe"
)

func main() {

	println(lengthOfLongestSubstring("i"))
}

func lengthOfLongestSubstring(s string) int {
	var m [256]int
	for i := range m {
		m[i] = -1
	}

	var max, cur int = 0, -1
	for i, c := range s {
		if cur < m[c] {
			cur = m[c]
		}

		if i-cur > max {
			max = i - cur
		}
		m[c] = i
	}
	return max
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

func intToList(v int) *list.List {
	var l = list.New()
	for ; v > 0; v /= 10 {
		l.PushFront(v % 10)
	}
	return l
}

func printIntList(l *list.List) {
	if l == nil {
		return
	}

	for e := l.Front(); e != nil; e = e.Next() {
		print(e.Value.(int))
	}

	fmt.Println()
}

func InspectSlice(slice []int) {

	// 数组的地址
	addr := unsafe.Pointer(&slice)

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

	for i := 0; i < *lenPtr; i++ {
		p := unsafe.Pointer(uintptr(arrPtr) + uintptr(i)*uintptr(unsafe.Sizeof(int(0))))
		fmt.Printf("[%d] %p %d\n", i, (*int)(p), *(*int)(p))
	}
}

func TryUnsafe() {
	type A struct {
		v uint16
	}

	var a = &A{v: 314} // 314 = 58 + 256

	// 将a从A结构体指针转换成uint16，
	// 由于unit16长度与A的长度相同
	// 转化后的值与a.v相同， 为 314
	p8 := *(*uint16)(unsafe.Pointer(a))
	fmt.Printf("A => uint16 \t[0x%08x] %v\n", p8, p8)

	// 将a从A结构体指针转换成uint8,
	// 由于uint8的长度小于A的长度
	// 转化后出现截断，结果为 a.v%2^8, 为58
	p16 := *(*uint8)(unsafe.Pointer(a))
	fmt.Printf("A => uint8 \t[0x%08x] %v\n", p16, p16)

	// 将a从A结构体转换成int
	// 由于int32的长度大于A，转化后的结果不可预知
	p32 := *(*uint32)(unsafe.Pointer(a))
	fmt.Printf("A => uint32 \t[0x%08x] %d\n", p32, p32)

	var b = &struct{ m, n int }{1, 2}
	// 与 p = unsafe.Pointer(&b.n)相同
	p := unsafe.Pointer(uintptr(unsafe.Pointer(b)) + uintptr(unsafe.Sizeof(b.m)))
	fmt.Printf("b addr %p, p addr %p, b.n addr %p\n", b, p, &b.n)

	p = unsafe.Pointer(b)
	//u := uintptr(unsafe.Pointer(p))
	p = unsafe.Pointer(uintptr(unsafe.Pointer(p)) + 2*uintptr(unsafe.Sizeof(b.n)))
	fmt.Printf("%d, %p\n", *(*int)(p), b)

	var s = []int{1, 2, 3}
	ip := unsafe.Pointer(uintptr(unsafe.Pointer(&s[0])) + 3*unsafe.Sizeof(int(0)))
	*(*int)(ip) = 123
	fmt.Printf("p is %d\n", *(*int)(ip))

	pf := (*int)(unsafe.Pointer(reflect.ValueOf(new(int)).Pointer()))
	*pf = 99 // 'c'

	var sf string
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&sf))
	hdr.Data = uintptr(unsafe.Pointer(pf))
	hdr.Len = 1
	fmt.Println("s => ", sf, ", len => ", len(sf))
}

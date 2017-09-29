package main

import "fmt"

func main() {

	// 底层array为 [0,1,2,3,4,5]
	var foo = []int{0,1,2,3,4,5}

	// 底层的array为[2,3,4,5]
	var bar = foo[2:3]

	fmt.Printf("cap(foo) => %d, cap(bar) => %d\n", cap(foo), cap(bar))

	// append操作会出现将foo里的值覆盖的情况
	fmt.Printf("before append to bar: \tfoo = %v\n", foo)
	bar = append(bar, 0)
	fmt.Printf("after append to bar: \tfoo = %v\n", foo)

	var array = []int{0,1,2,3,4,5}

	// 限制slice的capacity相当于限制他所允许处理的数据的范围
	slice := array[2:3:3]
	slice = append(slice, 0)

	fmt.Printf("before append to slice: array = %v\n", array)
	slice = append(slice, 0)
	fmt.Printf("after append to slice: \tarray = %v\n", array)

}

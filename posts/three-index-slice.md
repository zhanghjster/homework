---
title: Three-Index Slice
date: 2017-09-29 17:55:53
tags:
    - slice 
    - 切片
    - 三索引切片
---

 从Go1.2开始，对slice的切片操作支持设置capacity来限制新的slice容量，这样可以对底层的数组增加一层保护，对append操作进行更多的控制

slice的capacity是用来控制slice的最大容量，它反映的是底层数组的大小，如果对新的slice进行append操作，这会导致将源slice里的值覆盖的情况，这对于返回调用者一个slice的使用下是很危险的

Three-Index Slice的使用方法如下

```go
slice = array[i:j:k]
```

 ‘slice'的length为j-i, capacity为k-i, 其中k的取值范围可以为 j $<=$ k $<=$ len(array),  i可以为空表示从array的0开始切片，但j和k不能为空。append操作时，如果超出了’slice‘的capacity，完成后的slice的底层array就是一个新的数组

例子：

```go
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

benxYYTdeMac-mini:try ben$ go run try.go 
cap(foo) => 6, cap(bar) => 4
before append to bar:   foo = [0 1 2 3 4 5]
after append to bar:    foo = [0 1 2 0 4 5]
before append to slice: array = [0 1 2 3 4 5]
after append to slice:  array = [0 1 2 3 4 5]
```





参考：

https://blog.golang.org/go-slices-usage-and-internals

https://tip.golang.org/doc/go1.2#three_index

https://www.goinggo.net/2013/12/three-index-slices-in-go-12.html




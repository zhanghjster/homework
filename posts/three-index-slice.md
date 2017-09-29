---
title: Three-Index Slice
date: 2017-09-29 17:55:53
tags:
    - slice 
    - 切片
    - 三索引切片
---

 从Go1.2开始，对slice的切片操作支持设置capacity来限制新的slice容量，这样可以对底层的数组增加一层保护，并且对append操作进行更多的控制

slice的capacity是用来控制slice的最大容量，它反映的是底层数组的大小，如果对新的slice进行append操作，这回导致将源slice里的值覆盖的情况，这对于返回调用者一个slice的情境下是很危险的， 比如下面代码：

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

}

$ go run try.go 
ap(foo) => 6, cap(bar) => 4
before append to bar:   foo = [0 1 2 3 4 5]
after append to bar:    foo = [0 1 2 0 4 5]


```







参考：

https://blog.golang.org/go-slices-usage-and-internals

https://tip.golang.org/doc/go1.2#three_index

https://www.goinggo.net/2013/12/three-index-slices-in-go-12.html




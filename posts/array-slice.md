---
title: Array and Slice
date: 2017-09-29 17:55:53
tags:
    - slice
    - array
    - 数组 
    - 切片
    - 三索引切片
    - unsafe
    - go
categories:
    - go	
---

Go提供了array(数组)和slice(切片)两种类型化数据序列，array类似于其他语言的数组，具有定长的性质，slice建立在数组之上，但提供了比array更强大的功能和便利

<!-- more -->

#### Array

array通过指定长度和数据类型来定义，比如 '[4]int’定义了一个长度为4的整数array。array的长度是固定的，它的长度是类型的一部分([4]int和[5]int是不同类型)。array索引采用通常办法，从0开始，s[n]访问它的第n个元素

array在内存里是连续存储的同类型数据的序列，比如[4]int在内存里如下表示:

<img src="http://owo5nif4b.bkt.clouddn.com/go-slices-usage-and-internals_slice-array.png" width="400">

Go的array变量是值，表示整个array，它不像c那样数组名表示指向数组第一个元素的指针。这意味着当赋值或传递一个array时将复制其内容(为避免复制，可以传递指向array的指针)。可以将array看做成像结构体一样的固定长度复合值，但通过索引而不是字段名来访问里面元素

下面是c里对数组名进行指针运算

```c
 #include <stdio.h>
 
 int main() {
     int a[5] = {1,2,3,4,5};
 
     printf("a = %p, a+1 = %p\n", a, a+1);
     // a就是数组a的地址, 
     // 但地址位移则是按照类型a[5]的长度为单位
     printf("&a= %p, &a+1= %p, sizof(a)=%lu\n", &a, &a+1, sizeof(a));
 }
```

运行结果

```
a = 0x7fff5778fae0, a+1 = 0x7fff5778fae4
&a= 0x7fff5778fae0, &a+1= 0x7fff5778faf4, sizof(a)=20
```

类似上面的a+1的操作在golang里是不存在的

#### Slice

slice的类型规范是'[]T'，T表示里面元素的类型，不同于array的一个显著特点是它没有指定长度。

slice可以通过下面方法声明

```go
var slice = []byte{'a','b','c','d','e'}
```

还可以通过make函数来创建

```go
var foo = make([]int, 5, 10) // 5 表示length，10 表示capacity
var bar = make([]int, 5) // 5表示length，capacity被省略则和length相同
```

slice还可以从已经存在的slice或array‘切片’得来，通过使用冒号分隔两个索引指定半开范围(half-open)来完成，主表达式如下

```go
b = a[i:j] // 结果'b'由'a'从i开始j-i元素组成
```

* 如果'a'是array，则 $i$ $j$取值范围是 $0<=i<=j<=len(a)$,  $i$如果省略则默认值是0, $j$如果省略默认为$len(a)$, $cap(b)=j-i $ , $len(b) = j - i$
* 如果'a'是slice，则$i$ $j$取值范围是$0<=i<=j<=cap(a)$, $i$如果省略默认值是0，$j$如果省略默认值是$len(a)$ ,$cap(b) = cap(a) - i$, $len(b) = j - i$

#### Slice内部结构

slice实际上是一个描述符，它包含指向array的指针、段的长度和容量(段的最大长度)，如下图所示

<img src="http://owo5nif4b.bkt.clouddn.com/go-slices-usage-and-internals_slice-struct.png" width="400">



对于上文中定义的'slice', 它的结构如下图

<img src="http://owo5nif4b.bkt.clouddn.com/go-slices-usage-and-internals_slice-1.png" width="400">

对‘slice'进行切片

```go
var newSlice = slice[2:4]
```

'newSlice'的结构如下

<img src="http://owo5nif4b.bkt.clouddn.com/go-slices-usage-and-internals_slice-2-2.png" width="400">



切片操作不是复制slice的数据，而是创建一个新的slice，它的数组指针指向原slice里的数组，这使得切片操作非常高效。同时也造成了一个特性更新新的slice里的元素，原slice的元素也会更新。对于这个特性，在有些情况下是危险的，比如一个函数返回了一个slice的切片给调用者，但它只容忍调用者更新指定范围内的元素，如下面代码

```go
func caller() {
  src, s := callee()
  fmt.Printf("value of src before append 1 to s:\n%v\n",src)
  fmt.Printf("value of s before append 1 to it:\n%v\n",s)
  s = append(s, 1)
  fmt.Printf("value of src after append 1 to s:\n%v\n",src)
  fmt.Printf("value of s after append 1 to it:\n%v\n",s)  
}

func callee() (src, slice []int) {
  var foo = make([]int, 10)
  // 希望调用者使用索引5--7的元素
  // 调用者可以对新的slice更新或append
  // 不能影响foo里除5--7之外的元素
  return foo, foo[5:7] 
}
```

caller执行结果

```
value of src before append 1 to s:
[0 0 0 0 0 0 0 0 0 0]
value of s before append 1 to it:
[0 0]
value of src after append 1 to s:
[0 0 0 0 0 0 0 1 0 0]
value of s after append 1 to it:
[0 0 1]
```

从上面可以见到，原slice  ’src'的src[7]也被更新了，违背了callee的初衷，为了实现对caller的这种append操作的限制，从Go1.2开始，对slice的切片操作支持"Three-Index Slice"来限制切片的capacity，当caller对新切片执行append操作时，如果超出了capacity，append返回的slice会是扩展容量后slice，它内部的数组已经与callee的内部数组不是同一个了，这样就实现了callee部分数据对caller的限制

Three-Index Slice的表达式如下

```go
b = a[i:j:k] // 结果'b'由'a'从i开始j-i元素组成, b的容量 为k-i
```

- 如果'a'是array，则 $i$ $j$取值范围是 $0<=i<=j<=k<=len(a)$,  $i$如果省略则默认值是0,  $cap(b)=k - i$, $len(b) = j - i$
- 如果'a'是slice，则$i$ $j$取值范围是$0<=i<=j<=k<=cap(a)$, $i$如果省略则默认值是0， $cap(b) = k - i$, $len(b) = j - i$
- 只有$i$ 可以忽略，默认为0

下面让我们来看一下切片操作和append操作之后slice内部是如何变化的

 首先定义一个查看slice内部数据的函数，如下：

```go
// 查看slice的地址、length字段地址和值、capcity字段的地址和值
// 查看slice底层数组的每个元素的地址和值
// uintptr用于指针偏移运算
// unsafe.Pointer用于指针类型转换
func InspectSlice(slice []int) {

	// 数组的地址
	addr := unsafe.Pointer(&slice)

	// length字段的地址值和指针
	lenAddr := uintptr(addr) + uintptr(8)
  	lenPtr := (*int)(unsafe.Pointer(lenAddr))

	// capacity字段的地址值和指针
	capAddr := uintptr(addr) + uintptr(16)
	capPtr := (*int)(unsafe.Pointer(capAddr))

	// 取内部数组的地址
	// a = (*uintptr)(addr)，将addr从指向slice类型的指针转换成指向uintptr类型数据的指针
	// ptr = *(*uintptr)(addr))取a这个指针所指向的值, 
	// 为slice底层数据结构的第一个字段ptr的值
	// unsafe.Pointer(ptr)将ptr转换成指针
	arrayPtr := unsafe.Pointer(*(*uintptr)(addr))

	// 显示slice地址，length字段地址，capacity字段地址
	fmt.Printf("  Slice Addr[%p], Len Addr[0x%x] Cap Addr[0x%x]\n", addr, lenAddr,capAddr)
	
	// 显示slice的length和capacity
	fmt.Printf("  Slice length[%d] Cap[%d]\n", *lenPtr, *capPtr)

	for i :=0; i < *capPtr; i++ {
		//  根据数组地址+偏移量的办法取不同位置元素的指针和值
		p := unsafe.Pointer(uintptr(arrPtr) + uintptr(i) * uintptr(unsafe.Sizeof(int(0))))
		fmt.Printf("  [%d] %p %d\n", i, (*int)(p), *(*int)(p))
	}
}
```

进行没有设置capacity的切片操作

```go
// 底层array为 [0,1,2,3,4,5]
var foo = []int{0,1,2,3,4,5}

// 底层的array为[2,3,4,5]
var bar = foo[2:3]

fmt.Printf("\nfoo = []int{0,1,2,3,4,5}, bar = foo[2:3]\n")
fmt.Println("Inspect slice 'foo':")
InspectSlice(foo)

fmt.Println("Inspect slice 'bar':")
InspectSlice(bar)

fmt.Println()

// append操作会出现将foo里的值覆盖的情况
bar = append(bar, 0)
fmt.Println("Inspect slice 'foo' after append 0 to bar:")
InspectSlice(foo)

fmt.Println("Inspect slice 'bar' after append 0 to bar:")
InspectSlice(bar)
```

运行结果为：

```
foo = []int{0,1,2,3,4,5}, bar = foo[2:3]
Inspect slice 'foo':
  Slice Addr[0xc42000a060], Len Addr[0xc42000a068] Cap Addr[0xc42000a070]
  Slice length[6] Cap[6]
  [0] 0xc4200140c0 0
  [1] 0xc4200140c8 1
  [2] 0xc4200140d0 2
  [3] 0xc4200140d8 3
  [4] 0xc4200140e0 4
  [5] 0xc4200140e8 5
Inspect slice 'bar':
  Slice Addr[0xc42000a080], Len Addr[0xc42000a088] Cap Addr[0xc42000a090]
  Slice length[1] Cap[4]
  [0] 0xc4200140d0 2
  [1] 0xc4200140d8 3
  [2] 0xc4200140e0 4
  [3] 0xc4200140e8 5

Inspect slice 'foo' after append 0 to bar:
  Slice Addr[0xc42000a0a0], Len Addr[0xc42000a0a8] Cap Addr[0xc42000a0b0]
  Slice length[6] Cap[6]
  [0] 0xc4200140c0 0
  [1] 0xc4200140c8 1
  [2] 0xc4200140d0 2
  [3] 0xc4200140d8 0
  [4] 0xc4200140e0 4
  [5] 0xc4200140e8 5
Inspect slice 'bar' after append 0 to bar:
  Slice Addr[0xc42000a0c0], Len Addr[0xc42000a0c8] Cap Addr[0xc42000a0d0]
  Slice length[2] Cap[4]
  [0] 0xc4200140d0 2
  [1] 0xc4200140d8 0
  [2] 0xc4200140e0 4
  [3] 0xc4200140e8 5

```

从上面结果看，在执行append前，bar的底层数组就是foo的底层数组的**0xc4200140d0**--**0xc4200140e8**部分, append值‘0’到bar实际上是将0更新到bar的底层数组的**0xc4200140d8**， 同时也是foo的底层数组的值

下面进行three-index slice：

```go
foo = []int{0,1,2,3,4,5}

fmt.Printf("\nfoo = []int{0,1,2,3,4,5}, bar = foo[2:4:4]\n")
// 限制slice的capacity相当于限制被允许访问底层的数据的范围
bar = foo[2:4:4]

fmt.Println("Inspect slice 'foo':")
InspectSlice(foo)

fmt.Println("Inspect slice 'bar':")
InspectSlice(bar)

fmt.Println()

// append时超出了bar的capacity，bar底层新生成了一个容量更大的array来保存数据
bar = append(bar, 10)
fmt.Println("Inspect slice 'foo' after append 10 to bar:")
InspectSlice(foo)

fmt.Println("Inspect slice 'bar' after append 10 to bar:")
InspectSlice(bar)
```

运行结果为：

```
foo = []int{0,1,2,3,4,5}, bar = foo[2:4:4]
Inspect slice 'foo':
  Slice Addr[0xc42000a0e0], Len Addr[0xc42000a0e8] Cap Addr[0xc42000a0f0]
  Slice length[6] Cap[6]
  [0] 0xc420014120 0
  [1] 0xc420014128 1
  [2] 0xc420014130 2
  [3] 0xc420014138 3
  [4] 0xc420014140 4
  [5] 0xc420014148 5
Inspect slice 'bar':
  Slice Addr[0xc42000a100], Len Addr[0xc42000a108] Cap Addr[0xc42000a110]
  Slice length[2] Cap[2]
  [0] 0xc420014130 2
  [1] 0xc420014138 3

Inspect slice 'foo' after append 10 to bar:
  Slice Addr[0xc42000a120], Len Addr[0xc42000a128] Cap Addr[0xc42000a130]
  Slice length[6] Cap[6]
  [0] 0xc420014120 0
  [1] 0xc420014128 1
  [2] 0xc420014130 2
  [3] 0xc420014138 3
  [4] 0xc420014140 4
  [5] 0xc420014148 5
Inspect slice 'bar' after append 10 to bar:
  Slice Addr[0xc42000a140], Len Addr[0xc42000a148] Cap Addr[0xc42000a150]
  Slice length[3] Cap[4]
  [0] 0xc4200121c0 2
  [1] 0xc4200121c8 3
  [2] 0xc4200121d0 10
  [3] 0xc4200121d8 0
```

从上面结果看bar的底层数组被限制在了foo的底层数组的**0xc420014130**—**0xc420014138**位置，当执行append 10到bar时，由于超出了bar的容量，append会重新生成一个容量更大的array(原来的两倍)来保存bar的数据，地址为**0xc4200121c0**—**0xc4200121d8**. 这样就避免了append操作造成的对foo所允许bar访问的数据之外的数据的更新

参考：

https://golang.org/ref/spec#Slice_expressions

https://blog.golang.org/go-slices-usage-and-internals

https://tip.golang.org/doc/go1.2#three_index

https://www.goinggo.net/2013/12/three-index-slices-in-go-12.html




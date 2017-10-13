---
title: Golang unsafe Pointer
date: 2017-10-13 11:37:31
tags:
   - go
   - unsafe
   - pointer
categories: 基础
---

Golang的unsafe.Pointer表示指向任意类型的指针，它有其他类型不具备的四种特殊操作

- 一个任意类型的指针可以转化成Pointer
- 一个Pointer可以转化成任意类型的指针
- 一个uintptr可以转换成Pointer
- 一个Pointer可以转换成uintptr

通过使用Pointer可以让程序绕过类型系统，做到对任意内存地址的读取，但需要非常小心。下面介绍他的有效使用方式，除这些外其他均为无效的，即便是有效的模式，也有重要的注意事项

##### 将 *T1 转换成 *T2

假设T2的长度不大于T1，并且两者共享等效的内部布局，通过这个转换可以将一种类型数据重新解释为其他类型数据，比如math.Float64bits将float64类型数据转换成uint64

```go
func Float64bits(f float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&f))
}
```

如果T1的长度小于T2，会出现截断，大于T2是无效的会出现不可预期的结果， 比如下面的实验：

```go
type A struct {
  v uint16
}
var a = &A{v: 314} // 314 = 58 + 256

// 将a从A结构体指针转换成uint16，
// 由于unit16长度与A的长度相同
// 转化后的值与a.v相同， 为 314
p8 := *(*uint16)(unsafe.Pointer(a))
fmt.Printf("A => uint16 \t[0x%08x] %v\n", p8,  p8)

// 将a从A结构体指针转换成uint8,
// 由于uint8的长度小于A的长度
// 转化后出现截断，结果为 a.v%2^8, 为58
p16 := *(*uint8)(unsafe.Pointer(a))
fmt.Printf("A => uint8 \t[0x%08x] %v\n", p16,  p16)

// 将a从A结构体转换成int
// 由于int32的长度大于A，转化后的结果不可预知
p := *(*uint32)(unsafe.Pointer(a))
fmt.Printf("A => uint32 \t[0x%08x] %d\n", p,  p)
// 目前go1.9 用go vet还没有检测上面这行的错误用法，但不保证将来
```

运行结果

```
A => uint16     [0x0000013a] 314
A => uint8      [0x0000003a] 58
A => uint32     [0xbf30013a] 3207594298
```



从实验结果看， 第二行可见转化后是出现截断，第三行可见最终结果不可预测

##### 将Pointer转化为uintptr

将Pointer转化为uintptr后会产生一个保存Pointer指向的内存地址的整数值。通常情况下直接把uintptr转换为Pointer是无效的。当GC要删除对象时，不会更新uintptr里的值，所以uintprt所保存的地址已经不是原有内容了

接下列举几个唯一有效的将uintptr转换为Pointer的模式

##### 结合计算进行Pointer和uintptr之间互相转换

如果 p 指向一个分配的对象，他可以通过转换为uintptr后再加上一个偏移量，然后转化为Pointer实现指针的移动。比如

```go
p = unsafe.Pointer(uintptr(p) + offset)
```

这个办法通常用在访问一个struct的字段或者array的元素，如下面

```go
var b = &struct{m,n int}{1, 2}
// 与 p = unsafe.Pointer(&b.n)相同
p := unsafe.Pointer(uintptr(unsafe.Pointer(b)) + uintptr(unsafe.Sizeof(b.m)))
fmt.Printf("b addr %p, p addr %p, b.n addr %p\n", b, p, &b.n)
```

运行结果

```
b addr 0xc420014080, p addr 0xc420014088, b.n addr 0xc420014088
```

通过增加或减少偏移量来实现指针移动都是有效的，但不能超过原始分配空间的边界，否则无效，比如

```go
var s = []int{1,2,3}
// p 超出了s分配空间的大小，无效
p := unsafe.Pointer(uintptr(unsafe.Pointer(&s[0])) + 3 * unsafe.Sizeof(int(0)))
*(*int)(p) = 123
fmt.Printf("p is %d\n", *(*int)(p))
// 目前go1.9 用go vet没有检测这种错误用法，但不能保证将来是否会
```

从Pointer到uintptr经过计算互转应该放到一个表达式里，不要将uintptr保存到一个变量里，比如下面用法

```go
u := uintptr(p)
p = unsafe.Pointer(u + offset) 
// 目前go1.9 用go vet还不能检查出来这种错误用法，但不保证将来是否会
```

##### 调用syscall.Syscall时将Pointer转换成uintptr

Syscall函数直接将uiniptr参数传给操作系统，系统命令的实现会隐式的将其转换为Pointer，如果一个Pointer参数必须转化为uintptr作为参数，转化必须在表达式里，如下面所示：

```go
syscall.Syscall(SYS_READ, uintptr(fd), uintptr(unsafe.Pointer(p)), uintptr(n))
```

下面这种方式则是无效的

```go
u := uintptr(unsafe.Pointer(p))
syscall.Syscall(SYS_READ, uintptr(fd), u, uintptr(n))
```

##### 将reflect.Value.Pointer or reflect.Value.UnsafeAddr从UIuintptr转化为Pointer

reflect.Value的方法Pointer和UnsafeAddr返回了uintptr而不是unsafe.Pointer，目的是让调用者必须显示的导入unsafe包才能将返回结果转化成其他类型，但这要求在调用这连个方法时候立即将结果转化成Pointer, 如下

```go
p := (*int)(unsafe.Pointer(reflect.ValueOf(new(int)).Pointer()))
```

同样，也不能先把uintptr保存到一个变量之后再转换，下面方法无效

```go
u := reflect.ValueOf(new(int)).Pointer()
p := (*int)(unsafe.Pointer(u))
```

##### 将reflect.SliceHader或reflect.StringHeader的Data字段转化成Pointer

如上面情况类似，relfect.Sliceheader和refelect.Stringheader的Data字段是一个uintptr类型，目的是让调用者必须显示的导入unsafe才能将其转化成其他类型。然而这意味着SliceHeader和StringHeader只有在实际解析slice或者string时才有效

```go
p := (*int)(unsafe.Pointer(reflect.ValueOf(new(int)).Pointer()))
*p = 99 // 'c'

var s string
hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
hdr.Data = uintptr(unsafe.Pointer(p))
hdr.Len = 1
fmt.Println("s => ", s,", len => ", len(s))
// 输出结果为 s =>  c , len =>  1
```

通常 reflect.StringHeader和reflect.SliceHeader只能以*refelect.Sliceheader和 *reflect.StringHeader方式来使用，不能用作结果体

##### 总结

unsafe是用于Go compiler而不是Go runtime，使用它的时候必须小心并且遵循unsafe的说明文档，否则会有不可以预料的问题

参考：

https://golang.org/pkg/unsafe/

https://golang.org/cmd/vet/

http://www.tapirgames.com/blog/golang-unsafe
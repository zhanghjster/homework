package main

import (
	"fmt"
	"unsafe"
)

type A struct {
	c int8
	s int16
}

type B struct {
	s int16
	c int8
	i int
}

type C struct {
	c int8
	f float32
	i int
}

type D struct {
	f float32
	i int
	c int8
}

type F struct {
	c int8
	i int
	cc int16
}

type G struct {
	c int8
	cc int16
	i int
}

func main() {
	var a = A{}
	var b = B{}
	var c = C{}
	var d = D{}

	// size of A = 4
	// c: 1 字节 + 1 字节padding
	// s: 2 字节
	fmt.Printf("size of A = %d\n", unsafe.Sizeof(a))

	// size of B = 16
	// s: 2 字节
	// c: 1 字节 + 5 字节padding
	// i: 8 字节
	fmt.Printf("size of B = %d\n", unsafe.Sizeof(b))

	// size of c = 16
	// c: 1 字节 + 3 字节padding
	// f: 4 字节
	// i: 8 字节
	fmt.Printf("size of C = %d\n", unsafe.Sizeof(c))

	// size of d = 24
	// f: 4字节 + 4字节padding
	// i: 8 字节
	// c: 1字节 + 7 字节padding
	// c的尾部padding的原因是要保证是结构体自身也是8bit对齐的
	// 因为这样可以确保实现结构体数组时候里面每个元素也是对齐的
	fmt.Printf("size of D = %d\n", unsafe.Sizeof(d))


	// 由于有补齐，两个结构体即便有相同类型的字段，但前后顺序不同也可导致size不同
	var f = F{}
	var g = G{}

	// size of f = 24
	// c: 1字节 + 7字节padding
	// i: 8字节
	// cc: 2字节 + 4字节padding
	fmt.Printf("size of F = %d\n", unsafe.Sizeof(f))

	// size of g = 16
	// c: 1字节 + 1字节padding
	// cc: 2字节 + 4字节padding
	// i: 8字节
	fmt.Printf("size of G = %d\n", unsafe.Sizeof(g))
	fmt.Printf("offset of g.cc = %d\n", unsafe.Offsetof(g.cc))

	// 对齐原则：
	// 1. 数据类型自身对齐.
	// 		int8 byte bool	1字节对齐
	// 		int16 			2字节对齐
	// 		int31 float32 	4字节对齐
	// 		int int64 		8字节对齐
	// 		slice string	8字节对齐
	// 2. 结构体自身对齐，其大小要被最宽基本类型成员大小整除
	// 3. 结构体成员自身对齐
	// 结构体padding的原则就是要保证#2 和 #3

	// 为什么要对齐:
	// 每种数据类型都需要进行对齐操作，这是处理器结构决定的，而不是语言。
	// 现代计算机从内存读写数据，是按照'字长'（数据总线宽度）为单位的，将数据的地址设置为'字长'的倍数可以增加CPU读取数据的效率。
	//
	// 在32位机器上为4字节，64位则8字节。
	//
	//
	// 以32位机器的一个整数为例，CPU每次读取的数据长度是4，如果整数地址在4的倍数，比如 8 则一个读取周期能够完成。
	// 否则，如果在14，则需要读取两次。
	//
	// 如果64位机器则需要为 8 的倍数
	//
	// 参考:
	// http://www.geeksforgeeks.org/structure-member-alignment-padding-and-data-packing/
	// http://www.catb.org/esr/structure-packing/#_structure_alignment_and_padding
}

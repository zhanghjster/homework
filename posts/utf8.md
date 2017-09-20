---
title: UTF8
date: 2017-09-19 22:09:32
tags:
	- golang
	- utf8
	- unicode
---

### Unicode

在很久很久以前，计算机还处于中古时期的时候，人们用不多于127个数字就能表示所有用到的符号，用一个byte的7个bit就能完成此任务。其中 1-31 表示不可打印的字符，用于控制，比如7表示beep， 32-127为可打印。随着时间的推移，人们需要表示更多的符号，比如中文、日文、法文，但缺少统一的标准来规定哪个数字对应哪个符号。这样Unicode应运而生，将全世界的所有符号用统一的16-bit数字表示，最多能有65535个符号。

<!-- more -->

在unicode里，一个字符映射到一个理论上的概念Code Point，它只是一个用来表示字符的16bit数字，并不代表真实的存储格式。

举个栗子, Hello 用code point 表示为:

```
U+0048 U+0065 U+006C U+006C U+006F

```
#### 编码

在最早，人们直接用code point的值的方式进行编码，那么Hello的编码格式为

```
00 48 00 65 00 6C 00 6C 00 6F
```

还可以

```
48 00 65 00 6C 00 6C 00 6F 00
```

分别对应大端和小端。这种编码格式名为UCS-2.人们为了区分大端和小端这两种格式，在字符串的开始加上FE FF或FF FE来做出分别。

随着时间推移，65539个字符已经不够用了，UTF-16就出现了，它用一个或两个16-bit来表示。

这种办法用了一段时间后，人们开始觉得不方便，一是很多的00没有意义，空增加了字符的长度，二是那些ASCII编码的字符还得进行转码才能兼容。

这样UTF8蛋生了。它是另外一种保存code point的方式，能够兼容ASCII的编码格式。

### UTF8

UTF8是一种能够兼容ASCII的一种折中的编码格式，能够包含任何UNICODE字符，但文件长度会有一定的增加。UTF是Unicode Transformation Format的缩写, '8'表示用一至四个8bit的字节块来表示字符。

#### 详述

对于任何不大于127(0x7F)的字符，UTF8用一个byte表示，其中低7位表示值，最高为为0， 这与ASCII编码的值时相同的。

对于值为128--2047之间的字符，UTF8为两个byte。第一个byte为110XXXXX，第二个byte为10XXXXXX，共有11位(0x7FF)用来表示字符的值.

对于2018--65535之间的字符，UTF8位三个byte。三个byte格式分别为1110XXXX,10XXXXX和10XXXXXX, 共有16位（0xFFFF)来表示字符的值.

对于任何大于65536的字符，UTF8用四个byte表示。四个byte格式分别为11110XXX,10XXXXXX,10XXXXXX,10XXXXXX. 共有21位（0x1FFFFF), 但Unicode的字符并没有0x1FFFFF个，最大值为0x10FFFF

下面是字符对照表

1st Byte|2nd Byte|3rd Byte|4th Byte | Number of free Bits | Max unicode value
--------|--------|--------|---------|---------------------|-------------------
0xxxxxxx|        |        |         |       7             | 127
110xxxxx|10xxxxxx|        |         |     5+6 = 11        | 2047(0x7F)
1110xxxx|10xxxxxx|10xxxxxx|         |     4+6*2 = 16      | 65535(0xFFFF)
11110xxx|10xxxxxx|10xxxxxx|10xxxxxx |     3+6*3 = 21      | 1,114,111(0x10FFFF)


#### Golang里的string

在Go里，string是一个只读byte slice, 它里面存储的只是任意字节，并没有要求保存的是Unicode还是UTF8的文字还是其他预定义的格式

用Go写的代码是utf8的，所以在go代码里定义的字符串通常情况下是UTF8格式的。但如果字符串包含破坏UTF8格式的转义字符，则字符串里的内容就不全是utf8的了。比如:

```
package main

import "fmt"

func main() {
	var sample = "\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98"
	fmt.Printf("sample string : \t %s\n", sample)
	fmt.Printf("sample  byte: \t% x\n", sample)
	var in = `⌘`
	fmt.Printf("⌘ quoted: %+q, byte: % x\n", in, in)
}

localhost:run Ben$ go run main.go 
sample string :          ��=� ⌘
sample  byte:   bd b2 3d bc 20 e2 8c 98
⌘ quoted: "\u2318", byte: e2 8c 98

```
sample 字符串是乱码的，因为里面有部分转义字符破坏了字符utf8结构，无法解释

##### for range

在用for range遍历字符串时候，golang解码utf8字符，每次读出一个有效的utf字符而不是一个byte

```
package main

import "fmt"

func main() {
	var sample = "a中文"
	
	// 每次遍历均解析出utf8的字符
	// 如果遇到无法解析的就读取一个byte
	for i, v := range sample {
		fmt.Printf("%d %#U\n", i, v)
	}

	// 每次只遍历一个byte
	for i :=0; i< len(sample); i++ {
		fmt.Printf("%d, % x\n", i, sample[i])
	}
}

localhost:run Ben$ go run main.go 
0 U+0061 'a'
1 U+4E2D '中'
4 U+6587 '文'
0,  61
1,  e4
2,  b8
3,  ad
4,  e6
5,  96
6,  87

```

unicode/utf8 包提供了更多的utf8操作函数

#### 参考

https://blog.golang.org/strings










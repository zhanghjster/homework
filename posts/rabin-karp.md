---
title: String Match - Rabin Karp
date: 2017-10-17 21:15:50
tags:
    - Rabin Karp
    - 字符串匹配
categories: 算法
---

Rabin-Karp算法由Richard M.Karp和Michael O.Rabin创建，它的逻辑部分类似暴力匹配也是通过沿文本移动的窗口截取输入，但与暴力匹配不同的是，并不直接将窗口截出的文本子串与模式串做字符比对，而是通过hash函数将两者转换成数字后进行比较，逻辑如下面伪代码

```
function RabinKarp(string t[0..n], string p[0..m])
	hp := hash(p[0..m])
	for i from 0 to n-m
		hs := hash(s[i..i+m-1])
		if hs == hp and s[i..i+m-1] == p[0..m]
			return i
	return not found	
```
<!-- more -->
算法的核心部分为hash函数，从时间复杂度上，对于文本$t$的子串$s[i..i+m-1]$ hash运算需要时间复杂度是$O(m)$ ，所以如果对每个子串都进行hash，则程序的总的时间复杂度会到$O(n-m+1)m$ ，与暴力匹配相当。为了更快的计算速度，hash函数的时间必须是恒定的，为此Rabin-Karp使用了巧妙的Rolling hash办法

所谓rolling hash是指“下一次的hash使用上次hash的结果来进行计算”，算法将字符串$S[0..m]$ 转化一个以 $d$为基数的定位数系里的一个数，逻辑如下

$$\\ H=s[0]d^{m-1} + s[1]d^{m-2}+…+s[m-2]d + s[m-1]$$

对于文本$T[0..n]$，假设$t_s$表示$T$的长度为$m$的子串$T[s..s+m-1], (s=0,1,…,n-m)$  的哈希值，则有

* $T[s..m]$的hash值：

$$\\ t_s = T[s]d^{m-1} + T[s+1]d^{m-2}+…+T[s+m-2]d+T[s+m-1]$$

* $T[s+1...m+1]$的hash值：

$$\\ t_{s+1} = T[s+1]d^{m-1} + T[s+2]d^{m-2}+…+T[s+m-1]d + T[s+m] $$

所以有

$$\\ \begin{split} t_{s+1} = dt_s - T[s]d^m + T[s+m]\end{split}$$

可见$t_{s+1}, t_s$ 之间的关系是线性的，如果提前计算好$d^m$， 算法的总时间复杂度为$O(n-m+1)$

解决了时间复杂度问题后还需要解决哈希函数的冲突问题以及值过大问题，办法可以是$d$选择使用字符串所有用到字符种类数，然后将hash结果模一个素数$q$，逻辑如下

$$\\ t_{s+1} = (dt_s - T[s]d^m + T[s+m]) \mod q$$

如果提前计算出$h_d = (d \mod q)^m$ 则有

$$\\ t_{s+1} = (dt_s - T[s]h_d + T[s+m]) \mod q$$

模式$P$的哈希值则为

$$\\ p = (P[0]d^{m-1} + p[1]d^{m-2}+…+P[1]d+p[0]) \mod q$$

以在"hello from mars"中查找"mars"为例，$d$ 取256， $q$ 取103, 用字符的ASCII码作为字符的值

<img src="http://owo5nif4b.bkt.clouddn.com/kp9.png" width="400">

首先计算出模式$P$和文本$T[0..3]$的hash和$h_d$

$$\\ p = (m\times256^3+a\times256^2+r\times256 + s)\%103 = 39 $$

$$\\ t_0 = (h\times256^3+e\times256^2+l\times256+l)\%103=54$$

$$\\ h_d = (256\%103)^4 = 63$$

因为$p \neq t_0$，窗口右移

$$\\ t_1 = (256\times t_0 - h\times h_d + o)\%103 = 70$$

$p\neq t_1$ ， 继续右移直至检测到 'mars'的hash与p相同

下面是算法的详细实现：

```go
// t 为文本, p 为模式
// d 为基数, q 为模, 
// d为文本t的所有字符种类的和, q为满足d*q接近计算机字长的最大的素数
// 
// 实际应用中d如果选择合适则不需要进行模运算，直接利用无符号整数溢出来实现取模
// 这样既可以发挥现代计算机乘法效率高于除法优势，还可以减少除法的不必要运算
// Golang的strings库就采用了这个办法，d的选择是16777619(FNV算法也用到)
func RabinKarp(t, p string, d, q int) int {
	var n, m int = len(t), len(p)
	if n < m || m == 0 {
		return -1
	}

	// 计算 h = (d%q)^m
	// right to left binary exponentiation
	var h, s int = 1, d
	for i:=m; i>0;i>>=1 {
		if i&1 != 0 {
			h = (h*s)%q
		}
		s = (s*s)%q
	}

	// 计算t[0..m]和p的hash
	var ht, hp int
	for i:=0; i<m; i++ {
		ht = (d*ht + int(t[i]))%q
		hp = (d*hp + int(p[i]))%q
	}

	// 滚动检查
	for i:=0; i <= n-m; i++  {
		// 找到潜在匹配后做字符串验证
		if ht == hp && t[i:i+m] == p {
			return i
		}

		// 根据ht[i]计算ht[i+1]
		ht = (d*ht  - int(t[i])*h + int(t[i+m]))%q

		// ht为负数时转化成与之同余的正数以便于与p进行比较
		// 如果ht = (dt*ht + q - (int(t[i])*h)%q + int(t[i+m])%q)
		// 就不需要下面转化，因为ht肯定为正数，但它效率不高
		if ht < 0 {
			ht += q
		}
	}

	return -1
}
```

##### 资料

###### 模运算性质

$$\\ (a + b) \mod q = (a \mod q + b \mod q) \mod q$$

$$\\ (a\times b) \mod q = ((a \mod q)\times(b\mod q))\mod q$$

$$\\ (a + b) \mod q = (a + (b \mod q)) \mod q$$

$$\\ (a\times b)\mod q = (a \times  (b\mod q)) \mod q$$

###### 不同类型整数同余

假设，$n$位有符号整数$a, b, c, q$ ，其中$a < b$， 且有关系  $c \equiv (a - b) \mod q $， 那么它们是无符号整型时，这个同余关系还满足吗？答案是，看情况

在有符号整数情况下 

$$\\ (a - b) \mod q = (-(b-a)) \mod q $$

在无符号整数情况下， $a < b$ 导致溢出使得 $a -b = 2^n - (b - a)$ 所以

$$\\ (a - b) \mod q = (2^n - (b - a)) \mod q$$

所以当 $2^n \mod q = 0$ 时，同余关系才会满足，否则不成立

```go
package main

import "fmt"

func main() {
	var a, b, c, q int8
	a, b = 1, 2

	var A, B, C, Q uint8
	A, B = 1, 2

	for _, v := range []int8{7, 8} {
		q = int8(v)
		Q = uint8(v)

		c = (a-b)%q + q // + q使余数为正
		C = (A - B) % Q

		fmt.Printf("mod %d:\n", v)
		fmt.Printf("\ta-b = %d (%b), (a-b)%%q = %d\n", a-b, byte(a-b), c)
		fmt.Printf("\tA-B = %d (%b), (A-B)%%Q = %d\n", A-B, A-B, C)
	}
}
```

运行结果，可以看到当 q = 8 能够整除$2^n$时，同余成立

```
$ go run main.go 
mod 7:
        a-b = -1 (11111111), (a-b)%q = 6
        A-B = 255 (11111111), (A-B)%Q = 3
mod 8:
        a-b = -1 (11111111), (a-b)%q = 7
        A-B = 255 (11111111), (A-B)%Q = 7
```














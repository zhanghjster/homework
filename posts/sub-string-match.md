---
title: Sub-string match
date: 2017-10-08 21:15:50
tags:
    - brute force
    - rabin-karp
    - rolling hash
    - go
categories:
    - 算法
---

字符串匹配是IT们经常要遇到的问题，虽然很多现代语言和库来都提供了函数来处理，但研究一下他们的工作机理还是有必要的

<!-- more -->

#### Brute Force

一种最简单直接的匹配算法就是暴力匹配(Brute Force)，之所以称它暴力是因为要将文本与匹配串的每个字符都进行比较，耗时耗力

它的逻辑可以描述为，一个包含模式串的“模板”沿文本滑动，同时对每个偏移都要检查模板上的字符是否和文本中对应的字符相等。

举个栗子，假设有一段文本"Hello from mars"，要检查它是否包含 "mars"模式串，首先我们检查文本和模式串的第一个字符, 如下图

<img src="http://owo5nif4b.bkt.clouddn.com/bf1.png" width="400">

“H“和"m"不匹配，那么模板沿文本右移一个字符，如下图

<img src="http://owo5nif4b.bkt.clouddn.com/kp2.png" width="400">

 "e"和"m"也不匹配，模式串继续右移然后比较对应位置上的字符相等，一直到第10个字符‘m'出现了第一个字符的匹配，但模式的第二个字符则与文本上不匹配，如下图

<img src="http://owo5nif4b.bkt.clouddn.com/kp3.png" width="400">

模式串继续右移直至模式串上的字符与文本里相对应的字符都相等，如下图

<img src="http://owo5nif4b.bkt.clouddn.com/kp4.png" width="400">

最终找到了匹配的位置。如果文本长度为$n$ 模式长度为$m$ ，在最坏的情况下时间复杂度为$O(n-m+1)m$

下面是go的实现

```go
// 检查文本t中是否包含模式串p
// 返回出现模式串的位置，如不存在则返回-1
func BruteForce(t, p string) int {
	if len(t) < len(p) {
		return -1
	}

	for i := 0; i <= len(t) - len(p); i++ {
		// 检查t从i开始len(p)的是否与p相等
		if t[i:i+len(p)] == p {
			return i
		}

		// 更为原始的做法
		/* 
		j := 0
		for ;j<len(p);j++ {
			if t[i+j] != p[j] {
				break
			}
		}
		if j == len(p) {
			return i
		}
		 */
	}

	return -1
}
```

#### Rabin-Karp		

Rabin-Karp算法由Richard M.Karp和Michael O.Rabin创建。对于文本长度为$n$，模式串长度为$m$，平均时间复杂度为$O(m+n)$、最坏为$O(n-m+1)*m$， 虽然最坏的时间复杂度和暴力匹配相同，但在实际应用中要优于它

Rabin-Karp算法通过哈希算法将模式串和文本的子串转化成数值后进行比较，不过由于哈希有时会因为冲突将不同字符串转换成相同的数字，所以在出现数值相等时候只能认为是潜在的匹配，还要将文本子串与模式串做一次字符的比较，幸运的是选择合适的哈希函数能够有效减少冲突

算法的基本逻辑如下面伪码：

```go
function RabinKarp(string t[0..n], string p[0..m])
	hp := hash(p[0..m])
	for i from 0 to n-m
		hs := hash(s[i..i+m])
		if hs == hp and s[i..i+m] == p[0..m]
			return i
	return not found			
```

从上面的逻辑可以看到，算法的关键之处是hash函数，由于子串$s[i+1..i+m]$ hash运算需要时间复杂度是$O(m)$ ，所以如果对每个子串都进行hash，则程序的总的时间复杂度会到$O(n-m+1)m$ ，与暴力匹配相当。为了更快的计算速度，hash函数的时间必须是恒定时间的，为此Rabin-Karp使用了巧妙的hash函数(Rolling hash)

Rolling hash(滚动哈希)函数的输入为沿着文本滑动的窗口，每次窗口后移后的hash操作要用到上一次移动后hash的结果， 下面介绍一下它的逻辑

函数将字符串$S[0..m]$ 看成一个以 $d$为基数的定位数系里的一个数，hash函数就是计算出这个数，如下

$$H=s[0]d^{m-1} + s[1]d^{m-2}+…+s[m-2]d + s[m-1]$$

给定文本$T[0..n]$和模式$P[0..m]$，假设$t_s$表示$T$的长度为$m$的子串$T[s..s+m], (s=0,1,…,n-m)$  的哈希值，计算如下

$$t_s = T[s]d^{m-1} + T[s+1]d^{m-2}+…+T[s+m-1]d+T[s+m]$$

则$t_{s+1}$结果为

$$\begin{split}t_{s+1} \\ &= T[s+1]d^{m-1} + T[s+2]d^{m-2}+…+T[s+m]d + T[s+m+1]\\ &= d(t_s - T[s]d^{m-1}) + T[s+m+1] \end{split}$$

从上式可以看出，$t_{s+1}$的计算可以直接利用到$t_s$的值，如果预先计算了$d^{m-1}$的值，则hash运算的时间复杂度为$O(1)$

$d$的值可以是字符串所有字符种类数量，为避免hash值过大以及分布不均匀，还可以模一个素数。现实中，$d$会取一个很大的素数保证hash结果分布足够均匀，不需要取模的步骤


















---
title: String Match - Brute Force
date: 2017-10-08 21:15:50
tags:
    - brute force
    - 字符串匹配
    - 暴力匹配
categories: 算法
---

字符串匹配是IT们经常要遇到的问题，虽然很多现代语言和库都提供了函数来处理，但研究一下他们的工作机理还是有必要的

一种最简单直接的匹配算法就是暴力匹配(Brute Force)，之所以称它暴力是因为要将文本与匹配串的每个字符都进行比较，耗时耗力

它的逻辑可以描述为，一个包含模式串的“模板”沿文本滑动，同时对每个偏移都要检查模板上的字符是否和文本中对应的字符相等。

<!-- more -->

举个栗子，假设有一段文本"Hello from mars"，要检查它是否包含 "mars"模式串，首先我们检查文本和模式串的第一个字符, 如下图

<img src="http://owo5nif4b.bkt.clouddn.com/bf1.png" width="400">

“H“和"m"不匹配，那么模板沿文本右移一个字符，如下图

<img src="http://owo5nif4b.bkt.clouddn.com/kp2.png" width="400">

 "e"和"m"也不匹配，模式串继续右移然后比较对应位置上的字符相等，一直到第10个字符‘m'出现了第一个字符的匹配，但模式的第二个字符则与文本上不匹配，如下图

<img src="http://owo5nif4b.bkt.clouddn.com/kp3.png" width="400">

模式串继续右移直至模式串上的字符与文本里相对应的字符都相等，如下图

<img src="http://owo5nif4b.bkt.clouddn.com/kp4.png" width="400">

最终找到了匹配的位置。如果文本长度为$n$ 模式长度为$m$ ，在最坏的情况下时间复杂度为$O(n-m+1)m$

下面是实现

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

##### 总结

暴力匹配是最容易想到的，但由于时间复杂度问题，在字符串比较长时很少使用
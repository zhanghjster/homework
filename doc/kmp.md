---
title: String Match - KMP
date: 2017-10-17 21:49:04
tags:
   - 子串匹配
   - KMP
---

在字符串匹配的算法里，不管是Brute Force还是Rolling Hash都有个相同的问题，当出现失配时匹配串都是后移一个字符后重新进行匹配检查。有没有一种更好的办法来减少匹配次数呢？答案就是KMP

KMP算法由Donald Knuth、Vaughan Pratt、James H. Morris 三人于1977年联合发表，它通过对匹配串在出现失配时包含的信息来确定下一个匹配开始的位置来减少不必要的重复操作

首先来看一下Brute Force和Rolling Hash共有的的问题，以“hello is not hella”查找“hella”为例，从第一个字符开始对应位置的字符检查，如下图

<img src="http://owo5nif4b.bkt.clouddn.com/mismatch4.png" width="400">

当检查到第五个字符 'a'!='o'。出现失配，匹配串右移一个字符，如下图

<img src="http://owo5nif4b.bkt.clouddn.com/mismatch5.png" width="400">

匹配的位置回滚到文本串的第二个字符，不过看上去这完全是多余的，因为’ello'和‘hell'也是不匹配的，再重新进行比较就是重复的操作了，KMP算法就是解决这种多余的比较。

先观察一下下面的例子，文本串T和模式串P分别为"ab abcdab bcabcaba"和"abcdabd"

<img src="http://owo5nif4b.bkt.clouddn.com/kmp3.png" width="400">

此时匹配串的最后一个字符P[6]和文本串的T[9]不匹配，但可以看到T[7]和T[8]两个字符和P[0]P[1]相匹配，那么将匹配串右移动4个字符，如下图

<img src="http://owo5nif4b.bkt.clouddn.com/kmp4.png" width="400">



下一次匹配直接从P[2]与T[9]开始,不必再对P[0]P[1]进行比较，**跳过的字符长度是=匹配串失配字符之前字符串长度-最长相同的前缀和后缀的长度**。这既是KPM算法关键所在，如果为匹配串的每个字符都找到它之前子串的最长前缀后缀，那么就能知道当在这个字符失配时跳过的字符数，可以成为建立一个"最大长度表"来保存这些数字

以上面的模式串P为例，它的每个子串的前缀后缀分别如下

| 子串     | 前缀                  | 后缀                  | 最长相同前缀后缀 |
| ------ | ------------------- | ------------------- | -------- |
| a      | 空                   | 空                   | 0        |
| ab     | a                   | b                   | 0        |
| abc    | a,ab                | c,bc                | 0        |
| abca   | **a**,ab,abc        | **a**,ca,bca        | 1        |
| abcab  | a,**ab**,abc,abca   | b,**ab**,cab,bcab   | 2        |
| abcabd | a,ab,abc,abca,abcab | d,bd,abd,cabd,bcabd | 0        |

总结出每个字符对应的最大长度表如下

| 字符     | a    | b    | c    | a    | b    | d    |
| ------ | ---- | ---- | ---- | ---- | ---- | ---- |
| 最长相同长度 | 0    | 0    | 0    | 1    | 2    | 0    |

那么当出现失配时，匹配串移动的长度可以用下面的方法计算

```
匹配串移动的位数 = 
```



























参考：

1. [Computer Algorithms: Morris-Pratt String Searching](http://www.stoimen.com/blog/2012/04/09/computer-algorithms-morris-pratt-string-searching/)
2. [KMP 算法](http://wiki.jikexueyuan.com/project/kmp-algorithm/define.html)
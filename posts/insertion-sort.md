---
title: InsertionSort
date: 2017-09-29 00:16:32
tags:
    - 排序
    - 插入排序
    - go
categories:
    - 算法
---

插入排序是一种简单的适用于小数组的排序算法，时间复杂度为$O(n^2)$，逻辑类似玩扑克时把新拿到的牌按顺序插入到已有牌中，如下图：

<img src="http://owo5nif4b.bkt.clouddn.com/Insertion-Sort.jpg" width="400">

<!-- more -->

golang的实现

```go
// 假设数组为 a = {8，6，3，7，2，5}
// for循环执行过程
//   1. i = 1, j = -1, a = {6, 8, 3, 7, 2, 5}
//   2. i = 2, j = -1, a = {3, 6, 8, 7, 2, 5}
//   3. i = 3, j = 1,  a = {3, 6, 7, 8, 2, 5}
//   4. i = 4, j = 1,  a = {2, 3, 6, 7, 8, 5}
//   5. i = 5, j = 3,  a = {2, 3, 5, 6, 7, 8}
func InsertionSort(a []int) {
	if len(a) > 1 {
		for i, v := range a[1:] {
			j := i
			// 将i之前的所有大于v的值后移动
			// 停止在不大于v的位置
			for ;j>=0 && a[j] > v;j-- {
				a[j+1] = a[j]
			}
			// 将v插入到第一个不大于v的位置
			a[j+1] = v
		}
	}
}
```


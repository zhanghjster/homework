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
// 外循环每执行一步
//   1. i = 1, j = -1, a = {6, 8, 3, 7, 2, 5}
//   2. i = 2, j = -1, a = {3, 6, 8, 7, 2, 5}
//   3. i = 3, j = 1,  a = {3, 6, 7, 8, 2, 5}
//   4. i = 4, j = 1,  a = {2, 3, 6, 7, 8, 5}
//   5. i = 5, j = 3,  a = {2, 3, 5, 6, 7, 8}
func InsertionSort(nums []int) {
	if len(nums) < 2 {
		return
	}

	// 外循环
	for i := 1; i < len(nums); i++ {
		v := nums[i]
		// 向前遍历所有比v大的值并将其后移
		j := i - 1
		for j >= 0 && nums[j] > v {
			nums[j+1] = nums[j]
			j--
		}
		// 结束时的边界j已经是比v小的值
		nums[j+1] = v
	}
}
```


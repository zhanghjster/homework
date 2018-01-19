---
title: Counting Sort
date: 2017-12-03 20:34:50
categeries: 算法
tags: 
  - 排序
  - 计数排序
---

计数排序是一种根据小整数键对一组对象进行排序的算法，它通过计算具有不同整数键的对象的数量来操作，用这些计数来确定输出序列中每个键值对应的对象的位置。它的运行时间在对象总数和键值的最大最小之间是线性的，只试用键值的变化不明显大于对象数量的情况。

计数排序的基本思想是：假设输入有 $n$ 个元素，每个元素$p_i$具有一个非负的最大值为$k$的整数键$x_i $，对第个元素 $p_i$ 计算键值小于 $x_i$ 的元素的个数$c_i$，然后利用 $c_i$ 直接将元素$p_i$放到输出数组的位置上

##### 实现

~~~go
func CountingSort(a []int, k int) []int {
	// 数组的每一项保存键值小于或等于其索引的元素的个数
	var c = make([]int, k)

	// 将输入数组元素直接当做键值
	// 计算每个键值的出现次数
	for _, v := range a {
		c[v]++
	}

	// 计算不大于每个键值的键值个数
	// 也就是键值对应元素在输出数组里的索引
	for i := 1; i < len(c); i++ {
		c[i] += c[i-1]
	}

	// 按照c里计算的每个键值所在索引生成输出
	b := make([]int, len(a))
	for i := len(a) - 1; i >= 0; i-- {
		b[c[a[i]]-1] = a[i]
		c[a[i]]--
	}

	return b
}
~~~

##### 总结

计数排序的时间复杂度为 $O(k+n)$, 适用于当输入元素具有小整数的键值，小的意思是键值的大小范围不能远大于元素的个数

另外还有一个重要的特性是它是稳定的，也就是对于两个键值相同的元素，在输出数组中的顺序和在输入数组中是相同的











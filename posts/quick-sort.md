---
title: 快速排序
date: 2017-09-26 06:43:40
tags:
    	- 算法
	- 排序
	- quicksort
---

 快速排序是一个分治法排序，思想是取一个枢纽元“pivot”然后将数组元素以枢纽元为中心分割(partition)，递归执行。枢纽元可以从通过取第一个、取最后一个、随机、中位的办法从数组里去的。算法的关键是分割的过程，在这个过程里要讲数组里所有小于pivot的数放到它的前面，大于pivot的数放到它的后面

过程如下：

```go
func QuickSort(a []int) {
	if len(a) > 1 {
		pi := partition(a)
		QuickSort(a[:pi])
		QuickSort(a[pi+1:])
	}
}

func partition(a []int) int{
	// i 为小于pivot的值的最大的索引
	var i int = -1

	// 取最后一个元素为pivot
	var l = len(a) - 1
	var p = a[l]

	// j为要交换位置的元素的游标
	for j, v := range a[:l] {
		if v <= p {
			i++
			a[i], a[j] = a[j], a[i]
		}
	}

	// i为pivot的位置
	i++

	// 将pivot交换到i
	a[i], a[l] = a[l], a[i]

	return i
}
```



参考：

http://www.geeksforgeeks.org/quick-sort/




















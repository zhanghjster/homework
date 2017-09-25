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
func QuickSort(a []int, l, h int) {
	  if a[l] < a[h] {
      		pi = partition(a, l, h)
	  		QuickSort(arr, l, pi-1)
	  		QuickSort(arr, pi+1, h)
	  }	  
}

func partition(a []int, l, h int) int {
	pivot = arr[h]
    i = l-1
    for j := l; j < h-1; j++ {
   		if a[j] <= pivot {
          i++
          a[i], a[j] = a[j], a[i]
   		}
    }
    a[i+1], a[h] = a[h], a[i+1]
    
    return i+1    
}
```



参考：

http://www.geeksforgeeks.org/quick-sort/




















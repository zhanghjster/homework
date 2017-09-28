---
title: 快速排序
date: 2017-09-26 06:43:40
tags:
    	- 算法
	- 排序
	- quicksort
---

 快速排序是一个分治法排序，由[Tony Hoare](https://en.wikipedia.org/wiki/Tony_Hoare)在20世纪50年代提出，到现在还依然普遍使用，如果实现的好可以比其合并排序、堆排序等竞争算法快二到三倍。一般情况下它的时间复杂度为 $O(nlogn)$  , 最坏情况下为$O(n^2)$

算法的思想是取一个枢纽元“pivot”然后将数组元素以枢纽元为中心分割(partition)，递归执行。枢纽元可以从通过取第一个、取最后一个、随机、中位的办法从数组里取得。其中关分割是算法的关键，在这个过程里要将数组里所有小于pivot的数放到它的前面，大于pivot的数放到它的后面

递归调用的过程如下：

```go
func QuickSort(a []int) {
	if len(a) > 1 {
		pi := partition(a)
		QuickSort(a[:pi])
		QuickSort(a[pi+1:])
	}
}
```

#### Lomuto parition schema

算法由 Nico Lomuto提出，用数组的最后一个元素作为枢纽元，用$j$ 遍历数组，用$i$ 记录当前比pivot小的元素的最大索引，当 $j$ 处元素小于pivot，则将$i+1$处元素与$j$处元素互换，最终达到数组被分成大于和小于pivot元素的两部分。当数组是预排序的或者是都是相等元素，时间复杂度是$O(n^2)$ 

下面是一段用go实现代码，注释部分模拟了每一步的数组状态：

```go

// 假设数组为 a = {8，6，3，7，2，5，9，5}
// 初始i = -1, p = 5
// 循环 j <- 0 - len(a)-1
// j = 0: i = -1, (v = 8) > 5, skip, 				a = {8，6，3，7，2，5，9，5}
// j = 1: i = -1, (v = 6) > 5, skip, 				a = {8，6，3，7，2，5，9，5}
// j = 2: i = -1, (v = 3) < 5, i=0, a[0] <=> a[2], 	a = {3, 6，8，7，2，5，9，5}
// j = 3: i = 0,  (v = 7) > 5, skip				 	a = {3, 6，8，7，2，5，9，5}
// j = 4: i = 0,  (v = 2) < 5, i=1, a[1] <=> a[4],  a = {3, 2，8，7，6，5，9，5}
// j = 5: i = 1,  (v = 5) = 5, i=2, a[2] <=> a[5],  a = {3, 2，5，7，6，8，9，5}
// j = 6: i = 2,  (v = 9) > 5, skip					a = {3, 2，5，7，6，8，9，5}
// a[i+1] <=> a[len(a)-1], 							a = {3, 2，5，5，6，8，9，7}
func LomutoPartition(a []int) int {
	// i 为小于pivot的值的最大的索引
	var i int = -1

	var h = len(a) - 1
	// 取最后一个元素为pivot
	var p = a[h]

	// j为要交换位置的元素的游标
	for j, v := range a[:h] {
		if v <= p {
			i++
			a[i], a[j] = a[j], a[i]
		}
	}

	// i为pivot的位置
	i++

	// 将pivot从数组尾部交换到i
	a[i], a[h] = a[h], a[i]

	return i
}
```

#### Hoare partition scheme

算法由Hoare提出，两个索引$i$和$j$ 分别从数组的数位开始相向而行，当$i$ 处元素大于pivot并且$j$处小于pivot，交换两处的值，然后两个索引继续前行， 直到互换位置，此时交换两个索引处的值

下面是一段用go实现代码，注释部分模拟了每一步的数组状态：

```go
// 假设数组为 a = {8，6，3，7，2，5，9，5}
// 初始i, j = -1, 8
// p = 8
// for循环执行过程
//   1. i = 0, j=7, a[0]<=>a[7], a = {5, 6, 3, 7, 2, 5, 9, 8}
// 	 2. i = 1, j=5, a[1]<=>a[5], a = {5, 5, 3, 7, 2, 6, 9, 8}
//   3. i = 3, j=4, a[3]<=>a[4], a = {5, 5, 3, 2, 7, 6, 9, 8}
//   4. i = 4, j=3, return
func HoarePartition(a []int) int {
	var i, j = -1, len(a)
  
    // 第一个元素为pivot
	p := a[0]
	for {
      	// 遇到不小于pivot的值时i停止
		for i++;a[i] < p; i++ {} 
      
        // 遇到不大于pivot的值时j停止
		for j--;a[j] > p; j-- {}
      
        // i越过j后返回
		if i >= j {
			return j
		}
		
		a[i], a[j] = a[j], a[i]
	}
}
```



​	

参考：

https://en.wikipedia.org/wiki/Quicksort

http://www.geeksforgeeks.org/quick-sort/

https://www.khanacademy.org/computing/computer-science/algorithms/quick-sort/a/analysis-of-quicksort
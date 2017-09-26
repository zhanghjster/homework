package algrithm

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
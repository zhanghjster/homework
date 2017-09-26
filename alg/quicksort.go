package alg

func LomutoQuickSort(a []int) {
	if len(a) > 1 {
		pi := HoarePartition(a)
		LomutoQuickSort(a[:pi])
		LomutoQuickSort(a[pi+1:])
	}
}

func HoareQuickSort(a []int) {
	if len(a) > 1 {
		pi := HoarePartition(a)
		HoareQuickSort(a[:pi+1])
		HoareQuickSort(a[pi+1:])
	}
}

// 假设数组为 a = {8，6，3，7，2，5，9，5}
// i = -1, p = 5
// for j <- 0 - len(a)-1
// j = 0: i = -1, (v = 8) > 5, skip, 				a = {8，6，3，7，2，5，9，5}
// j = 1: i = -1, (v = 6) > 5, skip, 				a = {8，6，3，7，2，5，9，5}
// j = 2: i = -1, (v = 3) < 5, i=0, a[0] <=> a[2], 	a = {3, 6，8，7，2，5，9，5}
// j = 3: i = 0,  (v = 7) > 5, skip				 	a = {3, 6，8，7，2，5，9，5}
// j = 4: i = 0,  (v = 2) < 5, i=1, a[1] <=> a[4],  a = {3, 2，8，7，6，5，9，5}
// j = 5: i = 1,  (v = 5) = 5, i=2, a[2] <=> a[5],  a = {3, 2，5，7，6，8，9，5}
// j = 6: i = 2,  (v = 9) > 5, skip					a = {3, 2，5，7，6，8，9，5}
// a[i+1] <=> a[len(a)-1], 							a = {3, 2，5，5，6，8，9，7}
func LomutoPartition(a []int) int{
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

	// 将pivot交换到i
	a[i], a[h] = a[h], a[i]

	return i
}
// {9,4,1,5,2,8,7,3,5}
// {3,4,1,5,2,8,7,9,5}
func HoarePartition(a []int) int {
	var i, j = -1, len(a)
	p := a[0]
	for {
		for i++;a[i] < p; i++ {}
		for j--;a[j] > p; j-- {}
		if i >= j {
			return j
		}

		a[i], a[j] = a[j], a[i]
	}
}
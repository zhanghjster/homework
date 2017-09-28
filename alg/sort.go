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

// 对于元素都相同或者比较多的情况由于Hoare
// 对于预排序的数组时间复杂度(time complexity)和Hoare相同, O(n^2)
// 假设数组为 a = {8，6，3，7，2，5，9，5}
// i, j = -1, 8
// p = 8
// for循环执行过程
//   1. i = 0, j=7, a[0]<=>a[7], a = {5, 6, 3, 7, 2, 5, 9, 8}
// 	 2. i = 1, j=5, a[1]<=>a[5], a = {5, 5, 3, 7, 2, 6, 9, 8}
//   3. i = 3, j=4, a[3]<=>a[4], a = {5, 5, 3, 2, 7, 6, 9, 8}
//   4. i = 4, j=3, return
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
			for ;j>=0 && a[j] > v;j-- {
				a[j+1] = a[j]
			}
			a[j+1] = v
		}
	}
}
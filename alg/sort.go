package alg

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

func LomutoQuickSort(a []int) {
	if len(a) > 1 {
		pi := LomutoPartition(a)
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
func LomutoPartition(a []int) int {
	// i 为小于pivot的值的最大的索引
	var i = -1

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
		for i++; a[i] < p; i++ {
		}
		for j--; a[j] > p; j-- {
		}
		if i >= j {
			return j
		}

		a[i], a[j] = a[j], a[i]
	}
}

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

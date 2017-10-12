package alg

func TwoSum(a []int, t int) (int, int) {
	var m = make(map[int]int)
	for i, v := range a {
		if j, ok := m[v]; ok {
			return j, i
		} else {
			m[t-v] = i
		}
	}

	return -1, -1
}
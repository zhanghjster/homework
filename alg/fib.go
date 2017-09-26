package alg


func FibDesc(n int) int64 {
	if n < 3 {
		return 1
	}

	return FibDesc(n-1) + FibDesc(n-2)
}

func FibAsc(n int) uint64 {
	var a, b uint64 = 1, 1
	for i := 3; i <= n; i++ {
		a, b = b, a + b
	}
	return b
}
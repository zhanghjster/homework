package math


func FibDesc(n int) int64 {
	if n < 3 {
		return 1
	}

	return FibDesc(n-1) + FibDesc(n-2)
}

func FibAsc(n int) int64 {

	var fibM1, fibM2 = 1, 0
	for i := 2; i <= n; i++ {
	}

}
package math

import (
	"math"
)

const (
	MaxUint64 = ^uint64(0)
	MaxInt64 = MaxUint64 >> 1
)

func IsPrimitive(n uint64) bool {
	for i := uint64(2); i <= uint64(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func FindAllPrimitive(n uint64) []uint64 {
	var numbers = make([]bool, n+1)

	// mark all numbers primitive
	for i := uint64(2); i <= n; i++ {
		numbers[i] = true
	}

	for i := uint64(2); i*i <= n; i++ {
		if numbers[i] {
			for j := 2*i; j<=n; j += i {
				numbers[j] = false
			}
		}
	}

	var res []uint64
	for i, v := range numbers[2:] {
		if v {
			res = append(res, uint64(i))
		}
	}

	return res
}
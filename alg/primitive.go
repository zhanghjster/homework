package alg

import (
	"math"
)

const (
	MaxUint64 = ^uint64(0)
	MaxInt64  = MaxUint64 >> 1
)

func IsPrimitive(n uint64) bool {
	for i := uint64(2); i <= uint64(math.Sqrt(float64(n))); i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

// https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes
func EratosthenesSieve(n uint64) []uint64 {
	var numbers = make([]bool, n+1)

	// create a list of consecutive integer from 2 to n
	// 创建一个从2到n的连续整数集合
	for p := uint64(2); p <= n; p++ {
		numbers[p] = true
	}

	// initially, let p equal 2, the smallest primitive number
	// repeat the enumeration from 2 to square root of p
	// 重复遍历2到p的平方根次
	for p := uint64(2); p*p <= n; p++ {
		if numbers[p] {
			// enumerate the multiples of p from 2p to n in increment by p and mark them
			// 遍历所有p的倍数并标记
			for j := 2 * p; j <= n; j += p {
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

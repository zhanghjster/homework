package math

import (
	"math"
	"fmt"
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
	var np = make([]bool, n+1)
	var ch = make(chan uint64)

	go ProcessPrimitiveCheck(np, ch)

	for i := uint64(3); i<= n; i++ {
		if i%2 != 0 {
			ch <- i
		} else {
			np[i] = true
		}
	}

	close(ch)

	res := []uint64{}
	// print all primitive
	for i, v := range np {
		if i == 0 {
			continue
		}

		if !v {
			res = append(res, uint64(i))
			fmt.Println(i)
		}
	}

	return res
}

func ProcessPrimitiveCheck(np []bool, ch chan uint64) {
	nch := make(chan uint64)
	defer close(nch)

	n, ok := <- ch
	if !ok {
		return
	}

	go ProcessPrimitiveCheck(np, nch)

	for i := range ch {
		if i%n != 0 {
			nch <- i
		} else {
			np[i] = true
		}
	}
}


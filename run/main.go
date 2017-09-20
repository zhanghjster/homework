package main

import (
	"fmt"
)

func main() {
	var sample = "a中文"
	for i, v := range sample {
		fmt.Printf("%d %#U\n", i, v)
	}

	for i :=0; i< len(sample); i++ {
		fmt.Printf("%d, % x\n", i, sample[i])
	}


	// right to left pow
	var b, n, p = 3, 3, 1
	for ; n > 0; n >>=1 {
		if n&1 == 1 {
			p *= b
		}
		b *=b
	}

	println(p)
}
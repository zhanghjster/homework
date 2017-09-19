package main

import (
	"fmt"
	"golang.org/x/text/encoding/unicode"
)

func main() {
	var sample = "a中文"
	for i, v := range sample {
		fmt.Printf("%d %#U\n", i, v)
	}

	for i :=0; i< len(sample); i++ {
		fmt.Printf("%d, % x\n", i, sample[i])
	}
	
}
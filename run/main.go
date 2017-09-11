package main

import (
	"fmt"

	"math/rand"
)

func main() {

	var pk uint8 = 3
	var res uint8=^(uint8(0))
	rand.Seed(100)
	1 := rand.Int31()

	res = res * pk
	fmt.Printf("%b %v\n", res, res)

	fmt.Printf("%b %v\n", uint32(res) * uint32(pk), uint32(res) * uint32(pk))

}

package main

import (
<<<<<<< HEAD
	"fmt"

	"math/rand"
=======
	"container/list"
	"container/heap"
>>>>>>> 4ab8eda7bacc19f50aa7eec16df9b3c66a057725
)

func main() {

<<<<<<< HEAD
	var pk uint8 = 3
	var res uint8=^(uint8(0))
	rand.Seed(100)
	1 := rand.Int31()

	res = res * pk
	fmt.Printf("%b %v\n", res, res)

	fmt.Printf("%b %v\n", uint32(res) * uint32(pk), uint32(res) * uint32(pk))
=======
	var a uint8 = 255
	list.New()
	println(a * 4)
	heap.Init()
>>>>>>> 4ab8eda7bacc19f50aa7eec16df9b3c66a057725

}


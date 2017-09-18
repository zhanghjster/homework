package main

import (
<<<<<<< HEAD
	"fmt"
	"github.com/coocood/freecache"
	"log"
	"container/list"
	"os"
	list2 "github.com/zhanghjster/homework/list"
=======
<<<<<<< HEAD
	"fmt"

	"math/rand"
=======
	"container/list"
	"container/heap"
>>>>>>> 4ab8eda7bacc19f50aa7eec16df9b3c66a057725
>>>>>>> 91f1fbe08e4af9816f1d752f32250c829729305a
)

func main() {

<<<<<<< HEAD
	var b = list2.NewRingBuffer(7)

	var p = []byte{'a', 'b', 'c', 'd', 'e'}
	n, err := b.Write(p)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d written, start: %d, end: %d, length: %d", n, b.Start(), b.End(), b.Length())

	p = make([]byte, 8)
	n, err = b.Read(p)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d readed, start: %d, end: %d, length: %d", n, b.Start(), b.End(), b.Length())
	for i:=0; i <len(p); i++ {
		log.Printf("%s", string(p[i]))
	}

	// should be io.EOF
	_, err = b.Read(p)
	if err != nil {
		log.Println(err)
	}

	// write again
	_, err = b.Write([]byte{'f','g','h'})
	if err != nil {
		log.Fatal(err)
	}

	nn, err := b.Read(p[n:])
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d readed, start: %d, end: %d, length: %d", nn, b.Start(), b.End(), b.Length())
	for i:=n; i <len(p); i++ {
		log.Printf("%s", string(p[i]))
	}

	os.Exit(1)

	// 100M
	var cacheSize = 100*1024*1024
	cache := freecache.NewCache(cacheSize)

	k, v := []byte("foo"), []byte("bar")
	exp := 60

	cache.Set(k, v, exp)

	v1, err := cache.Get(k)
	if err != nil {
		log.Fatal(err)
	}


	list.New()
=======
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
>>>>>>> 91f1fbe08e4af9816f1d752f32250c829729305a

	println("Hit Count:", cache.HitCount())
	fmt.Printf("%s\n", string(v1))

}

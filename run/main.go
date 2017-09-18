package main

import (
	"log"
	list2 "github.com/zhanghjster/homework/list"
)

func main() {

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

}

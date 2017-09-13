package main

import "fmt"

type Item struct {
	next *Item
	value int
}

type List struct {
	sentinel Item
}

func ReverseLinkedList(list *List) {
	var head *Item
	for cur := list.sentinel.next; cur != nil; {
		// save the next item pointer
		next := cur.next

		// link head after the current item
		cur.next = head

		// move head to the current item
		head = cur

		// move current item to next
		cur = next
	}

	list.sentinel.next = head
}

func main() {
	//
	var xor = ^uint8(0)
	var n uint8 = 1
	fmt.Printf("%08b\n", byte(n^xor))
}

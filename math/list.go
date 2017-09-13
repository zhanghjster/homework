package math

type Item struct {
	next *Item
	value int
}

type List struct {
	sentinel Item
}

/*
 *	given a linked list, revert it in-place and one-pass
 *  a -> b -> c -> d -> e -> to e -> d -> c -> b ->a
 */
func RevertLink(list *List) {
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

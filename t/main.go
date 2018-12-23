package main

import "fmt"

func main() {
	var a = make([]int, 6)
	for i := 0; i < len(a); i++ {
		a[i] = i + 1
	}
	fmt.Printf("%#v\n", a)
	copy(a[1:], a[:])
	fmt.Printf("%#v\n", a)
}

func Copy(a []int, idx int) {
	copy(a[idx:], a[idx+1:])
	fmt.Printf("%#v\n", a)
}
func move(a []int, idx int) {
	for ; idx+1 < len(a); idx++ {
		a[idx] = a[idx+1]
	}
	fmt.Printf("%#v\n", a)
}

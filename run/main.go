package main

import (
	"net/url"
	"log"
	"fmt"
)

func myPow(x float64, n int) float64 {
	var ret float64 = 1
	for n > 0 {
		if n&1!=0 {
			ret *= x
		}
		x *= x
		n>>=1
	}
	return ret
}

func main() {
	var query = "a=1&b=2"
	es := url.QueryEscape(query)
	fmt.Printf("%s encoded to %s\n", query, es)

	us, err := url.QueryUnescape(es)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s decoded to %s\n", es, us)

	fmt.Println(myPow(2,3))

	var items = []int{1,2,5,5,5,3,5,4,6,5,5,5,5,5,11,1,1,1,1,1,1,1,1,1,1,1,}
	majority_item := items[0]
	count := 1
	for _, v := range items[1:] {
		if v == majority_item {
			count++
		} else {
			count--
			if count == 0 {
				majority_item = v
				count = 1
			}
		}
	}
	println(majority_item)
}
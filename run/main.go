package main

import (
	"net/url"
	"log"
	"fmt"
)

func main() {
	var query = "a=1&b=2"
	es := url.QueryEscape(query)
	fmt.Printf("%s encoded to %s\n", query, es)

	us, err := url.QueryUnescape(es)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s decoded to %s\n", es, us)
}
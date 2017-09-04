package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	const sample1 = "a中国b"

	for i, w := 0,0; i < len(sample1); i += w {
		v, width := utf8.DecodeRuneInString(sample1[i:])
		fmt.Printf("%#U\n", v)
		w = width
	}
}

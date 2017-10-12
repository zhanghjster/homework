package alg

import (
	"testing"
	"fmt"
)

func TestTwoSum(t *testing.T) {
	var a = []int{2,3,9,11,3}
	var i, j = TwoSum(a, 11)
	fmt.Println(i, j)
}
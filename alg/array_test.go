package alg

import "testing"

func TestRotateArray(t *testing.T) {
	var list = []int{1,2,3,4,5,6,7,8,9,10,11}
	RotateArray(list, 4)
	t.Logf("%v", list)
}

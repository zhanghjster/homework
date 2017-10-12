package alg

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestTwoSum(t *testing.T) {
	var a = []int{2,3,9,11,3}
	var i, j = TwoSum(a, 11)
	assert.Equal(t, i, 0)
	assert.Equal(t, j, 2)
}

func TestAddTwo(t *testing.T) {
	var a, b = 438, 92
	l1, l2 := IntToList(a), IntToList(b)

	l := AddTwo(l1, l2)

	assert.Equal(t, ListToInt(l), a+b)
}
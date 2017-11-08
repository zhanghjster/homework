package alg

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMaxSumSubsequence(t *testing.T) {
	var a = []int{-2, 1, -3, 4 , -1, 2, 1, -5, 4}
	s, m := MaxSumSubsequence(a)
	assert.Equal(t, 6, m)
	assert.ObjectsAreEqual([]int{4,-1,2,1}, s)

}

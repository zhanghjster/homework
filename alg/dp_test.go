package alg

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestLongestCommonSequence(t *testing.T) {
	var x, y = "ABCBDAB", "BDCABA"
	assert.Equal(t, "BCBA", LongestCommonSequence(x, y))
}

func BenchmarkLongestCommonSequence(b *testing.B) {
	for i:=0;i<b.N;i++ {
		LongestCommonSequence("ABCBDAB", "BDCABA")
	}
}

func TestLongestRepeatedSubsequenc(t *testing.T) {
	var x = "aabb"
	s := LongestRepeatedSubsequence(x)
	assert.Equal(t, "ab", s)
}

func TestCutRod(t *testing.T) {
	var p = []int{0,1,5,8,9,10,17,17,20,24,30}
	c, q := CutRod(p, 10)
	assert.ObjectsAreEqual([]int{}, c)
	assert.Equal(t, 30, q)

	c, q = CutRod(p, 4)
	assert.ObjectsAreEqual([]int{2}, c)
	assert.Equal(t, 10, q)
}

func TestMaxSumSubsequence(t *testing.T) {
	var a = []int{-2, 1, -3, 4 , -1, 2, 1, -5, 4}
	s, m := MaxSumSubsequence(a)
	assert.Equal(t, []int{4,-1,2,1,}, s)
	assert.Equal(t, m, 6)
}

func BenchmarkMaxSumSubsequence(b *testing.B) {
	var a = []int{-2, 1, -3, 4 , -1, 2, 1, -5, 4}
	for i:=0; i<b.N; i++ {
		MaxSumSubsequence(a)
	}
}

func TestLongestIncreasingSubsequence(t *testing.T) {
	var a = []int{10,22,9,33,21,50,41,60,80}
	assert.Equal(t, 6, LengthOfLIS(a))
}

func TestLengthOfLAS(t *testing.T) {
	assert.Equal(t, 3, LengthOfLAS([]int{1,5,4}))
	assert.Equal(t, 2, LengthOfLAS([]int{1,4,5}))
	assert.Equal(t, 1, LengthOfLAS([]int{1,1,1,1,1,1,1}))
	assert.Equal(t, 6, LengthOfLAS([]int{8, 9, 6, 4, 5, 7, 3, 2, 4}))
	assert.Equal(t, 6, LengthOfLAS([]int{10, 22, 9, 33, 49, 50, 31, 60}))
}
>>>>>>> 6ac4c74d3005170dafefd2486b3913dc0746c14d

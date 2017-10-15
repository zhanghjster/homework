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

func TestThreeSum(t *testing.T) {

}

func TestAddTwo(t *testing.T) {
	var a, b = 438, 92
	l1, l2 := IntToList(a), IntToList(b)

	l := AddTwoDoubly(l1, l2)

	assert.Equal(t, ListToInt(l), a+b)
}

func TestLengthOfLongestSubstring(t *testing.T) {
	for k, v := range map[string]int{
		"abcabcbb": 3, "bbbbb": 1, "pwwkew": 3,
	} {
		assert.Equal(t, LengthOfLongestSubstring(k), v)
	}
}

func TestPalindromeNumber(t *testing.T) {
	assert.False(t, PalindromeNumber(-2147447412))
	assert.True(t, PalindromeNumber(121))
}

func TestReverseInt(t *testing.T) {
	assert.Equal(t, ReverseInt(1534236469), 0)
	assert.Equal(t, ReverseInt(0), 0)
	assert.Equal(t, ReverseInt(234), 432)
	assert.Equal(t, ReverseInt(-234), -432)
}

func TestLongestPalindrome(t *testing.T) {
	assert.Equal(t, LongestPalindrome("abcda"),"a")
	assert.Equal(t, LongestPalindrome("abab"),"aba")
	assert.Equal(t, LongestPalindrome("abba"),"abba")
}
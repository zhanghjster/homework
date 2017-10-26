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

func TestRemoveElement(t *testing.T) {
	var a = []int{1,2,3,2,5,2,3,6}
	assert.Equal(t, RemoveElement(a, 2), 5)
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

func TestGcd(t *testing.T) {
	assert.Equal(t, 1, gcd(13, 3))
	assert.Equal(t, 3, gcd(9, 6))
}

func TestRotateJuggling(t *testing.T) {
	var a = []int{1,2,3,4,5,6,7,8,9,10,11,12}
	RotateJuggling(a, 6)
	t.Logf("%v", a)
}

func TestRotateReverse(t *testing.T) {
	var a = []int{1,2,3,4,5,6,7,8,9,10,11,12}
	RotateReverse(a, 3)
	t.Logf("%v", a)
}

func TestRotateBlockSwap(t *testing.T) {
	var a = []int{1,2,3,4,5,6,7,8,9,10,11,12}
	RotateBlockSwap(a, 9)
	t.Logf("%v", a)
}

func TestSearchInReversedSortedArray(t *testing.T) {
	assert.Equal(t, 0, SearchInReversedSortedArray(
		[]int{5,6,7,1,2,3,4}, 5,
	))
	assert.Equal(t, 2, SearchInReversedSortedArray(
		[]int{5,6,7,1,2,3,4}, 7,
	))
	assert.Equal(t, 3, SearchInReversedSortedArray(
		[]int{5,6,7,1,2,3,4}, 1,
	))
	assert.Equal(t, 0, SearchInReversedSortedArray(
		[]int{1,2,3,4,5,6,7}, 1,
	))
	assert.Equal(t, 6, SearchInReversedSortedArray(
		[]int{1,2,3,4,5,6,7}, 7,
	))
}

func TestBinarySearch(t *testing.T) {
	assert.Equal(t, 0, BinarySearch([]int{1,2,3,4,5}, 0, 4, 1))
	assert.Equal(t, 4, BinarySearch([]int{1,2,3,4,5}, 0, 4, 5))
	assert.Equal(t, -1, BinarySearch([]int{1,2,3,4,5}, 0, 4, 6))
	assert.Equal(t, -1, BinarySearch([]int{1,2,3,4,5}, 0, 4, 0))
}

func TestFindPivot(t *testing.T) {
	assert.Equal(t, 0, FindPivot([]int{5,4,3,2,1}, 0, 4))
	assert.Equal(t, 4, FindPivot([]int{1,2,3,4,5}, 0, 4))
	assert.Equal(t, 1, FindPivot([]int{4,5,1,2,3}, 0, 4))
	assert.Equal(t, 3, FindPivot([]int{2,3,4,5,1}, 0, 4))
}

func TestReArrangeArray(t *testing.T) {
	var a = []int{-1,-2,-3,-4,-5,-6,1,2,3}
	ReArrangeArray(a)
	var b = []int{-1,-2,-3,4,5,6,1,2,3}
	ReArrangeArray(b)
	t.Logf("%v", b)
	var c = []int{-1,-2,-3,4,5,6}
	ReArrangeArray(c)
	t.Logf("%v", c)
}

func TestReArrangePosNeg(t *testing.T) {
	var a = []int{-1,2,-3,4,-5,6}
	ReArrangePosNeg(a)
	assert.ObjectsAreEqual([]int{-1,-3,-6, 2, 4, 5}, a)
}

func TestReArrangePosNegExtraSpace(t *testing.T) {
	var a = []int{-1,2,-3,4,-5,6}
	ReArrangePosNegExtraSpace(a)
	assert.ObjectsAreEqual([]int{-1,-3,-6, 2, 4, 5}, a)
}

func TestReArrangeThreeWay(t *testing.T) {
	var a = []int{10,3,5,6,2,4,8,1,9,7}
	ReArrangeThreeWay(a, 4, 6)
	t.Log(a)
}

func TestFloydCycleDetect(t *testing.T) {
	var list = new(ListNode)
	var tail = list
	var e *ListNode // 环的入口点
	for i := 1; i < 8; i++ {
		list = frontPushList(list, i)
		if i == 5 {
			e = list
		}
	}

	// 7->6->5->4->3->2->1
	hasCycle, size, start := FloydCycleDetect(list)
	assert.False(t, hasCycle)
	assert.Equal(t, 0, size)
	assert.Equal(t, -1, start)

	tail.Next = e
	//       0<-1<-2
	//       |     |
	// 7->6->5->4->3
	hasCycle, size, start = FloydCycleDetect(list)
	assert.True(t, hasCycle)
	assert.Equal(t, 6, size)
	assert.Equal(t, 2, start)

	// first <-> second
	first, second := new(ListNode), new(ListNode)
	first.Next, second.Next = second, first
	hasCycle, size, start = FloydCycleDetect(first)
	assert.True(t, hasCycle)
	assert.Equal(t, 2, size)
	assert.Equal(t, 0, start)

	first.Next = first

	hasCycle, size, start = FloydCycleDetect(first)
	assert.True(t, hasCycle)
	assert.Equal(t, 1, size)
	assert.Equal(t, 0, start)
}

func BenchmarkReArrangePosNeg(b *testing.B) {
	for i:=0; i<b.N; i++ {
		var a = []int{1,2,3,4,5,6,7,8,9,10,11,12}
		RotateJuggling(a, 6)
	}
}

func BenchmarkRotateJuggling(b *testing.B) {
	for i:=0; i<b.N; i++ {
		var a = []int{-1,2,-3,4,-5,6}
		ReArrangePosNeg(a)
	}
}

func BenchmarkReArrangePosNegExtraSpace(b *testing.B) {
	for i:=0; i<b.N; i++ {
		var a = []int{-1,2,-3,4,-5,6}
		ReArrangePosNegExtraSpace(a)
	}
}

func BenchmarkRotateReverse(b *testing.B) {
	for i:=0; i<b.N; i++ {
		var a = []int{1,2,3,4,5,6,7,8,9,10,11,12}
		RotateReverse(a, 5)
	}
}

func BenchmarkRotateBlockSwap(b *testing.B) {
	for i:=0; i<b.N; i++ {
		var a = []int{1,2,3,4,5,6,7,8,9,10,11,12}
		RotateBlockSwap(a, 5)
	}
}

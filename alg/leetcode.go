package alg

import (
	"container/list"
	"fmt"
	"sort"
)

func TwoSum(a []int, t int) (int, int) {
	var m = make(map[int]int)
	for i, v := range a {
		if j, ok := m[v]; ok {
			return j, i
		} else {
			m[t-v] = i
		}
	}

	return -1, -1
}

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
type ListNode struct {
	Val  int
	Next *ListNode
}

func AddTwoSingly(l1, l2 *ListNode) *ListNode {
	var h = &ListNode{}
	var e = h

	var e1, e2 = l1, l2
	for d := 0; e1 != nil || e2 != nil || d > 0; d /= 10 {
		var v1, v2 int
		if e1 != nil {
			v1 = e1.Val
			e1 = e1.Next
		}
		if e2 != nil {
			v2 = e2.Val
			e2 = e2.Next
		}
		d += v1 + v2
		e.Next = &ListNode{
			Val: d % 10,
		}
		e = e.Next
	}

	return h.Next
}

func AddTwoDoubly(l1, l2 *list.List) *list.List {
	var l = list.New()

	var e1, e2 = l1.Back(), l2.Back()
	for d := 0; e1 != nil || e2 != nil; d /= 10 {
		var v1, v2 int
		if e1 != nil {
			v1 = e1.Value.(int)
			e1 = e1.Prev()
		}
		if e2 != nil {
			v2 = e2.Value.(int)
			e2 = e2.Prev()
		}

		d += v1 + v2
		l.PushFront(d % 10)
	}

	return l
}

func LengthOfLongestSubstring(s string) int {
	var m [256]int
	for i := range m {
		m[i] = -1
	}

	var max, cur int = 0, -1
	for i, c := range s {
		if cur < m[c] {
			cur = m[c]
		}

		if i-cur > max {
			max = i - cur
		}
		m[c] = i
	}
	return max
}

func LongestPalindrome(s string) string {
	if len(s) <= 1 {
		return s
	}

	var p, q int
	for j := 0; j < len(s)-1; j++ {
		// 以j和j+1为中心
		i, k := j, j+1
		for i>=0 && k < len(s) && s[i] == s[k] {
			if k-i > q-p {
				q, p = k, i
			}
			k++
			i--
		}

		// 以j为中心
		i, k = j-1, j+1
		for i>=0 && k < len(s) && s[i] == s[k] {
			if k-i > q-p {
				q, p = k, i
			}
			k++
			i--
		}
	}

	return s[p:q+1]
}

// 负数都不是回文的
// 个位数都是回文的
func PalindromeNumber(x int) bool {
	var y int
	for i := x; i > 0 && y >= 0; i /= 10 {
		y = 10*y + i%10
	}

	return x == y
}

func PowFloat(x float64, n int) float64 {
	var r float64 = 1.0
	if n < 0 {
		n = -n
		x = 1 / x
	}
	for n > 0 {
		if n&1 != 0 {
			r *= x
		}
		x *= x
		n >>= 1
	}

	return r
}

func ReverseInt(x int) int {
	var r, p int32
	for x != 0 {
		r = r*10 + int32(x%10)
		// check overflow
		if r/10 != p {
			return 0
		}
		p = r
		x /= 10
	}

	return int(r)
}

func ThreeSum(nums []int) [][]int {
	if len(nums) < 3 {
		return nil
	}

	sort.Ints(nums)

	var res = [][]int{}
	for i:=0; i< len(nums) - 2; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}

		var j = i + 1
		var k = len(nums)-1
		for j < k {
			if j != i+1 && nums[j] == nums[j-1] {
				j++
				continue
			}

			if k != len(nums) - 1 && nums[k] == nums[k+1] {
				k--
				continue
			}

			var p = nums[i] + nums[j] + nums[k]
			switch {
			case p == 0 :
				res = append(res, []int{nums[i], nums[j], nums[k]})
				k--
				j++
			case p > 0:
				k--
			case p <  0:
				j++
			}
		}
	}

	return res
}

//********************************工具函数**************************//

func IntToList(v int) *list.List {
	var l = list.New()
	for ; v > 0; v /= 10 {
		l.PushFront(v % 10)
	}
	return l
}

func ListToInt(l *list.List) int {
	var v int
	for e := l.Front(); e != nil; e = e.Next() {
		v = 10*v + e.Value.(int)
	}
	return v
}

func PrintIntList(l *list.List) {
	if l == nil {
		return
	}

	for e := l.Front(); e != nil; e = e.Next() {
		print(e.Value.(int))
	}

	fmt.Println()
}

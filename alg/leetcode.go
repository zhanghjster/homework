package alg

import (
	"container/list"
	"fmt"
	"sort"
)

// Given an array and a value, remove all instances of that value
// in place and return the new length.
//
// Do not allocate extra space for another array,
// you must do this in place with constant memory.
//
// The order of elements can be changed.
// It doesn't matter what you leave beyond the new length.
//
// Example:
// Given input array nums = [3,2,2,3], val = 3
//
// Your function should return length = 2,
// with the first two elements of nums being 2.
//
func RemoveElement(nums []int, val int) int {
	var ret int
	for i := 0; i < len(nums) ;i++ {
		if nums[i] != val {
			nums[ret] = nums[i]
			ret++
		}
	}
	return ret
}

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
	// 高效不通用版本
	var m [256]int
	for i := range m {
		m[i] = -1
	}

	// 最大长度
	var max int = 0
	// 上一次重复时的地址
	var last int = -1
	for i, c := range s {
		// 实现两个目的
		// 1. c 出现过
		// 2. 更新上一次重复出现的地址
		if last < m[c] {
			last = m[c]
		}

		if i-last > max {
			max = i - last
		}
		m[c] = i
	}

	return max

	/* 容易理解版本
	var m = make([]int, 256)
	// max 最大长度，cur当前长度
	var max, cur int
	for i, c := range s {
		// 上一次位置
		j := m[c]

		// 字符为出现过或上一次位置在 i - cur之前
		// cur++
		if j == 0  ||  i - cur + 1 > j {
			cur++
		} else {
			// 更新当前长度最新
			cur = i + 1 - j
		}

		if cur > max {
			max = cur
		}
		// 更新当前位置
		m[c] = i + 1
	}
	return max
	*/

	/* 低效率版
	var max int
	for i := 0; i < len(s) - max; i++ {
		var l = 0
		var m = make(map[byte]int)
		for j := i; j < len(s); j++ {
			c := s[j]
			idx, ok := m[c]
			if !ok {
				m[c] = j
				l++
			} else {
				i = idx
				break
			}
		}
		if max <= l {
			max = l
		}
	}
	return max*/
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

func myAtoi(str string) int {
	var ret int

	// 去掉首部无效字符
	for len(str) > 0 &&
		!(str[0] == '+' || str[0] == '-' || (str[0]>='0' && str[0] >='9')) {
		str = str[1:]
	}

	if len(str) == 0 {
		return 0
	}

	// 正负
	var neg bool
	if str[0] == '-' {
		neg = true
		str = str[1:]
	} else if str[0] == '+' {
		str = str[1:]
	}

	// 基数
	var base = 10
	if str[0] == '0' && len(str) > 1 {
		if str[1] == 'x' || str[1] == 'X' {
			base = 16
		} else {
			base = 8
		}
	}

	for i:= 0; i<len(str); i++ {
		var pre = ret
		var v int
		switch base {
		case 16:
			if str[i] > 'A' {
				v = int(str[i] - 'A') + 10
			} else {
				v = int(str[i] - '0')
			}
		case 8, 10:
			v = int(str[i] - '0')
		}

		if v < 0 || v >= base {
			continue
		}

		ret = base*ret + v

		// 溢出为0
		if  ret/base != pre {
			return 0
		}
	}

	if neg {
		ret = - ret
	}

	return ret
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

package alg

import (
	"container/list"
	"fmt"
	"sort"
)

// 将数组用low,high分成三段，
// 小于low的在最左边，
// 大于high的在最右边
// 其他在中间
// 不必考虑顺序
//
// 朴素的办法是排序，时间复杂度为O(nlogn)
// 可是使用Dutch National Flag算法O(n)
// 1. 三个指针start, cur, end
// 	start指向小于low的元素的最大的索引
// 	cur 指向当前要检查的元素
// 	end 指向大于high的元素的最低索引
//
// 2. start从-1开始，cur0开始，当cur小于low时start++，
// 	与cur指向的元素互换,之后cur指向下一个元素，
// 	因为此时cur左边的元素一定是小于high的
//
// 3. end从数组长度开始，cur大于high时，end--
//  与cur指向元素互换，因为不能保证现在的cur元素是否<low
//  需要再与low检查一次
func ReArrangeThreeWay(nums []int, low, high int) {
	var start, cur, end = -1, 0, len(nums)
	for cur < end {
		v := nums[cur]
		if v < low {
			start++
			nums[start], nums[cur] = nums[cur], nums[start]
			// cur的左边是已经处理过的元素，cur需要指向下一个
			cur++
		} else if v > high {
			end--
			nums[end], nums[cur] = nums[cur], nums[end]
			// cur 有可能是小于low，需要再检查一次
		} else {
			cur++
		}
	}
}

// 一个数组由正负数组成，将正负数在数组中分成左右两个集合
// 但负数之间，正数之间的原有顺序不能改变
// {-1, 2, 3, -4, 5, -6} 转化为 {-1, -4, -6, 2, 3, 5}
// 要求时空间复杂度O(1)
func ReArrangePosNeg(nums []int) {
	if len(nums) < 2 {
		return
	}

	// 使用插入排序的变种，以0为枢纽元
	for i := 1; i < len(nums); i++ {
		v := nums[i]

		if v > 0 {
			continue
		}

		j := i - 1
		// 将i之前的正数右移
		for j >= 0 && nums[j] >= 0 {
			nums[j+1] = nums[j]
			j--
		}

		nums[j+1] = v
	}
}

func ReArrangePosNegExtraSpace(nums []int) {
	var tmp = []int{}
	var i, j int
	for i < len(nums) {
		if nums[i] < 0 {
			nums[j] = nums[i]
			j++
		} else {
			tmp = append(tmp, nums[i])
		}
		i++
	}

	copy(nums[j:], tmp)
}

// 一个数组由正负数组成，要求将正负元素间隔重新排列
// 如果负数多，则多余的负数排在后面
// 如果正数多，则多余的证书排在后面
// {-1,2,3,-2,-4,5,6,7}重新排列后
// {-1,2,-2,3,-4,5,6,7}
// 要求时间复杂度O(n), 空间复杂度O(1)
// 步骤:
// 1. 用快速排序的办法将正负数分开
// 2. 负数隔位与每个正数交换
func ReArrangeArray(nums []int) {
	// 以0位枢纽元分割array
	var i, j = -1, len(nums)
	for {
		for i++; nums[i] < 0; i++ {
		}
		for j--; nums[j] > 0; j-- {
		}
		if i >= j {
			break
		}
		nums[i], nums[j] = nums[j], nums[i]
	}

	// i 正数最小位置
	// j 负数最大位置
	var k = 1
	for nums[k] < 0 && i < len(nums) {
		nums[k], nums[i] = nums[i], nums[k]
		k += 2
		i++
	}
}

// 在一个没有重复元素的经过翻转升序数组查找某个数字的位置，
// 如果不存在就返回 -1
// 步骤:
// 1. 找到最大的元素所在位置
// 2. 用最大元素将数组分割成两块，在要寻找元素落入的那个块里查找
func SearchInReversedSortedArray(nums []int, n int) int {
	var start, end = 0, len(nums) - 1
	var p int
	if p = FindPivot(nums, start, end); p == -1 {
		return -1
	}

	if n >= nums[0] {
		return BinarySearch(nums, 0, p, n)
	}

	if n < nums[0] {
		return BinarySearch(nums, p+1, end, n)
	}

	return -1
}

// 升序数组中查找一个元素
func BinarySearch(nums []int, start, end, n int) int {
	if start == end && nums[start] == n {
		return start
	}

	if start < end {
		i := (start + end) / 2
		if nums[i] < n {
			return BinarySearch(nums, i+1, end, n)
		} else {
			return BinarySearch(nums, start, i, n)
		}
	}

	return -1
}

// 在翻转过的升序数组里找到最大元素的的位置
// {4,5,6,1,2,3}
func FindPivot(nums []int, start, end int) int {
	if end-start == 1 {
		if nums[end] > nums[start] {
			return end
		} else {
			return start
		}
	}

	var i = (start + end) / 2

	if nums[start] < nums[i] {
		return FindPivot(nums, i, end)
	}

	if nums[start] > nums[i] {
		return FindPivot(nums, start, i)
	}

	return -1
}

// 将数组nums后d个元素翻转到尾部
// {1，2，3，4，5，6，7}的前3个元素翻转到尾部后
// {5，6，7，1，2，3，4}
func RotateReverse(nums []int, d int) {
	if d > 0 {
		d = len(nums) - d%len(nums)
	} else {
		d = -d % len(nums)
	}
	reverse(nums[:d])
	reverse(nums[d:])
	reverse(nums)
}
func reverse(nums []int) {
	var start, end = 0, len(nums) - 1
	for start < end {
		nums[start], nums[end] = nums[end], nums[start]
		start++
		end--
	}
}

// 将数组nums前d个元素翻转到尾部
// {1，2，3，4，5，6，7}的前3个元素翻转到尾部后
// {4，5，6，7，1，2，3}
func RotateBlockSwap(nums []int, d int) {
	var left, right = d, len(nums) - d

	if d == 0 || d == len(nums) {
		return
	}

	// AB => BA
	if left == right {
		swap(nums, d)
		return
	}

	// AlArB => BArAl => BAlAr
	if left > right {
		swap(nums, right)
		// d = Al = left - right
		RotateBlockSwap(nums[right:], left-right)
	}

	// ABlBr = BrBlA => BlBrA
	if left < right {
		swap(nums, left)
		// d = Br = left
		RotateBlockSwap(nums[:right], left)
	}
}

func swap(nums []int, d int) {
	var s1, s2 = 0, len(nums) - d
	for s1 < d {
		nums[s1], nums[s2] = nums[s2], nums[s1]
		s1++
		s2++
	}
}

func RotateJuggling(nums []int, d int) {
	var n = len(nums)

	// 假设gcd=(n,d)=g，则n=a*g,d=b*g
	// LCM(n,d) = a*b*g
	// 时间复杂度
	//  = gcd(n,d)*(LCM(n,d)/d)
	//  = g*(g*a*b/d)
	//  = (a*g)*(b*g)/d
	//  = a*g
	//  = n
	for i := 0; i < gcd(n, d); i++ {
		var tmp = nums[i]
		j := i
		// LCM(n,d)/d
		for {
			k := j + d
			if k >= n {
				k -= n
			}
			if k == i {
				break
			}
			nums[j] = nums[k]
			j = k
		}
		nums[j] = tmp
	}
}
func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

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
	for i := 0; i < len(nums); i++ {
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

func frontPushList(list *ListNode, v int) *ListNode {
	return &ListNode{
		Val:  v,
		Next: list,
	}
}

func printList(list *ListNode) {
	for list != nil {
		println(list.Val)
		list = list.Next
	}
}

// has: 是否有环
// size: 环的大小
// start: 环入口到链表头的距离
func FloydCycleDetect(list *ListNode) (has bool, size, start int) {
	start = -1

	// 环中相遇节点
	var m *ListNode

	// h 快指针， t慢指针
	var h, t = list, list
	for t != nil && h != nil && h.Next != nil {
		t = t.Next
		h = h.Next.Next

		// 环中相遇
		if h == t {
			m = h
			has = true
			break
		}
	}

	if has {
		h = list
		// h 指向链表头，t从环中相遇点开支向前移动
		// 两个指针每次都移动一步
		// 再次相遇时，h所走路程为环入口处距离
		for start++; h != t; start++ {
			t = t.Next
			h = h.Next
		}

		// t重新从环中相遇点前进
		// 再回到相遇点所走路径为环大小
		t = m.Next
		for size++; t != m; size++ {
			t = t.Next
		}
	}

	return
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
		for i >= 0 && k < len(s) && s[i] == s[k] {
			if k-i > q-p {
				q, p = k, i
			}
			k++
			i--
		}

		// 以j为中心
		i, k = j-1, j+1
		for i >= 0 && k < len(s) && s[i] == s[k] {
			if k-i > q-p {
				q, p = k, i
			}
			k++
			i--
		}
	}

	return s[p : q+1]
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
	for i := 0; i < len(nums)-2; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}

		var j = i + 1
		var k = len(nums) - 1
		for j < k {
			if j != i+1 && nums[j] == nums[j-1] {
				j++
				continue
			}

			if k != len(nums)-1 && nums[k] == nums[k+1] {
				k--
				continue
			}

			var p = nums[i] + nums[j] + nums[k]
			switch {
			case p == 0:
				res = append(res, []int{nums[i], nums[j], nums[k]})
				k--
				j++
			case p > 0:
				k--
			case p < 0:
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
		!(str[0] == '+' || str[0] == '-' || (str[0] >= '0' && str[0] >= '9')) {
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

	for i := 0; i < len(str); i++ {
		var pre = ret
		var v int
		switch base {
		case 16:
			if str[i] > 'A' {
				v = int(str[i]-'A') + 10
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
		if ret/base != pre {
			return 0
		}
	}

	if neg {
		ret = -ret
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

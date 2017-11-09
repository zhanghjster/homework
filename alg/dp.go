package alg

// dynamic programing

// count(c []int, n - c[m-1]) + count(c []int, c[0:m-1])
//
func ChangeOfCoins(coins []int, n int) int {

	return 0
}

// 用最少的硬币找零
// f(c[], n) = min(f(c[], n-c[0]),...,f(c[], n-c[len(c)-1]))
// 举例 c = {1,2,5}, n = 10
// f(0) = -1
// f(1) = min(f(0)) + 1 = 1
// f(2) = min(f(1), f(0)) + 1 = 1
// f(3) = min(f(2), f(1)) + 1 = 2
// f(4) = min(f(3), f(2)) + 1 = 2
// f(5) = min(f(4), f(3), f(0)) + 1 = 1
// f(6) = min(f(5), f(4), f(1)) + 1 = 2
// f(7) = min(f(6), f(5), f(2)) + 1 = 2
// f(8) = min(f(7), f(6), f(3)) + 1 = 3
// f(9) = min(f(8), f(7), f(4)) + 1 = 3
// f(10)= min(f(9), f(8), f(5)) + 1 = 2
func MinNumOfCoins(coins []int, n int) int {
	var r = make([]int, n+1)
	for i := 1; i <= n; i++ {
		r[i] = n + 1
		for _, c := range coins {
			if c <= i && r[i-c] < r[i] {
				r[i] = r[i-c]
			}
		}
		r[i]++
	}

	if r[n] > n {
		return -1
	}

	return r[n]
}

// 最长zigzag子序列
func LengthOfLAS(nums []int) int {
	var n = 1
	if len(nums) < 2 {
		return n
	}

	var up = nums[0] >= nums[1]
	for i := 0; i < len(nums)-1; i++ {
		if up != (nums[i] <= nums[i+1]) {
			up = !up
			n++
		}
	}

	return n
}

// 给定一个整数数组，找到最长的递增子序列的长度
// 动态规划法 O(n^2)
// 有一个O(nlogn)的算法
// http://www.geeksforgeeks.org/longest-monotonically-increasing-subsequence-size-n-log-n/
func LengthOfLIS(nums []int) (n int) {
	if len(nums) == 0 {
		return
	}

	// 以每个元素结尾的序列的最长递增子序列长度
	var max = make([]int, len(nums))

	for i := 0; i < len(nums); i++ {
		// L(i) = 1 + max(L(j)) (0 < j < i && a[j] < a[i])
		var m int
		for j := i - 1; j >= 0; j-- {
			if nums[j] < nums[i] && m < max[j] {
				m = max[j]
			}
		}

		max[i] = m + 1

		if max[i] > n {
			n = max[i]
		}
	}

	return
}

// 给定一个整数数组，找到和最大的连续子序列
// 给定数组 [-2, 1, -3, 4 , -1, 2, 1, -5, 4]
// 最大连续子序列为 [4, -1, 2, 1]
//
// seq: 最大子序列
// max: 最大值
func MaxSumSubsequence(nums []int) (seq []int, max int) {
	// max 最大子序列的和
	max = nums[0]

	// sum 以当前元素结尾的最大子序列的和
	// 初始化为第一个元素-1
	sum := max - 1

	// i 为当前元素j结尾的最大子序列的开始位置
	for i, j := 0, 1; j < len(nums); j++ {
		v := nums[j]

		// 构造当前元素为结尾的最大子序列
		// 换个角度，站在前一个元素nums[i-1]位置上，假设最大子序列为S[j-1],和为B[j-1]
		// 通过最S[j-1]和nums[j], 构造以nums[j]为结尾最大子序列S[j]，假设最大值为B[j]
		// B[j] = max(B[j-1]+nums[j], nums[j])

		// B[j-1] + nums[j] > nums[j]
		// 将nums[j]与S[j-1]组成的集合大于nums[j]组成的集合
		if sum+v >= v {
			sum += v
		} else {
			sum = v
			i = j
		}

		// 最大子序列更新
		if sum > max {
			max = sum
			seq = nums[i : j+1]
		}
	}

	return seq, max
}

// p: price
// n: rod length
func CutRod(p []int, n int) (cut []int, q int) {
	if n > len(p) {
		return
	}

	var r = make([]int, n+1)
	var c = make([]int, n+1)
	for j := 1; j <= n; j++ {
		for i := 1; i <= j; i++ {
			if q < p[i]+r[j-i] {
				q = p[i] + r[j-i]
				c[j] = i
			}
		}
		r[j] = q
	}

	cut = []int{}
	for n > 0 {
		if c[n] > 0 && c[n] < n {
			cut = append(cut, c[n])
		}
		n -= c[n]
	}

	return
}

// Given string "ABCBDAB" and "BDCABA",
// find the longest common sequence "BCAB"
func LongestCommonSequence(x, y string) string {
	var lx, ly = len(x), len(y)

	// store the lcs length
	var c = make([][]int, lx+1)
	for i := range c {
		c[i] = make([]int, ly+1)
	}

	for i := 1; i <= lx; i++ {
		for j := 1; j <= ly; j++ {
			switch {
			case x[i-1] == y[j-1]:
				c[i][j] = c[i-1][j-1] + 1
			case c[i-1][j] >= c[i][j-1]:
				c[i][j] = c[i-1][j]
			default:
				c[i][j] = c[i][j-1]
			}
		}
	}

	var b = make([]byte, c[lx][ly])
	genLCS(b, c, lx, ly, x)

	return string(b)
}

func genLCS(b []byte, c [][]int, i, j int, x string) {
	if i == 0 || j == 0 {
		return
	}

	switch {
	case c[i][j] > c[i-1][j] && c[i][j] > c[i][j-1]:
		genLCS(b, c, i-1, j-1, x)
		b[c[i][j]-1] = byte(x[i-1])
	case c[i-1][j] >= c[i][j-1]:
		genLCS(b, c, i-1, j, x)
	default:
		genLCS(b, c, i, j-1, x)
	}
}

func LongestRepeatedSubsequence(x string) string {
	var l = len(x)

	var c = make([][]int, l+1)
	for i := range c {
		c[i] = make([]int, l+1)
	}

	for i := 1; i <= l; i++ {
		for j := 1; j <= l; j++ {
			switch {
			case x[i-1] == x[j-1] && i != j:
				c[i][j] = c[i-1][j-1] + 1
			case c[i-1][j] >= c[i][j-1]:
				c[i][j] = c[i-1][j]
			default:
				c[i][j] = c[i][j-1]
			}
		}
	}

	var b = make([]byte, c[l][l])
	genLCS(b, c, l, l, x)
	return string(b)

}

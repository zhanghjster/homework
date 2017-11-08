---
title: Maximum contiguous subsequence
date: 2017-11-06 21:25:50
tags: 
  - 算法
  - 动态规划
  - 大连续子序列
  - dynamic programing
categories: 算法
---

给定一个由n个整数组成的数组A，找到和最大的连续子序列$A[i],A[i+1],\cdot\cdot\cdot,A[j]$ , 比如 A = [-2, 1, -3, 4 , -1, 2, 1, -5, 4]的子序列[4, -1, 2, 1]的和6，为最大的子序列，这个问题可以采用暴力、分治或动态规划解决，这里只总结动态规划的处理办法Kadane算法

Kandane算法开始于一个归纳问题：如果我们知道以第$i$元素结尾的最大子序列$S_i$,其和为

$B_i$，那么以$i+1$ 这个元素结尾的最长子序$S_{i+1}$及它和$B_{i+1}$是什么? 答案是很直接的，如果 $B_i + A_{i+1} > A_{i+1}$ 则$S_{i+1} = <S_i, A_{i+1}>$ 否则 $S_{i+1} = <A_{i+1}>$，所以 $B_{i+1} = max(B_i + A_{i+1}, A_{i+1}) $，可见其具有动态规划的最优子结构的性质

<!-- more -->

所以我们可以通过遍历一次数组就可以计算出以每个位置$i$的元素结尾的最大子序列，只需要追踪已经遇到的最大和，在遍历完数组后就能知道整个数组的最大子序列了

```go
// 给定一个正数数组，找到和最大的连续子序列
// 比如数组 [-2, 1, -3, 4 , -1, 2, 1, -5, 4]
// 最大连续子序列为 [4, -1, 2, 1]
func MaxSumSubsequence(nums []int) (subSeq []int, max int) {
	curMax := nums[0]

	start, end := 0, 1
	subSeq = nums[start:end]

	for i := 1; i < len(nums); i++ {
		v := nums[i]
		if curMax+v > v {
			curMax += v
			end++
		} else {
			curMax = v
			start, end = i, i+1
		}

		if curMax > max {
			max = curMax
			subSeq = nums[start:end]
		}
	}

	return
}
```

##### 总结：

使用动态规划的解决这个问题的核心是递归出最优子问题结构 $$B_{i+1} = max(B_i+A_{i+1}, A_{i+1})$$












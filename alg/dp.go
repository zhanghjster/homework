package alg

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

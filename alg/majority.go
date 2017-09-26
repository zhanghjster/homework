package alg

// 一个数组里有个元素重复次数超过半数，要求算法找出这个元素
// 时间复杂度在o(n),空间复杂度o(1)
// 思想源自Boyer–Moore的的投票算法，两两去除不同的元素，剩下的就是所求的
// https://en.wikipedia.org/wiki/Boyer–Moore_majority_vote_algorithm
func FindMajority(items []int) int {
	var count = 1
	var item = items[0]
	for _, v := range items {
		if v == item {
			count++
		} else {
			count--
			if count == 0 {
				item = v
				count = 1
			}
		}
	}
	return item
}

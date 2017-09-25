---
title: Find majority item
date: 2017-09-25 22:37:10
tags:
    - 算法
    - find majority
    - algorithm
    	
---

一次选美大赛，由N个评委投票来产生冠军，每个评委都将投票放到自己面前的小盒子里。在冠军的票数需要超过半数的情况下，实现算法从这N个小盒子里找出这个冠军。要求时间复杂度是O(n)，空间复杂度O(1)

这是一个经典的从数组中找出重复次数过半的元素的问题，可以使用Boyer–Moore投票算法解决，解决思想就是减掉两两不同的元素剩下的就是要求的

<!-- more -->

```go
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
```

参考：

https://en.wikipedia.org/wiki/Boyer–Moore_majority_vote_algorithm






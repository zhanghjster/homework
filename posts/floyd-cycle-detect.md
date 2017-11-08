---
title: Floyd Cycle Dectection
date: 2017-10-26 22:09:32
tags:
    - 判圈算法
    - 龟兔赛跑法
    - 算法
categories: 算法
---

Floyd判圈算法又称龟兔赛跑算法，可以在链表、状态机、迭代函数上判断是否存在环，并求出环路的起点与长度，算法逻辑如下

以链表为例，假设表头结点为S，有两个指针t和h均指向它，接着同时让t和h向前推进，t每次前进一步，h前进两步，二者一直前进直到到达尽头(没有后续结点)或两者再次相遇。如果两个指针再次相遇则说明有个环C，下面简要证明一下

 <!-- more -->

首先，如果两个指针相遇，一定是有个环，因为不是环的话快的指针就一路到链表尽头了。那么有环两个指针就一定相遇吗？不一定，但我们这里的t和h就一定相遇，原因如下

假设，t和h所在的环有n个元素，某一时刻h所在位置是0，t与h的在逆时针方向上距离是a，此时他们逆时针向前移动，如果t向前移动x步后两指针相遇，则有 

$$\\ 2x \% n \equiv (a+x)\%n$$

也就是

$$\\ x\%n\equiv a$$ 

所以t和h是否相遇是否成立转变成上面的模运算是否有解，显然x=a就是解。即，不管t再向前移动a步h就会追上

所以通过只要t和h能够相遇，就能证明充分证明有环，也不用担心有环但他们永远不会指向同一个指针情况

下面看如何计算环的入口位置以及环的大小

假设一个链表$\{e_1,e_2,\cdot\cdot\cdot\}$ 的 $\{ e_p,\cdot\cdot\cdot,e_q\}$ 组成了环，t和h在环上的m点相遇， 如下图

<img src="http://owo5nif4b.bkt.clouddn.com/cycle.png" width=400>

当两个指针到达相遇点m点 时

* 慢指针t移动了 $s+x$ 步，
* 快指针移动了$n(x+y)+ s + x$ 步，$n$为快指针在圈上空跑的圈数

又因为快指针速度是慢指针的两倍，所以

$$\\ 2(s+x) = n(x+y) + s + x$$

所以

$$\\  s+x = n(x+y)$$

也就是慢指针移动的距离是环的整数倍

#### 环的起点

从$s + x = n(x+y)$ 可以得到 

$$\\  s = (n-1)(x+y) + y$$ , 

两个指针相遇后，将快指针t指向链表起点, 以每次前进一步的速度前进, 慢指针从相遇点$m$前进，那么从上面的等式可以知道，两个指针继续相遇(指针相等)的地点就是p点, 环的起点，所以此时t指针走的步数就是s的大小，环起点也就求出来了

#### 环的大小

两个指针在$m$相遇时记录这个指针，慢指针继续向前移动，当再次移动到$m$时所走的步数就是环的大小

#### 举例

```go
// 检查是否存在环路，返回环路的大小，入口处距开口处位置
// has: 是否有环
// size: 环的大小
// start: 环入口到链表头的距离
type ListNode struct {
	Val  int
	Next *ListNode
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
```

#### 参考

1. [判圈算法](https://zh.wikipedia.org/wiki/Floyd判圈算法)
2. [https://stackoverflow.com/questions/3952805/proof-of-detecting-the-start-of-cycle-in-linked-list](https://stackoverflow.com/questions/3952805/proof-of-detecting-the-start-of-cycle-in-linked-list)
3. https://cs.stackexchange.com/questions/10360/floyds-cycle-detection-algorithm-determining-the-starting-point-of-cycle
4. https://www.quora.com/How-does-Floyds-cycle-finding-algorithm-work-How-does-moving-the-tortoise-to-the-beginning-of-the-linked-list-while-keeping-the-hare-at-the-meeting-place-followed-by-moving-both-one-step-at-a-time-make-them-meet-at-starting-point-of-the-cycle


Floyd判圈算法又称龟兔赛跑算法，可以在链表上判断是否存在环，并求出环路的起点与长度。

#### 环的存在

物理上，以不同速度向前移动的两个物体在环上是肯定会相遇的。对于一个由n个元素组成的环，上面有两个指针t和h，t每次向前移动一步，h向前移动两步，则不管他们的起点在环的哪个位置，如果向同一个方向前移动，t在一个环的周期内肯定会遇h。

以h所在位置为起点，t在h前面a步的位置，他们同时向前移动，如果t在向前移动了x步后t追上了h，则有 

​	$$\\  2x \% n \equiv (a+x)\%n$$

也就是

$$\\  x\%n\equiv a$$ 

有解，显然x=a就是解，也就是不管t在h前面哪个位置，t在一个n步内都会与h相遇

可见，如果两个指针相遇，一定有环，反过来，如果有环，两个指针一定相遇。所以用这两个指针是否相遇能够完全确定环的存在

一个链表$\{e_1,e_2,\cdot\cdot\cdot\}$ 假设$\{ e_m,\cdot\cdot\cdot,e_p\}$ 组成了环，求m和p，使用两个指针h和t，每次t前进一步，h前进两步。假设他们相遇于环上的$e_M$ , 如下图

<img src="http://owo5nif4b.bkt.clouddn.com/1.jpg" width=400>

到相遇点$e_M$ 时

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

两个指针相遇后，将快指针t指向链表起点$e_1$, 以每次前进一步的速度前进, 慢指针从相遇点$e_M$ 继续前进，那么从上面的等式可以知道，两个指针继续相遇(指针相等)的地点就是 $e_m$, 环的起点，所以此时t指针走的步数就是s的大小，环起点也就求出来了

#### 环的大小

两个指针在$e_M$相遇时记录这个指针，慢指针继续向前移动，当再次移动到$e_M$时所走的步数就是环的大小

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
func  FloydCycleDetect(list *ListNode) (has bool, size, start int) {
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
		// h,t 再次相遇时所走路程为环入口处距离
		for start++;h != t; start++{
			t = t.Next
			h = h.Next
		}

		// t从相遇点前进再回到相遇点所走路径为环大小
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


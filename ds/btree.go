package ds

import (
	"fmt"
)

type BTreeNode struct {
	n        int          // number of keys
	t        int          // minimum degree(minimum children)
	leaf     bool         // true for leaf node
	keys     []int        // key slice
	children []*BTreeNode // child node slice
}

func NewBTreeNode(t int, leaf bool) *BTreeNode {
	return &BTreeNode{
		n:        0,
		t:        t,
		leaf:     leaf,
		keys:     make([]int, 2*t-1),
		children: make([]*BTreeNode, 2*t),
	}
}

// 遍历所有以本节点为根的所有结点，打印它们的key
func (s *BTreeNode) Traverse() {
	var i int
	for i = 0; i < s.n; i++ {
		if !s.leaf {
			s.children[i].Traverse()
		}
		fmt.Printf("%d \n", s.keys[i])
	}

	if s.leaf {
		return
	}

	// 遍历最后一个子节点
	s.children[i].Traverse()
}

// 在以本节点为根的结点中搜索包含指定关键字的结点以及关键字所在位置
func (s *BTreeNode) Search(k int) (*BTreeNode, int) {
	var i int
	// 找到第一个不大于key的关键字所在位置
	// 结束条件，i = s.n || s.keys[i] >= k
	for i = 0; i < s.n && s.keys[i] < k; i++ {
	}

	if i < s.n && s.keys[i] == k {
		return s, i
	}

	if s.leaf {
		return nil, 0
	}

	// 在孩子中查找
	// 包含 i == s.n 或 s.keys[i] > k
	return s.children[i].Search(k)
}

func (s *BTreeNode) insertNonFull(k int) {
	// 找到要插入的位置
	var i int
	for i = 0; i < s.n && s.keys[i] < k; i++ {
	}

	// k 已经在本节点里存在
	if i < s.n && s.keys[i] == k {
		return
	}

	// 如果是叶节点则插入
	if s.leaf {
		// 将i之后的key右移
		copy(s.keys[i+1:], s.keys[i:s.n])

		// k放到位置i
		s.keys[i] = k
		s.n++

		return
	}

	var c = s.children[i]

	// 如果子节点已满，则分裂子节点
	if c.isFull() {
		s.splitChild(c, i)

		// 如果提升上来的key等于k则不插入
		if s.keys[i] == k {
			return
		}

		// 如果提升上来的新的key小于k则插入到它的右孩子
		if s.keys[i] < k {
			c = s.children[i+1]
		}
	}

	c.insertNonFull(k)
}

// c 为要分裂的子节点
// i 为父节点保存要提升的子节点key的位置
func (s *BTreeNode) splitChild(c *BTreeNode, i int) {
	// 生成新的结点
	z := NewBTreeNode(s.t, c.leaf)

	z.n, c.n = c.t-1, c.t-1

	// 将c的后t-1个key复制到新的结点
	// c中keys索引: 0..t-2, t-1, t..2t-2
	for j := c.t; j < 2*c.t-1; j++ {
		z.keys[j-c.t] = c.keys[j]
		c.keys[j] = 0
	}

	// 将s的keys里i位之后的元素右移
	for j := s.n - 1; j >= i; j-- {
		s.keys[j+1] = s.keys[j]
	}

	// 将中间的key复制到s的keys列表的位置i
	s.keys[i] = c.keys[c.t-1]
	s.n++
	c.keys[c.t-1] = 0

	// 将c的后t个孩子指针复制到新的结点
	// c中children索引: 0..t-1, t..2t-1
	if !c.leaf {
		for j := c.t; j < 2*c.t; j++ {
			z.children[j-c.t] = c.children[j]
			c.children[j] = nil
		}
	}

	// 将s的i+1之后的孩子指针在列表中后移
	for j := s.n; j >= i+1; j-- {
		s.children[j+1] = s.children[j]
	}

	// 将新的结点指针放到s的孩子列表位置i
	s.children[i+1] = z
}

func (s *BTreeNode) isFull() bool {
	return s.n == 2*s.t-1
}

type BTree struct {
	root *BTreeNode
	t    int
}

func NewBTree(t int) *BTree {
	return &BTree{t: t}
}

// 单程下行方式遍历树插入关键字。
// 关键方法是"主动分裂", 即，在遍历一个子节点前，如果子节点已满则对其进行分裂。
// 相反的"被动分裂"则是在要插入的时候遇满才分裂，会出现重复遍历的情况。
// 比如，从根节点到叶节点都是满的，当到达叶节点要发现其已满需进行分裂，
// 提升一个关键字到父节点，发现父节点也已满需要分裂父节点，重复下去一直到根节点。
// 这样就出现了从根节点到叶节点的重复遍历。而"主动分裂"则不出有出现这种情况，因为
// 在分裂一个子节点时候父节点是已经有足够空间容纳要提升的新的key了
//
// B数的增高依赖的是root结点分裂
// B数新增的关键字都增加到了叶节点上
func (t *BTree) Insert(k int) {
	if t.root == nil {
		t.root = NewBTreeNode(t.t, true)
	}

	var r = t.root
	// root已满则分裂root
	if r.isFull() {
		// 为新的root生成结点
		s := NewBTreeNode(t.t, false)

		// 老的root称为新结点的孩子
		s.children[0] = r

		// 分裂老的root，并将一个key提升到新的root
		s.splitChild(r, 0)

		// 新的root有两个孩子，决定将key插入到哪个孩子
		var i int
		if k > r.keys[0] {
			i++
		}
		s.children[i].insertNonFull(k)

		// 更新root
		t.root = s

		return
	}

	t.root.insertNonFull(k)
}

func (t *BTree) Traverse() {
	if t.root == nil {
		return
	}

	t.root.Traverse()
}

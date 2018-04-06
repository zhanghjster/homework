package tree

type BSTNode struct {
	key   int
	value int

	Left  *BSTNode
	Right *BSTNode
}

type BSTree struct {
	Root *BSTNode
}

func (t *BSTree) Find(key int) *BSTNode {
	var cur = t.Root
	for cur != nil && cur.key != key {
		if key < cur.key {
			cur = cur.Left
		} else {
			cur = cur.Right
		}
	}

	return cur
}

func (t *BSTree) Insert(key, value int) {
	if t.Root == nil {
		t.Root = &BSTNode{
			key:   key,
			value: value,
		}

		return
	}

	var cur = t.Root
	var parent = t.Root
	var isLeft bool

	for cur != nil {
		parent = cur

		if key < cur.key {
			cur = cur.Left
			isLeft = true
		} else {
			cur = cur.Right
			isLeft = false
		}
	}

	node := &BSTNode{key: key, value: value}
	if isLeft {
		parent.Left = node
	} else {
		parent.Right = node
	}
}

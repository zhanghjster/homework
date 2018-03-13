package tree

type BiTNode struct {
	Element interface{}
	Left    *BiTNode
	Right   *BiTNode
}

func CreateBiTree(t **BiTNode, i int, ele []interface{}) {
	if i >= len(ele) {
		return
	}

	*t = new(BiTNode)
	t.Element = ele[i]

	i++
	CreateBiTree(&t.Left, i, ele)

	i++
	CreateBiTree(&t.Right, i, ele)

}

func PreOrderTraverse(t *BiTNode, l int) {
	if t == nil {
		return
	}

	// process the data

	// traverse left
	PreOrderTraverse(t.Left, l+1)

	// traverse the right
	PreOrderTraverse(t.Right, l+1)
}

func InOrderTraverse(t *BiTNode, l int) {
	if t == nil {
		return
	}

	InOrderTraverse(t.Left, l+1)

	// process the data

	InOrderTraverse(t.Right, l+1)
}

func PostOrderTraverse(t *BiTNode, l int) {
	if t == nil {
		return
	}

	PostOrderTraverse(t.Left, l+1)
	PostOrderTraverse(t.Right, l+1)

	// process the data
}

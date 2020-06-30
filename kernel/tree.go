package kernel

import "fmt"

type Less = func(interface{}, interface{}) bool
type Equal = func(interface{}, interface{}) bool

type rbnode struct {
	black               bool
	key                 interface{}
	payload             interface{}
	left, right, parent *rbnode
}

type TreeIterator struct {
	visited     map[*rbnode]bool
	currentNode *rbnode
}

func NewIterator(t *RBTree) *TreeIterator {
	if t.root == nil {
		return &TreeIterator{currentNode: nil}
	}

	visited := make(map[*rbnode]bool)
	min := t.root.minimum()
	visited[min] = true

	return &TreeIterator{currentNode: min, visited: visited}
}

func (t *TreeIterator) HasMore() bool {
	return t.currentNode != nil
}

func (t *TreeIterator) Next() {
	right := t.currentNode.right
	visitedRight := t.visited[right]
	if !visitedRight && right != nil {
		t.currentNode = right.minimum()
		t.visited[t.currentNode] = true
	} else {
		for t.currentNode != nil && t.visited[t.currentNode] {
			t.currentNode = t.currentNode.parent
		}

		t.visited[t.currentNode] = true
	}
}

func (t *TreeIterator) Key() interface{} {
	return t.currentNode.key
}

func (t *TreeIterator) Value() interface{} {
	return t.currentNode.payload
}

type RBTree struct {
	root  *rbnode
	less  Less
	equal Equal
	size  int64 // to be done
}

func NewTreeMap(less Less, equal Equal) *RBTree {
	return &RBTree{less: less, equal: equal, size: 0}
}

func (tree *RBTree) Insert(key interface{}, value interface{}) {
	node := &rbnode{key: key, payload: value}
	if tree.root == nil {
		tree.root = node
		node.black = true
	} else {
		tree.insert(node)
	}
}

func (tree RBTree) findNode(key interface{}) *rbnode {
	current := tree.root
	for current != nil && !tree.equal(current.key, key) {
		fmt.Printf("current %+v\n", current)
		if !tree.less(current.key, key) {
			current = current.left
		} else {
			current = current.right
		}
	}

	return current
}

func (tree *RBTree) Empty() bool {
	return tree.root == nil
}

func (tree *RBTree) Clear() {
	tree.root = nil
}

func (tree *RBTree) Exists(key interface{}) bool {
	return tree.findNode(key) != nil
}

func (tree *RBTree) Get(key interface{}) interface{} {
	node := tree.findNode(key)
	if node != nil {
		return node.payload
	}

	return nil
}

func (tree *RBTree) Remove(key interface{}) {
	tree.size--

	x := tree.findNode(key)

	fmt.Printf("node %+v\n", x)
	if x != nil {
		tree.delete(x)
	}
}

func (tree *RBTree) leftRotate(x *rbnode) {
	y := x.right
	x.right = y.left

	if y.left != nil {
		y.left.parent = x
	}

	y.parent = x.parent
	if x.parent == nil {
		tree.root = y
	} else {
		if x == x.parent.left {
			x.parent.left = y
		} else {
			x.parent.right = y
		}
	}

	y.left = x
	x.parent = y
}

func (tree *RBTree) rightRotate(x *rbnode) {
	y := x.left
	x.left = y.right

	if y.right != nil {
		y.right.parent = x
	}

	y.parent = x.parent
	if x.parent == nil {
		tree.root = y
	} else {
		if x == x.parent.left {
			x.parent.left = y
		} else {
			x.parent.right = y
		}
	}

	y.right = x
	x.parent = y
}

func (node *rbnode) minimum() *rbnode {
	target := node
	for target.left != nil {
		target = target.left
	}

	return target
}

func (node *rbnode) successor() *rbnode {
	if node.right != nil {
		return node.right.minimum()
	}

	parent := node.parent
	target := node
	for parent != nil && node == parent.right {
		target = parent
		parent = parent.parent
	}

	return target
}

func (node *rbnode) isBlack() bool {
	if node == nil {
		return true
	}

	return node.black
}

func (tree *RBTree) deleteFixup(x *rbnode, parent *rbnode) {
	for tree.root != x && x.isBlack() {
		if x == parent.left {
			w := parent.right
			if !w.isBlack() {
				w.black = true
				parent.black = false
				tree.leftRotate(parent)

				w = parent.right
			}

			if w.left.isBlack() && w.right.isBlack() {
				w.black = false
				x = parent
				parent = x.parent
			} else {
				if w.right.isBlack() {
					w.left.black = false
					tree.rightRotate(w)
					w = w.parent.right
				}

				w.black = parent.black
				parent.black = true
				w.right.black = true

				tree.leftRotate(parent)
				x = tree.root
			}
		} else {
			if x == parent.right {
				w := parent.left
				if !w.black {
					w.black = true
					parent.black = false
					tree.leftRotate(parent)

					w = parent.left
				}

				if w.left.isBlack() && w.right.isBlack() {
					w.black = false
					x = parent
					parent = x.parent
				} else {
					if w.left.isBlack() {
						w.right.black = false
						tree.leftRotate(w)
						w = w.parent.left
					}

					w.black = parent.black
					parent.black = true
					w.left.black = true
					tree.rightRotate(parent)
					x = tree.root
				}
			}
		}
	}

	if x != nil {
		x.black = true
	}
}

func (tree *RBTree) delete(node *rbnode) {
	var y, x *rbnode

	if node.left == nil || node.right == nil {
		y = node
	} else {
		y = node.successor()
	}

	if y.left != nil {
		x = y.left
	} else {
		x = y.right
	}

	if x != nil {
		x.parent = y.parent
	}

	if y.parent == nil {
		tree.root = x
	} else if y == y.parent.left {
		y.parent.left = x
	} else {
		y.parent.right = x
	}

	if y != node {
		node.key = y.key
		node.payload = y.payload
	}

	if y.black {
		tree.deleteFixup(x, y.parent)
	}
}

func (tree *RBTree) insertNode(node *rbnode) {
	root := tree.root
	for {
		if root != nil && tree.less(node.key, root.key) {
			if root.left != nil {
				root = root.left
			} else {
				root.left = node
				break
			}
		} else if root != nil {
			if root.right != nil {
				root = root.right
			} else {
				root.right = node
				break
			}
		}
	}

	node.parent = root
	node.black = false
}

func (tree *RBTree) insert(x *rbnode) {
	tree.insertNode(x)

	for x != tree.root && !x.parent.black {
		if x.parent == x.parent.parent.left {
			y := x.parent.parent.right
			if !y.isBlack() {
				x.parent.black = true
				y.black = true
				x.parent.parent.black = false

				x = x.parent.parent
			} else {
				if x == x.parent.right {
					x = x.parent
					tree.leftRotate(x)
				}

				x.parent.black = true
				x.parent.parent.black = false
				tree.rightRotate(x.parent.parent)
			}
		} else {
			y := x.parent.parent.left
			if !y.isBlack() {
				x.parent.black = true
				y.black = true
				x.parent.parent.black = false

				x = x.parent.parent
			} else {
				if x == x.parent.left {
					x = x.parent
					tree.rightRotate(x)
				}

				x.parent.black = true
				x.parent.parent.black = false
				tree.leftRotate(x.parent.parent)
			}
		}
	}

	tree.root.black = true
}

func (tree *RBTree) GetIterator() *TreeIterator {
	return NewIterator(tree)
}

func (tree *RBTree) ToMap() map[interface{}]interface{} {
	it := tree.GetIterator()

	result := make(map[interface{}]interface{})

	for it.HasMore() {
		result[it.Key()] = it.Value()

		it.Next()
	}

	return result
}

package kernel

import (
	"fmt"
	"testing"
)

func lessInt(a interface{}, b interface{}) bool {
	return a.(int) < b.(int)
}

func equalInt(a interface{}, b interface{}) bool {
	return a.(int) == b.(int)
}

func TestRBTree_Insert(t *testing.T) {
	tree := NewTreeMap(lessInt, equalInt)

	for i := 0; i < 100; i++ {
		tree.Insert(i, i)
	}

	size := 0
	iterator := NewIterator(tree)
	for iterator.HasMore() {
		if iterator.Value() != size {
			t.Failed()
		}

		if iterator.Key() != size {
			t.Failed()
		}

		fmt.Printf("%d", iterator.Key())
		iterator.Next()
		size++
	}

	if size != 10 {
		t.Failed()
	}
}

func TestRBTree_Insert2(t *testing.T) {
	tree := NewTreeMap(lessInt, equalInt)

	for i := 99; i >= 0; i-- {
		tree.Insert(i, i)
	}

	size := 0
	iterator := NewIterator(tree)
	for iterator.HasMore() {
		if iterator.Value() != size {
			t.Failed()
		}

		if iterator.Key() != size {
			t.Failed()
		}

		fmt.Printf("%d", iterator.Key())
		iterator.Next()
		size++
	}

	if size != 10 {
		t.Failed()
	}
}

func TestRBTree_InsertDelete(t *testing.T) {
	tree := NewTreeMap(lessInt, equalInt)

	for i := 0; i < 100; i++ {
		tree.Insert(i, i)
	}

	tree.Remove(10)
	if tree.Exists(10) {
		t.Fatalf("key exists but it has been deletted")
	}

	for i := 200; i < 300; i++ {
		tree.Insert(i, i)
	}

	tree.Remove(50)
	tree.Remove(100)
	tree.Remove(150)
	tree.Remove(201)
	tree.Remove(200)
	tree.Remove(249)
	tree.Remove(199)
	tree.Remove(0)

	for i := 300; i < 400; i++ {
		tree.Insert(i, i)
	}

	iterator := NewIterator(tree)
	for iterator.HasMore() {
		fmt.Printf("%d\n", iterator.Key())
		iterator.Next()
	}
}

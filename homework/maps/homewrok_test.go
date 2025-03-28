package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
)

// go test -v homework_test.go

type Node[T constraints.Ordered, V any] struct {
	key   T
	value V
	lnode *Node[T, V]
	rnode *Node[T, V]
}

type OrderedMap[T constraints.Ordered, V any] struct {
	root   *Node[T, V]
	nodeSz int
}

func NewOrderedMap[T constraints.Ordered, V any]() OrderedMap[T, V] {
	return OrderedMap[T, V]{}
}

func NewNode[T constraints.Ordered, V any](key T, val V) *Node[T, V] {
	return &Node[T, V]{key: key, value: val}
}

func (m *OrderedMap[T, V]) Insert(key T, value V) {
	if m.root == nil {
		m.root = NewNode(key, value)
		m.nodeSz++
		return
	}

	var tmp *Node[T, V] = m.root
	for {
		if tmp.key == key {
			tmp.value = value
			break
		}

		if tmp.key > key {
			if tmp.lnode != nil {
				tmp = tmp.lnode
				continue
			}
			tmp.lnode = NewNode(key, value)
			m.nodeSz++
			break
		}

		if tmp.key < key {
			if tmp.rnode != nil {
				tmp = tmp.rnode
				continue
			}
			tmp.rnode = NewNode(key, value)
			m.nodeSz++
			break
		}
	}
}

func (m *OrderedMap[T, V]) Erase(key T) {
	var parentNode, deleteNode *Node[T, V] = m.root.RecursiveContains(key)
	if parentNode == nil && deleteNode == nil {
		return
	}

	if deleteNode.rnode == nil {
		if parentNode.lnode == deleteNode {
			parentNode.lnode = deleteNode.lnode
		} else {
			parentNode.rnode = deleteNode.lnode
		}
		deleteNode.lnode = nil
	} else {
		var parntmp, tmp *Node[T, V] = deleteNode, deleteNode.rnode
		for tmp.lnode != nil {
			parntmp = tmp
			tmp = tmp.lnode
		}
		deleteNode.key = tmp.key
		if parntmp.lnode == tmp {
			parntmp.lnode = tmp.rnode
		} else {
			parntmp.rnode = tmp.rnode
		}
	}
	m.nodeSz--
}

func (node *Node[T, V]) RecursiveContains(key T) (*Node[T, V], *Node[T, V]) {
	if node == nil {
		return nil, nil
	}

	if node.lnode != nil && node.lnode.key == key {
		return node, node.lnode
	}
	if node.rnode != nil && node.rnode.key == key {
		return node, node.rnode
	}

	if node.key > key {
		return node.lnode.RecursiveContains(key)
	} else {
		return node.rnode.RecursiveContains(key)
	}
}

func (m *OrderedMap[T, V]) Contains(key T) bool {
	var tmp *Node[T, V] = m.root
	for tmp != nil {
		if tmp.key == key {
			return true
		}
		if tmp.key > key {
			tmp = tmp.lnode
		} else {
			tmp = tmp.rnode
		}
	}
	return false
}

func (m *OrderedMap[T, V]) Size() int {
	return m.nodeSz
}

func (node *Node[T, V]) RecursiveBypass(action func(T, V)) {
	if node == nil {
		return
	}

	node.lnode.RecursiveBypass(action)
	action(node.key, node.value)
	node.rnode.RecursiveBypass(action)
}

func (m *OrderedMap[T, V]) ForEach(action func(T, V)) {
	m.root.RecursiveBypass(action)
}

func (node *Node[T, V]) PrintNode() {
	fmt.Printf("%p : %v : <%p, %p> \n", node, node.key, node.lnode, node.rnode)
	if node.lnode != nil {
		node.lnode.PrintNode()
	}
	if node.rnode != nil {
		node.rnode.PrintNode()
	}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap[int, int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	//data.root.PrintNode()

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	//data.root.PrintNode()

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}

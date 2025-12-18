package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go
type OrderedMap struct {
	tree *node
	size int
}

type node struct {
	key   int
	value int
	left  *node
	right *node
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
	currentNode := &m.tree

	for {
		if *currentNode == nil {
			m.size++
			*currentNode = &node{
				key:   key,
				value: value,
			}
			return
		}

		if (*currentNode).key == key {
			(*currentNode).value = value
			return
		}

		if (*currentNode).key > key {
			currentNode = &(*currentNode).left
			continue
		}

		currentNode = &(*currentNode).right
	}
}

func (m *OrderedMap) Erase(key int) {
	currentNode := &m.tree
	for {
		if *currentNode == nil {
			return
		}

		if (*currentNode).key > key {
			currentNode = &(*currentNode).left
			continue
		}

		if (*currentNode).key < key {
			currentNode = &(*currentNode).right
			continue
		}

		// if no left branch then replacing with right branch
		if (*currentNode).left == nil && (*currentNode).right != nil {
			m.size--
			*currentNode = (*currentNode).right
			return
		}

		// if no right branch then replacing with left branch
		if (*currentNode).right == nil && (*currentNode).left != nil {
			m.size--
			*currentNode = (*currentNode).left
			return
		}

		if (*currentNode).right == nil && (*currentNode).left == nil {
			m.size--
			*currentNode = nil
			return
		}

		// go to the right branch
		replacingNode := &(*currentNode).right
		for {
			// if there is no left branch
			if (*replacingNode).left == nil {
				m.size--
				*currentNode = *replacingNode
				return
			}

			nextLeftNode := &(*replacingNode).left

			// see if next left node will not contain left child
			if (*nextLeftNode).left == nil {
				m.size--
				(*nextLeftNode).right = (*currentNode).right // update relations right
				(*nextLeftNode).left = (*currentNode).left   // update relations left
				*currentNode = *nextLeftNode
				(*replacingNode).left = nil // remove relation with parent
				return
			}

			replacingNode = nextLeftNode
		}
	}
}

func (m *OrderedMap) Contains(key int) bool {
	currentNode := m.tree
	for {
		if currentNode == nil {
			return false
		}

		if key == currentNode.key {
			return true
		}

		if currentNode.key > key {
			currentNode = currentNode.left
			continue
		}

		if currentNode.key < key {
			currentNode = currentNode.right
			continue
		}
	}
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	if m.tree == nil {
		return
	}

	subTreeWalk(m.tree, action)
}

func subTreeWalk(tree *node, action func(int, int)) {
	if tree == nil {
		return
	}

	subTreeWalk(tree.left, action)
	action(tree.key, tree.value)
	subTreeWalk(tree.right, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

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

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}

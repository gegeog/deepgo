package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type CircularQueue struct {
	values []int
	front  *int
	rear   *int
}

func NewCircularQueue(size int) CircularQueue {
	return CircularQueue{
		values: make([]int, size),
	}
}

func (q *CircularQueue) borderReached(p *int) bool {
	if p == &q.values[cap(q.values)-1] {
		return true
	}

	return false
}

func (q *CircularQueue) Push(value int) bool {
	if q.Full() {
		return false
	}

	if q.Empty() {
		q.rear = &q.values[0]
		q.front = &q.values[0]
		q.values[0] = value

		return true
	}

	if q.borderReached(q.rear) {
		q.rear = &q.values[0]
		*q.rear = value

		return true
	}

	rearPointer := unsafe.Pointer(q.rear)
	next := (*int)(unsafe.Add(rearPointer, unsafe.Sizeof(value)))

	q.rear = next
	*q.rear = value

	return true
}

func (q *CircularQueue) Pop() bool {
	if q.Empty() {
		return false
	}

	if q.borderReached(q.front) {
		q.front = &q.values[0]
		return true
	}

	if q.front == q.rear {
		q.front = nil
		q.rear = nil

		return true
	}

	pointer := unsafe.Pointer(q.front)
	next := (*int)(unsafe.Add(pointer, unsafe.Sizeof(0)))
	q.front = next

	return true
}

func (q *CircularQueue) Front() int {
	if q.Empty() {
		return -1
	}

	return *q.front
}

func (q *CircularQueue) Back() int {
	if q.Empty() {
		return -1
	}

	return *q.rear
}

func (q *CircularQueue) Empty() bool {
	if q.rear == nil && q.front == nil {
		return true
	}

	return false
}

func (q *CircularQueue) Full() bool {
	first := &q.values[0]
	last := &q.values[cap(q.values)-1]
	if q.rear == last && q.front == first {
		return true
	}

	r := unsafe.Pointer(&q.rear)
	next := (*int)(unsafe.Add(r, unsafe.Sizeof(q.values[0])))
	if next == q.front {
		return true
	}

	return false
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue(queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, -1, queue.Front())
	assert.Equal(t, -1, queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, 1, queue.Front())
	assert.Equal(t, 3, queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, 2, queue.Front())
	assert.Equal(t, 4, queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}

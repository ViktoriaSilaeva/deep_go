package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go
const INIT_VALUE = -1
const FIRST_IDX = 0

type Integer interface {
	int | int8 | int16 | int32 | int64
}

type CircularQueue[T Integer] struct {
	values []T
	front  int
	rear   int
}

func NewCircularQueue[T Integer](size int) CircularQueue[T] {
	return CircularQueue[T]{values: make([]T, size),
		front: INIT_VALUE,
		rear:  INIT_VALUE}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}
	if q.Empty() {
		q.front = FIRST_IDX
	}
	q.rear = (q.rear + 1) % cap(q.values)
	q.values[q.rear] = value
	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}
	if q.front == q.rear {
		q.front = INIT_VALUE
		q.rear = INIT_VALUE
	} else {
		q.front = (q.front + 1) % cap(q.values)
	}
	return true
}

func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return -1
	}
	return q.values[q.front]
}

func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return -1
	}
	return q.values[q.rear]
}

func (q *CircularQueue[T]) Empty() bool {
	return q.front == INIT_VALUE
}

func (q *CircularQueue[T]) Full() bool {
	return q.front == (q.rear+1)%cap(q.values)
}

func TestCircularQueue(t *testing.T) {
	const queueSize = 3
	queue := NewCircularQueue[int](queueSize)

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())

	assert.Equal(t, int(-1), queue.Front())
	assert.Equal(t, int(-1), queue.Back())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Push(1))
	assert.True(t, queue.Push(2))
	assert.True(t, queue.Push(3))
	assert.False(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

	assert.False(t, queue.Empty())
	assert.True(t, queue.Full())

	assert.Equal(t, int(1), queue.Front())
	assert.Equal(t, int(3), queue.Back())

	assert.True(t, queue.Pop())
	assert.False(t, queue.Empty())
	assert.False(t, queue.Full())
	assert.True(t, queue.Push(4))

	assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

	assert.Equal(t, int(2), queue.Front())
	assert.Equal(t, int(4), queue.Back())

	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.True(t, queue.Pop())
	assert.False(t, queue.Pop())

	assert.True(t, queue.Empty())
	assert.False(t, queue.Full())
}

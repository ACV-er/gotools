package storage

import (
	"sync"
)

var list_node_pool *sync.Pool

func init() {
	list_node_pool = &sync.Pool{
		New: func() interface{} {
			return new(list_node)
		},
	}
}

type list_node struct {
	val  interface{}
	next *list_node
}

type list_queue struct {
	head *list_node
	tail *list_node
	len  int
}

func NewListQueue() Queue {
	return &list_queue{}
}

func (lq *list_queue) Push(v interface{}) {
	new_node := list_node_pool.Get().(*list_node)
	new_node.val = v
	new_node.next = nil

	lq.len++
	if lq.len == 1 {
		lq.head = new_node
		lq.tail = new_node
		return
	}

	lq.tail.next = new_node
	lq.tail = new_node
}

func (lq *list_queue) Pop() interface{} {
	if lq.len == 0 {
		return nil
	}

	rel := lq.head.val
	tmp := lq.head
	defer func() {
		list_node_pool.Put(tmp)
	}()
	lq.head = lq.head.next
	lq.len--
	if lq.len == 0 {
		lq.tail = nil
	}

	return rel
}

func (lq *list_queue) Front() interface{} {
	if lq.len == 0 {
		return nil
	}

	return lq.head.val
}

func (lq *list_queue) Back() interface{} {
	if lq.len == 0 {
		return nil
	}

	return lq.tail.val
}

func (lq *list_queue) Size() int {
	return lq.len
}

func (lq *list_queue) Clean() {
	st := lq.head
	// 异步清理节点
	go func() {
		for st != nil {
			tmp := st
			st = st.next
			list_node_pool.Put(tmp)
		}
	}()

	lq.head = nil
	lq.tail = nil
	lq.len = 0
}

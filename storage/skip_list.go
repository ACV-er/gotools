package storage

import (
	"math/rand"
)

const (
	MAX_LEVEL = 32
	P_FACTOR  = 0.25
)

type CmpReq int

const (
	LE CmpReq = iota
	EQ CmpReq = iota
	GT CmpReq = iota
)

func randomLevel() int {
	level := 1
	for rand.Float64() <= P_FACTOR && level < MAX_LEVEL {
		level++
	}
	return level
}

type Node struct {
	val  Interface
	next []*Node
	pre  []*Node
}

func (n *Node) Next() *Node {
	if n == nil {
		return nil
	}
	return n.next[0]
}

func (n *Node) Pre() *Node {
	if n == nil {
		return nil
	}
	return n.pre[0]
}

func (n *Node) Load() Interface {
	return n.val
}

type Skiplist struct {
	curLevel int
	levelMin *Node
}

type Interface interface {
	Cmp(interface{}) CmpReq
}

func NewSkipList() Skiplist {
	return Skiplist{
		curLevel: 1,
		levelMin: &Node{
			next: make([]*Node, MAX_LEVEL),
		},
	}
}

func (s *Skiplist) Front() *Node {
	return s.levelMin.next[0]
}

func (s *Skiplist) findPreNode(val Interface, level int) *Node {
	cur := s.levelMin
	for i := s.curLevel - 1; i >= level; i-- {
		for cur.next[i] != nil && cur.next[i].val.Cmp(val) == LE {
			cur = cur.next[i]
		}
	}
	return cur
}

func (s *Skiplist) Search(target Interface) bool {
	pre_node := s.findPreNode(target, 0)
	if pre_node.next[0] == nil || pre_node.next[0].val.Cmp(target) != EQ {
		return false
	}
	return true
}

func (s *Skiplist) Add(val Interface) {
	need_level := randomLevel()
	newNode := &Node{
		val:  val,
		next: make([]*Node, need_level),
		pre:  make([]*Node, need_level),
	}

	if need_level > s.curLevel {
		s.curLevel = need_level
	}

	for level := need_level - 1; level >= 0; level-- {
		pre_node := s.findPreNode(val, level)
		newNode.next[level] = pre_node.next[level]
		if newNode.next[level] != nil {
			newNode.next[level].pre[level] = newNode
		}
		newNode.pre[level] = pre_node
		pre_node.next[level] = newNode
	}
}

func (s *Skiplist) Erase(val Interface) bool {
	pre_node := s.findPreNode(val, 0)
	if pre_node.next[0] == nil || pre_node.next[0].val.Cmp(val) != EQ {
		return false
	}

	node := pre_node.next[0]
	for level := len(node.next) - 1; level >= 0; level-- {
		node.pre[level].next[level] = node.next[level]
		if node.next[level] != nil {
			node.next[level].pre[level] = node.pre[level]
		}
	}

	return true
}

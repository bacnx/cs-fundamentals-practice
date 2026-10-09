package lru

import "fmt"

type node struct {
	key          int
	value        int
	next         *node
	prev         *node
	relativeNode *node
}

type LinkedList struct {
	head *node
	tail *node
}

func NewLinkedList() *LinkedList {
	head := &node{}
	tail := &node{}
	head.next = tail
	tail.prev = head
	l := &LinkedList{head: head, tail: tail}
	return l
}

func (l *LinkedList) PutHead(key int, val int, relNode *node) {
	second := l.head.next // the node after new node
	newNode := &node{key: key, value: val, relativeNode: relNode}

	l.head.next = newNode
	newNode.prev = l.head

	newNode.next = second
	second.prev = newNode
}

func (l *LinkedList) PutTail(key int, val int, relativeNode *node) {
	first := l.tail.next // the node before new node
	newNode := &node{key: key, value: val, relativeNode: relativeNode}

	first.next = newNode
	newNode.prev = first

	newNode.next = l.tail
	l.tail.prev = newNode
}

func (l *LinkedList) FindNodeByKey(key int) *node {
	if l.head == nil || l.head.next == nil {
		return nil
	}

	n := l.head.next
	for n.next != nil && n.next != l.tail {
		if n.key == key {
			return n
		}

		n = n.next
	}
	return nil
}

func (l *LinkedList) GetTail() *node {
	if l.tail.prev == l.head {
		return nil
	}
	return l.tail.prev
}

func (l *LinkedList) RemoveNode(n *node) error {
	if n == nil || n.prev == nil || n.next == nil {
		return fmt.Errorf("invalid node")
	}
	first := n.prev
	second := n.next
	first.next = second
	second.prev = first

	n.prev = nil
	n.next = nil
	return nil
}

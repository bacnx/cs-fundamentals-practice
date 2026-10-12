package lru

import "fmt"

type LRUList struct {
	l *LinkedList
}

const RLU_NODE_TYPE = "rlu-node"

func NewLRUList() *LRUList {
	return &LRUList{l: NewLinkedList()}
}

func (l *LRUList) PutHead(key int, val int, relNode *node) *node {
	return l.l.PutHead(key, val, relNode, RLU_NODE_TYPE)
}

func (l *LRUList) GetHead() *node {
	return l.l.GetHead()
}

func (l *LRUList) PopTail() *node {
	tail := l.l.GetTail()
	if err := l.l.RemoveNode(tail); err != nil {
		return nil
	}
	return tail
}

func (l *LRUList) Remove(n *node) error {
	if n == nil {
		return fmt.Errorf("the node is nil")
	}
	if n.nodeType != RLU_NODE_TYPE {
		return fmt.Errorf("the node type is not correct")
	}
	return l.l.RemoveNode(n)
}

func (l *LRUList) UpdateToHead(n *node) error {
	if n == nil {
		return fmt.Errorf("node is nil")
	}
	if n.nodeType != RLU_NODE_TYPE {
		return fmt.Errorf("the node type is not correct")
	}

	err := l.Remove(n)
	if err != nil {
		return err
	}

	newNode := l.PutHead(n.key, n.value, n.relativeNode)
	newNode.relativeNode.relativeNode = newNode
	return nil
}

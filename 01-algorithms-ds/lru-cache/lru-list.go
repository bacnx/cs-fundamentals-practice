package lru

type LRUList struct {
	l *LinkedList
}

func NewLRUList() *LRUList {
	return &LRUList{l: NewLinkedList()}
}

func (l *LRUList) PutHead(key int, val int, relNode *node) {
	l.l.PutHead(key, val, relNode)
}

func (l *LRUList) PopTail() *node {
	tail := l.l.GetTail()
	if err := l.l.RemoveNode(tail); err != nil {
		return nil
	}
	return tail
}

func (l *LRUList) PopByKey(key int) *node {
	n := l.l.FindNodeByKey(key)
	if n == nil {
		return nil
	}
	if err := l.l.RemoveNode(n); err != nil {
		return nil
	}
	return n
}

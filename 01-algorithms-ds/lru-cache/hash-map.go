package lru

import "fmt"

const _MOD = 100005

type HashMap struct {
	cap     int
	len     int
	buckets []*LinkedList
}

const HASH_NODE_TYPE = "hash-node"

func NewHashMap(cap int) *HashMap {
	buckets := make([]*LinkedList, cap)
	for i := range cap {
		buckets[i] = NewLinkedList()
	}

	return &HashMap{
		cap:     cap,
		len:     0,
		buckets: buckets,
	}
}

func (h *HashMap) Put(key int, value int, relNode *node) (*node, error) {
	if _, err := h.Get(key); err == nil {
		return nil, fmt.Errorf("this key already existed")
	}
	if h.len >= h.cap {
		return nil, fmt.Errorf("len is equal or lager than cap")
	}

	bucketKey := h.generateBucketKey(key)
	n := h.buckets[bucketKey].PutTail(key, value, relNode, HASH_NODE_TYPE)
	h.len++
	return n, nil
}

func (h *HashMap) Get(key int) (*node, error) {
	bucketKey := h.generateBucketKey(key)
	n := h.buckets[bucketKey].FindNodeByKey(key)
	if n == nil {
		return nil, fmt.Errorf("key %v is not exist", key)
	}
	return n, nil
}

func (h *HashMap) Del(key int) error {
	bucketKey := h.generateBucketKey(key)
	n := h.buckets[bucketKey].FindNodeByKey(key)
	if n == nil {
		return fmt.Errorf("key %v is not exist", key)
	}

	if err := h.buckets[bucketKey].RemoveNode(n); err != nil {
		return err
	}
	h.len--

	return nil
}

func (h *HashMap) RemoveNode(n *node) error {
	if n == nil {
		return fmt.Errorf("node is nil")
	}
	if n.nodeType != HASH_NODE_TYPE {
		return fmt.Errorf("the node type is not correct")
	}

	bucketKey := h.generateBucketKey(n.key)
	err := h.buckets[bucketKey].RemoveNode(n)
	if err != nil {
		return err
	}
	h.len--
	return nil
}

func (h *HashMap) generateBucketKey(key int) int {
	if h.cap == 0 {
		return 0
	}
	if key < 0 {
		key = -key
	}
	return (key%_MOD + _MOD) * _MOD % h.cap
}

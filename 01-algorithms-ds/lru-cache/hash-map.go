package lru

import "fmt"

type HashMap struct {
	cap     int
	len     int
	buckets []*LinkedList
}

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

func (h *HashMap) Put(key int, value int, relNode *node) error {
	if _, err := h.Get(key); err != nil {
		return fmt.Errorf("this key already existed")
	}
	if h.len >= h.cap {
		return fmt.Errorf("len is equal or lager cap")
	}

	bucketKey := h.generateBucketKey(key)
	h.buckets[bucketKey].PutTail(key, value, relNode)
	h.len++
	return nil
}

func (h *HashMap) Get(key int) (int, error) {
	bucketKey := h.generateBucketKey(key)
	n := h.buckets[bucketKey].FindNodeByKey(key)
	if n == nil {
		return 0, fmt.Errorf("key %v is not exist", key)
	}
	return n.value, nil
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

	return nil
}

func (h *HashMap) generateBucketKey(key int) int {
	if h.cap == 0 {
		return 0
	}
	if key < 0 {
		key = -key
	}
	return (key + 100005) * 100005 % h.cap
}

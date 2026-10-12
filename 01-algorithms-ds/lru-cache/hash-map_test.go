package lru

import (
	"math"
	"testing"
)

// mustPut fails the test if Put returns an error or a nil node.
func mustPut(t *testing.T, h *HashMap, key, value int) *node {
	t.Helper()
	n, err := h.Put(key, value, nil)
	if err != nil {
		t.Fatalf("Put(%d) error: %v", key, err)
	}
	if n == nil {
		t.Fatalf("Put(%d) returned a nil node", key)
	}
	return n
}

// mustGet fails the test unless Get finds key holding want.
func mustGet(t *testing.T, h *HashMap, key, want int) *node {
	t.Helper()
	n, err := h.Get(key)
	if err != nil {
		t.Fatalf("Get(%d) error: %v", key, err)
	}
	if n.value != want {
		t.Fatalf("Get(%d) value = %d, want %d", key, n.value, want)
	}
	return n
}

func mustMiss(t *testing.T, h *HashMap, key int) {
	t.Helper()
	if n, err := h.Get(key); err == nil {
		t.Fatalf("Get(%d) found value %d, want a miss", key, n.value)
	}
}

// --- Get / Put ---

func TestHashMapGetOnEmpty(t *testing.T) {
	defer failOnPanic(t)
	mustMiss(t, NewHashMap(4), 1)
}

func TestHashMapPutThenGet(t *testing.T) {
	defer failOnPanic(t)
	h := NewHashMap(4)
	put := mustPut(t, h, 1, 100)
	if got := mustGet(t, h, 1, 100); got != put {
		t.Fatal("Get must return the same node Put returned")
	}
}

func TestHashMapPutStoresNodeFields(t *testing.T) {
	defer failOnPanic(t)
	h := NewHashMap(4)
	rel := &node{key: 99}
	n, err := h.Put(1, 100, rel)
	if err != nil {
		t.Fatalf("Put error: %v", err)
	}
	if n.key != 1 || n.value != 100 || n.relativeNode != rel || n.nodeType != HASH_NODE_TYPE {
		t.Fatalf("node = key %d value %d rel %p type %q, want 1, 100, %p, %q",
			n.key, n.value, n.relativeNode, n.nodeType, rel, HASH_NODE_TYPE)
	}
}

func TestHashMapGetMissingKey(t *testing.T) {
	defer failOnPanic(t)
	h := NewHashMap(4)
	mustPut(t, h, 1, 100)
	mustMiss(t, h, 2)
}

func TestHashMapPutDuplicateKeyFails(t *testing.T) {
	defer failOnPanic(t)
	h := NewHashMap(4)
	mustPut(t, h, 1, 100)
	if _, err := h.Put(1, 200, nil); err == nil {
		t.Fatal("Put of an existing key must return an error")
	}
	mustGet(t, h, 1, 100) // the old value must be untouched
}

func TestHashMapPutWhenFullFails(t *testing.T) {
	defer failOnPanic(t)
	h := NewHashMap(2)
	mustPut(t, h, 1, 1)
	mustPut(t, h, 2, 2)
	if _, err := h.Put(3, 3, nil); err == nil {
		t.Fatal("Put into a full map must return an error")
	}
	mustMiss(t, h, 3)
}

// Fill every slot. With cap keys spread over cap buckets, some buckets must
// hold several nodes, so this also checks lookup inside a shared bucket.
func TestHashMapFillToCapacity(t *testing.T) {
	defer failOnPanic(t)
	const size = 64
	h := NewHashMap(size)
	for k := range size {
		mustPut(t, h, k, k*10)
	}
	for k := range size {
		mustGet(t, h, k, k*10)
	}
}

// Capacity 1: one bucket, every key collides.
func TestHashMapCapacityOne(t *testing.T) {
	defer failOnPanic(t)
	h := NewHashMap(1)
	mustPut(t, h, 5, 50)
	mustGet(t, h, 5, 50)
	mustMiss(t, h, 6)
}

// --- key range ---

// Keys are any int. Each case runs in its own subtest so one failing key
// does not hide the others.
func TestHashMapAnyIntKey(t *testing.T) {
	for _, tc := range []struct {
		name string
		key  int
	}{
		{"zero", 0},
		{"negative", -1},
		{"large negative", -123456789},
		{"MinInt", math.MinInt},
		{"MaxInt", math.MaxInt},
		{"large positive (overflow in hash)", math.MaxInt / 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer failOnPanic(t)
			h := NewHashMap(7)
			mustPut(t, h, tc.key, 42)
			mustGet(t, h, tc.key, 42)
		})
	}
}

// key and -key must stay different keys even if they hash to the same bucket.
func TestHashMapKeyAndNegationAreDistinct(t *testing.T) {
	defer failOnPanic(t)
	h := NewHashMap(4)
	mustPut(t, h, 5, 1)
	mustPut(t, h, -5, 2)
	mustGet(t, h, 5, 1)
	mustGet(t, h, -5, 2)
}

// --- Del ---

func TestHashMapDel(t *testing.T) {
	defer failOnPanic(t)
	h := NewHashMap(4)
	mustPut(t, h, 1, 10)
	mustPut(t, h, 2, 20)
	if err := h.Del(1); err != nil {
		t.Fatalf("Del(1) error: %v", err)
	}
	mustMiss(t, h, 1)
	mustGet(t, h, 2, 20)
}

func TestHashMapDelMissingKey(t *testing.T) {
	defer failOnPanic(t)
	h := NewHashMap(4)
	if err := h.Del(1); err == nil {
		t.Fatal("Del of a key never inserted must return an error")
	}
}

// Deleting must give the slot back: a full map with one key deleted accepts a new key.
func TestHashMapDelFreesCapacity(t *testing.T) {
	defer failOnPanic(t)
	h := NewHashMap(2)
	mustPut(t, h, 1, 1)
	mustPut(t, h, 2, 2)
	if err := h.Del(1); err != nil {
		t.Fatalf("Del(1) error: %v", err)
	}
	mustPut(t, h, 3, 3)
	mustGet(t, h, 3, 3)
}

func TestHashMapDelThenReinsert(t *testing.T) {
	defer failOnPanic(t)
	h := NewHashMap(2)
	mustPut(t, h, 1, 1)
	_ = h.Del(1)
	mustPut(t, h, 1, 11)
	mustGet(t, h, 1, 11)
}

// --- RemoveNode ---

func TestHashMapRemoveNode(t *testing.T) {
	defer failOnPanic(t)
	h := NewHashMap(4)
	n := mustPut(t, h, 1, 10)
	mustPut(t, h, 2, 20)
	if err := h.RemoveNode(n); err != nil {
		t.Fatalf("RemoveNode error: %v", err)
	}
	mustMiss(t, h, 1)
	mustGet(t, h, 2, 20)
}

func TestHashMapRemoveNodeFreesCapacity(t *testing.T) {
	defer failOnPanic(t)
	h := NewHashMap(1)
	n := mustPut(t, h, 1, 1)
	if err := h.RemoveNode(n); err != nil {
		t.Fatalf("RemoveNode error: %v", err)
	}
	mustPut(t, h, 2, 2)
}

func TestHashMapRemoveNodeNil(t *testing.T) {
	defer failOnPanic(t)
	if err := NewHashMap(4).RemoveNode(nil); err == nil {
		t.Fatal("RemoveNode(nil) must return an error")
	}
}

// A node from the recency list must be rejected, and the map left unchanged.
func TestHashMapRemoveNodeWrongType(t *testing.T) {
	defer failOnPanic(t)
	h := NewHashMap(4)
	mustPut(t, h, 1, 10)
	foreign := NewLRUList().PutHead(1, 10, nil)
	if err := h.RemoveNode(foreign); err == nil {
		t.Fatal("RemoveNode of a recency-list node must return an error")
	}
	mustGet(t, h, 1, 10)
}

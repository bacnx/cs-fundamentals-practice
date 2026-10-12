package lru

import "testing"

// putKeys pushes keys at the head in order, so the last key ends up most recent.
func putKeys(l *LRUList, keys ...int) map[int]*node {
	nodes := map[int]*node{}
	for _, k := range keys {
		nodes[k] = l.PutHead(k, k*10, nil)
	}
	return nodes
}

// --- PutHead / PopTail ---

func TestLRUListPopTailOnEmpty(t *testing.T) {
	defer failOnPanic(t)
	if n := NewLRUList().PopTail(); n != nil {
		t.Fatalf("PopTail on empty list = key %d, want nil", n.key)
	}
}

func TestLRUListPutHeadStoresFields(t *testing.T) {
	defer failOnPanic(t)
	rel := &node{key: 99}
	n := NewLRUList().PutHead(1, 10, rel)
	if n.key != 1 || n.value != 10 || n.relativeNode != rel || n.nodeType != RLU_NODE_TYPE {
		t.Fatalf("node = key %d value %d rel %p type %q, want 1, 10, %p, %q",
			n.key, n.value, n.relativeNode, n.nodeType, rel, RLU_NODE_TYPE)
	}
}

// PopTail returns entries least recent first, then nil once empty.
func TestLRUListPopTailOrder(t *testing.T) {
	defer failOnPanic(t)
	l := NewLRUList()
	nodes := putKeys(l, 1, 2, 3)
	for _, want := range []int{1, 2, 3} {
		n := l.PopTail()
		if n != nodes[want] {
			t.Fatalf("PopTail did not return the node for key %d", want)
		}
	}
	if n := l.PopTail(); n != nil {
		t.Fatalf("PopTail after draining = key %d, want nil", n.key)
	}
	assertKeys(t, l.l, []int{})
}

// The popped node is handed back to the caller, so it must still carry its
// link to the hash node.
func TestLRUListPopTailKeepsRelativeNode(t *testing.T) {
	defer failOnPanic(t)
	l := NewLRUList()
	rel := &node{key: 1, nodeType: HASH_NODE_TYPE}
	l.PutHead(1, 10, rel)
	if n := l.PopTail(); n == nil || n.relativeNode != rel {
		t.Fatal("popped node lost its relativeNode")
	}
}

// --- Remove ---

func TestLRUListRemoveMiddle(t *testing.T) {
	defer failOnPanic(t)
	l := NewLRUList()
	nodes := putKeys(l, 1, 2, 3)
	if err := l.Remove(nodes[2]); err != nil {
		t.Fatalf("Remove error: %v", err)
	}
	assertKeys(t, l.l, []int{3, 1})
}

func TestLRUListRemoveNil(t *testing.T) {
	defer failOnPanic(t)
	if err := NewLRUList().Remove(nil); err == nil {
		t.Fatal("Remove(nil) must return an error")
	}
}

// A hash-map node must be rejected, and the list left unchanged.
func TestLRUListRemoveWrongType(t *testing.T) {
	defer failOnPanic(t)
	l := NewLRUList()
	putKeys(l, 1)
	foreign := &node{key: 1, nodeType: HASH_NODE_TYPE}
	if err := l.Remove(foreign); err == nil {
		t.Fatal("Remove of a hash-map node must return an error")
	}
	assertKeys(t, l.l, []int{1})
}

// --- UpdateToHead ---

func TestLRUListUpdateToHeadMovesToFront(t *testing.T) {
	defer failOnPanic(t)
	l := NewLRUList()
	nodes := putKeys(l, 1, 2, 3) // list: 3 2 1
	if err := l.UpdateToHead(nodes[1]); err != nil {
		t.Fatalf("UpdateToHead error: %v", err)
	}
	assertKeys(t, l.l, []int{1, 3, 2})
	if n := l.PopTail(); n == nil || n.key != 2 {
		t.Fatal("after refreshing 1, key 2 must be the least recent")
	}
}

func TestLRUListUpdateToHeadOnHeadIsNoOp(t *testing.T) {
	defer failOnPanic(t)
	l := NewLRUList()
	nodes := putKeys(l, 1, 2)
	if err := l.UpdateToHead(nodes[2]); err != nil {
		t.Fatalf("UpdateToHead error: %v", err)
	}
	assertKeys(t, l.l, []int{2, 1})
}

// The hash node keeps a pointer to its recency node. That pointer must still
// be the live node after a refresh, or the next refresh or removal through it
// operates on a detached node.
func TestLRUListUpdateToHeadKeepsNodeHandleValid(t *testing.T) {
	defer failOnPanic(t)
	l := NewLRUList()
	nodes := putKeys(l, 1, 2)
	n := nodes[1]
	if err := l.UpdateToHead(n); err != nil {
		t.Fatalf("first UpdateToHead error: %v", err)
	}
	if err := l.UpdateToHead(n); err != nil {
		t.Fatalf("second UpdateToHead on the same handle error: %v", err)
	}
	if err := l.Remove(n); err != nil {
		t.Fatalf("Remove after UpdateToHead error: %v", err)
	}
	assertKeys(t, l.l, []int{2})
}

func TestLRUListUpdateToHeadKeepsFields(t *testing.T) {
	defer failOnPanic(t)
	l := NewLRUList()
	rel := &node{key: 1, nodeType: HASH_NODE_TYPE}
	n := l.PutHead(1, 10, rel)
	l.PutHead(2, 20, nil)
	if err := l.UpdateToHead(n); err != nil {
		t.Fatalf("UpdateToHead error: %v", err)
	}
	l.PopTail() // drop key 2
	got := l.PopTail()
	if got == nil || got.key != 1 || got.value != 10 || got.relativeNode != rel {
		t.Fatal("refreshed entry lost its key, value or relativeNode")
	}
}

func TestLRUListUpdateToHeadNil(t *testing.T) {
	defer failOnPanic(t)
	if err := NewLRUList().UpdateToHead(nil); err == nil {
		t.Fatal("UpdateToHead(nil) must return an error")
	}
}

func TestLRUListUpdateToHeadWrongType(t *testing.T) {
	defer failOnPanic(t)
	l := NewLRUList()
	putKeys(l, 1, 2)
	foreign := &node{key: 1, nodeType: HASH_NODE_TYPE}
	if err := l.UpdateToHead(foreign); err == nil {
		t.Fatal("UpdateToHead of a hash-map node must return an error")
	}
	assertKeys(t, l.l, []int{2, 1})
}

// A node already removed from the list cannot be refreshed.
func TestLRUListUpdateToHeadDetachedNode(t *testing.T) {
	defer failOnPanic(t)
	l := NewLRUList()
	nodes := putKeys(l, 1, 2)
	_ = l.Remove(nodes[1])
	if err := l.UpdateToHead(nodes[1]); err == nil {
		t.Fatal("UpdateToHead of a removed node must return an error")
	}
	assertKeys(t, l.l, []int{2})
}

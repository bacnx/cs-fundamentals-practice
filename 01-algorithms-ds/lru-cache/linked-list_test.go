package lru

import (
	"slices"
	"testing"
)

// failOnPanic turns a panic into a test failure, so one crashing test does not
// abort the whole package run and hide every result after it.
// Use as the first line of a test: defer failOnPanic(t)
func failOnPanic(t *testing.T) {
	t.Helper()
	if r := recover(); r != nil {
		t.Fatalf("panicked: %v", r)
	}
}

// listKeys walks l from head to tail, checks that every next link is matched
// by the reverse prev link, and returns the keys in order.
func listKeys(t *testing.T, l *LinkedList) []int {
	t.Helper()
	keys := []int{}
	prev := l.head
	for n := l.head.next; n != l.tail; n = n.next {
		if n == nil {
			t.Fatal("reached nil before the tail sentinel")
		}
		if n.prev != prev {
			t.Fatalf("broken back link: node %d has prev != the node before it", n.key)
		}
		if len(keys) > 1000 {
			t.Fatal("list looks cyclic")
		}
		keys = append(keys, n.key)
		prev = n
	}
	if l.tail.prev != prev {
		t.Fatal("tail.prev does not point at the last node")
	}
	return keys
}

func assertKeys(t *testing.T, l *LinkedList, want []int) {
	t.Helper()
	if got := listKeys(t, l); !slices.Equal(got, want) {
		t.Fatalf("list = %v, want %v", got, want)
	}
}

// --- construction ---

func TestLinkedListNewIsEmpty(t *testing.T) {
	defer failOnPanic(t)
	l := NewLinkedList()
	assertKeys(t, l, []int{})
	if n := l.GetTail(); n != nil {
		t.Fatalf("GetTail on empty list = key %d, want nil", n.key)
	}
}

// --- PutHead / PutTail ---

func TestLinkedListPutHeadOrder(t *testing.T) {
	defer failOnPanic(t)
	l := NewLinkedList()
	l.PutHead(1, 10, nil, "t")
	l.PutHead(2, 20, nil, "t")
	l.PutHead(3, 30, nil, "t")
	assertKeys(t, l, []int{3, 2, 1})
}

func TestLinkedListPutTailOrder(t *testing.T) {
	defer failOnPanic(t)
	l := NewLinkedList()
	l.PutTail(1, 10, nil, "t")
	l.PutTail(2, 20, nil, "t")
	l.PutTail(3, 30, nil, "t")
	assertKeys(t, l, []int{1, 2, 3})
}

func TestLinkedListPutTailOnEmpty(t *testing.T) {
	defer failOnPanic(t)
	l := NewLinkedList()
	l.PutTail(1, 10, nil, "t")
	assertKeys(t, l, []int{1})
}

func TestLinkedListMixedPuts(t *testing.T) {
	defer failOnPanic(t)
	l := NewLinkedList()
	l.PutTail(2, 0, nil, "t")
	l.PutHead(1, 0, nil, "t")
	l.PutTail(3, 0, nil, "t")
	assertKeys(t, l, []int{1, 2, 3})
}

func TestLinkedListPutStoresFields(t *testing.T) {
	defer failOnPanic(t)
	l := NewLinkedList()
	rel := &node{key: 99}
	for name, put := range map[string]func(int, int, *node, string) *node{
		"PutHead": l.PutHead,
		"PutTail": l.PutTail,
	} {
		n := put(7, 70, rel, "kind")
		if n == nil {
			t.Fatalf("%s returned nil", name)
		}
		if n.key != 7 || n.value != 70 || n.relativeNode != rel || n.nodeType != "kind" {
			t.Fatalf("%s stored key=%d value=%d rel=%p type=%q, want 7,70,%p,\"kind\"",
				name, n.key, n.value, n.relativeNode, n.nodeType, rel)
		}
	}
}

// --- GetTail ---

func TestLinkedListGetTailReturnsLastNode(t *testing.T) {
	defer failOnPanic(t)
	l := NewLinkedList()
	l.PutHead(1, 0, nil, "t")
	l.PutHead(2, 0, nil, "t")
	if n := l.GetTail(); n == nil || n.key != 1 {
		t.Fatalf("GetTail = %v, want key 1", n)
	}
	assertKeys(t, l, []int{2, 1}) // GetTail must not modify the list
}

// --- FindNodeByKey ---

func TestLinkedListFindOnEmpty(t *testing.T) {
	defer failOnPanic(t)
	if n := NewLinkedList().FindNodeByKey(1); n != nil {
		t.Fatal("FindNodeByKey on empty list must return nil")
	}
}

func TestLinkedListFindSingleNode(t *testing.T) {
	defer failOnPanic(t)
	l := NewLinkedList()
	want := l.PutHead(1, 10, nil, "t")
	if got := l.FindNodeByKey(1); got != want {
		t.Fatal("FindNodeByKey(1) did not return the only node")
	}
}

// Every position must be searchable: first, middle and last.
func TestLinkedListFindEveryPosition(t *testing.T) {
	defer failOnPanic(t)
	l := NewLinkedList()
	nodes := map[int]*node{}
	for _, k := range []int{1, 2, 3} {
		nodes[k] = l.PutTail(k, k*10, nil, "t")
	}
	for _, k := range []int{1, 2, 3} {
		if got := l.FindNodeByKey(k); got != nodes[k] {
			t.Errorf("FindNodeByKey(%d) did not return the node holding key %d", k, k)
		}
	}
}

func TestLinkedListFindMissing(t *testing.T) {
	defer failOnPanic(t)
	l := NewLinkedList()
	l.PutTail(1, 0, nil, "t")
	l.PutTail(2, 0, nil, "t")
	if n := l.FindNodeByKey(3); n != nil {
		t.Fatal("FindNodeByKey(3) must return nil for a key never inserted")
	}
}

// The sentinels have key 0 because it is the zero value. Searching for 0 in a
// list that does not hold it must not return a sentinel.
func TestLinkedListFindZeroKeyIgnoresSentinels(t *testing.T) {
	defer failOnPanic(t)
	l := NewLinkedList()
	if n := l.FindNodeByKey(0); n != nil {
		t.Fatal("FindNodeByKey(0) on empty list returned a node (a sentinel?)")
	}
	l.PutTail(1, 0, nil, "t")
	if n := l.FindNodeByKey(0); n != nil {
		t.Fatal("FindNodeByKey(0) returned a node though key 0 was never inserted")
	}
}

// --- RemoveNode ---

func TestLinkedListRemoveEveryPosition(t *testing.T) {
	for _, tc := range []struct {
		name   string
		remove int
		want   []int
	}{
		{"first", 1, []int{2, 3}},
		{"middle", 2, []int{1, 3}},
		{"last", 3, []int{1, 2}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer failOnPanic(t)
			l := NewLinkedList()
			nodes := map[int]*node{}
			for _, k := range []int{1, 2, 3} {
				nodes[k] = l.PutTail(k, 0, nil, "t")
			}
			if err := l.RemoveNode(nodes[tc.remove]); err != nil {
				t.Fatalf("RemoveNode(%d) error: %v", tc.remove, err)
			}
			assertKeys(t, l, tc.want)
		})
	}
}

func TestLinkedListRemoveOnlyNodeLeavesEmpty(t *testing.T) {
	defer failOnPanic(t)
	l := NewLinkedList()
	n := l.PutHead(1, 0, nil, "t")
	if err := l.RemoveNode(n); err != nil {
		t.Fatalf("RemoveNode error: %v", err)
	}
	assertKeys(t, l, []int{})
	if l.GetTail() != nil {
		t.Fatal("GetTail must be nil after removing the only node")
	}
}

func TestLinkedListRemoveNil(t *testing.T) {
	defer failOnPanic(t)
	if err := NewLinkedList().RemoveNode(nil); err == nil {
		t.Fatal("RemoveNode(nil) must return an error")
	}
}

// Removing the same node twice must fail the second time and leave the list intact.
func TestLinkedListRemoveTwice(t *testing.T) {
	defer failOnPanic(t)
	l := NewLinkedList()
	n := l.PutTail(1, 0, nil, "t")
	l.PutTail(2, 0, nil, "t")
	if err := l.RemoveNode(n); err != nil {
		t.Fatalf("first RemoveNode error: %v", err)
	}
	if err := l.RemoveNode(n); err == nil {
		t.Fatal("second RemoveNode of the same node must return an error")
	}
	assertKeys(t, l, []int{2})
}

// A removed node can be put back and found again.
func TestLinkedListReuseAfterRemove(t *testing.T) {
	defer failOnPanic(t)
	l := NewLinkedList()
	n := l.PutTail(1, 0, nil, "t")
	_ = l.RemoveNode(n)
	l.PutTail(1, 0, nil, "t")
	assertKeys(t, l, []int{1})
	if l.FindNodeByKey(1) == nil {
		t.Fatal("re-inserted key 1 not found")
	}
}

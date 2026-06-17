package hashmap

import (
	"fmt"
	"testing"
)

func TestGetMissingKey(t *testing.T) {
	h := NewHash(16)
	_, ok := h.Get("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

func TestPutThenGet(t *testing.T) {
	h := NewHash(16)
	h.Put("foo", 42)
	v, ok := h.Get("foo")
	if !ok || v != 42 {
		t.Fatalf("expected 42 got %v (ok=%v)", v, ok)
	}
}

func TestPutUpdatesExistingKey(t *testing.T) {
	h := NewHash(16)
	h.Put("foo", 1)
	h.Put("foo", 2)
	v, ok := h.Get("foo")
	if !ok || v != 2 {
		t.Fatalf("expected 2 got %v", v)
	}
}

func TestDeleteRemovesKey(t *testing.T) {
	h := NewHash(16)
	h.Put("foo", 1)
	h.Delete("foo")
	_, ok := h.Get("foo")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestCollisionBothKeysAccessible(t *testing.T) {
	// Use a tiny capacity to force collisions
	h := NewHash(1)
	h.Put("a", 1)
	h.Put("b", 2)
	v1, ok1 := h.Get("a")
	v2, ok2 := h.Get("b")
	if !ok1 || v1 != 1 {
		t.Fatalf("expected a=1 got %v (ok=%v)", v1, ok1)
	}
	if !ok2 || v2 != 2 {
		t.Fatalf("expected b=2 got %v (ok=%v)", v2, ok2)
	}
}

func TestDeleteDoesNotAffectOtherKeys(t *testing.T) {
	h := NewHash(1) // force collision
	h.Put("a", 1)
	h.Put("b", 2)
	h.Delete("a")
	_, ok := h.Get("a")
	if ok {
		t.Fatal("expected a to be deleted")
	}
	v, ok := h.Get("b")
	if !ok || v != 2 {
		t.Fatalf("expected b=2 got %v", v)
	}
}

func TestMultipleKeys(t *testing.T) {
	h := NewHash(16)
	data := map[string]int{"x": 10, "y": 20, "z": 30}
	for k, v := range data {
		h.Put(k, v)
	}
	for k, want := range data {
		got, ok := h.Get(k)
		if !ok || got != want {
			t.Fatalf("key %q: expected %d got %v", k, want, got)
		}
	}
}

// Level 2 — Load factor and auto-resize

func TestResizeKeepsAllKeys(t *testing.T) {
	h := NewHash(4) // resize triggers at 75% of cap=4, i.e. after 3 inserts
	keys := []string{"a", "b", "c", "d", "e", "f"}
	for i, k := range keys {
		h.Put(k, i)
	}
	for i, k := range keys {
		got, ok := h.Get(k)
		if !ok || got != i {
			t.Fatalf("after resize: key %q expected %d got %v (ok=%v)", k, i, got, ok)
		}
	}
}

func TestMultipleResizesKeepAllKeys(t *testing.T) {
	h := NewHash(4)
	n := 100
	for i := 0; i < n; i++ {
		h.Put(fmt.Sprintf("%d", i), i)
	}
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("%d", i)
		got, ok := h.Get(key)
		if !ok || got != i {
			t.Fatalf("key %q: expected %d got %v (ok=%v)", key, i, got, ok)
		}
	}
}

func TestLenAccurateAfterUpdate(t *testing.T) {
	h := NewHash(16)
	h.Put("a", 1)
	h.Put("a", 2) // update, not a new key
	if got := h.Len(); got != 1 {
		t.Fatalf("Len should be 1 after updating same key, got %d", got)
	}
}

func TestLenDecreasesOnDelete(t *testing.T) {
	h := NewHash(16)
	h.Put("a", 1)
	h.Put("b", 2)
	h.Delete("a")
	if got := h.Len(); got != 1 {
		t.Fatalf("Len should be 1 after deleting one key, got %d", got)
	}
}

func TestLenUnchangedWhenDeletingMissingKey(t *testing.T) {
	h := NewHash(16)
	h.Put("a", 1)
	h.Delete("nonexistent")
	if got := h.Len(); got != 1 {
		t.Fatalf("Len should remain 1 after deleting nonexistent key, got %d", got)
	}
}

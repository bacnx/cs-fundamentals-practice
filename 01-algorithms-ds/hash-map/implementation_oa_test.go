package hashmap

import (
	"fmt"
	"testing"
)

// Level 3 — Open addressing (linear probing)
// Expects a HashOA type with NewHashOA(cap int) constructor.

func TestOAGetMissingKey(t *testing.T) {
	h := NewHashOA(16)
	_, ok := h.Get("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

func TestOAPutThenGet(t *testing.T) {
	h := NewHashOA(16)
	h.Put("foo", 42)
	v, ok := h.Get("foo")
	if !ok || v != 42 {
		t.Fatalf("expected 42 got %v (ok=%v)", v, ok)
	}
}

func TestOAUpdateExistingKey(t *testing.T) {
	h := NewHashOA(16)
	h.Put("foo", 1)
	h.Put("foo", 2)
	v, ok := h.Get("foo")
	if !ok || v != 2 {
		t.Fatalf("expected 2 got %v", v)
	}
}

func TestOADeleteRemovesKey(t *testing.T) {
	h := NewHashOA(16)
	h.Put("foo", 1)
	h.Delete("foo")
	_, ok := h.Get("foo")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

// The critical tombstone test: deleting a key in the middle of a probe chain
// must not make subsequent keys in that chain unreachable.
func TestOADeleteDoesNotBreakProbeChain(t *testing.T) {
	h := NewHashOA(16)
	keys := []string{"a", "b", "c", "d", "e", "f", "g"}
	for i, k := range keys {
		h.Put(k, i)
	}
	h.Delete("b")
	h.Delete("e")

	deleted := map[string]bool{"b": true, "e": true}
	for i, k := range keys {
		if deleted[k] {
			_, ok := h.Get(k)
			if ok {
				t.Fatalf("key %q should be deleted", k)
			}
			continue
		}
		got, ok := h.Get(k)
		if !ok || got != i {
			t.Fatalf("key %q: expected %d got %v (ok=%v)", k, i, got, ok)
		}
	}
}

// A tombstone slot must be reusable by a subsequent Put.
func TestOATombstoneSlotReusable(t *testing.T) {
	h := NewHashOA(16)
	h.Put("a", 1)
	h.Delete("a")
	h.Put("a", 2)
	v, ok := h.Get("a")
	if !ok || v != 2 {
		t.Fatalf("expected 2 after re-insert, got %v (ok=%v)", v, ok)
	}
	if l := h.Len(); l != 1 {
		t.Fatalf("Len should be 1 after delete+reinsert, got %d", l)
	}
}

func TestOALenAccurateAfterUpdate(t *testing.T) {
	h := NewHashOA(16)
	h.Put("a", 1)
	h.Put("a", 2) // update, not a new key
	if got := h.Len(); got != 1 {
		t.Fatalf("Len should be 1 after updating same key, got %d", got)
	}
}

func TestOALenUnchangedWhenDeletingMissingKey(t *testing.T) {
	h := NewHashOA(16)
	h.Put("a", 1)
	h.Delete("nonexistent")
	if got := h.Len(); got != 1 {
		t.Fatalf("Len should remain 1 after deleting nonexistent key, got %d", got)
	}
}

func TestOAResizeKeepsAllKeys(t *testing.T) {
	h := NewHashOA(4) // threshold 0.5: resize triggers after >2 keys
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

func TestOAMultipleResizesKeepAllKeys(t *testing.T) {
	h := NewHashOA(4)
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

package hashmap

import "testing"

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

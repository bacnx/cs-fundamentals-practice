package lru

import "testing"

// --- Level 1: basic get/put ---

func TestGetOnEmpty(t *testing.T) {
	c := NewLRU(2)
	if _, ok := c.Get(1); ok {
		t.Fatal("expected miss on empty cache")
	}
}

func TestPutThenGet(t *testing.T) {
	c := NewLRU(2)
	c.Put(1, 100)
	v, ok := c.Get(1)
	if !ok || v != 100 {
		t.Fatalf("expected 100 got %v (ok=%v)", v, ok)
	}
}

func TestGetMissingKey(t *testing.T) {
	c := NewLRU(2)
	c.Put(1, 100)
	if _, ok := c.Get(2); ok {
		t.Fatal("expected miss for key never inserted")
	}
}

func TestPutUpdatesValue(t *testing.T) {
	c := NewLRU(2)
	c.Put(1, 100)
	c.Put(1, 200) // update, same key
	v, ok := c.Get(1)
	if !ok || v != 200 {
		t.Fatalf("expected updated value 200 got %v", v)
	}
	if c.Len() != 1 {
		t.Fatalf("update must not grow size, Len=%d want 1", c.Len())
	}
}

// --- Level 2: capacity & eviction ---

func TestLenNeverExceedsCapacity(t *testing.T) {
	c := NewLRU(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(3, 3) // forces an eviction
	if c.Len() != 2 {
		t.Fatalf("Len=%d, must never exceed capacity 2", c.Len())
	}
}

func TestEvictsLeastRecentlyUsed(t *testing.T) {
	c := NewLRU(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(3, 3) // 1 is the LRU -> evicted
	if _, ok := c.Get(1); ok {
		t.Fatal("key 1 should have been evicted as least recently used")
	}
	if v, ok := c.Get(2); !ok || v != 2 {
		t.Fatalf("key 2 should survive, got %v (ok=%v)", v, ok)
	}
	if v, ok := c.Get(3); !ok || v != 3 {
		t.Fatalf("key 3 should be present, got %v (ok=%v)", v, ok)
	}
}

// The defining behavior: a Get must REFRESH recency, changing who gets evicted.
func TestGetRefreshesRecency(t *testing.T) {
	c := NewLRU(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Get(1)    // now 1 is most-recent, 2 is least-recent
	c.Put(3, 3) // should evict 2, NOT 1
	if _, ok := c.Get(2); ok {
		t.Fatal("key 2 should have been evicted after 1 was refreshed by Get")
	}
	if v, ok := c.Get(1); !ok || v != 1 {
		t.Fatalf("key 1 should survive because Get refreshed it, got %v (ok=%v)", v, ok)
	}
}

// Updating an existing key via Put must also refresh recency.
func TestPutRefreshesRecency(t *testing.T) {
	c := NewLRU(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(1, 10) // update 1 -> most-recent
	c.Put(3, 3)  // should evict 2
	if _, ok := c.Get(2); ok {
		t.Fatal("key 2 should be evicted; updating key 1 should have refreshed it")
	}
	if v, ok := c.Get(1); !ok || v != 10 {
		t.Fatalf("key 1 should survive with value 10, got %v (ok=%v)", v, ok)
	}
}

func TestCapacityOne(t *testing.T) {
	c := NewLRU(1)
	c.Put(1, 1)
	c.Put(2, 2) // evicts 1 immediately
	if _, ok := c.Get(1); ok {
		t.Fatal("with capacity 1, key 1 must be evicted by key 2")
	}
	if v, ok := c.Get(2); !ok || v != 2 {
		t.Fatalf("key 2 should be present, got %v (ok=%v)", v, ok)
	}
}

// --- Level 3: sequence (LeetCode 146 style) ---

func TestLeetCodeSequence(t *testing.T) {
	c := NewLRU(2)
	c.Put(1, 1)
	c.Put(2, 2)
	if v, ok := c.Get(1); !ok || v != 1 { // returns 1
		t.Fatalf("Get(1) = %v,%v want 1,true", v, ok)
	}
	c.Put(3, 3) // evicts key 2
	if _, ok := c.Get(2); ok {
		t.Fatal("Get(2) should miss (evicted)")
	}
	c.Put(4, 4) // evicts key 1
	if _, ok := c.Get(1); ok {
		t.Fatal("Get(1) should miss (evicted)")
	}
	if v, ok := c.Get(3); !ok || v != 3 {
		t.Fatalf("Get(3) = %v,%v want 3,true", v, ok)
	}
	if v, ok := c.Get(4); !ok || v != 4 {
		t.Fatalf("Get(4) = %v,%v want 4,true", v, ok)
	}
}

// Re-inserting an evicted key should behave like a fresh insert.
func TestReinsertEvictedKey(t *testing.T) {
	c := NewLRU(2)
	c.Put(1, 1)
	c.Put(2, 2)
	c.Put(3, 3) // evict 1
	c.Put(1, 11)
	if v, ok := c.Get(1); !ok || v != 11 {
		t.Fatalf("re-inserted key 1 should be 11, got %v (ok=%v)", v, ok)
	}
	if c.Len() != 2 {
		t.Fatalf("Len=%d want 2", c.Len())
	}
}

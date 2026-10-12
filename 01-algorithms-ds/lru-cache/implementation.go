package lru

// LRU is a fixed-capacity cache that evicts the least-recently-used entry
// when a new key would exceed capacity.
//
// Target complexity: BOTH Get and Put must run in O(1) worst case.
// Keys are any int, including negative values, 0 and math.MinInt.
type LRU struct {
	hash *HashMap
	lru  *LRUList
	cap  int
	len  int
}

// NewLRU returns an empty cache that holds at most `capacity` entries.
// You may assume capacity >= 1 (a test covers capacity == 1).
func NewLRU(capacity int) *LRU {
	if capacity < 1 {
		capacity = 1
	}

	h := NewHashMap(capacity)
	l := NewLRUList()
	return &LRU{
		hash: h,
		lru:  l,
		cap:  capacity,
		len:  0,
	}
}

// Get returns the value for key and true if present, marking it as the most
// recently used. If absent, it returns (0, false) and changes nothing.
func (l *LRU) Get(key int) (int, bool) {
	n, err := l.hash.Get(key)
	if err != nil {
		return 0, false
	}
	err = l.lru.UpdateToHead(n.relativeNode)
	if err != nil {
		return 0, false
	}
	return n.value, true
}

// Put inserts or updates key->value, marking it most recently used.
// If inserting a new key exceeds capacity, evict the least recently used entry
// first. Updating an existing key must NOT grow the size.
func (l *LRU) Put(key int, value int) {
	n, err := l.hash.Get(key)
	if err == nil {
		n.value = value
		err = l.lru.Remove(n.relativeNode)
		if err != nil {
			return
		}
		err = l.hash.RemoveNode(n)
		if err != nil {
			return
		}
		l.len--
	}

	if l.len >= l.cap {
		lastUsedNode := l.lru.PopTail()
		if lastUsedNode == nil {
			return
		}
		hashNode := lastUsedNode.relativeNode
		err = l.hash.RemoveNode(hashNode)
		if err != nil {
			return
		}
		l.len--
	}

	newHashNode, err := l.hash.Put(key, value, nil)
	if err != nil {
		return
	}
	newLRUNode := l.lru.PutHead(key, value, newHashNode)
	newHashNode.relativeNode = newLRUNode

	l.len++
}

// Len reports the current number of stored entries (for tests/invariants).
// It must never exceed capacity.
func (l *LRU) Len() int {
	return l.len
}

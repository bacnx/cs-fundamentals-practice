package lru

// LRU is a fixed-capacity cache that evicts the least-recently-used entry
// when a new key would exceed capacity.
//
// Target complexity: BOTH Get and Put must run in O(1) worst case.
// Keys are any int, including negative values, 0 and math.MinInt.
type LRU struct {
	// TODO: add the fields your design needs
}

// NewLRU returns an empty cache that holds at most `capacity` entries.
// You may assume capacity >= 1 (a test covers capacity == 1).
func NewLRU(capacity int) *LRU {
	// TODO: implement
	return &LRU{}
}

// Get returns the value for key and true if present, marking it as the most
// recently used. If absent, it returns (0, false) and changes nothing.
func (l *LRU) Get(key int) (int, bool) {
	// TODO: implement
	return 0, false
}

// Put inserts or updates key->value, marking it most recently used.
// If inserting a new key exceeds capacity, evict the least recently used entry
// first. Updating an existing key must NOT grow the size.
func (l *LRU) Put(key int, value int) {
	// TODO: implement
}

// Len reports the current number of stored entries (for tests/invariants).
// It must never exceed capacity.
func (l *LRU) Len() int {
	// TODO: implement
	return 0
}

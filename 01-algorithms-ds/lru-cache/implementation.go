package lru

// LRU is a fixed-capacity cache that evicts the least-recently-used entry
// when a new key would exceed capacity.
//
// Target complexity: BOTH Get and Put must run in O(1).
// That constraint is the whole puzzle — pick internal structures that make
// "look up by key", "move an entry to most-recent", and "drop the
// least-recent entry" all O(1).
//
// Suggested internals (you decide the final shape — fill these in):
//   - something to map key -> the entry's location in O(1)
//   - something to keep entries ordered by recency, where you can detach an
//     entry from the MIDDLE and re-attach it at one end in O(1)
//
// Hint already discussed: a hash map + a doubly linked list. The list node
// below is a starting scaffold; change it however you like.

type node struct {
	key  int
	val  int
	prev *node
	next *node
	// TODO: add fields if your design needs them
}

type LRU struct {
	capacity int
	// TODO: add the fields you need (the map, the list head/tail, a size counter...)
}

// NewLRU returns an empty cache that holds at most `capacity` entries.
// You may assume capacity >= 1 (a test covers capacity == 1).
func NewLRU(capacity int) *LRU {
	// TODO: initialize your internal structures here.
	return &LRU{capacity: capacity}
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

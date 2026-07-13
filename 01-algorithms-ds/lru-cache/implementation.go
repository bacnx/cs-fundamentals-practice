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
	key        int
	val        int
	valPointer *node
	prev       *node
	next       *node
}

type linkedList struct {
	head *node // head allway not nil
	tail *node // tail allway not nil
}

func newLinkedList() *linkedList {
	head := &node{}
	tail := &node{}
	head.next = tail
	tail.prev = head

	return &linkedList{
		head: head,
		tail: tail,
	}
}

func (l *linkedList) pushToHead(key int, value int, valPointer *node) *node {
	newNode := &node{
		key:        key,
		val:        value,
		valPointer: valPointer,
	}
	head := l.head
	next := head.next

	head.next = newNode
	newNode.prev = head
	newNode.next = next
	next.prev = newNode

	return newNode
}

func (l *linkedList) pushToTail(key int, value int, valPointer *node) *node {
	newNode := &node{
		key:        key,
		val:        value,
		valPointer: valPointer,
	}
	tail := l.tail
	prev := tail.prev

	prev.next = newNode
	newNode.prev = prev
	newNode.next = tail
	tail.prev = newNode

	return newNode
}

// remove remove a node in a linkedList
func remove(n *node) *node {
	if n.next == nil || n.prev == nil {
		return nil // invalid node
	}

	next := n.next
	prev := n.prev

	prev.next = next
	next.prev = prev

	n.next = nil
	n.prev = nil
	return n
}

func (l *linkedList) getLastItem() *node {
	if l.tail.prev == l.head {
		return nil
	}
	return l.tail.prev
}

type LRU struct {
	capacity int
	len      int
	sliCap   int
	lruList  *linkedList // head is newest
	sli      []*linkedList
}

// NewLRU returns an empty cache that holds at most `capacity` entries.
// You may assume capacity >= 1 (a test covers capacity == 1).
func NewLRU(capacity int) *LRU {
	sliCap := capacity * 3 // sliCap large for reduce collision
	return &LRU{
		capacity: capacity,
		sliCap:   sliCap,
		lruList:  newLinkedList(),
		sli:      make([]*linkedList, sliCap),
	}
}

// Get returns the value for key and true if present, marking it as the most
// recently used. If absent, it returns (0, false) and changes nothing.
func (l *LRU) Get(key int) (int, bool) {
	sliPosition := l.hash(key)
	pos := l.sli[sliPosition].searchInList(key)
	if pos != nil {
		remove(pos.valPointer)
		newNode := l.lruList.pushToHead(pos.key, pos.val, pos)
		pos.valPointer = newNode
		return pos.val, true
	}
	return 0, false
}

// Put inserts or updates key->value, marking it most recently used.
// If inserting a new key exceeds capacity, evict the least recently used entry
// first. Updating an existing key must NOT grow the size.
func (l *LRU) Put(key int, value int) {
	sliPosition := l.hash(key)
	if l.sli[sliPosition] == nil {
		l.sli[sliPosition] = newLinkedList()
	}
	posList := l.sli[sliPosition]

	if pos := posList.searchInList(key); pos != nil {
		l.removeByLruSliPos(pos)
	} else if l.len == l.capacity {
		lruNode := l.lruList.getLastItem()
		lruSliNode := lruNode.valPointer
		l.removeByLruSliPos(lruSliNode)
	}
	newNode := l.lruList.pushToHead(key, value, nil)
	newPosNode := l.sli[sliPosition].pushToHead(key, value, newNode)
	newNode.valPointer = newPosNode
	l.len++
}

func (l *LRU) removeByLruSliPos(n *node) {
	remove(n)
	remove(n.valPointer)
	l.len--
}

// Len reports the current number of stored entries (for tests/invariants).
// It must never exceed capacity.
func (l *LRU) Len() int {
	return l.len
}

func (l *LRU) hash(key int) int {
	prime := 100007
	return key * prime % l.capacity
}

func (list *linkedList) searchInList(key int) *node {
	if list == nil {
		return nil
	}

	for n := list.head; n != list.tail; n = n.next {
		if n.key == key {
			return n
		}
	}
	return nil
}

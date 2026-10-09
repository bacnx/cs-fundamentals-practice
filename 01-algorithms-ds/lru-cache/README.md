# LRU Cache

Applied problem 1 for `hash-map`.

## Real-world context

Browser cache, Redis with `maxmemory-policy allkeys-lru`, CPU instruction cache, ORM
identity maps, CDN edge caches. Any cache with a bounded size needs an eviction policy, and
"throw out whatever has gone longest without being touched" is the default one because it
approximates future reuse from past reuse at almost no cost.

## Theory

A hash map alone cannot do this. It answers *"where is key K"* in O(1), but it has no notion
of order, so *"which entry is the least recently used"* would be an O(n) scan — and that scan
would sit inside `Put`, destroying the O(1) guarantee.

The fix is to keep two structures over the same entries:

- a **hash map** key → the entry's location, for O(1) lookup
- a **doubly linked list** ordered by recency, for O(1) "detach from the middle, re-attach at
  the front" and O(1) "read the back"

```
 most recent                                     least recent
   head ⇄ (3,30) ⇄ (1,10) ⇄ (4,40) ⇄ ... ⇄ (2,20) ⇄ tail
                     ↑                        ↑
   map[1] ───────────┘      map[2] ───────────┘            evict from here
```

Why *doubly* linked: detaching a node in O(1) requires knowing its predecessor. A singly
linked list would have to walk from the head to find it — O(n) again, in the hot path.

Both `Get` and `Put` are "touch" operations: they move the entry to the front. The only
difference is that `Put` may first have to evict from the back.

## Contract

```go
func NewLRU(capacity int) *LRU
func (l *LRU) Get(key int) (int, bool)   // also marks key as most recently used
func (l *LRU) Put(key int, value int)    // insert or update, marks most recently used
func (l *LRU) Len() int                  // never exceeds capacity
```

`capacity >= 1` may be assumed. Keys are `int` and are **not** restricted to non-negative
values.

## Invariants

1. `Len()` never exceeds `capacity`.
2. Updating an existing key does not change `Len()`.
3. Every key in the hash map has exactly one node in the recency list, and vice versa.
4. The back of the list is always the least recently touched live entry.
5. A successful `Get` changes recency order but nothing else — not `Len()`, not any value.

## Complexity targets

| Operation | Target |
|---|---|
| Get | O(1) worst case |
| Put | O(1) worst case |
| Len | O(1) |

O(1) *worst case*, not amortized: there is no resize step here, so any O(n) path is a defect,
not a rare cost.

## Levels

Single level. Expected behaviour:

```
cache := NewLRU(3)
cache.Put(1, 10)
cache.Put(2, 20)
cache.Put(3, 30)
cache.Get(1)          // 1 is now most recently used
cache.Put(4, 40)      // evicts 2 (LRU), not 1
cache.Get(2)          // → 0, false
```

## Verification checklist

- [ ] Get on an empty cache returns `false`
- [ ] Put then Get returns the correct value
- [ ] Len never exceeds capacity
- [ ] Eviction removes the least recently used entry, not the oldest inserted
- [ ] Get refreshes recency; Put on an existing key refreshes recency
- [ ] capacity == 1 works
- [ ] Re-inserting an evicted key works
- [ ] Negative keys work, including `math.MinInt`
- [ ] Key `0` works (put, update, miss after eviction)

## Hint

O(1) eviction requires knowing which entry is LRU instantly. Hash map alone cannot do this —
what additional structure gives O(1) access to both ends?

## Why (deeper)

Where LRU sits among eviction policies, and what it gets wrong (scan resistance, the reason
Redis ships LFU too): second-brain → `topics/cs-fundamentals.md`, Nhóm 1.

# Hash Map

## Real-world context

The default "lookup by key" structure everywhere: language-level maps and dicts, routing
tables, symbol tables in compilers, in-memory caches, deduplication sets, request-scoped
indexes. When a system answers *"does this key exist, and what is its value"* in constant
time, this is almost always what is underneath.

## Theory

### Core mechanism

A hash map stores key-value pairs in an underlying array. A hash function converts any key
into an integer array index, enabling O(1) average-case access.

```
key "user:42"
    ↓ hash("user:42") % array_size
    → index 7
    ↓
array[7] = {key: "user:42", value: {...}}
```

### Collision handling

Two different keys can hash to the same index (collision). Two strategies:

**Chaining** — each bucket holds a linked list of entries:
```
array[7] → ("user:42", data1) → ("order:7", data2) → nil
```
Lookup: hash → index → scan the list for matching key.

**Open addressing** — on collision, probe for the next available slot:
```
array[7] occupied → try array[8] → try array[9] → ...
```
More cache-friendly but trickier to implement correctly (especially deletion).

### Load factor and resize

Load factor = `n / capacity` (elements / array size).

When load factor exceeds ~0.75, performance degrades. The map resizes: allocate a new array
(typically 2×) and rehash every entry into it. This is O(n) but happens rarely enough that
amortized cost per operation is still O(1).

### Why hash maps cannot be sorted

The array index is computed, not ordered — there is no inherent relationship between adjacent
slots. Use a balanced BST (TreeMap in Java, `btree` packages in Go) when you need ordered
iteration or range queries.

## Contract

Level 1–2, chaining (`Hash`):

```go
func NewHash(cap int) Hash
func (h *Hash) Put(key string, value any)
func (h *Hash) Get(key string) (any, bool)
func (h *Hash) Delete(key string)
func (h *Hash) Len() int
func (h *Hash) LoadFactor() float64   // required by Level 2, NOT implemented — see REVIEW-LOG.md
```

Level 3, open addressing (`HashOpenAddressing`): the same operations, built on a flat slot
array instead of per-bucket lists.

## Invariants

1. Every inserted key is retrievable until it is deleted.
2. All keys remain accessible after a resize.
3. Load factor stays below the threshold (0.75 chaining, 0.5 open addressing).
4. `Len()` equals the number of distinct live keys — updating an existing key does not change it.
5. Open addressing only: deleting a key never makes a *different* key unreachable.

## Complexity targets

| Operation | Average | Worst case |
|---|---|---|
| Get | O(1) | O(n) |
| Put | O(1) amortized | O(n) on resize |
| Delete | O(1) | O(n) |

Worst case occurs when many keys collide — a poor hash function, or adversarial input chosen
to collide on purpose (hash DoS).

## Levels

### Level 1 — Chaining hash map

Implement a hash map using separate chaining.

Required operations: `Put`, `Get` (second return indicates existence), `Delete`, `Len`.

Constraints:

- Fixed initial capacity (e.g. 16 buckets)
- Use a linked list (or slice) per bucket — no `map` built-in

### Level 2 — Load factor and auto-resize

Extend Level 1:

- Track load factor after every `Put`
- When load factor exceeds 0.75, resize to double capacity and rehash all entries
- `LoadFactor() float64` exposes the current load factor

### Level 3 — Open addressing (linear probing)

Implement a **separate** hash map using open addressing instead of chaining.

Same operations as Level 1, plus:

- `Put` finds the next available slot on collision
- `Delete` uses a **tombstone** marker (why? — figure it out during implementation)
- Resize when load factor exceeds 0.5 (a lower threshold than chaining — why?)

## Verification checklist

Level 1:

- [x] Get on a missing key returns `false`
- [x] Put then Get returns the correct value
- [x] Delete removes the key; a subsequent Get returns `false`
- [x] Two keys with the same hash are both accessible (collision works)

Level 2:

- [ ] Load factor is 0 on an empty map — **not verifiable, `LoadFactor()` is missing**
- [ ] Load factor increases correctly on each insert — **not verifiable, same reason**
- [x] After a resize, all previously inserted keys are still accessible
- [ ] Load factor drops below 0.75 after a resize — **not verifiable, same reason**

Level 3:

- [x] Collision resolved correctly (keys with the same hash both stored and retrievable)
- [x] Delete does not break a subsequent Get for other keys
- [ ] **Explain in a comment why tombstones are necessary** — still open, see `REVIEW-LOG.md`

## Hint

If a deleted slot looked exactly like a never-used slot, what would the probe loop do when it
reached it while searching for a key stored *after* that slot?

## Why (deeper)

Mechanism, the "access pattern" framing, hash DoS, and why databases index with B-trees rather
than hash tables: second-brain → `topics/cs-fundamentals.md`, Nhóm 1.

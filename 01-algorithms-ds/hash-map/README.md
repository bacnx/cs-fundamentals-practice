# Hash Map

## Theory

### Core mechanism

A hash map stores key-value pairs in an underlying array. A hash function converts any key into an integer array index, enabling O(1) average-case access.

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

When load factor exceeds ~0.75, performance degrades. The map resizes: allocate a new array (typically 2×) and rehash every entry into it. This is O(n) but happens rarely enough that amortized cost per operation is still O(1).

### Complexity

| Operation | Average | Worst case |
|---|---|---|
| Get | O(1) | O(n) |
| Put | O(1) amortized | O(n) on resize |
| Delete | O(1) | O(n) |

Worst case occurs when many keys collide (poor hash function or adversarial input).

### Why hash maps cannot be sorted

The array index is computed, not ordered — there is no inherent relationship between adjacent slots. Use a balanced BST (TreeMap in Java, `btree` packages in Go) when you need ordered iteration or range queries.

---

## Implementation tasks

Progress from basic to advanced. Each level must pass before moving on.

### Level 1 — Chaining hash map

Implement a hash map using separate chaining.

Required operations:
- `Put(key string, value any)` — insert or update
- `Get(key string) (any, bool)` — retrieve; second return indicates existence
- `Delete(key string)` — remove key if present
- `Len() int` — number of stored keys

Constraints:
- Fixed initial capacity (e.g., 16 buckets)
- Use a linked list (or slice) per bucket — no `map` built-in

**Verification checklist:**
- [ ] Get on missing key returns `false`
- [ ] Put then Get returns correct value
- [ ] Delete removes the key; subsequent Get returns `false`
- [ ] Two keys with same hash both accessible (collision works)

---

### Level 2 — Load factor and auto-resize

Extend Level 1:

- Track load factor after every `Put`
- When load factor exceeds 0.75, resize to double capacity and rehash all entries
- `LoadFactor() float64` — expose current load factor

**Verification checklist:**
- [ ] Load factor is 0 on empty map
- [ ] Load factor increases correctly on each insert
- [ ] After resize, all previously inserted keys still accessible
- [ ] Load factor drops below 0.75 after resize

---

### Level 3 — Open addressing (linear probing)

Implement a separate hash map using open addressing instead of chaining.

Required operations: same as Level 1, plus:
- `Put` must find next available slot on collision
- `Delete` must use a **tombstone** marker (why? — figure it out during implementation)
- Resize when load factor exceeds 0.5 (lower threshold needed — why?)

**Verification checklist:**
- [ ] Collision resolved correctly (keys with same hash both stored and retrievable)
- [ ] Delete does not break subsequent Get for other keys
- [ ] Explain in a comment why tombstones are necessary

---

## Applied problems

These problems use a hash map as the core data structure but require combining it with additional logic. Each reflects a real system design challenge.

---

### Problem 1 — LRU Cache

**Real-world context:** Browser cache, Redis `maxmemory-policy allkeys-lru`, CPU instruction cache.

Design a cache that:
- Stores up to `capacity` key-value pairs
- `Get(key)` returns the value and marks the entry as recently used
- `Put(key, value)` inserts the entry; if at capacity, evicts the **least recently used** entry first
- Both `Get` and `Put` must be O(1)

```
cache := NewLRUCache(3)
cache.Put("a", 1)
cache.Put("b", 2)
cache.Put("c", 3)
cache.Get("a")        // a is now most recently used
cache.Put("d", 4)     // evicts "b" (LRU), not "a"
cache.Get("b")        // → not found
```

**Hint:** O(1) eviction requires knowing which entry is LRU instantly. Hash map alone cannot do this — what additional structure gives O(1) access to both ends?

---

### Problem 2 — Rate limiter (sliding window)

**Real-world context:** API gateways (AWS API Gateway, Nginx), preventing abuse.

Implement a rate limiter:
- `Allow(userID string, now time.Time) bool`
- Returns `true` if the user has made fewer than `limit` requests in the last `windowSeconds`
- Returns `false` (and does not count the request) if the limit is exceeded

```
limiter := NewRateLimiter(limit=5, windowSeconds=60)
// user "u1" makes 5 requests at t=0..4s → all allowed
// user "u1" makes request at t=5s → denied
// user "u1" makes request at t=61s → allowed (window slid past early requests)
```

**Constraint:** Must work correctly when requests come in at arbitrary times — not just uniform intervals.

---

### Problem 3 — In-memory key-value store with TTL

**Real-world context:** Session stores, short-lived tokens, feature flags with expiry.

Implement a store where:
- `Set(key string, value any, ttl time.Duration)` — store with expiry
- `Get(key string) (any, bool)` — return value if not expired; `false` if missing or expired
- `DeleteExpired()` — remove all expired keys (called periodically by a background process)
- Expired keys must not be returned by `Get`, even if `DeleteExpired` hasn't run yet

```
store.Set("session:abc", userData, 30*time.Minute)
store.Get("session:abc")   // → userData, true (within TTL)
// 31 minutes later...
store.Get("session:abc")   // → nil, false (expired)
```

**Follow-up:** If `DeleteExpired` needs to efficiently find only the expired keys (not scan all keys), what secondary structure would you add?

---

### Problem 4 — Word frequency index

**Real-world context:** Search engine indexing, log analysis, autocomplete ranking.

Given a stream of words (one at a time, potentially millions):
- `Add(word string)` — record the word
- `TopK(k int) []string` — return the k most frequent words, in descending order of frequency
- Both operations should be efficient (not O(n) scan on every `TopK` call)

```
index.Add("go")
index.Add("rust")
index.Add("go")
index.Add("go")
index.Add("rust")
index.TopK(2)   // → ["go", "rust"]
```

**Hint:** `TopK` with efficient updates is a known combination of two data structures from this module.

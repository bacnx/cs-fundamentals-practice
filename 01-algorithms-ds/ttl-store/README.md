# In-memory KV store with TTL

Applied problem 3 for `hash-map`. **Spec staged — not started.**

## Real-world context

Session stores, short-lived auth tokens, feature flags with an expiry, idempotency keys,
cached API responses. Redis is this structure with a network protocol bolted on.

## Theory

Expiry looks like it needs a background job, but there are two independent mechanisms and a
correct store needs both:

- **Lazy expiry** — `Get` checks the entry's deadline and reports a miss if it has passed.
  Cheap, and guarantees correctness of reads. But an expired entry that is never read again
  occupies memory forever.
- **Active expiry** — a periodic sweep actually deletes expired entries, reclaiming memory.
  Needed for space, not for correctness.

A naive sweep scans every key, which is O(n) on a mostly-live store — wasted work proportional
to the data you are *keeping*. To sweep in time proportional to what actually expired, the
store needs a second structure that answers *"what is the next thing to expire"* — a different
access pattern from "look up by key", and therefore a different structure.

## Contract

```go
func NewStore() *Store
func (s *Store) Set(key string, value any, ttl time.Duration)
func (s *Store) Get(key string) (any, bool)   // false if missing OR expired
func (s *Store) DeleteExpired()               // called periodically by a background process
```

## Invariants

1. `Get` never returns an expired value, even if `DeleteExpired` has not run yet.
2. `DeleteExpired` removes every expired entry and no live one.
3. `Set` on an existing key replaces both value and deadline.
4. After `DeleteExpired`, no expired entry holds memory.

## Complexity targets

| Operation | Target |
|---|---|
| Set | O(log n) or better |
| Get | O(1) |
| DeleteExpired | O(k log n), where k = number of entries that actually expired — **not** O(n) |

## Levels

**Level 1 — lazy expiry only.** `Get` honours the deadline; `DeleteExpired` scans all keys.
Correct, and the O(n) sweep is accepted at this level.

**Level 2 — sweep proportional to expiries.** Add the secondary structure so `DeleteExpired`
touches only entries that have actually expired.

```
store.Set("session:abc", userData, 30*time.Minute)
store.Get("session:abc")   // → userData, true (within TTL)
// 31 minutes later...
store.Get("session:abc")   // → nil, false (expired)
```

## Verification checklist

Level 1:

- [ ] Get before the TTL returns the value
- [ ] Get after the TTL returns `false`, without `DeleteExpired` having run
- [ ] Set on an existing key resets its deadline
- [ ] A zero or negative TTL expires immediately

Level 2:

- [ ] `DeleteExpired` on a store with 1 expired and many live entries touches ~1 entry
- [ ] The secondary structure stays consistent when a key is overwritten with a new TTL
- [ ] The secondary structure stays consistent when a key is expired lazily by `Get` first

## Hint

If `DeleteExpired` needs to find only the expired keys without scanning all keys, what question
does it have to be able to ask in O(1) — and which structure answers exactly that question?

## Why (deeper)

Lazy vs active expiry, why Redis samples instead of sweeping exhaustively, and the memory/CPU
tradeoff that choice encodes: second-brain → `topics/cs-fundamentals.md`, Nhóm 1.

# Review log

Append-only. Every review adds an entry; an entry stays `OPEN` until the user fixes it
themselves. Agents record findings here and do not fix them.

A finding that generalises beyond this repo should also be promoted into second-brain
(`topics/on-tap/cs-fundamentals.md`) so it enters the spaced-repetition schedule. The column
below says whether that has happened.

---

## 2026-09-11 — hash-map (open addressing), lru-cache

Reviewed at the restructure of the repo. `go test ./...` passes for both packages, so every
finding below is something the existing test suite does not cover.

### OPEN — `hash-map`: Level 3 checklist item not done

`implementation_oa.go` has the `Deleted` flag and uses it on delete, but the Level 3 checklist
requires *"explain in a comment why tombstones are necessary"* and there is no such comment.
This is the understanding half of the level, not a formality: the answer is the reason the
probe loop cannot treat a deleted slot the same way it treats a never-used slot.

Promoted to second-brain: yes.

### OPEN — `hash-map`: Level 2 is missing the `LoadFactor()` accessor

The README's Level 2 requires `LoadFactor() float64`. Only the private
`loadFactorAndResize()` exists — the public accessor was dropped by commit `cc98535`
("remove unused methods"). Consequence: three of the four Level 2 checklist items cannot be
verified at all, because the quantity they assert about is not observable from outside the
package. The resize behaviour itself is tested and correct; the observability is not.

Lesson: "unused" measured against the current test file is not the same as "not required".
A spec item with no test looks like dead code.

Promoted to second-brain: yes.

### OPEN — `lru-cache`: panic on a negative key

Reproduced: `c := NewLRU(3); c.Put(-1, 100)` panics with `index out of range [-2]` at
`implementation.go:141`.

`hash()` computes `key * prime % l.capacity`. Go's `%` keeps the sign of the dividend, so a
negative key yields a negative index. The contract in the README does not restrict keys to
non-negative values, and no test passes a negative one.

Promoted to second-brain: yes — Go's `%` is remainder, not Euclidean modulo. Any
`hash(x) % n` over signed input has this hole.

### OPEN — `lru-cache`: two thirds of the bucket array is unreachable

`NewLRU` allocates `sli` with length `sliCap = capacity * 3`, commented
*"sliCap large for reduce collision"*. But `hash()` reduces modulo `l.capacity`, not
`l.sliCap`, so only indices `[0, capacity)` are ever produced. The extra memory is allocated
and never used, and the collision rate is exactly what it would be with `capacity` buckets.

Declared intent and actual behaviour disagree — review criterion 4. The test suite cannot
catch it because the map still behaves correctly, just with more collisions than designed.

Promoted to second-brain: yes — a comment stating an intent is not evidence the code
implements it.

---

## Earlier findings (reconstructed from git history)

These were found and fixed before this log existed. Kept because the lesson outlived the bug.

- **`hashFunc` hashed byte indices instead of rune values** (commits `ba5a56b`, `cc98535`).
  Ranging over a Go `string` yields `(byteIndex, rune)`; using the first element where the
  second was meant gives a hash function over positions, not over content — so it collides
  heavily and is unstable under multi-byte input. FIXED.
- **Both constructors accepted capacity 0** (commit `5d7ccca`), which makes every `% cap` a
  division by zero. FIXED.

## 2026-09-16 — hash-map (open addressing)

### CLOSED — `hash-map`: Level 3 checklist item not done

The missing explanation for tombstones was added. The comment correctly identifies that without it, the probe loop would terminate early and "ignore the items at the back".

### CLOSED — `hash-map`: Level 2 is missing the `LoadFactor()` accessor

The `LoadFactor() float64` method has been added and correctly implements `float64(len) / float64(cap)`. The missing tests have been added to verify that the load factor calculates correctly and properly triggers the resize logic.

### CLOSED — `lru-cache`: panic on a negative key

Fixed. Go's `%` preserves the sign, and `math.MinInt * -1` overflows and stays negative. The hash function now correctly applies the modulo *before* flipping the sign (and re-hashing), guaranteeing a positive index regardless of two's complement integer overflow.

## 2026-10-09 — lru-cache reset

The user cleared the implementation to redo it from scratch. `implementation.go` is back to a
signatures-only template; the test suite is kept.

### WITHDRAWN — `lru-cache`: two thirds of the bucket array is unreachable

The code it describes no longer exists. `TestBucketArrayUtilization` was removed with it: it
read the old design's private `sli` field, so it cannot compile against a fresh design. The
lesson stands — a comment stating an intent is not evidence the code implements it.

### CLOSED BY RESET — `lru-cache`: key 0 panics

Found before the reset: `NewLRU(1).Put(0, 7)` panicked with a nil pointer dereference. The
bucket search started at the sentinel head node, whose `key` field is the zero value `0`, so
key 0 matched the sentinel instead of a real entry. The three `TestZeroKey*` tests stay in the
suite so the rewrite is checked against it.

Promoted to second-brain: no — candidate: a sentinel node's zero-valued fields are real
values to any code that reads them; a search loop must skip sentinels, not trust their data.

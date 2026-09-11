# Word frequency index

Applied problem 4 for `hash-map`. **Spec staged — blocked on `heap`.**

## Real-world context

Search-engine indexing, log analysis ("top 10 error messages this hour"), autocomplete
ranking, trending topics, anomaly detection on request paths.

## Theory

Two operations pull in opposite directions. `Add` must be O(1) and is called millions of times.
`TopK` asks an *ordering* question, and a hash map has no order — so the naive answer is "scan
all distinct words and sort", which is O(n log n) per call.

The combination that fixes it: the hash map keeps the counts, and a **heap** keeps the current
best k. A heap answers "what is the smallest element of the k I am holding" in O(1), which is
exactly the question needed to decide whether an incoming word belongs in the top k at all.

Note that this is a *bounded* heap of size k, not a heap of every word. That is what makes the
memory cost independent of vocabulary size, and it is why the heap is a min-heap even though
the query asks for the largest: the thing you need instant access to is the weakest member of
the set you are defending.

## Prerequisite

The `heap` module must be built first. Do not start this with a `sort.Slice` stand-in — the
whole lesson is the structure.

## Contract

```go
func NewIndex() *Index
func (i *Index) Add(word string)
func (i *Index) TopK(k int) []string   // descending by frequency
```

## Invariants

1. `Add` never walks the vocabulary.
2. `TopK(k)` returns exactly `min(k, distinct words)` entries.
3. The result is ordered by descending frequency.
4. Ties are broken deterministically (pick a rule and state it — alphabetical is fine).

## Complexity targets

| Operation | Target |
|---|---|
| Add | O(1) amortized, or O(log k) |
| TopK | O(k log k) — **not** O(n log n) |

## Levels

Single level, but only after `heap` exists.

```
index.Add("go")
index.Add("rust")
index.Add("go")
index.Add("go")
index.Add("rust")
index.TopK(2)   // → ["go", "rust"]
```

## Verification checklist

- [ ] `TopK` on an empty index returns an empty slice, not nil-panic
- [ ] `k` larger than the vocabulary returns everything, in order
- [ ] `k == 0` returns empty
- [ ] A word that overtakes another is reflected in the next `TopK`
- [ ] Ties are broken by the stated rule, consistently across calls
- [ ] `TopK` does not scan the whole vocabulary (reason about it, or instrument it)

## Hint

`TopK` with efficient updates is a known combination of two structures from this module. If you
are holding only the best k so far, which single element do you need to compare an incoming
word against?

## Why (deeper)

Exact top-k vs streaming approximations (Count–Min Sketch, Space-Saving), and why real systems
usually choose the approximate one: second-brain → `topics/cs-fundamentals.md`, Nhóm 1.

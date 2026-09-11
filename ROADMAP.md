# Roadmap

Single source of truth for **practice status**. No theory here — theory lives in each module's
`README.md`. Learning status (what has been *taught*) lives in second-brain's
`indexes/topic-map.md`; this file tracks only what has been *built*.

## Next

**Close the open items on what is already built, then start `rate-limiter`.**
Four things are open — see `REVIEW-LOG.md`: the missing tombstone explanation in the
open-addressing map, the `LoadFactor()` accessor Level 2 requires but does not have, and two
defects in the LRU cache (panic on negative keys, two thirds of the bucket array unreachable).

## Legend

| | Meaning |
|---|---|
| ✅ | built, tests pass, checklist ticked |
| ⚠️ | built but has open findings in `REVIEW-LOG.md` |
| 🔜 | spec staged (README exists), implementation not started |
| ⏳ | queued, nothing scaffolded yet |

## Group order

Ordered by leverage for a fullstack developer, not by textbook order. Rationale is in
second-brain's `topics/cs-fundamentals.md`.

| # | Group | Order | Status |
|---|---|---|---|
| 01 | Algorithms & Data Structures | 1st | in progress |
| 03 | Database Internals | 2nd | ⏳ |
| 02 | Networking | 3rd | ⏳ |
| 04 | OS & Architecture | 4th | ⏳ |
| 05 | Distributed Systems | 5th | ⏳ |
| 06 | Type Systems & PLT | 6th | ⏳ |

## 01 — Algorithms & Data Structures

The group is organised around **access patterns**, not around algorithm names: pick the
structure the question demands. Lookup by key → hash map. Highest-priority element → heap.
Keys sharing a prefix → trie.

### Core structures

| Module | Level | What | Status |
|---|---|---|---|
| `hash-map` | 1 | Chaining: `Put` / `Get` / `Delete` / `Len`, fixed capacity | ✅ |
| `hash-map` | 2 | Load factor tracking + auto-resize at 0.75 | ⚠️ |
| `hash-map` | 3 | Open addressing with linear probing + tombstones | ⚠️ |
| `heap` | — | Binary heap: `Push` / `Pop` / `Peek`, sift up/down | ⏳ |
| `trie` | — | Prefix tree: `Insert` / `Search` / `StartsWith` | ⏳ |

### Applied problems

Each one combines a core structure with extra logic, and each mirrors a real system.

| Module | Problem | Needs | Status |
|---|---|---|---|
| `lru-cache` | LRU cache, O(1) `Get` and `Put` | hash map + doubly linked list | ⚠️ |
| `rate-limiter` | Sliding-window rate limiter | hash map | 🔜 |
| `ttl-store` | In-memory KV store with TTL | hash map (+ a second structure for efficient expiry) | 🔜 |
| `word-frequency` | Streaming word counts + `TopK` | hash map + **heap** | 🔜 — blocked on `heap` |

## 02–06

Nothing scaffolded. Each group gets its first module only when the corresponding theory has
been covered in second-brain — implementing against a spec you do not yet understand produces
code that passes tests and teaches nothing.

## What this file does not track

- **Theory** → the module's own `README.md`
- **Why it works / cross-topic links** → second-brain `topics/cs-fundamentals.md`
- **Revision schedule** → second-brain `indexes/review-schedule.md`
- **Open review findings** → `REVIEW-LOG.md`

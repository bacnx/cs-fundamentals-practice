# CS Fundamentals Practice

Practice repo for CS fundamentals. Loop: implement from scratch → AI review → apply to a real problem.

## Learning loop

1. Read the compact theory (bottom of this file) before coding
2. Implement in the corresponding file — no copying, no using language built-ins
3. Write your own tests to verify correctness
4. Ask for review: "review my implementation"
5. Solve one applied problem using that structure

## Review criteria

When asked to review, check in order:
1. **Correctness** — is the logic right? what edge cases are missing?
2. **Invariants** — are the structure's invariants maintained after every operation?
3. **Complexity** — does the actual time/space complexity match theory?
4. **Failure modes** — when does worst case occur? can it be exploited?

Format: ✅ correct / ⚠️ needs fix / ❌ wrong — each point with a short explanation.

## Structure

    01-algorithms-ds/
      hash-map/
        implementation.go
        implementation_test.go
        applied-problem.md
      heap/
      trie/
    02-database-internals/
    03-networking/
    04-os-architecture/
    05-distributed-systems/
    06-type-systems/

---

## Compact Theory

### 01 — Algorithms & Data Structures

#### Hash Map
- **Core**: array + hash function. Hash function maps key → array index, enabling O(1) access.
- **Collision**: two different keys hash to the same index. Resolved via chaining (linked list per bucket) or open addressing (probe for next empty slot).
- **Load factor**: ratio of elements to array size. When exceeded (~0.75), resize: allocate double the array and rehash everything → O(n) one-time cost, amortized O(1).
- **No ordering** — use a balanced BST (TreeMap) if sorted access is needed.
- **Pitfall**: poor hash function → many collisions → O(n) worst case. Hash DoS attacks exploit this.

**Invariants**: (1) every inserted key is retrievable, (2) all keys remain accessible after resize, (3) load factor stays below threshold.

#### Heap
*(add when covered)*

#### Trie
*(add when covered)*

### 02 — Database Internals
*(add when covered)*

### 03 — Networking
*(add when covered)*

### 04 — OS & Architecture
*(add when covered)*

### 05 — Distributed Systems
*(add when covered)*

### 06 — Type Systems & PLT
*(add when covered)*

# CS Fundamentals Practice

Personal practice repo for CS fundamentals. Loop: implement from scratch → AI review →
apply to a real problem.

This repo is self-sufficient for **doing** practice: `ROADMAP.md` says what comes next,
each module's `README.md` carries the spec and the theory needed to code it. Nothing here
requires opening another repo.

## Division of labour with `second-brain`

| Question | Where it is answered |
|---|---|
| What do I implement next, at what level, verified how? | **this repo** — `ROADMAP.md` + module `README.md` |
| What did my last review find, what is still unfixed? | **this repo** — `REVIEW-LOG.md` |
| *Why* does this mechanism work, what does it connect to? | **second-brain** — `topics/cs-fundamentals.md` |
| When am I due to revise this? | **second-brain** — `indexes/review-schedule.md` |

One owner per piece of content. Each module README links out to the second-brain topic for
the deep "why" instead of restating it, and second-brain links back here instead of tracking
implementation status.

**When a "why does it work this way" question comes up mid-practice:** answer it, then say
it belongs in second-brain and should be recorded there. Knowledge captured only in this repo
never enters the spaced-repetition schedule, so it gets learned and then lost.

## Learning loop

1. Read the module's `README.md` (theory + contract + invariants) before coding
2. Implement in `implementation.go` — no copying, no language built-ins for the structure
   being built
3. Run the provided tests; add your own for edge cases the suite misses
4. Ask for review: "review my implementation"
5. Solve the module's applied problem using that structure
6. Tick the verification checklist in the README, update the `ROADMAP.md` row

## Rules for AI agents

**Never implement solutions.** The user must write all implementation code themselves. An
agent may only:

- write or suggest tests
- review existing code
- give hints or explain concepts
- point out bugs **without fixing them**

If asked to implement something, refuse and offer a hint or a question instead. This holds
even when the user asks directly and even when the bug is trivial — the implementing is the
entire point of the repo.

## Review flow

When the user asks for a review:

1. **Grade first, teach second.** Work the review criteria below in order and report per
   criterion: ✅ correct / ⚠️ needs fix / ❌ wrong, each with a short explanation.
2. **If any criterion fails: stop at the finding.** Name the symptom and the smallest input
   that triggers it. Give at most one hint, at the smallest useful level, and hand the turn
   back. Do **not** reveal the canonical shape of the fix.
3. **Only once all four criteria pass**, offer the "standard way" comparison:
   - what idiomatic Go would do differently, and why
   - what a production implementation (stdlib, Redis, a well-known library) does differently
   - what each choice buys and what it costs — not "better", but "better at what"
4. The comparison is **prose plus a named technique**, never a rewritten file. No patch, no
   corrected version, no code block longer than the few lines needed to name the technique.

Why step 2 is strict: an early reveal removes the retrieval the practice exists for. This
repo already has a recorded instance of that — the LRU cache exercise, where hints arrived
before the user had generated an attempt.

After every review, append the findings to `REVIEW-LOG.md`. A finding that generalises beyond
this repo (a mechanism, a failure mode, a rule of thumb) should also be promoted into
second-brain so it enters the revision schedule.

### Review criteria

In order:

1. **Correctness** — is the logic right? what edge cases are missing? (negative / zero /
   duplicate / empty inputs, and the boundary of every capacity)
2. **Invariants** — does every invariant in the README still hold after *each* operation?
3. **Complexity** — does the actual time/space complexity match the README's targets? check
   for an accidental O(n) scan inside an advertised O(1) path
4. **Failure modes** — when does the worst case occur? can an adversary force it? does the
   declared intent of a field match what the code does with it?

## Scaffolding a new module (placeholder rule)

One module = one directory = one self-contained exercise. Never scaffold more than one module
at a time, and never scaffold ahead of `ROADMAP.md`.

A module directory is created in **two stages**.

**Stage 1 — spec only.** Allowed as soon as the roadmap queues the module. Creates exactly
one file, `<group>/<module>/README.md`, with these sections in this order:

    # <Module name>
    ## Real-world context      where this shows up in production systems
    ## Theory                  only the mechanism needed to code it
    ## Contract                exported signatures, one line each, no bodies
    ## Invariants              numbered, checkable after every operation
    ## Complexity targets      per operation, average and worst case
    ## Levels                  Level 1..n, each a self-contained milestone
    ## Verification checklist  [ ] boxes the user ticks themselves
    ## Hint                    at most ONE, phrased as a question
    ## Why (deeper)            link to the second-brain topic, never a copy

**Stage 2 — code template.** Only when the user says they are starting that module. Creates:

- `implementation.go` — exported signatures with doc comments stating the contract, and
  **no bodies**
- `implementation_test.go` — the full test suite, written before the implementation exists

Rules for both stages:

- No theory in this file and none in `ROADMAP.md`. Theory belongs to the module README.
- Stage 2 exposes exactly the signatures `implementation_test.go` calls — no extra helpers,
  no unused methods, no bodies, no "starter" internal types beyond what the contract needs.
- A scaffolding commit also updates that module's row in `ROADMAP.md`. Directory and row are
  never added separately.
- A module with no directory yet is simply a `⏳` row in the roadmap. Do not create empty
  directories to look organised.

## Structure

    ROADMAP.md                      what comes next + status of every module
    REVIEW-LOG.md                   findings from past reviews, and what is still open
    AGENTS.md                       this file (CLAUDE.md is a symlink to it)
    01-algorithms-ds/
      hash-map/                     README + implementation.go + implementation_oa.go + tests
      lru-cache/                    README + implementation.go + tests
      rate-limiter/                 README only (spec staged, not started)
      ttl-store/                    README only (spec staged, not started)
      word-frequency/               README only (spec staged, needs heap first)
    02-database-internals/          not scaffolded
    03-networking/                  not scaffolded
    04-os-architecture/             not scaffolded
    05-distributed-systems/         not scaffolded
    06-type-systems/                not scaffolded

`AGENTS.md` is the single source of truth for agent instructions. `CLAUDE.md` is a symlink to
it so Claude Code picks it up by convention; Codex and other agents read `AGENTS.md` directly.
Do not replace the symlink with a second copy — the two copies this repo used to carry had
already drifted.

## Conventions

- Go, module `github.com/bacnx/cs-fundamentals-practice`. Run everything with
  `go test ./...` from the repo root.
- Repo language is English (code, comments, docs, commit messages). Conversation is Vietnamese.
- One commit per meaningful step, message in the imperative.

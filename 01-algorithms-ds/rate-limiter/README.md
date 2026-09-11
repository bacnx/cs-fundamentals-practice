# Rate Limiter (sliding window)

Applied problem 2 for `hash-map`. **Spec staged — not started.**

## Real-world context

API gateways (AWS API Gateway, Nginx `limit_req`, Cloudflare), login-attempt throttling,
protecting an expensive downstream from a retry storm. The limiter is what stands between one
misbehaving client and everyone else's latency.

## Theory

A **fixed window** counter ("at most N per calendar minute") is easy and wrong at the seam: a
client can send N at 11:59:59 and N more at 12:00:00 — 2N inside one second, while every
window technically held. A **sliding window** instead asks "how many requests in the last
`windowSeconds`, measured from *now*", which has no seam to exploit.

That requires remembering *when* each recent request happened, per user — so the shape is a
hash map from user ID to a collection of timestamps, with old timestamps expiring off the
back as `now` advances.

Note which access pattern that second collection has: you always drop from the oldest end and
add at the newest. That is a different structure from the one that answers "which is largest".

## Contract

```go
func NewRateLimiter(limit int, windowSeconds int) *RateLimiter
func (r *RateLimiter) Allow(userID string, now time.Time) bool
```

`Allow` returns `true` if the user has made fewer than `limit` requests within the last
`windowSeconds`, and counts this one. It returns `false` if the limit is exceeded, and a
denied request is **not** counted.

`now` is a parameter rather than read from the clock so that tests can control time.

## Invariants

1. A denied request never affects the decision on any later request.
2. Memory per user is bounded by `limit`, not by the number of requests ever made.
3. The decision depends only on `now` and the retained timestamps — calling `Allow` twice with
   the same `now` and a limit that is not yet reached allows both.

## Complexity targets

| Operation | Target |
|---|---|
| Allow | O(1) amortized per call |

Amortized, because one call may expire several stale timestamps at once — but each timestamp
is expired exactly once over its lifetime.

## Levels

Single level.

```
limiter := NewRateLimiter(limit=5, windowSeconds=60)
// user "u1" makes 5 requests at t=0..4s → all allowed
// user "u1" makes a request at t=5s     → denied
// user "u1" makes a request at t=61s    → allowed (window slid past the early requests)
```

**Constraint:** must be correct when requests arrive at arbitrary times, not just uniform
intervals.

## Verification checklist

- [ ] The `limit`-th request in a window is allowed; the `limit+1`-th is denied
- [ ] A denied request does not consume quota (still denied, but does not push the window)
- [ ] Quota recovers gradually as individual timestamps age out, not all at once
- [ ] Two different users do not interfere
- [ ] Requests at the exact window boundary are handled consistently
- [ ] Retained timestamps per user never exceed `limit`

## Hint

What has to be true about the *oldest* retained timestamp before a new request can be allowed?

## Why (deeper)

Fixed vs sliding vs token bucket vs leaky bucket, and why distributed rate limiting is a
different problem: second-brain → `topics/cs-fundamentals.md`, Nhóm 1 (and later Nhóm 5).

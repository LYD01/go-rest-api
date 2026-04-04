# Load test goals, scenarios, and success criteria

Use this document before running [scripts/load](../scripts/load/README.md). Update thresholds to match your environment.

## Goal

Establish how the go-api service behaves under increasing concurrent HTTP traffic: latency percentiles, throughput (RPS), error rate, and stability over time. Identify the first failure modes on **your** hardware (e.g. CPU saturation, timeouts, connection limits) so tuning or scaling decisions are data-driven.

## Success criteria (initial defaults)

Adjust these after a baseline run if your machine is much faster or slower.

| Criterion | Target | Notes |
|-----------|--------|--------|
| Availability | 0% HTTP 5xx under **baseline** load | 401/429 are policy, not server bugs |
| Latency (albums GET) | p95 &lt; 100 ms at **sustained** tier | Localhost; remote adds RTT |
| Latency (health) | p95 &lt; 50 ms at **sustained** tier | No API key path |
| Stability | RSS memory drift &lt; ~10% over **sustained** window | No obvious leak on this codebase |
| Stress | Document RPS/latency at **first** sustained error or timeout | Defines practical ceiling on test host |

## Scenarios

### 1. Baseline

- **Purpose:** Happy-path latency with minimal concurrency.
- **Load:** 1 virtual user (VU), ~10 iterations, or 10 RPS for 30s.
- **Endpoints:** `GET /healthz`, `GET /albums` (with `X-API-Key`).

### 2. Ramp

- **Purpose:** Find where latency or errors begin to climb.
- **Load:** Increase VUs or RPS in stages (e.g. 0 → 50 VUs over 2 minutes), hold briefly at each step, or use k6 stages.
- **Endpoints:** Primarily `GET /albums` (authenticated).

### 3. Sustained

- **Purpose:** Stability: memory, goroutines, steady error rate at ~70–80% of ramp “knee.”
- **Load:** Fixed concurrency (e.g. 30–50 VUs) for 5–15 minutes.
- **Endpoints:** Mix `GET /albums` and optional `GET /albums/{id}`.

### 4. Spike

- **Purpose:** Burst behavior and recovery (e.g. traffic 10× normal for 30s).
- **Load:** Short high-VU stage, then return to low VUs.
- **Endpoints:** `GET /albums`.

### 5. Stress (optional)

- **Purpose:** Deliberately exceed capacity; record breaking point **on your test host only**.
- **Load:** Push VUs or RPS until errors/timeouts &gt; 1% or tool limits hit.
- **Endpoints:** Same as sustained; optionally add `POST /albums` to stress write locks.

### 6. Write mix (optional)

- **Purpose:** Compare read-heavy vs write-heavy contention on the in-memory `RWMutex` store.
- **Load:** High concurrency `POST /albums` with small JSON bodies vs `GET /albums` only.

## Non-goals (for this repo’s first report)

- TLS / HTTP/2 client behavior (plain HTTP on localhost is fine).
- Multi-instance or load balancer (single process only).
- Production or third-party targets (never load-test systems you do not own or lack permission to test).

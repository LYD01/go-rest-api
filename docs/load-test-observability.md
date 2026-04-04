# Observability during load tests

Pair k6 (or any load generator) with **server-side** signals so you can explain latency spikes and errors in a benchmark report.

## Process metrics (always useful)

- **CPU and memory:** `htop`, `top`, or `ps` on the PID running `go run ./cmd/server`.
- **Note:** If k6 and the server run on the **same** machine, CPU% is shared; call that out in the report as a limiting factor.

## k6 output (client-side)

k6 prints:

- **http_req_duration** percentiles (p50, p95, p99)
- **http_req_failed** rate
- **iterations** and **vus**

Use `--summary-export=path.json` for a machine-readable artifact (see [scripts/load/README.md](../scripts/load/README.md)).

## Logs

The server uses `slog` in [internal/middleware/middleware.go](../internal/middleware/middleware.go) for each request. Under load, log volume can itself cost CPU; for maximum throughput experiments you might temporarily lower log level (not required for first report).

## Optional: Go pprof (non-production only)

When `PPROF_ADDR` is set, the server starts a **separate** HTTP listener on that address with the standard Go pprof endpoints (same as `import _ "net/http/pprof"` on `DefaultServeMux`).

1. Start the API with pprof (example):

   ```bash
   export PPROF_ADDR=127.0.0.1:6060
   export API_KEY=your-secret
   go run ./cmd/server
   ```

2. While a load test runs, capture profiles (from another terminal):

   ```bash
   go tool pprof -http=:0 http://127.0.0.1:6060/debug/pprof/profile?seconds=30
   ```

   Or save CPU profile to file:

   ```bash
   curl -s 'http://127.0.0.1:6060/debug/pprof/profile?seconds=15' > cpu.pb.gz
   go tool pprof -top cpu.pb.gz
   ```

3. **Mutex block profile** (optional): run the server with mutex profiling enabled, e.g. `GODEBUG=mutexprofilefraction=1`, then open `http://127.0.0.1:6060/debug/pprof/mutex` while under load — useful for comparing `GET`-heavy vs `POST`-heavy scenarios against the store lock.

**Do not** expose `PPROF_ADDR` on a public network. Use `127.0.0.1` and disable in production.

## What to paste into a report

- Table or screenshot of k6 summary (RPS, p95, error %).
- One line on server RSS before/after sustained test.
- If pprof was used: one top-N stack sample or a sentence on where time went (JSON, mutex, etc.).

# Load tests (k6)

[k6](https://k6.io/) drives HTTP scenarios against **your** server. Do not point these scripts at hosts you do not own or lack permission to test.

## Prerequisites

- [Install k6](https://grafana.com/docs/k6/latest/set-up/install-k6/) (e.g. `sudo snap install k6` on Ubuntu, or download a release).
- For **manual** runs: start the server with the same `API_KEY` as k6 (see [.env.example](../../.env.example)).

## One-command matrix (recommended)

From the repo root:

```bash
./scripts/load/run-all.sh
```

This frees `PORT` (default **18081**), starts `go run ./cmd/server`, sources `.env` so the key matches `config.Load()`, runs all standard scenarios, and writes JSON to [`docs/load-results/`](../../docs/load-results/).

Override port/key:

```bash
PORT=18082 API_KEY=my-secret ./scripts/load/run-all.sh
```

## Environment variables (manual k6)

Scripts read `BASE_URL` and `API_KEY` from the environment (`__ENV` in k6). Passing explicit `-e` flags is reliable:

```bash
k6 run -e API_KEY='your-secret' -e BASE_URL='http://127.0.0.1:8080' scripts/load/albums-get.js
```

| Variable | Default in scripts | Description |
|----------|-------------------|-------------|
| `BASE_URL` | `http://127.0.0.1:8080` | Origin only (no trailing slash) |
| `API_KEY` | *(required for album scripts)* | Must match server after `.env` is loaded |

## Scripts

| File | Purpose |
|------|---------|
| `health.js` | `GET /healthz` — no API key |
| `albums-get.js` | `GET /albums` with `X-API-Key` |
| `scenario-baseline.js` | Light load: health + albums |
| `scenario-ramp.js` | Staged VU ramp on albums |
| `scenario-sustained.js` | Fixed VUs for several minutes |
| `scenario-spike.js` | Burst then cooldown |
| `scenario-stress.js` | High VUs until failures (host-only) |
| `scenario-write-mix.js` | `GET /albums` + `POST /albums` |

## Commands

From repo root (or `scripts/load` with paths adjusted):

```bash
k6 run -e API_KEY="$API_KEY" -e BASE_URL="$BASE_URL" scripts/load/health.js
k6 run -e API_KEY="$API_KEY" -e BASE_URL="$BASE_URL" scripts/load/albums-get.js
k6 run -e API_KEY="$API_KEY" -e BASE_URL="$BASE_URL" scripts/load/scenario-baseline.js
```

JSON summary (for reports / CI):

```bash
k6 run -e API_KEY="$API_KEY" -e BASE_URL="$BASE_URL" \
  --summary-export=docs/load-results/summary-baseline.json scripts/load/scenario-baseline.js
```

## Reproducibility

Record in your report: k6 version (`k6 version`), `BASE_URL`, `API_KEY` redacted, script name, and machine specs (CPU, RAM, OS).

#!/usr/bin/env bash
# Run all k6 scenarios against a fresh local server. Use only on hosts you control.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT"

PORT="${PORT:-18081}"

# config.Load() reads .env and overwrites os env; align k6 with the same API_KEY.
set -a
if [[ -f "$ROOT/.env" ]]; then
	# shellcheck disable=SC1091
	source "$ROOT/.env"
fi
set +a
KEY="${API_KEY:-loadtest-dev-key}"

K6="${K6_BIN:-k6}"
if ! command -v "$K6" &>/dev/null; then
  if [[ -x "$ROOT/.tools/k6" ]]; then
    K6="$ROOT/.tools/k6"
  else
    echo "Install k6 (https://k6.io/docs/getting-started/installation/) or place binary at .tools/k6" >&2
    exit 1
  fi
fi

mkdir -p docs/load-results

fuser -k "${PORT}/tcp" 2>/dev/null || true
sleep 0.5

( API_KEY="$KEY" ADDR="127.0.0.1:$PORT" go run ./cmd/server ) &
SRV_PID=$!
cleanup() { kill -TERM "$SRV_PID" 2>/dev/null || true; wait "$SRV_PID" 2>/dev/null || true; }
trap cleanup EXIT

READY=0
for _ in $(seq 1 100); do
  if curl -sf "http://127.0.0.1:${PORT}/healthz" >/dev/null 2>&1; then
    READY=1
    break
  fi
  sleep 0.05
done
if [[ "$READY" -ne 1 ]]; then
  echo "server did not become ready on 127.0.0.1:${PORT}" >&2
  exit 1
fi

BASE_URL="http://127.0.0.1:${PORT}"
export API_KEY="$KEY" BASE_URL

# -e ensures __ENV is set even if the runner strips inherited env.
_k6() {
	"$K6" run -e "API_KEY=${KEY}" -e "BASE_URL=${BASE_URL}" "$@"
}

_k6 --summary-export=docs/load-results/summary-baseline.json scripts/load/scenario-baseline.js
_k6 --summary-export=docs/load-results/summary-health.json scripts/load/health.js
_k6 --summary-export=docs/load-results/summary-albums-get.json scripts/load/albums-get.js
_k6 --summary-export=docs/load-results/summary-ramp.json scripts/load/scenario-ramp.js
SUSTAIN_DURATION="${SUSTAIN_DURATION:-45s}" SUSTAIN_VUS="${SUSTAIN_VUS:-35}" \
	_k6 --summary-export=docs/load-results/summary-sustained.json scripts/load/scenario-sustained.js
_k6 --summary-export=docs/load-results/summary-spike.json scripts/load/scenario-spike.js

echo "Summaries written to docs/load-results/"

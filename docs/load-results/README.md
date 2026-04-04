# k6 summary exports

JSON files in this directory are produced by:

```bash
./scripts/load/run-all.sh
```

or manually, for example:

```bash
k6 run -e API_KEY=your-key -e BASE_URL=http://127.0.0.1:8080 \
  --summary-export=docs/load-results/summary-custom.json scripts/load/albums-get.js
```

They are suitable to diff in PRs or feed into spreadsheets. **Do not commit secrets**; summaries contain rates and latencies only.

See the narrative report: [`../perf-report-2026-04-05.md`](../perf-report-2026-04-05.md).

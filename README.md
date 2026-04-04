# go-api

A small JSON HTTP service in Go: in-memory albums, API-key protection for album routes, middleware (request ID, logging, panic recovery), and graceful shutdown.

## Requirements

- [Go](https://go.dev/dl/) 1.22 or newer (uses `net/http` method + path patterns, for example `GET /albums/{id}`)

## Configuration

Environment variables (optional `.env` in the project root is loaded automatically if present):

| Variable | Required | Description |
|----------|----------|-------------|
| `API_KEY` | Yes | Secret sent by clients as `X-API-Key` for album endpoints |
| `ADDR` | No | Listen address (default `:8080`) |

Copy the example file and edit values:

```bash
cp .env.example .env
```

## Run

From the repository root:

```bash
go run ./cmd/server
```

The server logs the address on startup. Stop with `Ctrl+C` for a graceful shutdown (up to about 10 seconds).

## HTTP API

### Health (no API key)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/healthz` | Liveness check; JSON `{"status":"ok"}` |

### Albums (requires `X-API-Key: <your API_KEY>`)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/albums` | List all albums |
| `GET` | `/albums/{id}` | Get one album by `id` |
| `POST` | `/albums` | Create an album (JSON body) |

Responses include an `X-Request-Id` header for tracing.

**Album JSON**

```json
{
  "id": "optional; generated if omitted",
  "title": "string (required)",
  "artist": "string (required)",
  "price": 0
}
```

**Example**

```bash
export API_KEY='your-secret'
go run ./cmd/server
```

```bash
# health (no key)
curl -s localhost:8080/healthz

# list albums
curl -s -H "X-API-Key: your-secret" localhost:8080/albums

# create
curl -s -X POST -H "Content-Type: application/json" -H "X-API-Key: your-secret" \
  -d '{"title":"Blue Train","artist":"John Coltrane","price":56.99}' \
  localhost:8080/albums
```

## Project layout

```
cmd/server/          # main: routing, server, signal handling
internal/albums/     # models, in-memory store, HTTP handlers
internal/config/     # env / .env loading
internal/middleware/ # request ID, logging, recover, API key gate
internal/respond/    # JSON helpers
```


# Listello monorepo architecture

Target layout for shipping Listello as a Go API and a UI — the UI talks to the API over HTTP, not to Go directly.

## Packages

```text
listello/
  api/          # Go module: domain, application, adapters, HTTP server, CLI
  ui/           # UI app; calls api over HTTP
```

Optional root `package.json` with npm/pnpm workspaces:

```json
{
  "private": true,
  "workspaces": ["ui"]
}
```

Go remains its own module (`api/go.mod`). JS workspaces do not replace the Go module.

## Responsibilities

| Package | Owns | Does not own |
|---------|------|--------------|
| `api` | Domain, application, ports/adapters, HTTP API, native CLI | Presentation, npm metadata |
| `ui` | Presentation and UX | Persistence, domain logic |

Dependency direction:

```text
ui  →  HTTP  →  api
```

## Build flow

1. **Run the API** (Go HTTP server + composition root):

   ```bash
   # from api/
   go run ./cmd/listello serve   # or equivalent entrypoint
   ```

2. **Run the UI** (Vite dev server + optional dev proxy to the API):

   ```bash
   # from ui/
   npm run dev
   ```

```text
api/
  cmd/listello        → CLI + HTTP server (composition root)
  internal/
    listello-domain   → business logic and invariants
    listello-application → use cases; defines ports
    listello-adapter  → port implementations (SQLite, event log)

ui/
  src/client/         → React app
```

## HTTP API surface

The API exposes REST (or RPC-style) endpoints over HTTP. Example shape:

```http
POST /api/lists
Content-Type: application/json

{ "name": "Next actions" }
```

```json
{ "ID": "...", "Name": "Next actions" }
```

- Request and response bodies are JSON; field naming follows one convention (camelCase or PascalCase) and is documented in the API.
- Errors use HTTP status codes and a consistent error body (e.g. `{ "error": "..." }`).
- The UI (or any HTTP client) is the consumer; no Go types cross the boundary.

`internal/` packages stay under `api`. The HTTP handlers live in the same Go module (e.g. `cmd/listello` or a dedicated `internal/listello-http` package) and may import `internal/...`.

## Persistence

File-backed SQLite (e.g. ncruces) and file event log, as today. The API is the single composition root for adapters in this layout.

## Local development

Typical workflow:

1. Start `api` on a fixed port (e.g. `:8080`).
2. Start `ui` with Vite; configure the dev server or client to proxy `/api` to the Go server.
3. Use Bruno, `api.http`, or similar to hit the API directly when debugging backend behavior.

## Migration from the current repo

The repo uses this layout: `api/` (Go module) and `ui/` (npm workspace). The UI calls the API over HTTP; there is no Node/WASM SDK.

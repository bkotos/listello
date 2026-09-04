---
name: create-api-client
description: >-
  Scaffolds UI API client functions in ui/src/lib/api/ that call HTTP
  endpoints. Use when adding or extending API clients, fetch wrappers,
  list-client/item-client, or UI integration with api-types following
  Listello architecture. For React Query hooks, use create-api-queries.
---

# Create API Client

Guide for adding UI client code in `ui/src/lib/api/`. Read this skill before adding functions that call the HTTP API.

Architecture context: see [README.md](../../../README.md) (UI client calls the API via a shared `request` helper; types come from `api-types`). Layer order: [LAYER-ORDER.md](../LAYER-ORDER.md). TDD workflow: see [.cursor/rules/tdd.mdc](../../rules/tdd.mdc).

## Upstream dependencies

**Stop and do not proceed** until the API handler and types exist. See [LAYER-ORDER.md](../LAYER-ORDER.md).

Before starting, verify:

| Check | How |
|-------|-----|
| Handler function exists | `handlers.{Handler}(...)` in `api/cmd/server/handlers/` |
| Route registered | `mux.HandleFunc("{METHOD} /api/...", ...)` in `server.go` |
| Response type in api-types | e.g. `{Resource}Response` in `api-types/index.ts` |
| Request type in api-types (writes) | e.g. `{Action}{Resource}Request` in `api-types/index.ts` |

If handler or route missing → **stop**. Tell the user to use `create-api-handler` first (which requires `create-application-service` first).

If types missing → **stop**. Tell the user to add view DTOs via `create-api-handler` and run `make api-types`.

Do not write client tests or fetch wrappers until the handler is implemented.

## Scope

**In scope:** `{resource}-client.ts`, `{resource}-client.spec.ts`.

**Out of scope** (mention as follow-ups only; do not implement unless asked):

- React Query hooks (`{resource}-queries.ts`) — use [create-api-queries](../create-api-queries/SKILL.md)

- Go API handlers and DTOs (`api/cmd/server/`, `api/internal/view-dtos/`) — use [create-api-handler](../create-api-handler/SKILL.md)
- Application / adapter layers
- React pages, components, routes (`ui/src/pages/`, `ui/src/components/`) — use [create-ui-component](../create-ui-component/SKILL.md)
- App context (`ui/src/contexts/`) — wire queries into context separately if needed
- Bruno collection (`bruno/`)

## Architecture constraints

- Client functions call the shared `request` helper from `util.ts` — never call `fetch` directly.
- **Import request and response types from `api-types`** (generated from Go view DTOs via tygo).
- One client file per API resource (`list-client.ts`, `item-client.ts`).
- Client functions are thin: build path, method, body; return typed `Promise<T>`.
- Pass `RequestInit` (especially `signal`) through on reads for query cancellation.
- Do not duplicate API path strings across files — add a new function in the resource client.

## Preconditions

These duplicate the upstream gate — all must pass:

1. **HTTP handler implemented** — not a stub; route wired in `server.go`.
2. **`api-types` current** — run `make api-types` after the handler's view DTOs were added.

If the endpoint or types are missing → **stop** and use `create-api-handler` first.

## Decision tree

1. **New resource?** → Create `{resource}-client.ts` + spec.
2. **New operation on existing resource?** → Add function to existing client.
3. **GET (read)?** → Client function; then use [create-api-queries](../create-api-queries/SKILL.md) for the hook if the page does not already query that data.
4. **POST/PUT/PATCH/DELETE (write)?** → Client function; then [create-ui-component](../create-ui-component/SKILL.md) (page calls the client + `invalidateQueries`). Do **not** add a mutation hook unless the user asks for [create-api-queries](../create-api-queries/SKILL.md).
5. **Types missing in `api-types`?** → Stop; add Go DTOs and run `make api-types`.

## Scaffold checklist

```
Task progress:
- [ ] Verify upstream: handler + route + api-types exist (stop if not — see LAYER-ORDER.md)
- [ ] Confirm API route, method, path, request/response types in api-types
- [ ] Run make api-types if Go DTOs changed recently
- [ ] Add client function(s) in ui/src/lib/api/{resource}-client.ts
- [ ] Write failing client test (mock request) in {resource}-client.spec.ts
- [ ] Run tests — confirm failure is missing behavior
- [ ] STOP — summarize spec and failure for user review
- [ ] (After approval) Implement client function(s)
- [ ] Re-run tests — confirm green
- [ ] (Follow-up) Pages/components — use [create-ui-component](../create-ui-component/SKILL.md) (client + invalidate; no mutation hook)
- [ ] (Follow-up) React Query **read** hooks — use [create-api-queries](../create-api-queries/SKILL.md) only for new reads
```

## Naming conventions

| Artifact | Convention |
|----------|------------|
| Client file | `{resource}-client.ts` (e.g. `list-client.ts`, `item-client.ts`) |
| Client test | `{resource}-client.spec.ts` |
| Client function | `{verb}{Resource}` camelCase (e.g. `getAllLists`, `createList`, `defineItem`) |
| Types import | `import type { ListResponse, CreateListRequest } from "api-types"` |

For query hooks, see [create-api-queries](../create-api-queries/SKILL.md).

## Code templates

### GET collection

```ts
import type { {Resource}Response } from "api-types";
import { request } from "./util.ts";

export async function getAll{Resources}(init?: RequestInit): Promise<{Resource}Response[]> {
  return request<{Resource}Response[]>("/api/{resources}", init);
}
```

### GET by ID

```ts
export async function get{Resource}(id: string, init?: RequestInit): Promise<{Resource}Response> {
  return request<{Resource}Response>(`/api/{resources}/${id}`, init);
}
```

### POST with request body

```ts
import type { {Resource}Response, {Action}{Resource}Request } from "api-types";

export async function {action}{Resource}(body: {Action}{Resource}Request): Promise<{Resource}Response> {
  return request<{Resource}Response>("/api/{resources}", {
    method: "POST",
    body: JSON.stringify(body),
  });
}
```

For simple single-field bodies where no request type exists yet, accept typed parameters and build the body inline — but prefer `api-types` request interfaces once the Go DTO exists.

## Test templates

### Client test — mock `request`

```ts
import { afterEach, describe, expect, it, vi } from "vitest";

vi.mock("./util.ts", () => ({
  request: vi.fn(),
}));

import { {functionName} } from "./{resource}-client.ts";
import { request } from "./util.ts";

afterEach(() => {
  vi.clearAllMocks();
});

describe("{functionName}", () => {
  it("requests {description} from the API", async () => {
    // Arrange
    const expected = { /* typed response */ };
    vi.mocked(request).mockResolvedValue(expected);

    // Act
    const result = await {functionName}(/* args */);

    // Assert
    expect(request).toHaveBeenCalledWith("/api/...", /* expected init */);
    expect(result).toEqual(expected);
  });
});
```

Assert exact path, HTTP method, and JSON body for writes.

## TDD workflow

Follow `.cursor/rules/tdd.mdc`:

1. **Red** — Write failing client test first (function may not exist). Confirm failure is missing behavior.
2. **Stop** — Summarize spec and failure. **Do not implement** until the user approves.
3. **Green** — Implement client function. Re-run tests.

Do not write client implementation in the same turn as a new failing spec.

For pages after the client is green, use [create-ui-component](../create-ui-component/SKILL.md) (call the client + `invalidateQueries`). Use [create-api-queries](../create-api-queries/SKILL.md) for **new reads**, or for a mutation hook only if the user asks.

## Verification

```bash
# Client tests
cd ui && npm test -- src/lib/api/{resource}-client.spec.ts

# Or from repo root
make test

# Regenerate types if Go DTOs changed
make api-types
```

## Further reading

Annotated walkthrough of `list-client`: [examples.md](examples.md). For query hooks, see [create-api-queries](../create-api-queries/examples.md).

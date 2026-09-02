---
name: create-api-client
description: >-
  Scaffolds UI API client functions and React Query hooks in ui/src/lib/api/
  that call HTTP endpoints. Use when adding or extending API clients, fetch
  wrappers, list-client/item-client, React Query hooks, or UI integration with
  api-types following Listello architecture.
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

**In scope:** `{resource}-client.ts`, `{resource}-client.spec.ts`, `{resource}-queries.ts`, `{resource}-queries.spec.ts`.

**Out of scope** (mention as follow-ups only; do not implement unless asked):

- Go API handlers and DTOs (`api/cmd/server/`, `api/internal/view-dtos/`) — use [create-api-handler](../create-api-handler/SKILL.md)
- Application / adapter layers
- React pages, components, routes (`ui/src/pages/`, `ui/src/components/`)
- App context (`ui/src/contexts/`) — wire queries into context separately if needed
- Bruno collection (`bruno/`)

## Architecture constraints

- Client functions call the shared `request` helper from `util.ts` — never call `fetch` directly.
- **Import request and response types from `api-types`** (generated from Go view DTOs via tygo).
- One client file per API resource (`list-client.ts`, `item-client.ts`).
- Client functions are thin: build path, method, body; return typed `Promise<T>`.
- React Query hooks live in `{resource}-queries.ts` — query keys, `useQuery`, `useMutation` with cache invalidation.
- Pass `RequestInit` (especially `signal`) through on reads for query cancellation.
- Do not duplicate API path strings across files — add a new function in the resource client.

## Preconditions

These duplicate the upstream gate — all must pass:

1. **HTTP handler implemented** — not a stub; route wired in `server.go`.
2. **`api-types` current** — run `make api-types` after the handler's view DTOs were added.

If the endpoint or types are missing → **stop** and use `create-api-handler` first.

## Decision tree

1. **New resource?** → Create `{resource}-client.ts` + spec; `{resource}-queries.ts` + spec.
2. **New operation on existing resource?** → Add function to existing client; add hook if needed.
3. **GET (read)?** → Client function + `useQuery` hook with query key.
4. **POST/PUT/PATCH (write)?** → Client function + `useMutation` hook; invalidate related query keys on success.
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
- [ ] Add query keys + useQuery/useMutation in {resource}-queries.ts
- [ ] Write failing query hook test in {resource}-queries.spec.ts
- [ ] STOP — summarize hook spec if written in same batch
- [ ] (After approval) Implement hooks
- [ ] Re-run tests — confirm green
```

## Naming conventions

| Artifact | Convention |
|----------|------------|
| Client file | `{resource}-client.ts` (e.g. `list-client.ts`, `item-client.ts`) |
| Client test | `{resource}-client.spec.ts` |
| Queries file | `{resource}-queries.ts` |
| Queries test | `{resource}-queries.spec.ts` |
| Client function | `{verb}{Resource}` camelCase (e.g. `getAllLists`, `createList`, `defineItem`) |
| Query keys | `{resource}QueryKeys` object (e.g. `listQueryKeys`) |
| Query hook | `use{Resource}Query`, `useAll{Resources}Query` |
| Mutation hook | `use{Action}{Resource}Mutation` (e.g. `useCreateListMutation`) |
| Types import | `import type { ListResponse, CreateListRequest } from "api-types"` |

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

### Query keys and useQuery

```ts
import { useQuery } from "@tanstack/react-query";
import { getAll{Resources}, get{Resource} } from "./{resource}-client.ts";

export const {resource}QueryKeys = {
  all: ["{resources}"] as const,
  detail: (id: string) => ["{resources}", id] as const,
};

export function useAll{Resources}Query() {
  return useQuery({
    queryKey: {resource}QueryKeys.all,
    queryFn: ({ signal }) => getAll{Resources}({ signal }),
  });
}

export function use{Resource}Query(id: string | undefined) {
  return useQuery({
    queryKey: {resource}QueryKeys.detail(id ?? ""),
    queryFn: ({ signal }) => get{Resource}(id!, { signal }),
    enabled: Boolean(id),
  });
}
```

### useMutation with invalidation

```ts
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { {action}{Resource} } from "./{resource}-client.ts";

export function use{Action}{Resource}Mutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: {Action}{Resource}Request) => {action}{Resource}(body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: {resource}QueryKeys.all });
    },
  });
}
```

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

### Query hook test — mock client

```ts
import { renderHook, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { createQueryWrapper } from "../../test/renderWithQueryClient.tsx";

vi.mock("./{resource}-client.ts", () => ({
  {action}{Resource}: vi.fn(),
  getAll{Resources}: vi.fn(),
}));

import { {action}{Resource}, getAll{Resources} } from "./{resource}-client.ts";
import { use{Action}{Resource}Mutation, useAll{Resources}Query } from "./{resource}-queries.ts";

describe("use{Action}{Resource}Mutation", () => {
  it("{action}s a {resource} and invalidates the list query", async () => {
    // Arrange — mock client, renderHook with createQueryWrapper()
    // Act — mutation.mutate(...)
    // Assert — client called, invalidate refetches, data updated
  });
});
```

One behavioral concern per test. Client tests mock `request`; query tests mock the client module.

## TDD workflow

Follow `.cursor/rules/tdd.mdc`:

1. **Red** — Write failing client test first (function may not exist). Confirm failure is missing behavior.
2. **Stop** — Summarize spec and failure. **Do not implement** until the user approves.
3. **Green** — Implement client function. Re-run tests.
4. Repeat for query hooks if in scope.

Do not write client implementation in the same turn as a new failing spec.

## Verification

```bash
# Client + query tests
cd ui && npm test -- src/lib/api/

# Or from repo root
make test

# Regenerate types if Go DTOs changed
make api-types
```

## Further reading

Annotated walkthrough of `list-client` and `list-queries`: [examples.md](examples.md)

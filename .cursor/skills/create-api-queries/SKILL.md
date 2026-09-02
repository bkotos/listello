---
name: create-api-queries
description: >-
  Scaffolds React Query hooks in ui/src/lib/api/{resource}-queries.ts with
  query keys, useQuery, useMutation, and cache invalidation. Use when adding
  or extending query hooks, item-queries, list-queries, or wiring client
  functions into TanStack Query following Listello architecture.
---

# Create API Queries

Guide for adding React Query hooks in `ui/src/lib/api/`. Read this skill before adding `{resource}-queries.ts` hooks.

Architecture context: see [README.md](../../../README.md). Layer order: [LAYER-ORDER.md](../LAYER-ORDER.md). TDD workflow: see [.cursor/rules/tdd.mdc](../../rules/tdd.mdc).

## Upstream dependencies

**Stop and do not proceed** until the client function exists. See [LAYER-ORDER.md](../LAYER-ORDER.md).

Before starting, verify:

| Check | How |
|-------|-----|
| Client function exists | `{verb}{Resource}(...)` in `ui/src/lib/api/{resource}-client.ts` |
| Client test passes | `npm test -- src/lib/api/{resource}-client.spec.ts` green |
| Types from `api-types` | Client return type matches generated DTO |

If client function missing → **stop**. Tell the user to use [create-api-client](../create-api-client/SKILL.md) first (which requires [create-api-handler](../create-api-handler/SKILL.md) for the HTTP endpoint).

Do not write query hook tests or hooks until the client function is implemented.

## Scope

**In scope:** `{resource}-queries.ts`, `{resource}-queries.spec.ts` — query keys, `useQuery`, `useMutation`, cache invalidation.

**Out of scope** (mention as follow-ups only; do not implement unless asked):

- Client functions (`{resource}-client.ts`) — use [create-api-client](../create-api-client/SKILL.md)
- Go API handlers and DTOs — use [create-api-handler](../create-api-handler/SKILL.md)
- React pages, components, routes (`ui/src/pages/`, `ui/src/components/`)
- App context (`ui/src/contexts/`)

## Architecture constraints

- Hooks call **client functions** from `{resource}-client.ts` — never call `request` or `fetch` directly.
- One queries file per API resource (`list-queries.ts`, `item-queries.ts`).
- Export a `{resource}QueryKeys` object — all invalidation uses these keys.
- Forward `{ signal }` from `queryFn` to client reads for cancellation.
- Use `enabled: Boolean(id)` when the hook accepts an optional ID param.
- Mutations invalidate affected query keys in `onSuccess` — only keys the spec requires.
- Do not duplicate API paths or client logic in hooks.

## Decision tree

1. **New resource?** → Create `{resource}-queries.ts` + spec.
2. **New read on existing resource?** → Add query key + `useQuery` hook; add test.
3. **New write on existing resource?** → Add `useMutation` hook; invalidate related keys on success.
4. **Collection read (no param)?** → `useAll{Resources}Query()` with `{resource}QueryKeys.all`.
5. **Read by ID?** → `use{Resource}Query(id: string | undefined)` with `{resource}QueryKeys.detail(id)`.
6. **Read scoped to parent (e.g. items by list)?** → `useAll{Resources}Query(parentId)` with `{resource}QueryKeys.byList(parentId)` or similar.

## Scaffold checklist

```
Task progress:
- [ ] Verify upstream: client function exists and client test green (stop if not — see LAYER-ORDER.md)
- [ ] Read client function signature and return type
- [ ] Decide: useQuery vs useMutation; query key shape; invalidation targets
- [ ] Write failing hook test in {resource}-queries.spec.ts (mock client module)
- [ ] Run tests — confirm failure is missing behavior
- [ ] STOP — summarize spec and failure for user review
- [ ] (After approval) Add query keys + hook(s) in {resource}-queries.ts
- [ ] Re-run tests — confirm green
```

## Naming conventions

| Artifact | Convention |
|----------|------------|
| Queries file | `{resource}-queries.ts` |
| Queries test | `{resource}-queries.spec.ts` |
| Query keys | `{resource}QueryKeys` object |
| Collection query | `useAll{Resources}Query()` |
| Detail query | `use{Resource}Query(id: string \| undefined)` |
| Scoped collection | `useAll{Resources}Query(scopeId: string \| undefined)` |
| Mutation hook | `use{Action}{Resource}Mutation` (e.g. `useCreateListMutation`) |
| Test helper | `createQueryWrapper()` from `ui/src/test/renderWithQueryClient.tsx` |
| Client mock | `vi.mock("./{resource}-client.ts", () => ({ ... }))` |

## Code templates

### Query keys

```ts
export const {resource}QueryKeys = {
  all: ["{resources}"] as const,
  detail: (id: string) => ["{resources}", id] as const,
  byList: (listId: string) => ["items", listId] as const, // scoped reads
};
```

Pick key shapes that match how data is fetched and invalidated. Existing: `listQueryKeys.all`, `listQueryKeys.detail`, `itemQueryKeys.byList`.

### useQuery — collection

```ts
import { useQuery } from "@tanstack/react-query";
import { getAll{Resources} } from "./{resource}-client.ts";

export function useAll{Resources}Query() {
  return useQuery({
    queryKey: {resource}QueryKeys.all,
    queryFn: ({ signal }) => getAll{Resources}({ signal }),
  });
}
```

### useQuery — by ID

```ts
export function use{Resource}Query(id: string | undefined) {
  return useQuery({
    queryKey: {resource}QueryKeys.detail(id ?? ""),
    queryFn: ({ signal }) => get{Resource}(id!, { signal }),
    enabled: Boolean(id),
  });
}
```

### useQuery — scoped collection

```ts
export function useAll{Resources}Query(scopeId: string | undefined) {
  return useQuery({
    queryKey: {resource}QueryKeys.byList(scopeId ?? ""),
    queryFn: ({ signal }) => getAll{Resources}(scopeId!, { signal }),
    enabled: Boolean(scopeId),
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
    mutationFn: (/* args */) => {action}{Resource}(/* args */),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: {resource}QueryKeys.all });
    },
  });
}
```

For mutations scoped to a parent (e.g. define item on a list), pass the scope ID and invalidate `{resource}QueryKeys.byList(listId)`.

## Test templates

Mock the **client module**, not `request` or `fetch`. Use `createQueryWrapper()`.

### useQuery test

```ts
import { renderHook, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { createQueryWrapper } from "../../test/renderWithQueryClient.tsx";

vi.mock("./{resource}-client.ts", () => ({
  getAll{Resources}: vi.fn(),
}));

import { getAll{Resources} } from "./{resource}-client.ts";
import { useAll{Resources}Query } from "./{resource}-queries.ts";

afterEach(() => {
  vi.clearAllMocks();
});

describe("useAll{Resources}Query", () => {
  it("loads {resources} from the client", async () => {
    const expected = [/* typed DTOs */];
    vi.mocked(getAll{Resources}).mockResolvedValue(expected);

    const { QueryWrapper } = createQueryWrapper();
    const { result } = renderHook(() => useAll{Resources}Query(), {
      wrapper: QueryWrapper,
    });

    await waitFor(() => {
      expect(result.current.data).toEqual(expected);
    });
    expect(getAll{Resources}).toHaveBeenCalledWith(
      expect.objectContaining({ signal: expect.any(AbortSignal) }),
    );
  });
});
```

For scoped reads, assert the client was called with the scope ID and `{ signal }`.

### useMutation + invalidation test

```ts
describe("use{Action}{Resource}Mutation", () => {
  it("{action}s a {resource} and invalidates the list query", async () => {
    const initial = [/* ... */];
    const updated = [/* ... */];
    vi.mocked({action}{Resource}).mockResolvedValue(/* single result */);
    vi.mocked(getAll{Resources})
      .mockResolvedValueOnce(initial)
      .mockResolvedValueOnce(updated);

    const { QueryWrapper } = createQueryWrapper();
    const { result } = renderHook(
      () => ({
        query: useAll{Resources}Query(),
        mutation: use{Action}{Resource}Mutation(),
      }),
      { wrapper: QueryWrapper },
    );

    await waitFor(() => expect(result.current.query.data).toEqual(initial));

    result.current.mutation.mutate(/* args */);

    await waitFor(() => expect(result.current.mutation.isSuccess).toBe(true));
    expect({action}{Resource}).toHaveBeenCalledWith(/* args */);
    await waitFor(() => expect(result.current.query.data).toEqual(updated));
    expect(getAll{Resources}).toHaveBeenCalledTimes(2);
  });
});
```

One behavioral concern per test. Do not test client HTTP details here — that belongs in `{resource}-client.spec.ts`.

## TDD workflow

Follow `.cursor/rules/tdd.mdc`:

1. **Red** — Write failing hook test first. Hook may not exist (compile error) or return no data. Confirm failure is missing behavior.
2. **Stop** — Summarize spec and failure. **Do not implement** until the user approves.
3. **Green** — Implement query keys + hook(s). Re-run tests.

Do not write hook implementation in the same turn as a new failing spec.

## Verification

```bash
# Query tests only
cd ui && npm test -- src/lib/api/{resource}-queries.spec.ts

# All API lib tests
cd ui && npm test -- src/lib/api/
```

## Further reading

Annotated walkthroughs of `list-queries` and `item-queries`: [examples.md](examples.md)

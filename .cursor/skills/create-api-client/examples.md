# API Client Examples

Annotated references from the Listello UI. Read when implementing client functions.

For React Query hooks, see [create-api-queries/examples.md](../create-api-queries/examples.md).

## 1. Client module: `list-client.ts`

**File:** `ui/src/lib/api/list-client.ts`

```ts
import type { ListResponse } from "api-types";
import { request } from "./util.ts";

export async function getAllLists(init?: RequestInit): Promise<ListResponse[]> {
  return request<ListResponse[]>("/api/lists", init);
}

export async function getList(id: string, init?: RequestInit): Promise<ListResponse> {
  return request<ListResponse>(`/api/lists/${id}`, init);
}

export async function createList(name: string): Promise<ListResponse> {
  return request<ListResponse>("/api/lists", {
    method: "POST",
    body: JSON.stringify({ name }),
  });
}
```

Key points:

- Types from `api-types` — matches Go `ListResponse` DTO.
- Shared `request` helper handles JSON headers, error parsing (`ApiError`).
- Reads accept optional `RequestInit` for `signal` (query cancellation).
- **Target pattern for writes:** use `CreateListRequest` from `api-types` once the Go DTO exists:

```ts
import type { CreateListRequest, ListResponse } from "api-types";

export async function createList(body: CreateListRequest): Promise<ListResponse> {
  return request<ListResponse>("/api/lists", {
    method: "POST",
    body: JSON.stringify(body),
  });
}
```

## 2. Shared `request` helper

**File:** `ui/src/lib/api/util.ts`

```ts
export async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...init?.headers,
    },
  });
  if (!response.ok) {
    // parses { error: string } from API
    throw new ApiError(response.status, message);
  }
  return response.json() as Promise<T>;
}
```

Do not bypass this helper. Vite dev proxy forwards `/api` to the Go server on `:8080`.

## 3. Client tests

**File:** `ui/src/lib/api/list-client.spec.ts`

```ts
vi.mock("./util.ts", () => ({
  request: vi.fn(),
}));

describe("getAllLists", () => {
  it("requests lists from the API", async () => {
    const lists = [{ ID: "LS_1", Name: "Work" }];
    vi.mocked(request).mockResolvedValue(lists);

    const result = await getAllLists();

    expect(request).toHaveBeenCalledWith("/api/lists", undefined);
    expect(result).toEqual(lists);
  });
});

describe("createList", () => {
  it("posts a new list to the API", async () => {
    const list = { ID: "LS_1", Name: "Next actions" };
    vi.mocked(request).mockResolvedValue(list);

    const result = await createList("Next actions");

    expect(request).toHaveBeenCalledWith("/api/lists", {
      method: "POST",
      body: JSON.stringify({ name: "Next actions" }),
    });
    expect(result).toEqual(list);
  });
});
```

Key points:

- Mock `request`, not `fetch`.
- Assert path and `RequestInit` exactly.
- Use `afterEach(() => vi.clearAllMocks())`.

## 4. Generated types

**File:** `api-types/index.ts` (from `make api-types`)

```ts
export interface ListResponse {
  ID: string;
  Name: string;
}
```

Request types (e.g. `CreateListRequest`, `DefineItemRequest`) appear here after adding Go DTOs and running `make api-types`. Client functions must use these types — do not hand-write duplicate interfaces.

## 5. App context (out of skill scope)

**File:** `ui/src/contexts/AppContext.tsx`

Context consumes query hooks:

```ts
import { useAllListsQuery } from "../lib/api/list-queries.ts";
```

Wire new resources into context or pages separately unless the user asks.

## 6. Next client example: define item

Given `POST /api/lists/{id}/items` with `DefineItemRequest` / `ItemResponse` in `api-types`:

1. **Client** (`item-client.ts`):

```ts
export async function defineItem(
  listId: string,
  body: DefineItemRequest,
): Promise<ItemResponse> {
  return request<ItemResponse>(`/api/lists/${listId}/items`, {
    method: "POST",
    body: JSON.stringify(body),
  });
}
```

2. **Client test:** assert path `/api/lists/LS_1/items`, method POST, body `{ title: "Buy milk" }`.

3. **UI:** wire the page via [create-ui-component](../create-ui-component/SKILL.md) — call the client, then `invalidateQueries`. Do not add a mutation hook unless asked.

4. **Prerequisite:** API handler + DTOs via `create-api-handler`, then `make api-types`.

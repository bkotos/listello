# API Queries Examples

Annotated references from the Listello UI. Read when implementing React Query hooks.

## 1. List queries — collection, detail, mutation

**File:** `ui/src/lib/api/list-queries.ts`

```ts
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createList, getAllLists, getList } from "./list-client.ts";

export const listQueryKeys = {
  all: ["lists"] as const,
  detail: (id: string) => ["lists", id] as const,
};

export function useAllListsQuery() {
  return useQuery({
    queryKey: listQueryKeys.all,
    queryFn: ({ signal }) => getAllLists({ signal }),
  });
}

export function useListQuery(listId: string | undefined) {
  return useQuery({
    queryKey: listQueryKeys.detail(listId ?? ""),
    queryFn: ({ signal }) => getList(listId!, { signal }),
    enabled: Boolean(listId),
  });
}

export function useCreateListMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (name: string) => createList(name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: listQueryKeys.all });
    },
  });
}
```

Key points:

- `listQueryKeys` centralizes keys for reads and invalidation.
- `signal` forwarded from `queryFn` for abort on unmount.
- `enabled: Boolean(listId)` skips fetch when route param is missing.
- `useCreateListMutation` invalidates `listQueryKeys.all` after success.

## 2. Item queries — scoped collection

**File:** `ui/src/lib/api/item-queries.ts`

```ts
import { useQuery } from "@tanstack/react-query";
import { getAllItems } from "./item-client.ts";

export const itemQueryKeys = {
  byList: (listId: string) => ["items", listId] as const,
};

export function useAllItemsQuery(listId: string | undefined) {
  return useQuery({
    queryKey: itemQueryKeys.byList(listId ?? ""),
    queryFn: ({ signal }) => getAllItems(listId!, { signal }),
    enabled: Boolean(listId),
  });
}
```

Key points:

- Items are scoped by list — key is `["items", listId]`, not a global `all`.
- Same `enabled` + `signal` pattern as `useListQuery`.
- A future `useDefineItemMutation(listId)` would invalidate `itemQueryKeys.byList(listId)`.

## 3. Query hook test — useQuery

**File:** `ui/src/lib/api/item-queries.spec.ts`

```ts
vi.mock("./item-client.ts", () => ({
  getAllItems: vi.fn(),
}));

import { getAllItems } from "./item-client.ts";
import { useAllItemsQuery } from "./item-queries.ts";
import { createQueryWrapper } from "../../test/renderWithQueryClient.tsx";

describe("useAllItemsQuery", () => {
  it("loads items for a list", async () => {
    const expected = [{ ID: "IT_1", ListID: "LS_1", Title: "Buy milk", /* ... */ }];
    vi.mocked(getAllItems).mockResolvedValue(expected);

    const { QueryWrapper } = createQueryWrapper();
    const { result } = renderHook(() => useAllItemsQuery("LS_1"), {
      wrapper: QueryWrapper,
    });

    await waitFor(() => {
      expect(result.current.data).toEqual(expected);
    });
    expect(getAllItems).toHaveBeenCalledWith(
      "LS_1",
      expect.objectContaining({ signal: expect.any(AbortSignal) }),
    );
  });
});
```

Key points:

- Mock `./item-client.ts`, not `request`.
- Use `createQueryWrapper()` from `ui/src/test/renderWithQueryClient.tsx`.
- Assert client called with scope ID and abort `signal`.

## 4. Mutation + invalidation test

**File:** `ui/src/lib/api/list-queries.spec.ts`

```ts
describe("useCreateListMutation", () => {
  it("creates a list and invalidates the lists query", async () => {
    const initialLists = [{ ID: "LS_1", Name: "Work" }];
    const updatedLists = [...initialLists, { ID: "LS_2", Name: "Reading" }];
    vi.mocked(createList).mockResolvedValue({ ID: "LS_2", Name: "Reading" });
    vi.mocked(getAllLists)
      .mockResolvedValueOnce(initialLists)
      .mockResolvedValueOnce(updatedLists);

    const { QueryWrapper } = createQueryWrapper();
    const { result } = renderHook(
      () => ({
        listsQuery: useAllListsQuery(),
        createMutation: useCreateListMutation(),
      }),
      { wrapper: QueryWrapper },
    );

    await waitFor(() => {
      expect(result.current.listsQuery.data).toEqual(initialLists);
    });

    result.current.createMutation.mutate("Reading");

    await waitFor(() => {
      expect(result.current.createMutation.isSuccess).toBe(true);
    });
    expect(createList).toHaveBeenCalledWith("Reading");
    await waitFor(() => {
      expect(result.current.listsQuery.data).toEqual(updatedLists);
    });
    expect(getAllLists).toHaveBeenCalledTimes(2);
  });
});
```

Key points:

- Render query + mutation together to test invalidation.
- Mock client twice (`mockResolvedValueOnce`) to simulate refetch after invalidate.
- Assert `getAllLists` called twice (initial load + post-mutation refetch).

## 5. Using hooks in pages

**File:** `ui/src/pages/ListPage.tsx`

```ts
import { useAllItemsQuery } from "../lib/api/item-queries.ts";
import { useListQuery } from "../lib/api/list-queries.ts";

function ListPage() {
  const { listId } = useParams();
  const { data: list } = useListQuery(listId);
  const { data: items } = useAllItemsQuery(listId);
  // ...
}
```

Page tests mock the **client module** (not the queries module) because hooks call client functions under the hood:

```ts
vi.mock("../lib/api/item-client.ts", () => ({
  getAllItems: vi.fn(),
}));
```

## 6. Next example: define item mutation

Given `defineItem(listId, body)` in `item-client.ts`:

```ts
export function useDefineItemMutation(listId: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (body: DefineItemRequest) => defineItem(listId, body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: itemQueryKeys.byList(listId) });
    },
  });
}
```

Test: mock `defineItem` + `getAllItems` twice; mutate; assert items refetch.

Prerequisite: client function via [create-api-client](../create-api-client/SKILL.md).

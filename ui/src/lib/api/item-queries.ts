import { useQuery } from "@tanstack/react-query";
import { getAllItems, getItem, getItemComments } from "./item-client";

export const itemQueryKeys = {
  byList: (listId: string) => ["items", listId] as const,
  detail: (id: string) => ["items", id] as const,
  comments: (itemId: string) => ["items", itemId, "comments"] as const,
};

export function useAllItemsQuery(listId: string | undefined) {
  return useQuery({
    queryKey: itemQueryKeys.byList(listId ?? ""),
    queryFn: ({ signal }) => getAllItems(listId!, { signal }),
    enabled: Boolean(listId),
  });
}

export function useItemQuery(id: string | undefined) {
  return useQuery({
    queryKey: itemQueryKeys.detail(id ?? ""),
    queryFn: ({ signal }) => getItem(id!, { signal }),
    enabled: Boolean(id),
  });
}

export function useItemCommentsQuery(itemId: string | undefined) {
  return useQuery({
    queryKey: itemQueryKeys.comments(itemId ?? ""),
    queryFn: ({ signal }) => getItemComments(itemId!, { signal }),
    enabled: Boolean(itemId),
  });
}

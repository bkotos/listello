import { useQuery } from "@tanstack/react-query";
import { getAllItems, getItem } from "./item-client";

export const itemQueryKeys = {
  byList: (listId: string) => ["items", listId] as const,
  detail: (id: string) => ["items", id] as const,
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

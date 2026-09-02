import { useQuery } from "@tanstack/react-query";
import { getAllItems } from "./item-client";

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

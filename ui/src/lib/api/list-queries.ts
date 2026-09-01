import { useQuery } from "@tanstack/react-query";
import { getAllLists, getList } from "./list-client.ts";

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

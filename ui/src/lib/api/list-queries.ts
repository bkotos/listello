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

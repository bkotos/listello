import { useQuery } from "@tanstack/react-query";
import { getInstance } from "./instance-client";

export const instanceQueryKeys = {
  current: ["instance"] as const,
};

export function useInstanceQuery() {
  return useQuery({
    queryKey: instanceQueryKeys.current,
    queryFn: ({ signal }) => getInstance({ signal }),
  });
}

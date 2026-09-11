import { useQuery } from "@tanstack/react-query";
import { getDefaultPersistenceLocation, getInstance } from "./instance-client";

export const instanceQueryKeys = {
  current: ["instance"] as const,
  defaultPersistenceLocation: ["instance", "default-persistence-location"] as const,
};

export function useInstanceQuery() {
  return useQuery({
    queryKey: instanceQueryKeys.current,
    queryFn: ({ signal }) => getInstance({ signal }),
  });
}

export function useDefaultPersistenceLocationQuery() {
  return useQuery({
    queryKey: instanceQueryKeys.defaultPersistenceLocation,
    queryFn: ({ signal }) => getDefaultPersistenceLocation({ signal }),
  });
}

import { createContext, useMemo, type ReactNode } from "react";
import type { ListResponse } from "api-types";
import { useAllListsQuery } from "../lib/api/list-queries";

export type AppContextValue = {
  lists: ListResponse[];
  isLoadingLists: boolean;
  listsError: string | null;
  refreshLists: () => Promise<void>;
};

export const AppContext = createContext<AppContextValue | null>(null);

type AppProviderProps = {
  children: ReactNode;
};

export function AppProvider({ children }: AppProviderProps) {
  const { data, isLoading, error, refetch } = useAllListsQuery();

  const value = useMemo(
    () => ({
      lists: data ?? [],
      isLoadingLists: isLoading,
      listsError: error instanceof Error ? error.message : error ? "Failed to load lists" : null,
      refreshLists: async () => {
        await refetch();
      },
    }),
    [data, isLoading, error, refetch],
  );

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
}

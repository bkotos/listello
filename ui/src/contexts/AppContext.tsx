import {
  createContext,
  useCallback,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import type { ListResponse } from "api-types";
import { getAllLists } from "../lib/api/list-client.ts";

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
  const [lists, setLists] = useState<ListResponse[]>([]);
  const [isLoadingLists, setIsLoadingLists] = useState(true);
  const [listsError, setListsError] = useState<string | null>(null);

  const refreshLists = useCallback(async () => {
    setIsLoadingLists(true);
    setListsError(null);

    try {
      const fetchedLists = await getAllLists();
      console.log("lists", fetchedLists);
      setLists(fetchedLists);
    } catch (error) {
      setListsError(error instanceof Error ? error.message : "Failed to load lists");
    } finally {
      setIsLoadingLists(false);
    }
  }, []);

  useEffect(() => {
    void refreshLists();
  }, [refreshLists]);

  const value = useMemo(
    () => ({
      lists,
      isLoadingLists,
      listsError,
      refreshLists,
    }),
    [lists, isLoadingLists, listsError, refreshLists],
  );

  return <AppContext.Provider value={value}>{children}</AppContext.Provider>;
}

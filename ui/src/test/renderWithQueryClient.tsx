import { QueryClientProvider } from "@tanstack/react-query";
import { createElement, type ReactNode } from "react";
import { createTestQueryClient } from "../lib/query-client.ts";

export function createQueryWrapper() {
  const queryClient = createTestQueryClient();

  function QueryWrapper({ children }: { children: ReactNode }) {
    return createElement(QueryClientProvider, { client: queryClient }, children);
  }

  return { queryClient, QueryWrapper };
}

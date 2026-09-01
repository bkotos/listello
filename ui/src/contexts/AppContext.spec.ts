import { renderHook, waitFor } from "@testing-library/react";
import { createElement, type ReactNode } from "react";
import { afterEach, describe, expect, it, vi } from "vitest";

vi.mock("../lib/api/list-client.ts", () => ({
  getAllLists: vi.fn(),
}));

import { AppProvider } from "./AppContext.tsx";
import { getAllLists } from "../lib/api/list-client.ts";
import { useAppContext } from "./useAppContext.ts";

function wrapper({ children }: { children: ReactNode }) {
  return createElement(AppProvider, null, children);
}

afterEach(() => {
  vi.clearAllMocks();
});

describe("AppProvider", () => {
  it("loads lists on mount", async () => {
    // Arrange
    const lists = [{ ID: "LS_1", Name: "Work" }];
    vi.mocked(getAllLists).mockResolvedValue(lists);
    const logSpy = vi.spyOn(console, "log").mockImplementation(() => {});

    // Act
    const { result } = renderHook(() => useAppContext(), { wrapper });

    // Assert
    await waitFor(() => {
      expect(result.current.isLoadingLists).toBe(false);
    });
    expect(getAllLists).toHaveBeenCalledOnce();
    expect(result.current.lists).toEqual(lists);
    expect(result.current.listsError).toBeNull();
    expect(logSpy).toHaveBeenCalledWith("lists", lists);
    logSpy.mockRestore();
  });

  it("stores an error when loading lists fails", async () => {
    // Arrange
    vi.mocked(getAllLists).mockRejectedValue(new Error("network error"));

    // Act
    const { result } = renderHook(() => useAppContext(), { wrapper });

    // Assert
    await waitFor(() => {
      expect(result.current.isLoadingLists).toBe(false);
    });
    expect(result.current.lists).toEqual([]);
    expect(result.current.listsError).toBe("network error");
  });

  it("refreshes lists when refreshLists is called", async () => {
    // Arrange
    const initialLists = [{ ID: "LS_1", Name: "Work" }];
    const updatedLists = [
      { ID: "LS_1", Name: "Work" },
      { ID: "LS_2", Name: "Personal" },
    ];
    vi.mocked(getAllLists)
      .mockResolvedValueOnce(initialLists)
      .mockResolvedValueOnce(updatedLists);
    const logSpy = vi.spyOn(console, "log").mockImplementation(() => {});

    const { result } = renderHook(() => useAppContext(), { wrapper });
    await waitFor(() => {
      expect(result.current.lists).toEqual(initialLists);
    });

    // Act
    await result.current.refreshLists();

    // Assert
    expect(getAllLists).toHaveBeenCalledTimes(2);
    await waitFor(() => {
      expect(result.current.lists).toEqual(updatedLists);
    });
    logSpy.mockRestore();
  });
});

import { renderHook, waitFor } from "@testing-library/react";
import { createElement, type ReactNode } from "react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { createQueryWrapper } from "../test/renderWithQueryClient";

vi.mock("../lib/api/list-client", () => ({
  getAllLists: vi.fn(),
}));

import { AppProvider } from "./AppContext";
import { getAllLists } from "../lib/api/list-client";
import { useAppContext } from "./useAppContext";

function wrapper({ children }: { children: ReactNode }) {
  const { QueryWrapper } = createQueryWrapper();

  return createElement(
    QueryWrapper,
    null,
    createElement(AppProvider, null, children),
  );
}

afterEach(() => {
  vi.clearAllMocks();
});

describe("AppProvider", () => {
  it("loads lists on mount", async () => {
    // Arrange
    const lists = [{ ID: "LS_1", Name: "Work" }];
    vi.mocked(getAllLists).mockResolvedValue(lists);

    // Act
    const { result } = renderHook(() => useAppContext(), { wrapper });

    // Assert
    await waitFor(() => {
      expect(result.current.isLoadingLists).toBe(false);
    });
    expect(getAllLists).toHaveBeenCalledOnce();
    expect(result.current.lists).toEqual(lists);
    expect(result.current.listsError).toBeNull();
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
  });
});

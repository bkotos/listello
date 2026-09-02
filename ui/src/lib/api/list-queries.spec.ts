import { renderHook, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { createQueryWrapper } from "../../test/renderWithQueryClient";

vi.mock("./list-client", () => ({
  createList: vi.fn(),
  getAllLists: vi.fn(),
  getList: vi.fn(),
}));

import { createList, getAllLists } from "./list-client";
import { useAllListsQuery, useCreateListMutation } from "./list-queries";

afterEach(() => {
  vi.clearAllMocks();
});

describe("useCreateListMutation", () => {
  it("creates a list and invalidates the lists query", async () => {
    // Arrange
    const initialLists = [{ ID: "LS_1", Name: "Work" }];
    const updatedLists = [...initialLists, { ID: "LS_2", Name: "Reading" }];
    vi.mocked(createList).mockResolvedValue({ ID: "LS_2", Name: "Reading" });
    vi.mocked(getAllLists)
      .mockResolvedValueOnce(initialLists)
      .mockResolvedValueOnce(updatedLists);

    const { QueryWrapper } = createQueryWrapper();
    const { result } = renderHook(
      () => ({
        listsQuery: useAllListsQuery(),
        createMutation: useCreateListMutation(),
      }),
      { wrapper: QueryWrapper },
    );

    await waitFor(() => {
      expect(result.current.listsQuery.data).toEqual(initialLists);
    });

    // Act
    result.current.createMutation.mutate("Reading");

    // Assert
    await waitFor(() => {
      expect(result.current.createMutation.isSuccess).toBe(true);
    });
    expect(createList).toHaveBeenCalledWith("Reading");
    await waitFor(() => {
      expect(result.current.listsQuery.data).toEqual(updatedLists);
    });
    expect(getAllLists).toHaveBeenCalledTimes(2);
  });
});

import { renderHook, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { createQueryWrapper } from "../../test/renderWithQueryClient";

vi.mock("./item-client", () => ({
  getAllItems: vi.fn(),
  getItem: vi.fn(),
}));

import { getAllItems, getItem } from "./item-client";
import { useAllItemsQuery, useItemQuery } from "./item-queries";

afterEach(() => {
  vi.clearAllMocks();
});

describe("useAllItemsQuery", () => {
  it("loads items for a list", async () => {
    // Arrange
    const expected = [
      {
        ID: "IT_1",
        ListID: "LS_1",
        ParentID: "",
        Title: "Buy milk",
        Description: "",
        DueDate: "",
        Tags: [],
        Priority: "",
        State: "outstanding",
      },
    ];
    vi.mocked(getAllItems).mockResolvedValue(expected);

    const { QueryWrapper } = createQueryWrapper();
    const { result } = renderHook(() => useAllItemsQuery("LS_1"), {
      wrapper: QueryWrapper,
    });

    // Assert
    await waitFor(() => {
      expect(result.current.data).toEqual(expected);
    });
    expect(getAllItems).toHaveBeenCalledWith(
      "LS_1",
      expect.objectContaining({ signal: expect.any(AbortSignal) }),
    );
  });
});

describe("useItemQuery", () => {
  it("loads an item by id from the client", async () => {
    const expected = {
      ID: "IT_1",
      ListID: "LS_2",
      ParentID: "",
      Title: "Plan trip",
      Description: "",
      DueDate: "",
      Tags: [],
      Priority: "",
      State: "outstanding",
    };
    vi.mocked(getItem).mockResolvedValue(expected);

    const { QueryWrapper } = createQueryWrapper();
    const { result } = renderHook(() => useItemQuery("IT_1"), {
      wrapper: QueryWrapper,
    });

    // Assert
    await waitFor(() => {
      expect(result.current.data).toEqual(expected);
    });
    expect(getItem).toHaveBeenCalledWith(
      "IT_1",
      expect.objectContaining({ signal: expect.any(AbortSignal) }),
    );
  });
});

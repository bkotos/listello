import { renderHook, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { createQueryWrapper } from "../../test/renderWithQueryClient";

vi.mock("./item-client", () => ({
  getAllItems: vi.fn(),
  getItem: vi.fn(),
  getItemComments: vi.fn(),
}));

import { getAllItems, getItem, getItemComments } from "./item-client";
import { useAllItemsQuery, useItemCommentsQuery, useItemQuery } from "./item-queries";

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

describe("useItemCommentsQuery", () => {
  it("loads comments for an item from the client", async () => {
    const expected = [
      {
        ID: "CM_1",
        ItemID: "IT_1",
        UserID: "US_1",
        UserName: "Alex",
        Body: "Need 2%",
        CreatedAt: "2026-10-17T00:00:00Z",
      },
    ];
    vi.mocked(getItemComments).mockResolvedValue(expected);

    const { QueryWrapper } = createQueryWrapper();
    const { result } = renderHook(() => useItemCommentsQuery("IT_1"), {
      wrapper: QueryWrapper,
    });

    await waitFor(() => {
      expect(result.current.data).toEqual(expected);
    });
    expect(getItemComments).toHaveBeenCalledWith(
      "IT_1",
      expect.objectContaining({ signal: expect.any(AbortSignal) }),
    );
  });
});

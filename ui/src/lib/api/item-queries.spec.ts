import { renderHook, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { createQueryWrapper } from "../../test/renderWithQueryClient.tsx";

vi.mock("./item-client.ts", () => ({
  getAllItems: vi.fn(),
}));

import { getAllItems } from "./item-client.ts";
import { useAllItemsQuery } from "./item-queries.ts";

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

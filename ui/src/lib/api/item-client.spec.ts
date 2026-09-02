import { afterEach, describe, expect, it, vi } from "vitest";

vi.mock("./util.ts", () => ({
  request: vi.fn(),
}));

import { defineItem } from "./item-client.ts";
import { request } from "./util.ts";

afterEach(() => {
  vi.clearAllMocks();
});

describe("defineItem", () => {
  it("posts a new item on a list to the API", async () => {
    // Arrange
    const item = {
      ID: "IT_1",
      ListID: "LS_1",
      ParentID: "",
      Title: "Buy milk",
      Description: "",
      DueDate: "",
      Tags: [],
      Priority: "",
      State: "outstanding",
    };
    vi.mocked(request).mockResolvedValue(item);

    // Act
    const result = await defineItem("LS_1", { title: "Buy milk" });

    // Assert
    expect(request).toHaveBeenCalledWith("/api/lists/LS_1/items", {
      method: "POST",
      body: JSON.stringify({ title: "Buy milk" }),
    });
    expect(result).toEqual(item);
  });
});

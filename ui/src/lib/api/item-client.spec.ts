import { afterEach, describe, expect, it, vi } from "vitest";
import { completeItem, defineItem, getAllItems } from "./item-client";
import { request } from "./util";

vi.mock(import("./util"), () => ({
  request: vi.fn(),
}));

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

describe("getAllItems", () => {
  it("requests all items for a list from the API", async () => {
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
      {
        ID: "IT_2",
        ListID: "LS_1",
        ParentID: "",
        Title: "Call dentist",
        Description: "",
        DueDate: "",
        Tags: [],
        Priority: "",
        State: "outstanding",
      },
    ];
    const init = { signal: new AbortController().signal };
    vi.mocked(request).mockResolvedValue(expected);

    // Act
    const result = await getAllItems("LS_1", init);

    // Assert
    expect(request).toHaveBeenCalledWith("/api/lists/LS_1/items", init);
    expect(result).toEqual(expected);
  });
});

describe("completeItem", () => {
  it("posts a complete action for an item to the API", async () => {
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
      State: "complete",
    };
    vi.mocked(request).mockResolvedValue(item);

    // Act
    const result = await completeItem("IT_1");

    // Assert
    expect(request).toHaveBeenCalledWith("/api/items/IT_1/complete", {
      method: "POST",
    });
    expect(result).toEqual(item);
  });
});

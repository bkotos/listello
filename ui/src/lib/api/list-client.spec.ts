import { afterEach, describe, expect, it, vi } from "vitest";

vi.mock("./util", () => ({
  request: vi.fn(),
}));

import { createFirstList, createList, deleteList, getAllLists, getList } from "./list-client";
import { request } from "./util";

afterEach(() => {
  vi.clearAllMocks();
});

describe("getAllLists", () => {
  it("requests lists from the API", async () => {
    // Arrange
    const lists = [{ ID: "LS_1", Name: "Work" }];
    vi.mocked(request).mockResolvedValue(lists);

    // Act
    const result = await getAllLists();

    // Assert
    expect(request).toHaveBeenCalledWith("/api/lists", undefined);
    expect(result).toEqual(lists);
  });
});

describe("getList", () => {
  it("requests a list by id from the API", async () => {
    // Arrange
    const list = { ID: "LS_1", Name: "Work" };
    vi.mocked(request).mockResolvedValue(list);

    // Act
    const result = await getList("LS_1");

    // Assert
    expect(request).toHaveBeenCalledWith("/api/lists/LS_1", undefined);
    expect(result).toEqual(list);
  });
});

describe("createList", () => {
  it("posts a new list to the API", async () => {
    // Arrange
    const list = { ID: "LS_1", Name: "Next actions" };
    vi.mocked(request).mockResolvedValue(list);

    // Act
    const result = await createList("Next actions");

    // Assert
    expect(request).toHaveBeenCalledWith("/api/lists", {
      method: "POST",
      body: JSON.stringify({ name: "Next actions" }),
    });
    expect(result).toEqual(list);
  });
});

describe("createFirstList", () => {
  it("posts the first list to the API", async () => {
    // Arrange
    const list = { ID: "LS_1", Name: "Errands" };
    vi.mocked(request).mockResolvedValue(list);

    // Act
    const result = await createFirstList("Errands");

    // Assert
    expect(request).toHaveBeenCalledWith("/api/lists/first", {
      method: "POST",
      body: JSON.stringify({ name: "Errands" }),
    });
    expect(result).toEqual(list);
  });
});

describe("deleteList", () => {
  it("deletes a list via the API", async () => {
    // Arrange
    vi.mocked(request).mockResolvedValue(undefined);

    // Act
    const result = await deleteList("LS_1");

    // Assert
    expect(request).toHaveBeenCalledWith("/api/lists/LS_1", {
      method: "DELETE",
    });
    expect(result).toBeUndefined();
  });
});

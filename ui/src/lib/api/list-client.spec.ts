import { afterEach, describe, expect, it, vi } from "vitest";

vi.mock("./util.ts", () => ({
  request: vi.fn(),
}));

import { createList, getAllLists } from "./list-client.ts";
import { request } from "./util.ts";

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
    expect(request).toHaveBeenCalledWith("/api/lists");
    expect(result).toEqual(lists);
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

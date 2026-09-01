import { afterEach, describe, expect, it, vi } from "vitest";
import { ApiError, createList, getAllLists } from "./client.ts";

afterEach(() => {
  vi.unstubAllGlobals();
});

describe("getAllLists", () => {
  it("fetches lists from the API", async () => {
    // Arrange
    const lists = [{ ID: "LS_1", Name: "Work" }];
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => lists,
    });
    vi.stubGlobal("fetch", fetchMock);

    // Act
    const result = await getAllLists();

    // Assert
    expect(fetchMock).toHaveBeenCalledWith("/api/lists", {
      headers: { "Content-Type": "application/json" },
    });
    expect(result).toEqual(lists);
  });
});

describe("createList", () => {
  it("posts a new list to the API", async () => {
    // Arrange
    const list = { ID: "LS_1", Name: "Next actions" };
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 201,
      json: async () => list,
    });
    vi.stubGlobal("fetch", fetchMock);

    // Act
    const result = await createList("Next actions");

    // Assert
    expect(fetchMock).toHaveBeenCalledWith("/api/lists", {
      method: "POST",
      body: JSON.stringify({ name: "Next actions" }),
      headers: { "Content-Type": "application/json" },
    });
    expect(result).toEqual(list);
  });
});

describe("ApiError", () => {
  it("uses the API error message when the response body includes error", async () => {
    // Arrange
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: false,
        status: 400,
        statusText: "Bad Request",
        json: async () => ({ error: "name is required" }),
      }),
    );

    // Act
    const error = await createList("").catch((caught: unknown) => caught);

    // Assert
    expect(error).toBeInstanceOf(ApiError);
    expect(error).toMatchObject({
      status: 400,
      message: "name is required",
    });
  });

  it("falls back to status text when the error body is not JSON", async () => {
    // Arrange
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({
        ok: false,
        status: 500,
        statusText: "Internal Server Error",
        json: async () => {
          throw new Error("invalid json");
        },
      }),
    );

    // Act
    const error = await getAllLists().catch((caught: unknown) => caught);

    // Assert
    expect(error).toBeInstanceOf(ApiError);
    expect(error).toMatchObject({
      status: 500,
      message: "Internal Server Error",
    });
  });
});

import { afterEach, describe, expect, it, vi } from "vitest";
import { ApiError, request } from "./util";

afterEach(() => {
  vi.unstubAllGlobals();
});

describe("request", () => {
  it("returns parsed JSON for successful responses", async () => {
    // Arrange
    const payload = { ID: "LS_1", Name: "Work" };
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => payload,
    });
    vi.stubGlobal("fetch", fetchMock);

    // Act
    const result = await request<typeof payload>("/api/lists");

    // Assert
    expect(fetchMock).toHaveBeenCalledWith("/api/lists", {
      headers: { "Content-Type": "application/json" },
    });
    expect(result).toEqual(payload);
  });

  it("forwards request init to fetch", async () => {
    // Arrange
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      status: 201,
      json: async () => ({}),
    });
    vi.stubGlobal("fetch", fetchMock);

    // Act
    await request("/api/lists", {
      method: "POST",
      body: JSON.stringify({ name: "Next actions" }),
    });

    // Assert
    expect(fetchMock).toHaveBeenCalledWith("/api/lists", {
      method: "POST",
      body: JSON.stringify({ name: "Next actions" }),
      headers: { "Content-Type": "application/json" },
    });
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
    const error = await request("/api/lists").catch((caught: unknown) => caught);

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
    const error = await request("/api/lists").catch((caught: unknown) => caught);

    // Assert
    expect(error).toBeInstanceOf(ApiError);
    expect(error).toMatchObject({
      status: 500,
      message: "Internal Server Error",
    });
  });
});

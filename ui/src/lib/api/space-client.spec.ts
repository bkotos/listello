import { afterEach, describe, expect, it, vi } from "vitest";

vi.mock("./util.ts", () => ({
  request: vi.fn(),
}));

import { createSpace } from "./space-client";
import { request } from "./util";

describe("createSpace", () => {
  afterEach(() => {
    vi.clearAllMocks();
  });

  it("calls POST /api/spaces with the space data", async () => {
    // Arrange
    const body = { name: "Personal" };
    const expected = { ID: "SP_1", Name: "Personal" };
    vi.mocked(request).mockResolvedValue(expected);

    // Act
    const result = await createSpace(body);

    // Assert
    expect(request).toHaveBeenCalledWith("/api/spaces", {
      method: "POST",
      body: JSON.stringify(body),
    });
    expect(result).toEqual(expected);
  });
});

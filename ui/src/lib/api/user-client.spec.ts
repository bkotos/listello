import { afterEach, describe, expect, it, vi } from "vitest";

vi.mock("./util", () => ({
  request: vi.fn(),
}));

import { createUser } from "./user-client";
import { request } from "./util";

afterEach(() => {
  vi.clearAllMocks();
});

describe("createUser", () => {
  it("posts a new user to the API", async () => {
    // Arrange
    const user = { ID: "US_1", Name: "Alex" };
    vi.mocked(request).mockResolvedValue(user);

    // Act
    const result = await createUser({ name: "Alex" });

    // Assert
    expect(request).toHaveBeenCalledWith("/api/users", {
      method: "POST",
      body: JSON.stringify({ name: "Alex" }),
    });
    expect(result).toEqual(user);
  });
});

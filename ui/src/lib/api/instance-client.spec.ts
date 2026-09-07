import { afterEach, describe, expect, it, vi } from "vitest";

vi.mock("./util", () => ({
  request: vi.fn(),
}));

import { createInstance, getInstance } from "./instance-client";
import { request } from "./util";

afterEach(() => {
  vi.clearAllMocks();
});

describe("createInstance", () => {
  it("posts a new instance to the API", async () => {
    // Arrange
    const instance = {
      HostingMode: "",
      PersistenceLocation: "",
      PersistenceState: "",
      SetupState: "",
    };
    vi.mocked(request).mockResolvedValue(instance);

    // Act
    const result = await createInstance();

    // Assert
    expect(request).toHaveBeenCalledWith("/api/instance", {
      method: "POST",
    });
    expect(result).toEqual(instance);
  });
});

describe("getInstance", () => {
  it("requests the instance from the API", async () => {
    // Arrange
    const instance = {
      HostingMode: "",
      PersistenceLocation: "",
      PersistenceState: "",
      SetupState: "",
    };
    vi.mocked(request).mockResolvedValue(instance);

    // Act
    const result = await getInstance();

    // Assert
    expect(request).toHaveBeenCalledWith("/api/instance", undefined);
    expect(result).toEqual(instance);
  });

  it("returns null when the instance does not exist yet", async () => {
    // Arrange
    vi.mocked(request).mockResolvedValue(null);

    // Act
    const result = await getInstance();

    // Assert
    expect(result).toBeNull();
  });
});

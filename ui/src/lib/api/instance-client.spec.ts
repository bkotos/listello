import { afterEach, describe, expect, it, vi } from "vitest";

vi.mock("./util", () => ({
  request: vi.fn(),
}));

import { createInstance, getDefaultPersistenceLocation, getInstance, selectHostingMode, selectPersistenceLocation } from "./instance-client";
import { request } from "./util";

afterEach(() => {
  vi.clearAllMocks();
});

describe("getDefaultPersistenceLocation", () => {
  it("requests the default persistence location from the API", async () => {
    // Arrange
    const expected = { Location: "/Users/me/Library/Application Support/listello" };
    vi.mocked(request).mockResolvedValue(expected);

    // Act
    const result = await getDefaultPersistenceLocation();

    // Assert
    expect(request).toHaveBeenCalledWith("/api/instance/default-persistence-location", undefined);
    expect(result).toEqual(expected);
  });
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

describe("selectPersistenceLocation", () => {
  it("posts the persistence location to the API", async () => {
    // Arrange
    const instance = {
      HostingMode: "local",
      PersistenceLocation: "/var/listello",
      PersistenceState: "",
      SetupState: "",
    };
    vi.mocked(request).mockResolvedValue(instance);

    // Act
    const result = await selectPersistenceLocation({ location: "/var/listello" });

    // Assert
    expect(request).toHaveBeenCalledWith("/api/instance/persistence-location", {
      method: "POST",
      body: JSON.stringify({ location: "/var/listello" }),
    });
    expect(result).toEqual(instance);
  });
});

describe("selectHostingMode", () => {
  it("posts the hosting mode to the API", async () => {
    // Arrange
    const instance = {
      HostingMode: "local",
      PersistenceLocation: "",
      PersistenceState: "",
      SetupState: "",
    };
    vi.mocked(request).mockResolvedValue(instance);

    // Act
    const result = await selectHostingMode({ mode: "local" });

    // Assert
    expect(request).toHaveBeenCalledWith("/api/instance/hosting-mode", {
      method: "POST",
      body: JSON.stringify({ mode: "local" }),
    });
    expect(result).toEqual(instance);
  });
});

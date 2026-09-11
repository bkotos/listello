import { renderHook, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { createQueryWrapper } from "../../test/renderWithQueryClient";

vi.mock("./instance-client", () => ({
  getInstance: vi.fn(),
  getDefaultPersistenceLocation: vi.fn(),
}));

import { getDefaultPersistenceLocation, getInstance } from "./instance-client";
import { useDefaultPersistenceLocationQuery, useInstanceQuery } from "./instance-queries";

afterEach(() => {
  vi.clearAllMocks();
});

describe("useDefaultPersistenceLocationQuery", () => {
  it("loads the default persistence location from the client", async () => {
    // Arrange
    const expected = { Location: "/Users/me/Library/Application Support/listello" };
    vi.mocked(getDefaultPersistenceLocation).mockResolvedValue(expected);

    const { QueryWrapper } = createQueryWrapper();
    const { result } = renderHook(() => useDefaultPersistenceLocationQuery(), {
      wrapper: QueryWrapper,
    });

    // Assert
    await waitFor(() => {
      expect(result.current.data).toEqual(expected);
    });
    expect(getDefaultPersistenceLocation).toHaveBeenCalledWith(
      expect.objectContaining({ signal: expect.any(AbortSignal) }),
    );
  });
});

describe("useInstanceQuery", () => {
  it("loads the instance from the client", async () => {
    // Arrange
    const expected = {
      HostingMode: "",
      PersistenceLocation: "",
      PersistenceState: "",
      SetupState: "",
    };
    vi.mocked(getInstance).mockResolvedValue(expected);

    const { QueryWrapper } = createQueryWrapper();
    const { result } = renderHook(() => useInstanceQuery(), {
      wrapper: QueryWrapper,
    });

    // Assert
    await waitFor(() => {
      expect(result.current.data).toEqual(expected);
    });
    expect(getInstance).toHaveBeenCalledWith(
      expect.objectContaining({ signal: expect.any(AbortSignal) }),
    );
  });

  it("returns null when the instance does not exist yet", async () => {
    // Arrange
    vi.mocked(getInstance).mockResolvedValue(null);

    const { QueryWrapper } = createQueryWrapper();
    const { result } = renderHook(() => useInstanceQuery(), {
      wrapper: QueryWrapper,
    });

    // Assert
    await waitFor(() => {
      expect(result.current.data).toBeNull();
    });
  });
});

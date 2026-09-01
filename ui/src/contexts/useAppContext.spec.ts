import { renderHook } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import { useAppContext } from "./useAppContext.ts";

describe("useAppContext", () => {
  it("throws when used outside AppProvider", () => {
    // Arrange
    const consoleError = vi.spyOn(console, "error").mockImplementation(() => {});

    // Act & Assert
    expect(() => renderHook(() => useAppContext())).toThrow(
      "useAppContext must be used within AppProvider",
    );

    consoleError.mockRestore();
  });
});

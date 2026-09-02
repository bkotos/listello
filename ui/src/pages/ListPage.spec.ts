import { cleanup, fireEvent, screen, waitFor } from "@testing-library/react";
import { createElement } from "react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { renderPageWithShellContext } from "../test/renderPageWithShellContext.ts";
import ListPage from "./ListPage.tsx";

vi.mock("../lib/api/list-client.ts", () => ({
  getList: vi.fn(),
}));

vi.mock("../lib/api/item-client.ts", () => ({
  getAllItems: vi.fn(),
}));

import { getAllItems } from "../lib/api/item-client.ts";
import { getList } from "../lib/api/list-client.ts";

afterEach(() => {
  cleanup();
  vi.clearAllMocks();
});

describe("ListPage", () => {
  it("renders the list title from the API", async () => {
    // Arrange
    vi.mocked(getList).mockResolvedValue({ ID: "LS_1", Name: "Work" });
    vi.mocked(getAllItems).mockResolvedValue([]);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_1",
    });

    // Assert
    await waitFor(() => {
      expect(screen.getByRole("heading", { name: "Work" })).toBeInTheDocument();
    });
    expect(getList).toHaveBeenCalledWith(
      "LS_1",
      expect.objectContaining({ signal: expect.any(AbortSignal) }),
    );
    expect(screen.getByPlaceholderText("Add to Work…")).toBeInTheDocument();
    expect(screen.getByText("No items yet.")).toBeInTheDocument();
  });

  it("renders a generic title when the list cannot be loaded", async () => {
    // Arrange
    vi.mocked(getList).mockRejectedValue(new Error("not found"));
    vi.mocked(getAllItems).mockResolvedValue([]);
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_missing",
    });

    // Assert
    await waitFor(() => {
      expect(getList).toHaveBeenCalledWith(
        "LS_missing",
        expect.objectContaining({ signal: expect.any(AbortSignal) }),
      );
    });
    expect(screen.getByRole("heading", { name: "List" })).toBeInTheDocument();
    expect(screen.getByPlaceholderText("Add to List…")).toBeInTheDocument();
  });

  it("calls openSidebar when the open menu button is clicked", async () => {
    // Arrange
    vi.mocked(getList).mockResolvedValue({ ID: "LS_2", Name: "Personal" });
    vi.mocked(getAllItems).mockResolvedValue([]);
    const { openSidebar } = renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/LS_2",
    });
    await waitFor(() => {
      expect(screen.getByRole("heading", { name: "Personal" })).toBeInTheDocument();
    });

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Open menu" }));

    // Assert
    expect(openSidebar).toHaveBeenCalledOnce();
  });
});

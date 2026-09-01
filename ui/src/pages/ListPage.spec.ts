import { cleanup, fireEvent, screen } from "@testing-library/react";
import { createElement } from "react";
import { afterEach, describe, expect, it } from "vitest";
import { renderPageWithShellContext } from "../test/renderPageWithShellContext.ts";
import ListPage from "./ListPage.tsx";

afterEach(() => {
  cleanup();
});

describe("ListPage", () => {
  it("renders the list title from the route param", () => {
    // Arrange
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/work",
    });

    // Assert
    expect(screen.getByRole("heading", { name: "Work" })).toBeInTheDocument();
    expect(screen.getByPlaceholderText("Add to Work…")).toBeInTheDocument();
    expect(screen.getByText("No items yet.")).toBeInTheDocument();
  });

  it("renders a generic title for an unknown list id", () => {
    // Arrange
    renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/unknown",
    });

    // Assert
    expect(screen.getByRole("heading", { name: "List" })).toBeInTheDocument();
    expect(screen.getByPlaceholderText("Add to List…")).toBeInTheDocument();
  });

  it("calls openSidebar when the open menu button is clicked", () => {
    // Arrange
    const { openSidebar } = renderPageWithShellContext(createElement(ListPage), {
      path: "lists/:listId",
      initialEntry: "/lists/personal",
    });

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Open menu" }));

    // Assert
    expect(openSidebar).toHaveBeenCalledOnce();
  });
});

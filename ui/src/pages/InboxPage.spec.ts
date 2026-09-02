import { cleanup, fireEvent, screen } from "@testing-library/react";
import { createElement } from "react";
import { afterEach, describe, expect, it } from "vitest";
import { renderPageWithShellContext } from "../test/renderPageWithShellContext";
import InboxPage from "./InboxPage";

afterEach(() => {
  cleanup();
});

describe("InboxPage", () => {
  it("renders the inbox content area", () => {
    // Arrange
    renderPageWithShellContext(createElement(InboxPage), {
      path: "inbox",
      initialEntry: "/inbox",
    });

    // Assert
    expect(screen.getByRole("heading", { name: "Inbox" })).toBeInTheDocument();
    expect(screen.getByPlaceholderText("Capture something on your mind…")).toBeInTheDocument();
    expect(screen.getByText("No captured items yet.")).toBeInTheDocument();
  });

  it("calls openSidebar when the open menu button is clicked", () => {
    // Arrange
    const { openSidebar } = renderPageWithShellContext(createElement(InboxPage), {
      path: "inbox",
      initialEntry: "/inbox",
    });

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Open menu" }));

    // Assert
    expect(openSidebar).toHaveBeenCalledOnce();
  });
});

import { cleanup, fireEvent, render, screen } from "@testing-library/react";
import { createElement } from "react";
import { Inbox } from "lucide-react";
import { afterEach, describe, expect, it, vi } from "vitest";
import { ContentArea, type ContentAreaProps } from "./ContentArea.tsx";

type RenderContentAreaOptions = Partial<ContentAreaProps>;

afterEach(() => {
  cleanup();
});

function renderContentArea(
  {
    children = "No captured items yet.",
    title = "Inbox",
    capturePlaceholder = "Capture something on your mind…",
    ...rest
  }: RenderContentAreaOptions = {},
) {
  return render(
    createElement(ContentArea, {
      title,
      capturePlaceholder,
      children,
      ...rest,
    }),
  );
}

describe("ContentArea", () => {
  it("renders the title, capture input, and children", () => {
    // Arrange
    renderContentArea();

    // Assert
    expect(screen.getByRole("heading", { name: "Inbox" })).toBeInTheDocument();
    expect(screen.getByPlaceholderText("Capture something on your mind…")).toBeInTheDocument();
    expect(screen.getByText("No captured items yet.")).toBeInTheDocument();
  });

  it("renders an optional icon", () => {
    // Arrange
    const { container } = renderContentArea({
      icon: createElement(Inbox, { "aria-label": "Inbox icon" }),
    });

    // Assert
    expect(container.querySelector('[aria-label="Inbox icon"]')).toBeInTheDocument();
  });

  it("does not show the open menu button without onOpenSidebar", () => {
    // Arrange
    renderContentArea();

    // Assert
    expect(screen.queryByRole("button", { name: "Open menu" })).not.toBeInTheDocument();
  });

  it("shows the open menu button when onOpenSidebar is provided", () => {
    // Arrange
    renderContentArea({ onOpenSidebar: vi.fn() });

    // Assert
    expect(screen.getByRole("button", { name: "Open menu" })).toBeInTheDocument();
  });

  it("calls onOpenSidebar when the open menu button is clicked", () => {
    // Arrange
    const onOpenSidebar = vi.fn();
    renderContentArea({ onOpenSidebar });

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Open menu" }));

    // Assert
    expect(onOpenSidebar).toHaveBeenCalledOnce();
  });
});

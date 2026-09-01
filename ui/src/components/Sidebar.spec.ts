import { cleanup, fireEvent, render, screen } from "@testing-library/react";
import { createElement } from "react";
import { MemoryRouter } from "react-router-dom";
import { afterEach, describe, expect, it, vi } from "vitest";
import { Sidebar } from "./Sidebar.tsx";

type RenderSidebarOptions = {
  initialEntry?: string;
  onNavigate?: () => void;
};

afterEach(() => {
  cleanup();
});

function renderSidebar({ initialEntry = "/inbox", onNavigate }: RenderSidebarOptions = {}) {
  render(
    createElement(
      MemoryRouter,
      { initialEntries: [initialEntry] },
      createElement(Sidebar, { onNavigate }),
    ),
  );
}

describe("Sidebar", () => {
  it("renders branding, inbox, lists, and account menu", () => {
    // Arrange
    renderSidebar();

    // Assert
    expect(screen.getByText("Listello")).toBeInTheDocument();
    expect(screen.getByText("Lists")).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Inbox" })).toHaveAttribute("href", "/inbox");
    expect(screen.getByRole("link", { name: "Work" })).toHaveAttribute("href", "/lists/work");
    expect(screen.getByRole("link", { name: "Personal" })).toHaveAttribute("href", "/lists/personal");
    expect(screen.getByRole("link", { name: "Reading" })).toHaveAttribute("href", "/lists/reading");
    expect(screen.getByRole("button", { name: "Add list" })).toBeInTheDocument();
    expect(screen.getByText("Signed in")).toBeInTheDocument();
  });

  it("highlights the inbox link on the inbox route", () => {
    // Arrange
    renderSidebar({ initialEntry: "/inbox" });

    // Assert
    expect(screen.getByRole("link", { name: "Inbox" })).toHaveClass("is-active");
    expect(screen.getByRole("link", { name: "Work" })).not.toHaveClass("is-active");
  });

  it("highlights the active list link on a list route", () => {
    // Arrange
    renderSidebar({ initialEntry: "/lists/work" });

    // Assert
    expect(screen.getByRole("link", { name: "Inbox" })).not.toHaveClass("is-active");
    expect(screen.getByRole("link", { name: "Work" })).toHaveClass("is-active");
  });

  it("calls onNavigate when the inbox link is clicked", () => {
    // Arrange
    const onNavigate = vi.fn();
    renderSidebar({ onNavigate });

    // Act
    fireEvent.click(screen.getByRole("link", { name: "Inbox" }));

    // Assert
    expect(onNavigate).toHaveBeenCalledOnce();
  });

  it("calls onNavigate when a list link is clicked", () => {
    // Arrange
    const onNavigate = vi.fn();
    renderSidebar({ onNavigate });

    // Act
    fireEvent.click(screen.getByRole("link", { name: "Personal" }));

    // Assert
    expect(onNavigate).toHaveBeenCalledOnce();
  });
});

import { cleanup, fireEvent, render, screen } from "@testing-library/react";
import { createElement } from "react";
import { MemoryRouter } from "react-router-dom";
import { afterEach, describe, expect, it, vi } from "vitest";
import { Sidebar } from "./Sidebar.tsx";

vi.mock("../contexts/useAppContext.ts", () => ({
  useAppContext: () => ({
    lists: [
      { ID: "work", Name: "Work" },
      { ID: "personal", Name: "Personal" },
      { ID: "reading", Name: "Reading" },
    ],
    isLoadingLists: false,
    listsError: null,
    refreshLists: vi.fn(),
  }),
}));

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

  it("shows a list name input at the bottom of the lists section when Add list is clicked", () => {
    // Arrange
    renderSidebar();

    // Assert
    expect(screen.queryByPlaceholderText("List name")).not.toBeInTheDocument();

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Add list" }));

    // Assert
    const input = screen.getByPlaceholderText("List name");
    expect(input).toBeInTheDocument();
    expect(input).toHaveClass("input", "is-small", "mt-2");
  });

  it("focuses the list name input when Add list is clicked", () => {
    // Arrange
    renderSidebar();

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Add list" }));

    // Assert
    expect(screen.getByPlaceholderText("List name")).toHaveFocus();
  });

  it("hides the list name input when Add list is clicked again", () => {
    // Arrange
    renderSidebar();
    const addButton = screen.getByRole("button", { name: "Add list" });
    fireEvent.click(addButton);

    // Act
    fireEvent.click(addButton);

    // Assert
    expect(screen.queryByPlaceholderText("List name")).not.toBeInTheDocument();
  });

  it("hides the list name input when Escape is pressed", () => {
    // Arrange
    renderSidebar();
    fireEvent.click(screen.getByRole("button", { name: "Add list" }));

    // Act
    fireEvent.keyDown(screen.getByPlaceholderText("List name"), { key: "Escape" });

    // Assert
    expect(screen.queryByPlaceholderText("List name")).not.toBeInTheDocument();
  });
});

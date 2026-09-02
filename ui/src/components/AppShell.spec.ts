import { cleanup, fireEvent, render, screen, within } from "@testing-library/react";
import { createElement, type ReactNode } from "react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { afterEach, describe, expect, it, vi } from "vitest";
import { createQueryWrapper } from "../test/renderWithQueryClient";
import { AppShell } from "./AppShell";
import InboxPage from "../pages/InboxPage";

vi.mock("../contexts/useAppContext", () => ({
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

afterEach(() => {
  cleanup();
});

function renderAppShell(routes: ReactNode, initialEntry = "/inbox") {
  const { QueryWrapper } = createQueryWrapper();

  render(
    createElement(
      QueryWrapper,
      null,
      createElement(
        MemoryRouter,
        { initialEntries: [initialEntry] },
        createElement(
          Routes,
          null,
          createElement(Route, { element: createElement(AppShell) }, routes),
        ),
      ),
    ),
  );
}

function mobileDrawer() {
  return document.querySelector(".app-drawer");
}

describe("AppShell", () => {
  it("renders the shell layout and outlet content", () => {
    // Arrange
    renderAppShell(createElement(Route, { path: "inbox", element: createElement(InboxPage) }));

    // Assert
    expect(document.querySelector(".app-shell")).toBeInTheDocument();
    expect(document.querySelector(".app-sidebar")).toBeInTheDocument();
    expect(document.querySelector("main.app-main")).toBeInTheDocument();
    expect(screen.getByRole("heading", { name: "Inbox" })).toBeInTheDocument();
    expect(screen.getByText("No captured items yet.")).toBeInTheDocument();
    expect(screen.getAllByText("Listello")).toHaveLength(1);
  });

  it("does not show the mobile drawer initially", () => {
    // Arrange
    renderAppShell(createElement(Route, { path: "inbox", element: createElement(InboxPage) }));

    // Assert
    expect(mobileDrawer()).not.toBeInTheDocument();
    expect(document.querySelector(".app-overlay")).not.toBeInTheDocument();
  });

  it("opens the mobile drawer when the outlet calls openSidebar", () => {
    // Arrange
    renderAppShell(createElement(Route, { path: "inbox", element: createElement(InboxPage) }));

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Open menu" }));

    // Assert
    expect(mobileDrawer()).toBeInTheDocument();
    expect(document.querySelector(".app-overlay")).toBeInTheDocument();
    expect(screen.getAllByText("Listello")).toHaveLength(2);
  });

  it("closes the mobile drawer when the overlay is clicked", () => {
    // Arrange
    renderAppShell(createElement(Route, { path: "inbox", element: createElement(InboxPage) }));
    fireEvent.click(screen.getByRole("button", { name: "Open menu" }));

    // Act
    fireEvent.click(document.querySelector(".app-overlay")!);

    // Assert
    expect(mobileDrawer()).not.toBeInTheDocument();
    expect(document.querySelector(".app-overlay")).not.toBeInTheDocument();
  });

  it("closes the mobile drawer when a sidebar link is clicked", () => {
    // Arrange
    renderAppShell(createElement(Route, { path: "inbox", element: createElement(InboxPage) }));
    fireEvent.click(screen.getByRole("button", { name: "Open menu" }));
    const drawer = mobileDrawer() as HTMLElement;

    // Act
    fireEvent.click(within(drawer).getByRole("link", { name: "Work" }));

    // Assert
    expect(mobileDrawer()).not.toBeInTheDocument();
  });
});

import { cleanup, fireEvent, render, screen } from "@testing-library/react";
import { createElement } from "react";
import { afterEach, describe, expect, it } from "vitest";
import { AccountMenu } from "./AccountMenu";

afterEach(() => {
  cleanup();
});

function renderAccountMenu() {
  render(createElement(AccountMenu));
  return screen.getByRole("button", { name: "Account options" });
}

function dropdown() {
  return screen.getByRole("menu").closest(".dropdown");
}

describe("AccountMenu", () => {
  it("shows the signed-in user", () => {
    // Act
    render(createElement(AccountMenu));

    // Assert
    expect(screen.getByText("Signed in")).toBeInTheDocument();
    expect(screen.getByText("You · Owner")).toBeInTheDocument();
  });

  it("starts closed", () => {
    // Arrange
    const trigger = renderAccountMenu();

    // Assert
    expect(trigger).toHaveAttribute("aria-expanded", "false");
    expect(dropdown()).not.toHaveClass("is-active");
  });

  it("opens when the account options button is clicked", () => {
    // Arrange
    const trigger = renderAccountMenu();

    // Act
    fireEvent.click(trigger);

    // Assert
    expect(trigger).toHaveAttribute("aria-expanded", "true");
    expect(dropdown()).toHaveClass("is-active");
    expect(screen.getByRole("menuitem", { name: "Profile" })).toBeInTheDocument();
    expect(screen.getByRole("menuitem", { name: "Settings" })).toBeInTheDocument();
    expect(screen.getByRole("menuitem", { name: "Delegation policy" })).toBeInTheDocument();
    expect(screen.getByRole("menuitem", { name: "Sign out" })).toBeInTheDocument();
  });

  it("closes when a menu item is clicked", () => {
    // Arrange
    const trigger = renderAccountMenu();
    fireEvent.click(trigger);

    // Act
    fireEvent.click(screen.getByRole("menuitem", { name: "Settings" }));

    // Assert
    expect(trigger).toHaveAttribute("aria-expanded", "false");
    expect(dropdown()).not.toHaveClass("is-active");
  });

  it("closes when Escape is pressed", () => {
    // Arrange
    const trigger = renderAccountMenu();
    fireEvent.click(trigger);

    // Act
    fireEvent.keyDown(document, { key: "Escape" });

    // Assert
    expect(trigger).toHaveAttribute("aria-expanded", "false");
    expect(dropdown()).not.toHaveClass("is-active");
  });

  it("closes when clicking outside", () => {
    // Arrange
    const trigger = renderAccountMenu();
    fireEvent.click(trigger);

    // Act
    fireEvent.mouseDown(document.body);

    // Assert
    expect(trigger).toHaveAttribute("aria-expanded", "false");
    expect(dropdown()).not.toHaveClass("is-active");
  });
});

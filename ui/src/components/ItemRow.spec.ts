import { cleanup, fireEvent, render, screen } from "@testing-library/react";
import { createElement } from "react";
import type { ItemDto } from "api-types";
import { afterEach, describe, expect, it, vi } from "vitest";
import { ItemRow } from "./ItemRow";

const baseItem: ItemDto = {
  ID: "IT_1",
  ListID: "LS_1",
  ParentID: "",
  Title: "Idea: weekly review template",
  Description: "",
  DueDate: "",
  Tags: [],
  Priority: "",
  State: "outstanding",
};

afterEach(() => {
  cleanup();
});

describe("ItemRow", () => {
  it("renders an outstanding item without completed styling", () => {
    // Arrange
    render(createElement(ItemRow, { item: baseItem, onComplete: vi.fn(), onUncomplete: vi.fn() }));

    // Assert
    const toggle = screen.getByRole("button", { name: "Mark complete" });
    expect(toggle).not.toHaveClass("is-checked");
    expect(toggle.querySelector("svg")).not.toBeInTheDocument();

    const title = screen.getByText("Idea: weekly review template");
    expect(title).toHaveClass("is-size-6");
    expect(title).not.toHaveClass("muted");
    expect(title).not.toHaveClass("line-through");
  });

  it("renders a completed item with a checked toggle", () => {
    // Arrange
    render(
      createElement(ItemRow, {
        item: { ...baseItem, State: "complete" },
        onComplete: vi.fn(),
        onUncomplete: vi.fn(),
      }),
    );

    // Assert
    const toggle = screen.getByRole("button", { name: "Mark incomplete" });
    expect(toggle).toHaveClass("is-checked");
    expect(toggle.querySelector("svg")).toBeInTheDocument();
  });

  it("renders a completed item with a struck-through title", () => {
    // Arrange
    render(
      createElement(ItemRow, {
        item: { ...baseItem, State: "complete" },
        onComplete: vi.fn(),
        onUncomplete: vi.fn(),
      }),
    );

    // Assert
    const title = screen.getByText("Idea: weekly review template");
    expect(title).toHaveClass("is-size-6");
    expect(title).toHaveClass("muted");
    expect(title).toHaveClass("line-through");
  });

  it("calls onComplete when the checkbox is clicked", () => {
    // Arrange
    const onComplete = vi.fn();
    render(createElement(ItemRow, { item: baseItem, onComplete, onUncomplete: vi.fn() }));

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Mark complete" }));

    // Assert
    expect(onComplete).toHaveBeenCalledWith("IT_1");
  });

  it("calls onUncomplete when the checked checkbox is clicked", () => {
    // Arrange
    const onUncomplete = vi.fn();
    render(
      createElement(ItemRow, {
        item: { ...baseItem, State: "complete" },
        onComplete: vi.fn(),
        onUncomplete,
      }),
    );

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Mark incomplete" }));

    // Assert
    expect(onUncomplete).toHaveBeenCalledWith("IT_1");
  });

  it("immediately marks the checkbox as checked when clicked", () => {
    // Arrange
    const onComplete = vi.fn(() => new Promise<void>(() => {}));
    render(createElement(ItemRow, { item: baseItem, onComplete, onUncomplete: vi.fn() }));

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Mark complete" }));

    // Assert
    const toggle = screen.getByRole("button", { name: "Mark incomplete" });
    expect(toggle).toHaveClass("is-checked");
    expect(toggle.querySelector("svg")).toBeInTheDocument();
  });

  it("immediately strikes through the title when the checkbox is clicked", () => {
    // Arrange
    const onComplete = vi.fn(() => new Promise<void>(() => {}));
    render(createElement(ItemRow, { item: baseItem, onComplete, onUncomplete: vi.fn() }));

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Mark complete" }));

    // Assert
    const title = screen.getByText("Idea: weekly review template");
    expect(title).toHaveClass("muted");
    expect(title).toHaveClass("line-through");
  });

  it("immediately unchecks the checkbox when clicked after completing", () => {
    // Arrange
    const onComplete = vi.fn();
    const onUncomplete = vi.fn();
    render(createElement(ItemRow, { item: baseItem, onComplete, onUncomplete }));

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Mark complete" }));
    fireEvent.click(screen.getByRole("button", { name: "Mark incomplete" }));

    // Assert
    expect(onUncomplete).toHaveBeenCalledWith("IT_1");
    const toggle = screen.getByRole("button", { name: "Mark complete" });
    expect(toggle).not.toHaveClass("is-checked");
    expect(toggle.querySelector("svg")).not.toBeInTheDocument();
  });

  it("immediately removes the strikethrough when the checkbox is clicked after completing", () => {
    // Arrange
    const onComplete = vi.fn();
    const onUncomplete = vi.fn();
    render(createElement(ItemRow, { item: baseItem, onComplete, onUncomplete }));

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Mark complete" }));
    fireEvent.click(screen.getByRole("button", { name: "Mark incomplete" }));

    // Assert
    const title = screen.getByText("Idea: weekly review template");
    expect(title).not.toHaveClass("muted");
    expect(title).not.toHaveClass("line-through");
  });
});

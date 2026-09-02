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
    render(createElement(ItemRow, { item: baseItem, onComplete: vi.fn() }));

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
    render(createElement(ItemRow, { item: baseItem, onComplete }));

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Mark complete" }));

    // Assert
    expect(onComplete).toHaveBeenCalledWith("IT_1");
  });
});

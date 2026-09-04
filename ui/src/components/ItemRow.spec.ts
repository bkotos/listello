import { cleanup, fireEvent, render, screen } from "@testing-library/react";
import { createElement } from "react";
import type { ItemDto } from "api-types";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
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

  describe("hover actions", () => {
    beforeEach(() => {
      render(
        createElement(ItemRow, {
          item: { ...baseItem, Title: "Reply to the venue about the offsite" },
          onComplete: vi.fn(),
          onUncomplete: vi.fn(),
        }),
      );
    });

    it("renders the row as a hover-parent button", () => {
      // Assert
      const row = screen.getByText("Reply to the venue about the offsite").closest(".task-row");
      expect(row).toHaveClass("hover-parent", "is-flex", "p-3");
      expect(row).toHaveAttribute("role", "button");
      expect(row).toHaveAttribute("tabindex", "0");
    });

    it("renders a Set date icon button", () => {
      // Assert
      const setDate = screen.getByRole("button", { name: "Set date" });
      expect(setDate).toHaveClass("icon-btn");
      expect(setDate).toHaveAttribute("title", "Set date");
      expect(setDate.querySelector("svg.lucide-calendar")).toBeInTheDocument();
    });

    it("renders an Add comment icon button", () => {
      // Assert
      const addComment = screen.getByRole("button", { name: "Add comment" });
      expect(addComment).toHaveClass("icon-btn");
      expect(addComment).toHaveAttribute("title", "Add comment");
      expect(addComment.querySelector("svg.lucide-message-square")).toBeInTheDocument();
    });

    it("groups date and comment actions in a hover-reveal row", () => {
      // Assert
      const setDate = screen.getByRole("button", { name: "Set date" });
      const hoverActions = setDate.parentElement;
      expect(hoverActions).toHaveClass("is-inline-flex", "is-align-items-center", "hover-reveal");
      expect(hoverActions?.parentElement).toHaveClass(
        "is-flex",
        "is-flex-wrap-wrap",
        "is-align-items-center",
        "mt-2",
        "is-size-7",
        "muted",
      );
    });

    it("renders a closed Task options dropdown", () => {
      // Assert
      const options = screen.getByRole("button", { name: "Task options" });
      expect(options).toHaveClass("icon-btn", "hover-reveal");
      expect(options).toHaveAttribute("aria-haspopup", "menu");
      expect(options).toHaveAttribute("aria-expanded", "false");
      expect(options.querySelector("svg.lucide-ellipsis")).toBeInTheDocument();
      expect(options.parentElement).toHaveClass("dropdown-trigger");
      expect(options.closest(".dropdown")).toHaveClass("is-right");
    });
  });

  describe("when the Task options button is clicked", () => {
    beforeEach(() => {
      render(
        createElement(ItemRow, {
          item: { ...baseItem, Title: "Reply to the venue about the offsite" },
          onComplete: vi.fn(),
          onUncomplete: vi.fn(),
        }),
      );
      fireEvent.click(screen.getByRole("button", { name: "Task options" }));
    });

    it("opens the dropdown", () => {
      // Assert
      const trigger = screen.getByRole("button", { name: "Task options" });
      expect(trigger).toHaveAttribute("aria-expanded", "true");
      expect(trigger.closest(".dropdown")).toHaveClass("is-right", "is-active");
    });

    it("marks the trigger as active", () => {
      // Assert
      const trigger = screen.getByRole("button", { name: "Task options" });
      expect(trigger).toHaveClass("icon-btn", "is-active");
      expect(trigger).not.toHaveClass("hover-reveal");
    });

    it("renders a Rename menuitem", () => {
      // Assert
      const rename = screen.getByRole("menuitem", { name: "Rename" });
      expect(rename).toHaveClass("dropdown-item", "is-flex", "is-align-items-center");
      expect(rename.querySelector("svg.lucide-pencil")).toBeInTheDocument();
    });

    it("renders a Delete menuitem", () => {
      // Assert
      const remove = screen.getByRole("menuitem", { name: "Delete" });
      expect(remove).toHaveClass("dropdown-item", "is-flex", "is-align-items-center");
      expect(remove.querySelector("svg.lucide-trash-2")).toBeInTheDocument();
    });

    it("separates Rename and Delete with a divider", () => {
      // Assert
      const content = screen.getByRole("menu").querySelector(".dropdown-content");
      const children = [...(content?.children ?? [])];
      expect(children[0]).toHaveTextContent("Rename");
      expect(children[1]).toHaveClass("dropdown-divider");
      expect(children[2]).toHaveTextContent("Delete");
    });

    it("closes when clicking outside", () => {
      // Act
      fireEvent.mouseDown(document.body);

      // Assert
      const trigger = screen.getByRole("button", { name: "Task options" });
      expect(trigger).toHaveAttribute("aria-expanded", "false");
      expect(trigger.closest(".dropdown")).not.toHaveClass("is-active");
      expect(screen.queryByRole("menu")).not.toBeInTheDocument();
    });

    it("closes when the Task options button is clicked again", () => {
      // Act
      fireEvent.click(screen.getByRole("button", { name: "Task options" }));

      // Assert
      const trigger = screen.getByRole("button", { name: "Task options" });
      expect(trigger).toHaveAttribute("aria-expanded", "false");
      expect(trigger.closest(".dropdown")).not.toHaveClass("is-active");
      expect(screen.queryByRole("menu")).not.toBeInTheDocument();
    });
  });
});

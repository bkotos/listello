import { cleanup, fireEvent, render, screen } from "@testing-library/react";
import { createElement } from "react";
import type { ItemCommentDto, ItemDto } from "api-types";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { ItemDetail } from "./ItemDetail";

const baseItem: ItemDto = {
  ID: "IT_1",
  ListID: "LS_1",
  ParentID: "",
  Title: "Buy windshield wipers for truck",
  Description: "",
  DueDate: "",
  Tags: [],
  Priority: "",
  State: "outstanding",
};

afterEach(() => {
  cleanup();
});

describe("ItemDetail", () => {
  beforeEach(() => {
    render(createElement(ItemDetail, { item: baseItem, listName: "Work", onClose: vi.fn(), onModifyTitle: vi.fn(), onComment: vi.fn(), onAddSubtask: vi.fn(), onCompleteSubtask: vi.fn(), onUncompleteSubtask: vi.fn(), onDeleteSubtask: vi.fn() }));
  });

  it("renders a right-hand detail aside", () => {
    const panel = document.querySelector("aside.app-detail");
    expect(panel).toBeInTheDocument();
    expect(panel).toHaveClass("is-hidden-touch");
  });

  it("renders the list name as a tag", () => {
    const tag = screen.getByText("Work");
    expect(tag).toHaveClass("tag", "is-primary", "is-light", "is-family-code", "is-uppercase");
  });

  it("renders an unclarified status", () => {
    const status = screen.getByText("unclarified");
    expect(status).toHaveClass("is-family-code", "is-uppercase");
  });

  it("renders a Delete item button", () => {
    const remove = screen.getByRole("button", { name: "Delete item" });
    expect(remove).toHaveClass("icon-btn", "is-danger-hover");
    expect(remove.querySelector("svg.lucide-trash-2")).toBeInTheDocument();
  });

  it("renders a Close button", () => {
    const close = screen.getByRole("button", { name: "Close" });
    expect(close).toHaveClass("icon-btn");
    expect(close.querySelector("svg.lucide-x")).toBeInTheDocument();
  });

  it("renders a Mark complete toggle", () => {
    const toggle = screen.getByRole("button", { name: "Mark complete" });
    expect(toggle).toHaveClass("check-toggle");
  });

  it("renders the item title in a textarea", () => {
    const title = screen.getByDisplayValue("Buy windshield wipers for truck");
    expect(title.tagName).toBe("TEXTAREA");
    expect(title).toHaveClass("textarea", "is-shadowless");
  });

  it("prevents Enter from inserting a newline in the title", () => {
    const title = screen.getByDisplayValue("Buy windshield wipers for truck");
    const proceeded = fireEvent.keyDown(title, { key: "Enter" });
    expect(proceeded).toBe(false);
  });

  it("strips newlines from pasted title content", () => {
    const title = screen.getByDisplayValue("Buy windshield wipers for truck") as HTMLTextAreaElement;
    fireEvent.paste(title, {
      clipboardData: {
        getData: () => "Line one\nLine two",
      },
    });
    expect(title.value).toBe("Line oneLine two");
  });

  it("renders a Due date field", () => {
    expect(screen.getByText("Due date").querySelector("svg.lucide-calendar")).toBeInTheDocument();
    const dateInput = document.querySelector('input[type="date"]');
    expect(dateInput).toHaveClass("input", "is-small");
  });

  it("renders a No priority button", () => {
    const none = screen.getByRole("button", { name: "No priority" });
    expect(none).toHaveClass("button", "is-small", "is-primary", "is-light");
  });

  it("renders a Low priority button", () => {
    const low = screen.getByRole("button", { name: "Low" });
    expect(low).toHaveClass("button", "is-small", "is-white");
  });

  it("renders a Medium priority button", () => {
    const medium = screen.getByRole("button", { name: "Medium" });
    expect(medium).toHaveClass("button", "is-small", "is-white");
  });

  it("renders a High priority button", () => {
    const high = screen.getByRole("button", { name: "High" });
    expect(high).toHaveClass("button", "is-small", "is-white");
  });

  it("renders a List select", () => {
    expect(screen.getByText("List").querySelector("svg.lucide-chevron-right")).toBeInTheDocument();
    expect(screen.getByRole("combobox")).toBeInTheDocument();
  });

  it("renders an add tag button", () => {
    const add = screen.getByRole("button", { name: "add" });
    expect(add).toHaveClass("button", "is-small", "is-white");
    expect(add.querySelector("svg.lucide-plus")).toBeInTheDocument();
  });

  it("renders a Notes heading and textarea", () => {
    expect(screen.getByRole("heading", { name: "Notes" })).toHaveClass(
      "is-family-code",
      "is-uppercase",
      "muted",
    );
    expect(screen.getByPlaceholderText("Add notes, links, or context…")).toHaveClass("textarea");
  });

  it("renders a Subtasks heading and add input", () => {
    expect(screen.getByRole("heading", { name: "Subtasks" })).toHaveClass(
      "is-family-code",
      "is-uppercase",
      "muted",
    );
    expect(screen.getByPlaceholderText("Add a subtask…")).toHaveClass("input", "is-small", "is-shadowless");
  });

  it("renders a Comments heading and composer", () => {
    expect(screen.getByRole("heading", { name: "Comments" })).toHaveClass(
      "is-family-code",
      "is-uppercase",
      "muted",
    );
    expect(screen.getByPlaceholderText("Write a comment…")).toHaveClass("input", "is-small");
  });

  it("renders a disabled Send comment button", () => {
    const send = screen.getByRole("button", { name: "Send comment" });
    expect(send).toBeDisabled();
    expect(send).toHaveClass("button", "is-small", "is-primary");
    expect(send.querySelector("svg.lucide-send")).toBeInTheDocument();
  });
});

describe("ItemDetail subtasks", () => {
  const onCompleteSubtask = vi.fn();
  const onDeleteSubtask = vi.fn();

  beforeEach(() => {
    onCompleteSubtask.mockReset();
    onDeleteSubtask.mockReset();
    render(
      createElement(ItemDetail, {
        item: {
          ...baseItem,
          Subtasks: [{ ID: "ST_1", Title: "Get 2%", Completed: false }],
        },
        listName: "Work",
        onClose: vi.fn(),
        onModifyTitle: vi.fn(),
        onComment: vi.fn(),
        onAddSubtask: vi.fn(),
        onCompleteSubtask,
        onUncompleteSubtask: vi.fn(),
        onDeleteSubtask,
      }),
    );
  });

  it("renders a Mark subtask complete toggle", () => {
    const toggle = screen.getByRole("button", { name: "Mark subtask complete" });
    expect(toggle).toHaveClass("check-toggle");
  });

  it("renders the subtask title", () => {
    expect(screen.getByText("Get 2%")).toHaveClass("is-size-6");
  });

  it("renders a Delete subtask button", () => {
    const remove = screen.getByRole("button", { name: "Delete subtask" });
    expect(remove).toHaveClass("icon-btn", "is-danger-hover", "hover-reveal");
    expect(remove.querySelector("svg.lucide-trash-2")).toBeInTheDocument();
  });

  it("calls onCompleteSubtask when the toggle is clicked", () => {
    fireEvent.click(screen.getByRole("button", { name: "Mark subtask complete" }));
    expect(onCompleteSubtask).toHaveBeenCalledWith("IT_1", "ST_1");
  });

  it("calls onDeleteSubtask when Delete subtask is clicked", () => {
    fireEvent.click(screen.getByRole("button", { name: "Delete subtask" }));
    expect(onDeleteSubtask).toHaveBeenCalledWith("IT_1", "ST_1");
  });
});

describe("ItemDetail completed subtask", () => {
  const onUncompleteSubtask = vi.fn();

  beforeEach(() => {
    onUncompleteSubtask.mockReset();
    render(
      createElement(ItemDetail, {
        item: {
          ...baseItem,
          Subtasks: [{ ID: "ST_1", Title: "Get 2%", Completed: true }],
        },
        listName: "Work",
        onClose: vi.fn(),
        onModifyTitle: vi.fn(),
        onComment: vi.fn(),
        onAddSubtask: vi.fn(),
        onCompleteSubtask: vi.fn(),
        onUncompleteSubtask,
        onDeleteSubtask: vi.fn(),
      }),
    );
  });

  it("renders the completed subtask as checked", () => {
    const toggle = screen.getByRole("button", { name: "Mark subtask incomplete" });
    expect(toggle).toHaveClass("check-toggle", "is-checked");
    expect(toggle.querySelector("svg.lucide-check")).toBeInTheDocument();
  });

  it("renders the completed subtask title with strikethrough", () => {
    expect(screen.getByText("Get 2%")).toHaveClass("is-size-6", "muted", "line-through");
  });

  it("calls onUncompleteSubtask when the toggle is clicked", () => {
    fireEvent.click(screen.getByRole("button", { name: "Mark subtask incomplete" }));
    expect(onUncompleteSubtask).toHaveBeenCalledWith("IT_1", "ST_1");
  });
});

describe("ItemDetail comments", () => {
  const onDeleteComment = vi.fn();
  const comments: ItemCommentDto[] = [
    {
      ID: "CM_1",
      ItemID: "IT_1",
      UserID: "US_1",
      UserName: "Alex",
      Body: "Need 2%",
      CreatedAt: new Date().toISOString(),
    },
    {
      ID: "CM_2",
      ItemID: "IT_1",
      UserID: "US_2",
      UserName: "Sam",
      Body: "Get oat milk",
      CreatedAt: new Date().toISOString(),
    },
  ];

  beforeEach(() => {
    onDeleteComment.mockReset();
    render(
      createElement(ItemDetail, {
        item: baseItem,
        listName: "Work",
        comments,
        currentUserId: "US_1",
        onClose: vi.fn(),
        onModifyTitle: vi.fn(),
        onComment: vi.fn(),
        onAddSubtask: vi.fn(),
        onCompleteSubtask: vi.fn(),
        onUncompleteSubtask: vi.fn(),
        onDeleteSubtask: vi.fn(),
        onDeleteComment,
      }),
    );
  });

  it("renders the comment author name", () => {
    expect(screen.getByText("Alex")).toHaveClass("is-size-7", "has-text-weight-medium");
  });

  it("renders the comment body", () => {
    expect(screen.getByText("Need 2%")).toHaveClass("is-size-7", "muted");
  });

  it("renders the comment author initial", () => {
    expect(screen.getByText("A")).toHaveClass(
      "is-flex",
      "is-align-items-center",
      "is-justify-content-center",
      "is-family-code",
      "has-text-primary",
      "has-text-weight-semibold",
    );
  });

  it("renders a Delete comment button for the current user's comment", () => {
    const remove = screen.getByRole("button", { name: "Delete comment" });
    expect(remove).toHaveClass("icon-btn", "is-danger-hover", "hover-reveal");
    expect(remove.querySelector("svg.lucide-trash-2")).toBeInTheDocument();
  });

  it("does not render a Delete comment button for another user's comment", () => {
    expect(screen.getByText("Sam")).toBeInTheDocument();
    expect(screen.getAllByRole("button", { name: "Delete comment" })).toHaveLength(1);
  });

  it("calls onDeleteComment when Delete comment is clicked", () => {
    fireEvent.click(screen.getByRole("button", { name: "Delete comment" }));
    expect(onDeleteComment).toHaveBeenCalledWith("IT_1", "CM_1");
  });
});

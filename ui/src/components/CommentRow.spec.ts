import { cleanup, fireEvent, render, screen } from "@testing-library/react";
import { createElement } from "react";
import type { ItemCommentDto } from "api-types";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { CommentRow } from "./CommentRow";

const ownComment: ItemCommentDto = {
  ID: "CM_1",
  ItemID: "IT_1",
  UserID: "US_1",
  UserName: "Alex",
  Body: "Need 2%",
  CreatedAt: new Date().toISOString(),
};

const otherComment: ItemCommentDto = {
  ID: "CM_2",
  ItemID: "IT_1",
  UserID: "US_2",
  UserName: "Sam",
  Body: "Get oat milk",
  CreatedAt: new Date().toISOString(),
};

afterEach(() => {
  cleanup();
});

describe("CommentRow", () => {
  const onDelete = vi.fn();

  beforeEach(() => {
    onDelete.mockReset();
    render(
      createElement(CommentRow, {
        comment: ownComment,
        currentUserId: "US_1",
        onDelete,
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

  it("renders the comment relative time", () => {
    expect(screen.getByText("less than a minute ago")).toHaveClass("is-size-7", "muted");
  });

  it("renders a Delete comment button for the current user's comment", () => {
    const remove = screen.getByRole("button", { name: "Delete comment" });
    expect(remove).toHaveClass("icon-btn", "is-danger-hover", "hover-reveal");
    expect(remove.querySelector("svg.lucide-trash-2")).toBeInTheDocument();
  });

  it("calls onDelete when Delete comment is clicked", () => {
    fireEvent.click(screen.getByRole("button", { name: "Delete comment" }));
    expect(onDelete).toHaveBeenCalledWith("CM_1");
  });
});

describe("CommentRow for another user's comment", () => {
  beforeEach(() => {
    render(
      createElement(CommentRow, {
        comment: otherComment,
        currentUserId: "US_1",
        onDelete: vi.fn(),
      }),
    );
  });

  it("does not render a Delete comment button", () => {
    expect(screen.getByText("Sam")).toBeInTheDocument();
    expect(screen.queryByRole("button", { name: "Delete comment" })).not.toBeInTheDocument();
  });
});

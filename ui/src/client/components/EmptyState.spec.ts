import { render, screen } from "@testing-library/react";
import { createElement } from "react";
import { describe, expect, it } from "vitest";
import { EmptyState } from "./EmptyState.tsx";

describe("EmptyState", () => {
  it("renders the title and description", () => {
    render(
      createElement(EmptyState, {
        title: "No lists yet",
        description: "Create a list to get started.",
      }),
    );

    expect(screen.getByRole("heading", { name: "No lists yet" })).toBeInTheDocument();
    expect(screen.getByText("Create a list to get started.")).toBeInTheDocument();
  });
});

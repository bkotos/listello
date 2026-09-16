import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import { act, render, screen, cleanup } from "@testing-library/react";
import { createElement } from "react";
import { SystemSetupStep, SystemSetupFooter } from "./Step7-SystemSetup";

afterEach(() => {
  cleanup();
  vi.useRealTimers();
});

describe("SystemSetupStep", () => {
  const onComplete = vi.fn();

  beforeEach(() => {
    vi.useFakeTimers();
    onComplete.mockReset();
    render(
      createElement(SystemSetupStep, {
        spaceName: "Personal",
        userName: "Alex",
        onComplete,
      }),
    );
  });

  it("renders the Workspace · Automatic eyebrow", () => {
    // Assert
    const eyebrow = screen.getByText("Workspace · Automatic");
    expect(eyebrow).toHaveClass("step-eyebrow");
  });

  it("renders the Getting things ready heading", () => {
    // Assert
    const heading = screen.getByRole("heading", { name: "Getting things ready" });
    expect(heading).toHaveClass("step-title", "text-balance");
  });

  it("renders the setup lead with space name", () => {
    // Assert
    const lead = screen.getByText(/Listello is wiring up the essentials for/);
    expect(lead).toHaveClass("step-lead", "text-pretty");
    expect(lead).toHaveTextContent("Listello is wiring up the essentials for Personal.");
    const strong = lead.querySelector("strong");
    expect(strong).toHaveTextContent("Personal");
  });

  it("renders a Create Inbox setup check in progress", () => {
    // Assert
    const label = screen.getByText("Create Inbox");
    expect(label).toHaveClass("setup-check-label");
    const row = label.closest(".setup-check");
    expect(row).not.toHaveClass("is-pending");
    expect(row).not.toHaveClass("is-done");
    const loader = row?.querySelector("svg.lucide-loader-circle");
    expect(loader).toBeInTheDocument();
    expect(loader).toHaveClass("spin");
  });

  it("renders an Assign space to user setup check as pending", () => {
    // Assert
    const label = screen.getByText("Assign Personal to Alex");
    expect(label).toHaveClass("setup-check-label");
    const row = label.closest(".setup-check");
    expect(row).toHaveClass("is-pending");
    expect(row?.querySelector(".check-toggle")).toBeInTheDocument();
  });

  describe("after 500ms", () => {
    beforeEach(() => {
      act(() => {
        vi.advanceTimersByTime(500);
      });
    });

    it("renders a Create Inbox setup check as done", () => {
      // Assert
      const label = screen.getByText("Create Inbox");
      expect(label).toHaveClass("setup-check-label");
      const row = label.closest(".setup-check");
      expect(row).toHaveClass("is-done");
      expect(row?.querySelector("svg.lucide-circle-check")).toBeInTheDocument();
    });

    it("renders an Assign space to user setup check in progress", () => {
      // Assert
      const label = screen.getByText("Assign Personal to Alex");
      expect(label).toHaveClass("setup-check-label");
      const row = label.closest(".setup-check");
      expect(row).not.toHaveClass("is-pending");
      expect(row).not.toHaveClass("is-done");
      const loader = row?.querySelector("svg.lucide-loader-circle");
      expect(loader).toBeInTheDocument();
      expect(loader).toHaveClass("spin");
    });

    describe("after another 500ms", () => {
      beforeEach(() => {
        act(() => {
          vi.advanceTimersByTime(500);
        });
      });

      it("renders an Assign space to user setup check as done", () => {
        // Assert
        const label = screen.getByText("Assign Personal to Alex");
        expect(label).toHaveClass("setup-check-label");
        const row = label.closest(".setup-check");
        expect(row).toHaveClass("is-done");
        expect(row?.querySelector("svg.lucide-circle-check")).toBeInTheDocument();
      });

      it("calls onComplete", () => {
        // Assert
        expect(onComplete).toHaveBeenCalledOnce();
      });
    });
  });
});

describe("SystemSetupFooter", () => {
  beforeEach(() => {
    render(
      createElement(SystemSetupFooter, {
        continueEnabled: false,
        onContinue: () => {},
      })
    );
  });

  it("renders a disabled Continue button", () => {
    // Assert
    const button = screen.getByRole("button", { name: "Continue" });
    expect(button).toHaveClass("button", "is-primary", "footer-grow");
    expect(button).toBeDisabled();
    expect(button.querySelector("svg.lucide-arrow-right")).toBeInTheDocument();
  });

  describe("when continueEnabled is true", () => {
    beforeEach(() => {
      cleanup();
      render(
        createElement(SystemSetupFooter, {
          continueEnabled: true,
          onContinue: () => {},
        })
      );
    });

    it("enables the Continue button", () => {
      // Assert
      const button = screen.getByRole("button", { name: "Continue" });
      expect(button).not.toBeDisabled();
    });
  });
});

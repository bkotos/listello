import { describe, it, expect, beforeEach, afterEach } from "vitest";
import { render, screen, cleanup } from "@testing-library/react";
import { createElement } from "react";
import { SystemSetupStep, SystemSetupFooter } from "./Step7-SystemSetup";

afterEach(cleanup);

describe("SystemSetupStep", () => {
  beforeEach(() => {
    render(
      createElement(SystemSetupStep, {
        spaceName: "Personal",
        userName: "Alex",
        onComplete: () => {},
      })
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

  it("renders Create Inbox setup check", () => {
    // Assert
    const check = screen.getByText("Create Inbox");
    expect(check).toHaveClass("setup-check-label");
    const row = check.closest(".setup-check");
    expect(row).toBeInTheDocument();
  });

  it("renders Assign space to user setup check", () => {
    // Assert
    const check = screen.getByText("Assign Personal to Alex");
    expect(check).toHaveClass("setup-check-label");
    const row = check.closest(".setup-check");
    expect(row).toBeInTheDocument();
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

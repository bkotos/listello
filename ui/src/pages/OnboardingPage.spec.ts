import { cleanup, render, screen } from "@testing-library/react";
import { createElement } from "react";
import { afterEach, beforeEach, describe, expect, it } from "vitest";
import OnboardingPage from "./OnboardingPage";

afterEach(() => {
  cleanup();
});

describe("OnboardingPage", () => {
  beforeEach(() => {
    render(createElement(OnboardingPage));
  });

  it("renders the Welcome to Listello heading", () => {
    // Assert
    const heading = screen.getByRole("heading", { name: "Welcome to Listello" });
    expect(heading).toHaveClass("step-title", "text-balance");
    expect(heading.closest(".onboarding-page")).toBeInTheDocument();
  });

  it("renders the Listello brand mark with a check icon", () => {
    // Assert
    const brand = screen.getByText("Listello");
    expect(brand).toHaveClass("brand-mark");
    expect(brand.querySelector(".brand-dot")).toBeInTheDocument();
    expect(brand.querySelector("svg.lucide-check")).toBeInTheDocument();
    expect(brand.closest(".onboarding-topbar")).toBeInTheDocument();
  });

  it("renders a First-time setup badge", () => {
    // Assert
    const badge = screen.getByText("First-time setup");
    expect(badge).toHaveClass("setup-badge");
  });

  it("hides phase progress from assistive tech", () => {
    // Assert
    const progress = screen.getByText("Instance", { hidden: true }).closest(".phase-progress");
    expect(progress).toHaveAttribute("aria-hidden", "true");
  });

  it("renders the Instance phase as active", () => {
    // Assert
    const label = screen.getByText("Instance", { hidden: true });
    expect(label).toHaveClass("phase-seg-label");
    const seg = label.closest(".phase-seg");
    expect(seg).toHaveClass("is-active");
    expect(seg?.querySelector(".phase-seg-index")).toHaveTextContent("1");
    expect(seg?.querySelector(".phase-seg-fill")).toHaveStyle({ width: "25%" });
  });

  it("renders the Workspace phase as upcoming", () => {
    // Assert
    const label = screen.getByText("Workspace", { hidden: true });
    expect(label).toHaveClass("phase-seg-label");
    const seg = label.closest(".phase-seg");
    expect(seg).toHaveClass("is-upcoming");
    expect(seg?.querySelector(".phase-seg-index")).toHaveTextContent("2");
    expect(seg?.querySelector(".phase-seg-fill")).toHaveStyle({ width: "0%" });
  });

  it("renders the Ready phase as upcoming", () => {
    // Assert
    const label = screen.getByText("Ready", { hidden: true });
    expect(label).toHaveClass("phase-seg-label");
    const seg = label.closest(".phase-seg");
    expect(seg).toHaveClass("is-upcoming");
    expect(seg?.querySelector(".phase-seg-index")).toHaveTextContent("3");
    expect(seg?.querySelector(".phase-seg-fill")).toHaveStyle({ width: "0%" });
  });

  it("renders a sparkles icon", () => {
    // Assert
    const heading = screen.getByRole("heading", { name: "Welcome to Listello" });
    const icon = heading.parentElement?.querySelector(".big-icon");
    expect(icon).toBeInTheDocument();
    expect(icon?.querySelector("svg.lucide-sparkles")).toBeInTheDocument();
  });

  it("renders the welcome lead", () => {
    // Assert
    const lead = screen.getByText(/Let's create your instance and set up a calm, GTD-style workspace/);
    expect(lead).toHaveClass("step-lead", "text-pretty");
    expect(lead).toHaveTextContent(
      "Let's create your instance and set up a calm, GTD-style workspace. It only takes a minute, and you can change everything later.",
    );
  });

  it("renders a Choose how Listello is hosted setup check", () => {
    // Assert
    const label = screen.getByText("Choose how Listello is hosted");
    expect(label).toHaveClass("setup-check-label");
    const row = label.closest(".setup-check");
    expect(row).toHaveClass("is-done");
    expect(row?.querySelector("svg.lucide-server")).toBeInTheDocument();
  });

  it("renders a Create your space and profile setup check", () => {
    // Assert
    const label = screen.getByText("Create your space and profile");
    expect(label).toHaveClass("setup-check-label");
    const row = label.closest(".setup-check");
    expect(row).toHaveClass("is-done");
    expect(row?.querySelector("svg.lucide-layers")).toBeInTheDocument();
  });

  it("renders a Start your first list setup check", () => {
    // Assert
    const label = screen.getByText("Start your first list");
    expect(label).toHaveClass("setup-check-label");
    const row = label.closest(".setup-check");
    expect(row).toHaveClass("is-done");
    expect(row?.querySelector("svg.lucide-list-checks")).toBeInTheDocument();
  });

  it("renders a Create instance button", () => {
    // Assert
    const button = screen.getByRole("button", { name: "Create instance" });
    expect(button).toHaveClass("button", "is-primary", "footer-grow");
    expect(button.querySelector("svg.lucide-arrow-right")).toBeInTheDocument();
    expect(button.closest(".onboarding-footer")).toBeInTheDocument();
  });
});

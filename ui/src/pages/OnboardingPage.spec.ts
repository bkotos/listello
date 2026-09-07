import {
  cleanup,
  fireEvent,
  render,
  screen,
  waitFor,
} from "@testing-library/react";
import { createElement } from "react";
import type { ListelloInstanceResponse } from "api-types/listello-instance";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import OnboardingPage from "./OnboardingPage";

vi.mock("../lib/api/instance-client", () => ({
  createInstance: vi.fn(),
  getInstance: vi.fn(),
}));

import { createInstance } from "../lib/api/instance-client";

const createdInstance: ListelloInstanceResponse = {
  HostingMode: "",
  PersistenceLocation: "",
  PersistenceState: "",
  SetupState: "",
};

afterEach(() => {
  cleanup();
  vi.clearAllMocks();
});

describe("OnboardingPage", () => {
  beforeEach(() => {
    render(createElement(OnboardingPage));
  });

  it("renders the Welcome to Listello heading", () => {
    // Assert
    const heading = screen.getByRole("heading", {
      name: "Welcome to Listello",
    });
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
    const progress = document.querySelector(".phase-progress");
    expect(progress).toHaveAttribute("aria-hidden", "true");
  });

  it("renders the Instance phase as active", () => {
    // Assert
    const label = document.querySelector(
      ".phase-seg.is-active .phase-seg-label",
    );
    expect(label).toHaveClass("phase-seg-label");
    expect(label).toHaveTextContent("Instance");
    const seg = label?.closest(".phase-seg");
    expect(seg).toHaveClass("is-active");
    expect(seg?.querySelector(".phase-seg-index")).toHaveTextContent("1");
    expect(seg?.querySelector(".phase-seg-fill")).toHaveStyle({ width: "25%" });
  });

  it("renders the Workspace phase as upcoming", () => {
    // Assert
    const label = [...document.querySelectorAll(".phase-seg-label")].find(
      (el) => el.textContent === "Workspace",
    );
    expect(label).toHaveClass("phase-seg-label");
    const seg = label?.closest(".phase-seg");
    expect(seg).toHaveClass("is-upcoming");
    expect(seg?.querySelector(".phase-seg-index")).toHaveTextContent("2");
    expect(seg?.querySelector(".phase-seg-fill")).toHaveStyle({ width: "0%" });
  });

  it("renders the Ready phase as upcoming", () => {
    // Assert
    const label = [...document.querySelectorAll(".phase-seg-label")].find(
      (el) => el.textContent === "Ready",
    );
    expect(label).toHaveClass("phase-seg-label");
    const seg = label?.closest(".phase-seg");
    expect(seg).toHaveClass("is-upcoming");
    expect(seg?.querySelector(".phase-seg-index")).toHaveTextContent("3");
    expect(seg?.querySelector(".phase-seg-fill")).toHaveStyle({ width: "0%" });
  });

  it("renders a sparkles icon", () => {
    // Assert
    const heading = screen.getByRole("heading", {
      name: "Welcome to Listello",
    });
    const icon = heading.parentElement?.querySelector(".big-icon");
    expect(icon).toBeInTheDocument();
    expect(icon?.querySelector("svg.lucide-sparkles")).toBeInTheDocument();
  });

  it("renders the welcome lead", () => {
    // Assert
    const lead = screen.getByText(
      /Let's create your instance and set up a calm, GTD-style workspace/,
    );
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

  it("calls createInstance when Create instance is clicked", () => {
    // Arrange
    vi.mocked(createInstance).mockResolvedValue(createdInstance);

    // Act
    fireEvent.click(screen.getByRole("button", { name: "Create instance" }));

    // Assert
    expect(createInstance).toHaveBeenCalledOnce();
  });

  describe("when the instance is created", () => {
    beforeEach(async () => {
      vi.mocked(createInstance).mockResolvedValue(createdInstance);
      fireEvent.click(screen.getByRole("button", { name: "Create instance" }));
      await waitFor(() => {
        expect(
          screen.getByRole("heading", {
            name: "How should Listello be hosted?",
          }),
        ).toBeInTheDocument();
      });
    });

    it("renders the How should Listello be hosted heading", () => {
      // Assert
      const heading = screen.getByRole("heading", {
        name: "How should Listello be hosted?",
      });
      expect(heading).toHaveClass("step-title", "text-balance");
    });

    it("renders the Instance · Hosting eyebrow", () => {
      // Assert
      const eyebrow = screen.getByText("Instance · Hosting");
      expect(eyebrow).toHaveClass("step-eyebrow");
    });

    it("renders the hosting lead", () => {
      // Assert
      const lead = screen.getByText(
        /Your hosting mode decides where this instance runs/,
      );
      expect(lead).toHaveClass("step-lead", "text-pretty");
      expect(lead).toHaveTextContent(
        "Your hosting mode decides where this instance runs and how your data is stored. You can migrate later.",
      );
    });

    it("renders the Instance phase fill at 50%", () => {
      // Assert
      const label = document.querySelector(
        ".phase-seg.is-active .phase-seg-label",
      );
      expect(label).toHaveTextContent("Instance");
      expect(
        label?.closest(".phase-seg")?.querySelector(".phase-seg-fill"),
      ).toHaveStyle({
        width: "50%",
      });
    });

    it("renders a selected Local choice card", () => {
      // Assert
      const card = screen.getByRole("button", { name: /Local/ });
      expect(card).toHaveClass("choice-card", "is-selected");
      expect(card).toHaveAttribute("aria-pressed", "true");
      expect(card.querySelector("svg.lucide-hard-drive")).toBeInTheDocument();
      expect(card.querySelector(".choice-title")).toHaveTextContent("Local");
      expect(card.querySelector(".choice-desc")).toHaveTextContent(
        "Runs directly on this machine, with local filesystem-backed persistence.",
      );
      expect(card.querySelector(".choice-meta")).toHaveTextContent(
        "Uses SQLite",
      );
      expect(card.querySelector("svg.lucide-database")).toBeInTheDocument();
      expect(
        card.querySelector(".choice-check svg.lucide-circle-check"),
      ).toBeInTheDocument();
    });

    it("renders a Standalone Web choice card", () => {
      // Assert
      const card = screen.getByRole("button", { name: /Standalone Web/ });
      expect(card).toHaveClass("choice-card");
      expect(card).not.toHaveClass("is-selected");
      expect(card).toHaveAttribute("aria-pressed", "false");
      expect(card.querySelector("svg.lucide-globe")).toBeInTheDocument();
      expect(card.querySelector(".choice-title")).toHaveTextContent(
        "Standalone Web",
      );
      expect(card.querySelector(".choice-desc")).toHaveTextContent(
        "Runs independently in the browser as a PWA or standalone app.",
      );
      expect(card.querySelector(".choice-meta")).toHaveTextContent(
        "Uses IndexedDB / OPFS",
      );
      expect(card.querySelector("svg.lucide-database")).toBeInTheDocument();
      expect(card.querySelector(".choice-check")).not.toBeInTheDocument();
    });

    it("renders a Back button", () => {
      // Assert
      const back = screen.getByRole("button", { name: "Back" });
      expect(back).toHaveClass("button", "is-light");
      expect(back.querySelector("svg.lucide-arrow-left")).toBeInTheDocument();
      expect(back.closest(".onboarding-footer")).toBeInTheDocument();
    });

    it("renders a Continue button", () => {
      // Assert
      const button = screen.getByRole("button", { name: "Continue" });
      expect(button).toHaveClass("button", "is-primary", "footer-grow");
      expect(
        button.querySelector("svg.lucide-arrow-right"),
      ).toBeInTheDocument();
      expect(button.closest(".onboarding-footer")).toBeInTheDocument();
    });

    describe("when the Standalone Web choice card is clicked", () => {
      beforeEach(() => {
        fireEvent.click(screen.getByRole("button", { name: /Standalone Web/ }));
      });

      it("selects the Standalone Web choice card", () => {
        // Assert
        const card = screen.getByRole("button", { name: /Standalone Web/ });
        expect(card).toHaveClass("choice-card", "is-selected");
        expect(card).toHaveAttribute("aria-pressed", "true");
        expect(
          card.querySelector(".choice-check svg.lucide-circle-check"),
        ).toBeInTheDocument();
      });

      it("deselects the Local choice card", () => {
        // Assert
        const card = screen.getByRole("button", { name: /Local/ });
        expect(card).toHaveClass("choice-card");
        expect(card).not.toHaveClass("is-selected");
        expect(card).toHaveAttribute("aria-pressed", "false");
        expect(card.querySelector(".choice-check")).not.toBeInTheDocument();
      });

      it("renders the Instance phase fill at 67%", () => {
        // Assert
        const label = document.querySelector(
          ".phase-seg.is-active .phase-seg-label",
        );
        expect(label).toHaveTextContent("Instance");
        expect(
          label?.closest(".phase-seg")?.querySelector(".phase-seg-fill"),
        ).toHaveStyle({
          width: "67%",
        });
      });
    });
  });
});

import {
  act,
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
  selectHostingMode: vi.fn(),
}));

import { createInstance, selectHostingMode } from "../lib/api/instance-client";

const createdInstance: ListelloInstanceResponse = {
  HostingMode: "",
  PersistenceLocation: "",
  PersistenceState: "",
  SetupState: "",
};

function itRendersInitializeStepBeforeTimeouts() {
  it("renders the Setting up persistence heading", () => {
    // Assert
    const heading = screen.getByRole("heading", {
      name: "Setting up persistence",
    });
    expect(heading).toHaveClass("step-title", "text-balance");
  });

  it("renders the Instance · Initialize eyebrow", () => {
    // Assert
    const eyebrow = screen.getByText("Instance · Initialize");
    expect(eyebrow).toHaveClass("step-eyebrow");
  });

  it("renders the initialize lead", () => {
    // Assert
    const lead = screen.getByText(/We're initializing the/);
    expect(lead).toHaveClass("step-lead", "text-pretty");
    expect(lead).toHaveTextContent(
      "We're initializing the SQLite store for your local instance.",
    );
    expect(lead.querySelector("strong")).toHaveTextContent("SQLite");
    expect(lead.querySelectorAll("strong")[1]).toHaveTextContent("local");
  });

  it("renders the Instance phase fill at 100%", () => {
    // Assert
    const label = document.querySelector(
      ".phase-seg.is-active .phase-seg-label",
    );
    expect(label).toHaveTextContent("Instance");
    expect(
      label?.closest(".phase-seg")?.querySelector(".phase-seg-fill"),
    ).toHaveStyle({
      width: "100%",
    });
  });

  it("renders a Create instance setup check in progress", () => {
    // Assert
    const label = screen.getByText("Create instance");
    expect(label).toHaveClass("setup-check-label");
    const row = label.closest(".setup-check");
    expect(row).not.toHaveClass("is-pending");
    expect(row).not.toHaveClass("is-done");
    const loader = row?.querySelector("svg.lucide-loader-circle");
    expect(loader).toBeInTheDocument();
    expect(loader).toHaveClass("spin");
  });

  it("renders a Prepare SQLite store setup check as pending", () => {
    // Assert
    const label = screen.getByText("Prepare SQLite store");
    expect(label).toHaveClass("setup-check-label");
    const row = label.closest(".setup-check");
    expect(row).toHaveClass("is-pending");
    expect(row?.querySelector(".check-toggle")).toBeInTheDocument();
  });

  it("renders an Initialize persistence setup check as pending", () => {
    // Assert
    const label = screen.getByText("Initialize persistence");
    expect(label).toHaveClass("setup-check-label");
    const row = label.closest(".setup-check");
    expect(row).toHaveClass("is-pending");
    expect(row?.querySelector(".check-toggle")).toBeInTheDocument();
  });

  it("renders a Back button", () => {
    // Assert
    const back = screen.getByRole("button", { name: "Back" });
    expect(back).toHaveClass("button", "is-light");
    expect(back.querySelector("svg.lucide-arrow-left")).toBeInTheDocument();
    expect(back.closest(".onboarding-footer")).toBeInTheDocument();
  });

  it("renders a disabled Continue button", () => {
    // Assert
    const button = screen.getByRole("button", { name: "Continue" });
    expect(button).toHaveClass("button", "is-primary", "footer-grow");
    expect(button).toBeDisabled();
    expect(button.querySelector("svg.lucide-arrow-right")).toBeInTheDocument();
    expect(button.closest(".onboarding-footer")).toBeInTheDocument();
  });
}

afterEach(() => {
  cleanup();
  vi.clearAllMocks();
  vi.useRealTimers();
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

    it("calls selectHostingMode with local when Continue is clicked", () => {
      // Arrange
      vi.mocked(selectHostingMode).mockResolvedValue(createdInstance);

      // Act
      fireEvent.click(screen.getByRole("button", { name: "Continue" }));

      // Assert
      expect(selectHostingMode).toHaveBeenCalledWith({
        mode: "local",
      });
    });

    describe("when the Back button is clicked", () => {
      beforeEach(async () => {
        fireEvent.click(screen.getByRole("button", { name: "Back" }));
        await waitFor(() => {
          expect(
            screen.getByRole("heading", { name: "Welcome to Listello" }),
          ).toBeInTheDocument();
        });
      });

      it("renders the Welcome to Listello heading", () => {
        // Assert
        const heading = screen.getByRole("heading", {
          name: "Welcome to Listello",
        });
        expect(heading).toHaveClass("step-title", "text-balance");
      });
    });

    describe("when the Continue button is clicked", () => {
      beforeEach(async () => {
        fireEvent.click(screen.getByRole("button", { name: "Continue" }));
        await waitFor(() => {
          expect(
            screen.getByRole("heading", { name: "Choose a data directory" }),
          ).toBeInTheDocument();
        });
      });

      it("renders the Choose a data directory heading", () => {
        // Assert
        const heading = screen.getByRole("heading", {
          name: "Choose a data directory",
        });
        expect(heading).toHaveClass("step-title", "text-balance");
      });

      it("renders the Instance · Local eyebrow", () => {
        // Assert
        const eyebrow = screen.getByText("Instance · Local");
        expect(eyebrow).toHaveClass("step-eyebrow");
      });

      it("renders the data directory lead", () => {
        // Assert
        const lead = screen.getByText(
          /This is where Listello keeps its data for this instance/,
        );
        expect(lead).toHaveClass("step-lead", "text-pretty");
        expect(lead).toHaveTextContent(
          "This is where Listello keeps its data for this instance. You can move it later.",
        );
      });

      it("renders the Instance phase fill at 75%", () => {
        // Assert
        const label = document.querySelector(
          ".phase-seg.is-active .phase-seg-label",
        );
        expect(label).toHaveTextContent("Instance");
        expect(
          label?.closest(".phase-seg")?.querySelector(".phase-seg-fill"),
        ).toHaveStyle({
          width: "75%",
        });
      });

      it("renders a Data directory field", () => {
        // Assert
        const input = screen.getByLabelText("Data directory");
        expect(input).toHaveClass("input");
        expect(input).toHaveAttribute("id", "onb-location");
        expect(input).toHaveAttribute("placeholder", "~/listello");
        expect(input).toHaveValue(
          "/Users/jdoe/Library/Application Support/listello",
        );
        expect(screen.getByText("Data directory")).toHaveClass("label");
      });

      it("renders a folder icon in the Data directory field", () => {
        // Assert
        const control = screen.getByLabelText("Data directory").closest(".control");
        expect(control).toHaveClass("has-icons-left");
        expect(control?.querySelector("svg.lucide-folder")).toBeInTheDocument();
      });

      it("renders the data directory help", () => {
        // Assert
        const help = screen.getByText(
          "Logs, the SQLite database, and config are stored here.",
        );
        expect(help).toHaveClass("help");
      });

      describe("when the Back button is clicked", () => {
        beforeEach(async () => {
          fireEvent.click(screen.getByRole("button", { name: "Back" }));
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
      });

      describe("when the Continue button is clicked", () => {
        beforeEach(async () => {
          fireEvent.click(screen.getByRole("button", { name: "Continue" }));
          await waitFor(() => {
            expect(
              screen.getByRole("heading", { name: "Setting up persistence" }),
            ).toBeInTheDocument();
          });
        });

        itRendersInitializeStepBeforeTimeouts();
      });

      describe("after 500ms", () => {
        beforeEach(() => {
          vi.useFakeTimers();
          fireEvent.click(screen.getByRole("button", { name: "Continue" }));
          act(() => {
            vi.advanceTimersByTime(500);
          });
        });

        it("renders a Create instance setup check as done", () => {
          // Assert
          const label = screen.getByText("Create instance");
          expect(label).toHaveClass("setup-check-label");
          const row = label.closest(".setup-check");
          expect(row).toHaveClass("is-done");
          expect(row?.querySelector("svg.lucide-circle-check")).toBeInTheDocument();
        });

        it("renders a Prepare SQLite store setup check in progress", () => {
          // Assert
          const label = screen.getByText("Prepare SQLite store");
          expect(label).toHaveClass("setup-check-label");
          const row = label.closest(".setup-check");
          expect(row).not.toHaveClass("is-pending");
          expect(row).not.toHaveClass("is-done");
          const loader = row?.querySelector("svg.lucide-loader-circle");
          expect(loader).toBeInTheDocument();
          expect(loader).toHaveClass("spin");
        });

        it("renders an Initialize persistence setup check as pending", () => {
          // Assert
          const label = screen.getByText("Initialize persistence");
          expect(label).toHaveClass("setup-check-label");
          const row = label.closest(".setup-check");
          expect(row).toHaveClass("is-pending");
          expect(row?.querySelector(".check-toggle")).toBeInTheDocument();
        });

        describe("after another 500ms", () => {
          beforeEach(() => {
            act(() => {
              vi.advanceTimersByTime(500);
            });
          });

          it("renders a Prepare SQLite store setup check as done", () => {
            // Assert
            const label = screen.getByText("Prepare SQLite store");
            expect(label).toHaveClass("setup-check-label");
            const row = label.closest(".setup-check");
            expect(row).toHaveClass("is-done");
            expect(
              row?.querySelector("svg.lucide-circle-check"),
            ).toBeInTheDocument();
          });

          it("renders an Initialize persistence setup check in progress", () => {
            // Assert
            const label = screen.getByText("Initialize persistence");
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

            it("renders an Initialize persistence setup check as done", () => {
              // Assert
              const label = screen.getByText("Initialize persistence");
              expect(label).toHaveClass("setup-check-label");
              const row = label.closest(".setup-check");
              expect(row).toHaveClass("is-done");
              expect(
                row?.querySelector("svg.lucide-circle-check"),
              ).toBeInTheDocument();
            });

            it("enables the Continue button", () => {
              // Assert
              const button = screen.getByRole("button", { name: "Continue" });
              expect(button).toHaveClass("button", "is-primary", "footer-grow");
              expect(button).toBeEnabled();
              expect(
                button.querySelector("svg.lucide-arrow-right"),
              ).toBeInTheDocument();
              expect(button.closest(".onboarding-footer")).toBeInTheDocument();
            });

            describe("when the Back button is clicked", () => {
              beforeEach(() => {
                fireEvent.click(screen.getByRole("button", { name: "Back" }));
              });

              describe("when the Continue button is clicked", () => {
                beforeEach(() => {
                  fireEvent.click(screen.getByRole("button", { name: "Continue" }));
                });

                itRendersInitializeStepBeforeTimeouts();
              });
            });
          });
        });
      });
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

      describe("when the Local choice card is clicked", () => {
        beforeEach(() => {
          fireEvent.click(screen.getByRole("button", { name: /Local/ }));
        });

        it("selects the Local choice card", () => {
          // Assert
          const card = screen.getByRole("button", { name: /Local/ });
          expect(card).toHaveClass("choice-card", "is-selected");
          expect(card).toHaveAttribute("aria-pressed", "true");
          expect(
            card.querySelector(".choice-check svg.lucide-circle-check"),
          ).toBeInTheDocument();
        });

        it("deselects the Standalone Web choice card", () => {
          // Assert
          const card = screen.getByRole("button", { name: /Standalone Web/ });
          expect(card).toHaveClass("choice-card");
          expect(card).not.toHaveClass("is-selected");
          expect(card).toHaveAttribute("aria-pressed", "false");
          expect(card.querySelector(".choice-check")).not.toBeInTheDocument();
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
      });

      describe("when the Continue button is clicked", () => {
        beforeEach(() => {
          fireEvent.click(screen.getByRole("button", { name: "Continue" }));
        });

        // TODO this is temporary until we implement this flow for standalone web
        it("stays on the How should Listello be hosted heading", () => {
          // Assert
          expect(
            screen.getByRole("heading", {
              name: "How should Listello be hosted?",
            }),
          ).toBeInTheDocument();
        });
      });
    });
  });
});

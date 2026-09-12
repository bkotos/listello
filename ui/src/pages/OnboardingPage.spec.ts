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
import { createQueryWrapper } from "../test/renderWithQueryClient";
import OnboardingPage from "./OnboardingPage";

vi.mock("../lib/api/instance-client", () => ({
  createInstance: vi.fn(),
  getInstance: vi.fn(),
  getDefaultPersistenceLocation: vi.fn(),
  selectHostingMode: vi.fn(),
  selectPersistenceLocation: vi.fn(),
  initializePersistence: vi.fn(),
}));

vi.mock("../lib/api/user-client", () => ({
  createUser: vi.fn(),
}));

vi.mock("../lib/api/space-client", () => ({
  createSpace: vi.fn(),
}));

import {
  createInstance,
  getDefaultPersistenceLocation,
  initializePersistence,
  selectHostingMode,
  selectPersistenceLocation,
} from "../lib/api/instance-client";
import { createUser } from "../lib/api/user-client";
import { createSpace } from "../lib/api/space-client";

const createdInstance: ListelloInstanceResponse = {
  HostingMode: "",
  PersistenceLocation: "",
  PersistenceState: "",
  SetupState: "",
};

const apiDefaultPersistenceLocation =
  "/Users/api/Library/Application Support/listello";

function renderOnboardingPage() {
  const { QueryWrapper } = createQueryWrapper();
  return render(createElement(QueryWrapper, null, createElement(OnboardingPage)));
}

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
    vi.mocked(getDefaultPersistenceLocation).mockResolvedValue({
      Location: apiDefaultPersistenceLocation,
    });
    renderOnboardingPage();
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

      it("renders a Data directory field", async () => {
        // Assert
        const input = screen.getByLabelText("Data directory");
        expect(input).toHaveClass("input");
        expect(input).toHaveAttribute("id", "onb-location");
        expect(input).toHaveAttribute("placeholder", "~/listello");
        await waitFor(() => {
          expect(input).toHaveValue(apiDefaultPersistenceLocation);
        });
        expect(getDefaultPersistenceLocation).toHaveBeenCalledWith(
          expect.objectContaining({ signal: expect.any(AbortSignal) }),
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

      it("calls selectPersistenceLocation when Continue is clicked", async () => {
        // Arrange
        vi.mocked(selectPersistenceLocation).mockResolvedValue(createdInstance);
        await waitFor(() => {
          expect(screen.getByLabelText("Data directory")).toHaveValue(
            apiDefaultPersistenceLocation,
          );
        });

        // Act
        fireEvent.click(screen.getByRole("button", { name: "Continue" }));

        // Assert
        expect(selectPersistenceLocation).toHaveBeenCalledWith({
          location: apiDefaultPersistenceLocation,
        });
      });

      it("calls initializePersistence when Continue is clicked", async () => {
        // Arrange
        vi.mocked(selectPersistenceLocation).mockResolvedValue(createdInstance);
        vi.mocked(initializePersistence).mockResolvedValue(createdInstance);
        await waitFor(() => {
          expect(screen.getByLabelText("Data directory")).toHaveValue(
            apiDefaultPersistenceLocation,
          );
        });

        // Act
        fireEvent.click(screen.getByRole("button", { name: "Continue" }));

        // Assert
        await waitFor(() => {
          expect(initializePersistence).toHaveBeenCalledOnce();
        });
        expect(initializePersistence).toHaveBeenCalledWith();
      });

      it("calls initializePersistence after selectPersistenceLocation", async () => {
        // Arrange
        vi.mocked(selectPersistenceLocation).mockResolvedValue(createdInstance);
        vi.mocked(initializePersistence).mockResolvedValue(createdInstance);
        await waitFor(() => {
          expect(screen.getByLabelText("Data directory")).toHaveValue(
            apiDefaultPersistenceLocation,
          );
        });

        // Act
        fireEvent.click(screen.getByRole("button", { name: "Continue" }));

        // Assert
        await waitFor(() => {
          expect(initializePersistence).toHaveBeenCalledOnce();
        });
        expect(vi.mocked(selectPersistenceLocation).mock.invocationCallOrder[0]).toBeLessThan(
          vi.mocked(initializePersistence).mock.invocationCallOrder[0],
        );
      });

      describe("when the Data directory field is changed", () => {
        beforeEach(() => {
          fireEvent.change(screen.getByLabelText("Data directory"), {
            target: { value: "/tmp/listello" },
          });
        });

        it("calls selectPersistenceLocation with the updated location when Continue is clicked", () => {
          // Arrange
          vi.mocked(selectPersistenceLocation).mockResolvedValue(createdInstance);

          // Act
          fireEvent.click(screen.getByRole("button", { name: "Continue" }));

          // Assert
          expect(selectPersistenceLocation).toHaveBeenCalledWith({
            location: "/tmp/listello",
          });
        });
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

            describe("when the Continue button is clicked", () => {
              beforeEach(() => {
                fireEvent.click(screen.getByRole("button", { name: "Continue" }));
              });

              it("renders the Name your space heading", () => {
                // Assert
                const heading = screen.getByRole("heading", {
                  name: "Name your space",
                });
                expect(heading).toHaveClass("step-title", "text-balance");
              });

              it("renders the Workspace · Space eyebrow", () => {
                // Assert
                const eyebrow = screen.getByText("Workspace · Space");
                expect(eyebrow).toHaveClass("step-eyebrow");
              });

              it("renders the space lead", () => {
                // Assert
                const lead = screen.getByText(
                  /A space groups your lists together/,
                );
                expect(lead).toHaveClass("step-lead", "text-pretty");
                expect(lead).toHaveTextContent(
                  "A space groups your lists together. Most people start with a single personal space.",
                );
              });

              it("renders the Instance phase as done", () => {
                // Assert
                const label = [...document.querySelectorAll(".phase-seg-label")].find(
                  (el) => el.textContent === "Instance",
                );
                expect(label).toHaveClass("phase-seg-label");
                const seg = label?.closest(".phase-seg");
                expect(seg).toHaveClass("is-done");
                expect(
                  seg?.querySelector(".phase-seg-index svg.lucide-check"),
                ).toBeInTheDocument();
                expect(seg?.querySelector(".phase-seg-fill")).toHaveStyle({
                  width: "100%",
                });
              });

              it("renders the Workspace phase as active", () => {
                // Assert
                const label = document.querySelector(
                  ".phase-seg.is-active .phase-seg-label",
                );
                expect(label).toHaveTextContent("Workspace");
                const seg = label?.closest(".phase-seg");
                expect(seg).toHaveClass("is-active");
                expect(seg?.querySelector(".phase-seg-index")).toHaveTextContent(
                  "2",
                );
                expect(seg?.querySelector(".phase-seg-fill")).toHaveStyle({
                  width: "25%",
                });
              });

              it("renders the Ready phase as upcoming", () => {
                // Assert
                const label = [...document.querySelectorAll(".phase-seg-label")].find(
                  (el) => el.textContent === "Ready",
                );
                expect(label).toHaveClass("phase-seg-label");
                const seg = label?.closest(".phase-seg");
                expect(seg).toHaveClass("is-upcoming");
                expect(seg?.querySelector(".phase-seg-index")).toHaveTextContent(
                  "3",
                );
                expect(seg?.querySelector(".phase-seg-fill")).toHaveStyle({
                  width: "0%",
                });
              });

              it("renders a Space name field", () => {
                // Assert
                const input = screen.getByLabelText("Space name");
                expect(input).toHaveClass("input");
                expect(input).toHaveAttribute("id", "onb-space");
                expect(input).toHaveAttribute("placeholder", "Personal");
                expect(input).toHaveValue("Personal");
                expect(screen.getByText("Space name")).toHaveClass("label");
              });

              it("renders a layers icon in the Space name field", () => {
                // Assert
                const control = screen
                  .getByLabelText("Space name")
                  .closest(".control");
                expect(control).toHaveClass("has-icons-left");
                expect(
                  control?.querySelector("svg.lucide-layers"),
                ).toBeInTheDocument();
              });

              it("renders a Back button", () => {
                // Assert
                const back = screen.getByRole("button", { name: "Back" });
                expect(back).toHaveClass("button", "is-light");
                expect(
                  back.querySelector("svg.lucide-arrow-left"),
                ).toBeInTheDocument();
                expect(back.closest(".onboarding-footer")).toBeInTheDocument();
              });

              it("renders a Continue button", () => {
                // Assert
                const button = screen.getByRole("button", { name: "Continue" });
                expect(button).toHaveClass("button", "is-primary", "footer-grow");
                expect(button).toBeEnabled();
                expect(
                  button.querySelector("svg.lucide-arrow-right"),
                ).toBeInTheDocument();
                expect(button.closest(".onboarding-footer")).toBeInTheDocument();
              });

              describe("when the Continue button is clicked", () => {
                beforeEach(async () => {
                  vi.mocked(createSpace).mockResolvedValue({
                    ID: "SP_1",
                    Name: "Personal",
                  });
                  
                  await act(async () => {
                    fireEvent.click(screen.getByRole("button", { name: "Continue" }));
                  });
                });

                it("renders the What should we call you heading", () => {
                  // Assert
                  const heading = screen.getByRole("heading", {
                    name: "What should we call you?",
                  });
                  expect(heading).toHaveClass("step-title", "text-balance");
                });

                it("renders the Workspace · You eyebrow", () => {
                  // Assert
                  const eyebrow = screen.getByText("Workspace · You");
                  expect(eyebrow).toHaveClass("step-eyebrow");
                });

                it("renders the name lead", () => {
                  // Assert
                  const lead = screen.getByText(
                    /Your name shows up on comments and activity/,
                  );
                  expect(lead).toHaveClass("step-lead", "text-pretty");
                  expect(lead).toHaveTextContent(
                    "Your name shows up on comments and activity. It's just for you — no account needed.",
                  );
                });

                it("renders the Instance phase as done", () => {
                  // Assert
                  const label = [
                    ...document.querySelectorAll(".phase-seg-label"),
                  ].find((el) => el.textContent === "Instance");
                  expect(label).toHaveClass("phase-seg-label");
                  const seg = label?.closest(".phase-seg");
                  expect(seg).toHaveClass("is-done");
                  expect(
                    seg?.querySelector(".phase-seg-index svg.lucide-check"),
                  ).toBeInTheDocument();
                  expect(seg?.querySelector(".phase-seg-fill")).toHaveStyle({
                    width: "100%",
                  });
                });

                it("renders the Workspace phase fill at 50%", () => {
                  // Assert
                  const label = document.querySelector(
                    ".phase-seg.is-active .phase-seg-label",
                  );
                  expect(label).toHaveTextContent("Workspace");
                  expect(
                    label?.closest(".phase-seg")?.querySelector(".phase-seg-fill"),
                  ).toHaveStyle({
                    width: "50%",
                  });
                });

                it("renders the Ready phase as upcoming", () => {
                  // Assert
                  const label = [
                    ...document.querySelectorAll(".phase-seg-label"),
                  ].find((el) => el.textContent === "Ready");
                  expect(label).toHaveClass("phase-seg-label");
                  const seg = label?.closest(".phase-seg");
                  expect(seg).toHaveClass("is-upcoming");
                  expect(seg?.querySelector(".phase-seg-index")).toHaveTextContent(
                    "3",
                  );
                  expect(seg?.querySelector(".phase-seg-fill")).toHaveStyle({
                    width: "0%",
                  });
                });

                it("renders a Your name field", () => {
                  // Assert
                  const input = screen.getByLabelText("Your name");
                  expect(input).toHaveClass("input");
                  expect(input).toHaveAttribute("id", "onb-user");
                  expect(input).toHaveAttribute("placeholder", "e.g. Alex");
                  expect(input).toHaveValue("");
                  expect(screen.getByText("Your name")).toHaveClass("label");
                });

                it("renders a user icon in the Your name field", () => {
                  // Assert
                  const control = screen
                    .getByLabelText("Your name")
                    .closest(".control");
                  expect(control).toHaveClass("has-icons-left");
                  expect(
                    control?.querySelector("svg.lucide-user"),
                  ).toBeInTheDocument();
                });

                it("renders a Back button", () => {
                  // Assert
                  const back = screen.getByRole("button", { name: "Back" });
                  expect(back).toHaveClass("button", "is-light");
                  expect(
                    back.querySelector("svg.lucide-arrow-left"),
                  ).toBeInTheDocument();
                  expect(back.closest(".onboarding-footer")).toBeInTheDocument();
                });

                it("renders a disabled Continue button", () => {
                  // Assert
                  const button = screen.getByRole("button", { name: "Continue" });
                  expect(button).toHaveClass("button", "is-primary", "footer-grow");
                  expect(button).toBeDisabled();
                  expect(
                    button.querySelector("svg.lucide-arrow-right"),
                  ).toBeInTheDocument();
                  expect(button.closest(".onboarding-footer")).toBeInTheDocument();
                });

                describe("when a name is entered", () => {
                  beforeEach(() => {
                    fireEvent.change(screen.getByLabelText("Your name"), {
                      target: { value: "Alex" },
                    });
                  });

                  it("enables the Continue button", () => {
                    // Assert
                    const continueBtn = screen.getByRole("button", { name: "Continue" });
                    expect(continueBtn).not.toBeDisabled();
                  });

                  describe("when the name is deleted", () => {
                    beforeEach(() => {
                      fireEvent.change(screen.getByLabelText("Your name"), {
                        target: { value: "" },
                      });
                    });

                    it("disables the Continue button again", () => {
                      // Assert
                      const continueBtn = screen.getByRole("button", { name: "Continue" });
                      expect(continueBtn).toBeDisabled();
                    });
                  });

                  describe("when the Continue button is clicked", () => {
                    beforeEach(async () => {
                      vi.mocked(createUser).mockResolvedValue({
                        ID: "US_1",
                        Name: "Alex",
                      });
                      
                      await act(async () => {
                        fireEvent.click(screen.getByRole("button", { name: "Continue" }));
                      });
                    });

                    it("calls createUser with the entered name", () => {
                      // Assert
                      expect(createUser).toHaveBeenCalledWith({ name: "Alex" });
                    });

                    describe("when the Continue button is clicked", () => {
                      beforeEach(() => {
                        fireEvent.click(screen.getByRole("button", { name: "Continue" }));
                      });

                      it("renders the Getting things ready heading", () => {
                        // Assert
                        const heading = screen.getByRole("heading", {
                          name: "Getting things ready",
                        });
                        expect(heading).toHaveClass("step-title", "text-balance");
                      });

                      it("renders the Workspace · Automatic eyebrow", () => {
                        // Assert
                        const eyebrow = screen.getByText("Workspace · Automatic");
                        expect(eyebrow).toHaveClass("step-eyebrow");
                      });

                      it("renders the setup lead with the space name", () => {
                        // Assert
                        const lead = screen.getByText(/Listello is wiring up the essentials for/);
                        expect(lead).toHaveClass("step-lead", "text-pretty");
                        expect(lead).toHaveTextContent("Listello is wiring up the essentials for Personal.");
                      });

                      it("renders the Workspace phase fill at 75%", () => {
                        // Assert
                        const label = document.querySelector(
                          ".phase-seg.is-active .phase-seg-label",
                        );
                        expect(label).toHaveTextContent("Workspace");
                        expect(
                          label?.closest(".phase-seg")?.querySelector(".phase-seg-fill"),
                        ).toHaveStyle({
                          width: "75%",
                        });
                      });

                      it("renders the Create Inbox setup check", () => {
                        // Assert
                        const check = screen.getByText("Create Inbox");
                        expect(check).toHaveClass("setup-check-label");
                      });

                      it("renders the Assign space to user setup check", () => {
                        // Assert
                        const check = screen.getByText("Assign Personal to Alex");
                        expect(check).toHaveClass("setup-check-label");
                      });

                      it("renders a disabled Continue button", () => {
                        // Assert
                        const button = screen.getByRole("button", { name: "Continue" });
                        expect(button).toHaveClass("button", "is-primary", "footer-grow");
                        expect(button).toBeDisabled();
                      });
                    });
                  });
                });
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

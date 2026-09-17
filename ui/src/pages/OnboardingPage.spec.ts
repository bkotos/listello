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
  pairSpace: vi.fn(),
}));

vi.mock("../lib/api/user-client", () => ({
  createUser: vi.fn(),
}));

import {
  createInstance,
  getDefaultPersistenceLocation,
  getInstance,
  initializePersistence,
  pairSpace,
  selectHostingMode,
  selectPersistenceLocation,
} from "../lib/api/instance-client";
import { createUser } from "../lib/api/user-client";

const createdInstance: ListelloInstanceResponse = {
  HostingMode: "",
  PersistenceLocation: "",
  PersistenceState: "",
  SetupState: "",
  Space: { ID: "", Name: "" },
  User: { ID: "", Name: "" },
};

const apiDefaultPersistenceLocation =
  "/Users/api/Library/Application Support/listello";

const instancePersistenceLocation = "/Users/me/Documents/listello";

function renderOnboardingPage() {
  const { QueryWrapper } = createQueryWrapper();
  return render(createElement(QueryWrapper, null, createElement(OnboardingPage)));
}

function expectStepHeading(name: string) {
  const heading = screen.getByRole("heading", { name });
  expect(heading).toHaveClass("step-title", "text-balance");
}

function itRendersStepHeading(name: string) {
  it(`renders the ${name} heading`, () => {
    // Assert
    expectStepHeading(name);
  });
}

function itRendersStepHeadingOnPage(name: string) {
  it(`renders the ${name} heading`, () => {
    // Assert
    const heading = screen.getByRole("heading", { name });
    expect(heading).toHaveClass("step-title", "text-balance");
    expect(heading.closest(".onboarding-page")).toBeInTheDocument();
  });
}

function expectStepEyebrow(text: string) {
  const eyebrow = screen.getByText(text);
  expect(eyebrow).toHaveClass("step-eyebrow");
}

function itRendersStepEyebrow(text: string) {
  it(`renders the ${text} eyebrow`, () => {
    // Assert
    expectStepEyebrow(text);
  });
}

function expectStepLead(matcher: string | RegExp, fullText: string) {
  const lead = screen.getByText(matcher);
  expect(lead).toHaveClass("step-lead", "text-pretty");
  expect(lead).toHaveTextContent(fullText);
}

function itRendersStepLead(
  title: string,
  matcher: string | RegExp,
  fullText: string,
) {
  it(title, () => {
    // Assert
    expectStepLead(matcher, fullText);
  });
}

function itRendersStepLeadWithStrongs(
  title: string,
  matcher: string | RegExp,
  fullText: string,
  firstStrong: string,
  secondStrong: string,
) {
  it(title, () => {
    // Assert
    const lead = screen.getByText(matcher);
    expect(lead).toHaveClass("step-lead", "text-pretty");
    expect(lead).toHaveTextContent(fullText);
    expect(lead.querySelector("strong")).toHaveTextContent(firstStrong);
    expect(lead.querySelectorAll("strong")[1]).toHaveTextContent(secondStrong);
  });
}

function phaseSeg(name: string) {
  const label = [...document.querySelectorAll(".phase-seg-label")].find(
    (el) => el.textContent === name,
  );
  return { label, seg: label?.closest(".phase-seg") };
}

function expectActivePhase(name: string, fill: string, index: string) {
  const label = document.querySelector(
    ".phase-seg.is-active .phase-seg-label",
  );
  expect(label).toHaveClass("phase-seg-label");
  expect(label).toHaveTextContent(name);
  const seg = label?.closest(".phase-seg");
  expect(seg).toHaveClass("is-active");
  expect(seg?.querySelector(".phase-seg-index")).toHaveTextContent(index);
  expect(seg?.querySelector(".phase-seg-fill")).toHaveStyle({ width: fill });
}

function itRendersActivePhase(name: string, fill: string, index: string) {
  it(`renders the ${name} phase as active`, () => {
    // Assert
    expectActivePhase(name, fill, index);
  });
}

function expectUpcomingPhase(name: string, fill: string, index: string) {
  const { label, seg } = phaseSeg(name);
  expect(label).toHaveClass("phase-seg-label");
  expect(seg).toHaveClass("is-upcoming");
  expect(seg?.querySelector(".phase-seg-index")).toHaveTextContent(index);
  expect(seg?.querySelector(".phase-seg-fill")).toHaveStyle({ width: fill });
}

function itRendersUpcomingPhase(name: string, fill: string, index: string) {
  it(`renders the ${name} phase as upcoming`, () => {
    // Assert
    expectUpcomingPhase(name, fill, index);
  });
}

function expectDonePhase(name: string, fill: string) {
  const { label, seg } = phaseSeg(name);
  expect(label).toHaveClass("phase-seg-label");
  expect(seg).toHaveClass("is-done");
  expect(
    seg?.querySelector(".phase-seg-index svg.lucide-check"),
  ).toBeInTheDocument();
  expect(seg?.querySelector(".phase-seg-fill")).toHaveStyle({ width: fill });
}

function itRendersDonePhase(name: string, fill: string) {
  it(`renders the ${name} phase as done`, () => {
    // Assert
    expectDonePhase(name, fill);
  });
}

function expectPhaseFill(name: string, fill: string) {
  const label = document.querySelector(
    ".phase-seg.is-active .phase-seg-label",
  );
  expect(label).toHaveTextContent(name);
  expect(
    label?.closest(".phase-seg")?.querySelector(".phase-seg-fill"),
  ).toHaveStyle({
    width: fill,
  });
}

function itRendersPhaseFill(name: string, fill: string) {
  it(`renders the ${name} phase fill at ${fill}`, () => {
    // Assert
    expectPhaseFill(name, fill);
  });
}

function setupCheckRow(label: string) {
  const labelEl = screen.getByText(label);
  expect(labelEl).toHaveClass("setup-check-label");
  return labelEl.closest(".setup-check");
}

function expectPendingSetupCheck(label: string) {
  const row = setupCheckRow(label);
  expect(row).toHaveClass("is-pending");
  expect(row?.querySelector(".check-toggle")).toBeInTheDocument();
}

function expectInProgressSetupCheck(label: string) {
  const row = setupCheckRow(label);
  expect(row).not.toHaveClass("is-pending");
  expect(row).not.toHaveClass("is-done");
  const loader = row?.querySelector("svg.lucide-loader-circle");
  expect(loader).toBeInTheDocument();
  expect(loader).toHaveClass("spin");
}

function expectDoneSetupCheck(label: string) {
  const row = setupCheckRow(label);
  expect(row).toHaveClass("is-done");
  expect(row?.querySelector("svg.lucide-circle-check")).toBeInTheDocument();
}

function expectDoneSetupCheckWithIcon(label: string, icon: string) {
  const row = setupCheckRow(label);
  expect(row).toHaveClass("is-done");
  expect(row?.querySelector(`svg.${icon}`)).toBeInTheDocument();
}

function itRendersPendingSetupCheck(title: string, label: string) {
  it(title, () => {
    // Assert
    expectPendingSetupCheck(label);
  });
}

function itRendersInProgressSetupCheck(title: string, label: string) {
  it(title, () => {
    // Assert
    expectInProgressSetupCheck(label);
  });
}

function itRendersDoneSetupCheck(title: string, label: string) {
  it(title, () => {
    // Assert
    expectDoneSetupCheck(label);
  });
}

function itRendersDoneSetupCheckWithIcon(label: string, icon: string) {
  it(`renders a ${label} setup check`, () => {
    // Assert
    expectDoneSetupCheckWithIcon(label, icon);
  });
}

function expectBackButton() {
  const back = screen.getByRole("button", { name: "Back" });
  expect(back).toHaveClass("button", "is-light");
  expect(back.querySelector("svg.lucide-arrow-left")).toBeInTheDocument();
  expect(back.closest(".onboarding-footer")).toBeInTheDocument();
}

function itRendersBackButton() {
  it("renders a Back button", () => {
    // Assert
    expectBackButton();
  });
}

function continueButton() {
  return screen.getByRole("button", { name: "Continue" });
}

function expectContinueButton() {
  const button = continueButton();
  expect(button).toHaveClass("button", "is-primary", "footer-grow");
  expect(button.querySelector("svg.lucide-arrow-right")).toBeInTheDocument();
  expect(button.closest(".onboarding-footer")).toBeInTheDocument();
}

function expectEnabledContinueButton() {
  const button = continueButton();
  expect(button).toHaveClass("button", "is-primary", "footer-grow");
  expect(button).toBeEnabled();
  expect(button.querySelector("svg.lucide-arrow-right")).toBeInTheDocument();
  expect(button.closest(".onboarding-footer")).toBeInTheDocument();
}

function expectDisabledContinueButton() {
  const button = continueButton();
  expect(button).toHaveClass("button", "is-primary", "footer-grow");
  expect(button).toBeDisabled();
  expect(button.querySelector("svg.lucide-arrow-right")).toBeInTheDocument();
  expect(button.closest(".onboarding-footer")).toBeInTheDocument();
}

function expectDisabledContinueButtonWithoutChrome() {
  const button = continueButton();
  expect(button).toHaveClass("button", "is-primary", "footer-grow");
  expect(button).toBeDisabled();
}

function itRendersContinueButton() {
  it("renders a Continue button", () => {
    // Assert
    expectContinueButton();
  });
}

function itRendersEnabledContinueButton() {
  it("renders a Continue button", () => {
    // Assert
    expectEnabledContinueButton();
  });
}

function itRendersDisabledContinueButton() {
  it("renders a disabled Continue button", () => {
    // Assert
    expectDisabledContinueButton();
  });
}

function itRendersDisabledContinueButtonWithoutChrome() {
  it("renders a disabled Continue button", () => {
    // Assert
    expectDisabledContinueButtonWithoutChrome();
  });
}

function itEnablesContinueButton() {
  it("enables the Continue button", () => {
    // Assert
    expectEnabledContinueButton();
  });
}

function expectField(
  label: string,
  options: { id: string; placeholder: string; value: string },
) {
  const input = screen.getByLabelText(label);
  expect(input).toHaveClass("input");
  expect(input).toHaveAttribute("id", options.id);
  expect(input).toHaveAttribute("placeholder", options.placeholder);
  expect(input).toHaveValue(options.value);
  expect(screen.getByText(label)).toHaveClass("label");
}

function itRendersField(
  label: string,
  options: { id: string; placeholder: string; value: string },
) {
  it(`renders a ${label} field`, () => {
    // Assert
    expectField(label, options);
  });
}

function expectFieldIcon(label: string, iconName: string) {
  const control = screen.getByLabelText(label).closest(".control");
  expect(control).toHaveClass("has-icons-left");
  expect(control?.querySelector(`svg.lucide-${iconName}`)).toBeInTheDocument();
}

function itRendersFieldIcon(label: string, iconName: string) {
  it(`renders a ${iconName} icon in the ${label} field`, () => {
    // Assert
    expectFieldIcon(label, iconName);
  });
}

type ChoiceCardDetails = {
  name: string | RegExp;
  title: string;
  desc: string;
  meta: string;
  icon: string;
};

function choiceCard(name: string | RegExp) {
  return screen.getByRole("button", { name });
}

function expectSelectedChoiceCardChrome(card: HTMLElement) {
  expect(card).toHaveClass("choice-card", "is-selected");
  expect(card).toHaveAttribute("aria-pressed", "true");
  expect(
    card.querySelector(".choice-check svg.lucide-circle-check"),
  ).toBeInTheDocument();
}

function expectUnselectedChoiceCardChrome(card: HTMLElement) {
  expect(card).toHaveClass("choice-card");
  expect(card).not.toHaveClass("is-selected");
  expect(card).toHaveAttribute("aria-pressed", "false");
  expect(card.querySelector(".choice-check")).not.toBeInTheDocument();
}

function expectChoiceCardDetails(card: HTMLElement, options: ChoiceCardDetails) {
  expect(card.querySelector(`svg.lucide-${options.icon}`)).toBeInTheDocument();
  expect(card.querySelector(".choice-title")).toHaveTextContent(options.title);
  expect(card.querySelector(".choice-desc")).toHaveTextContent(options.desc);
  expect(card.querySelector(".choice-meta")).toHaveTextContent(options.meta);
  expect(card.querySelector("svg.lucide-database")).toBeInTheDocument();
}

function itRendersSelectedChoiceCard(
  title: string,
  options: ChoiceCardDetails,
) {
  it(title, () => {
    // Assert
    const card = choiceCard(options.name);
    expectSelectedChoiceCardChrome(card);
    expectChoiceCardDetails(card, options);
  });
}

function itRendersUnselectedChoiceCard(
  title: string,
  options: ChoiceCardDetails,
) {
  it(title, () => {
    // Assert
    const card = choiceCard(options.name);
    expectUnselectedChoiceCardChrome(card);
    expectChoiceCardDetails(card, options);
  });
}

function itSelectsChoiceCard(name: string) {
  it(`selects the ${name} choice card`, () => {
    // Assert
    expectSelectedChoiceCardChrome(choiceCard(new RegExp(name)));
  });
}

function itDeselectsChoiceCard(name: string) {
  it(`deselects the ${name} choice card`, () => {
    // Assert
    expectUnselectedChoiceCardChrome(choiceCard(new RegExp(name)));
  });
}

async function clickAndWaitForHeading(
  buttonName: string,
  heading: string,
) {
  fireEvent.click(screen.getByRole("button", { name: buttonName }));
  await waitFor(() => {
    expect(screen.getByRole("heading", { name: heading })).toBeInTheDocument();
  });
}

async function clickContinueAndWaitForHeading(heading: string) {
  await clickAndWaitForHeading("Continue", heading);
}

async function clickBackAndWaitForHeading(heading: string) {
  await clickAndWaitForHeading("Back", heading);
}

function advanceBy500ms() {
  act(() => {
    vi.advanceTimersByTime(500);
  });
}

async function waitForDataDirectoryValue() {
  await waitFor(() => {
    expect(screen.getByLabelText("Data directory")).toHaveValue(
      apiDefaultPersistenceLocation,
    );
  });
}

function itRendersInitializeStepBeforeTimeouts() {
  itRendersStepHeading("Setting up persistence");
  itRendersStepEyebrow("Instance · Initialize");
  itRendersStepLeadWithStrongs(
    "renders the initialize lead",
    /We're initializing the/,
    "We're initializing the SQLite store for your local instance.",
    "SQLite",
    "local",
  );
  itRendersPhaseFill("Instance", "100%");
  itRendersInProgressSetupCheck(
    "renders a Create instance setup check in progress",
    "Create instance",
  );
  itRendersPendingSetupCheck(
    "renders a Prepare SQLite store setup check as pending",
    "Prepare SQLite store",
  );
  itRendersPendingSetupCheck(
    "renders an Initialize persistence setup check as pending",
    "Initialize persistence",
  );
  itRendersBackButton();
  itRendersDisabledContinueButton();
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

  itRendersStepHeadingOnPage("Welcome to Listello");

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

  itRendersActivePhase("Instance", "25%", "1");
  itRendersUpcomingPhase("Workspace", "0%", "2");
  itRendersUpcomingPhase("Ready", "0%", "3");

  it("renders a sparkles icon", () => {
    // Assert
    const heading = screen.getByRole("heading", {
      name: "Welcome to Listello",
    });
    const icon = heading.parentElement?.querySelector(".big-icon");
    expect(icon).toBeInTheDocument();
    expect(icon?.querySelector("svg.lucide-sparkles")).toBeInTheDocument();
  });

  itRendersStepLead(
    "renders the welcome lead",
    /Let's create your instance and set up a calm, GTD-style workspace/,
    "Let's create your instance and set up a calm, GTD-style workspace. It only takes a minute, and you can change everything later.",
  );
  itRendersDoneSetupCheckWithIcon(
    "Choose how Listello is hosted",
    "lucide-server",
  );
  itRendersDoneSetupCheckWithIcon(
    "Create your space and profile",
    "lucide-layers",
  );
  itRendersDoneSetupCheckWithIcon("Start your first list", "lucide-list-checks");

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
      await clickAndWaitForHeading(
        "Create instance",
        "How should Listello be hosted?",
      );
    });

    itRendersStepHeading("How should Listello be hosted?");
    itRendersStepEyebrow("Instance · Hosting");
    itRendersStepLead(
      "renders the hosting lead",
      /Your hosting mode decides where this instance runs/,
      "Your hosting mode decides where this instance runs and how your data is stored. You can migrate later.",
    );
    itRendersPhaseFill("Instance", "50%");
    itRendersSelectedChoiceCard("renders a selected Local choice card", {
      name: /Local/,
      title: "Local",
      desc: "Runs directly on this machine, with local filesystem-backed persistence.",
      meta: "Uses SQLite",
      icon: "hard-drive",
    });
    itRendersUnselectedChoiceCard("renders a Standalone Web choice card", {
      name: /Standalone Web/,
      title: "Standalone Web",
      desc: "Runs independently in the browser as a PWA or standalone app.",
      meta: "Uses IndexedDB / OPFS",
      icon: "globe",
    });
    itRendersBackButton();
    itRendersContinueButton();

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
        await clickBackAndWaitForHeading("Welcome to Listello");
      });

      itRendersStepHeading("Welcome to Listello");
    });

    describe("when the Continue button is clicked", () => {
      beforeEach(async () => {
        await clickContinueAndWaitForHeading("Choose a data directory");
      });

      itRendersStepHeading("Choose a data directory");
      itRendersStepEyebrow("Instance · Local");
      itRendersStepLead(
        "renders the data directory lead",
        /This is where Listello keeps its data for this instance/,
        "This is where Listello keeps its data for this instance. You can move it later.",
      );
      itRendersPhaseFill("Instance", "75%");

      it("renders a Data directory field", async () => {
        // Assert
        const input = screen.getByLabelText("Data directory");
        expect(input).toHaveClass("input");
        expect(input).toHaveAttribute("id", "onb-location");
        expect(input).toHaveAttribute("placeholder", "~/listello");
        await waitForDataDirectoryValue();
        expect(getDefaultPersistenceLocation).toHaveBeenCalledWith(
          expect.objectContaining({ signal: expect.any(AbortSignal) }),
        );
        expect(screen.getByText("Data directory")).toHaveClass("label");
      });

      itRendersFieldIcon("Data directory", "folder");

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
        await waitForDataDirectoryValue();

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
        await waitForDataDirectoryValue();

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
        await waitForDataDirectoryValue();

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
          await clickBackAndWaitForHeading("How should Listello be hosted?");
        });

        itRendersStepHeading("How should Listello be hosted?");
      });

      describe("when the Continue button is clicked", () => {
        beforeEach(async () => {
          await clickContinueAndWaitForHeading("Setting up persistence");
        });

        itRendersInitializeStepBeforeTimeouts();
      });

      describe("after 500ms", () => {
        beforeEach(() => {
          vi.useFakeTimers();
          fireEvent.click(screen.getByRole("button", { name: "Continue" }));
          advanceBy500ms();
        });

        itRendersDoneSetupCheck(
          "renders a Create instance setup check as done",
          "Create instance",
        );
        itRendersInProgressSetupCheck(
          "renders a Prepare SQLite store setup check in progress",
          "Prepare SQLite store",
        );
        itRendersPendingSetupCheck(
          "renders an Initialize persistence setup check as pending",
          "Initialize persistence",
        );

        describe("after another 500ms", () => {
          beforeEach(() => {
            advanceBy500ms();
          });

          itRendersDoneSetupCheck(
            "renders a Prepare SQLite store setup check as done",
            "Prepare SQLite store",
          );
          itRendersInProgressSetupCheck(
            "renders an Initialize persistence setup check in progress",
            "Initialize persistence",
          );

          describe("after another 500ms", () => {
            beforeEach(() => {
              advanceBy500ms();
            });

            itRendersDoneSetupCheck(
              "renders an Initialize persistence setup check as done",
              "Initialize persistence",
            );
            itEnablesContinueButton();

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

              itRendersStepHeading("Name your space");
              itRendersStepEyebrow("Workspace · Space");
              itRendersStepLead(
                "renders the space lead",
                /A space groups your lists together/,
                "A space groups your lists together. Most people start with a single personal space.",
              );
              itRendersDonePhase("Instance", "100%");
              itRendersActivePhase("Workspace", "25%", "2");
              itRendersUpcomingPhase("Ready", "0%", "3");
              itRendersField("Space name", {
                id: "onb-space",
                placeholder: "Personal",
                value: "Personal",
              });
              itRendersFieldIcon("Space name", "layers");
              itRendersBackButton();
              itRendersEnabledContinueButton();

              describe("when the Back button is clicked", () => {
                beforeEach(() => {
                  fireEvent.click(screen.getByRole("button", { name: "Back" }));
                });

                itRendersStepHeading("Choose a data directory");
              });

              describe("when the Continue button is clicked", () => {
                beforeEach(async () => {
                  vi.mocked(pairSpace).mockResolvedValue(createdInstance);

                  await act(async () => {
                    fireEvent.click(screen.getByRole("button", { name: "Continue" }));
                  });
                });

                itRendersStepHeading("What should we call you?");
                itRendersStepEyebrow("Workspace · You");
                itRendersStepLead(
                  "renders the name lead",
                  /Your name shows up on comments and activity/,
                  "Your name shows up on comments and activity. It's just for you — no account needed.",
                );
                itRendersDonePhase("Instance", "100%");
                itRendersPhaseFill("Workspace", "50%");
                itRendersUpcomingPhase("Ready", "0%", "3");
                itRendersField("Your name", {
                  id: "onb-user",
                  placeholder: "e.g. Alex",
                  value: "",
                });
                itRendersFieldIcon("Your name", "user");
                itRendersBackButton();
                itRendersDisabledContinueButton();

                it("calls pairSpace with the space name", () => {
                  // Assert
                  expect(pairSpace).toHaveBeenCalledWith({ name: "Personal" });
                });

                describe("when the Back button is clicked", () => {
                  beforeEach(() => {
                    fireEvent.click(screen.getByRole("button", { name: "Back" }));
                  });

                  itRendersStepHeading("Name your space");
                });

                describe("when a name is entered", () => {
                  beforeEach(() => {
                    fireEvent.change(screen.getByLabelText("Your name"), {
                      target: { value: "Alex" },
                    });
                  });

                  it("enables the Continue button", () => {
                    // Assert
                    const continueBtn = screen.getByRole("button", {
                      name: "Continue",
                    });
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
                      const continueBtn = screen.getByRole("button", {
                        name: "Continue",
                      });
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
                        fireEvent.click(
                          screen.getByRole("button", { name: "Continue" }),
                        );
                      });
                    });

                    it("calls createUser with the entered name", () => {
                      // Assert
                      expect(createUser).toHaveBeenCalledWith({ name: "Alex" });
                    });

                    describe("when the Continue button is clicked", () => {
                      beforeEach(() => {
                        fireEvent.click(
                          screen.getByRole("button", { name: "Continue" }),
                        );
                      });

                      itRendersStepHeading("Getting things ready");
                      itRendersStepEyebrow("Workspace · Automatic");
                      itRendersStepLead(
                        "renders the setup lead with the space name",
                        /Listello is wiring up the essentials for/,
                        "Listello is wiring up the essentials for Personal.",
                      );
                      itRendersPhaseFill("Workspace", "75%");
                      itRendersInProgressSetupCheck(
                        "renders a Create Inbox setup check in progress",
                        "Create Inbox",
                      );
                      itRendersPendingSetupCheck(
                        "renders an Assign space to user setup check as pending",
                        "Assign Personal to Alex",
                      );
                      itRendersBackButton();
                      itRendersDisabledContinueButtonWithoutChrome();

                      describe("when the Back button is clicked", () => {
                        beforeEach(() => {
                          fireEvent.click(screen.getByRole("button", { name: "Back" }));
                        });

                        itRendersStepHeading("What should we call you?");
                      });

                      describe("after 500ms", () => {
                        beforeEach(() => {
                          advanceBy500ms();
                        });

                        itRendersDoneSetupCheck(
                          "renders a Create Inbox setup check as done",
                          "Create Inbox",
                        );
                        itRendersInProgressSetupCheck(
                          "renders an Assign space to user setup check in progress",
                          "Assign Personal to Alex",
                        );

                        it("keeps the Continue button disabled", () => {
                          // Assert
                          const button = screen.getByRole("button", {
                            name: "Continue",
                          });
                          expect(button).toBeDisabled();
                        });

                        describe("after another 500ms", () => {
                          beforeEach(() => {
                            advanceBy500ms();
                          });

                          itRendersDoneSetupCheck(
                            "renders an Assign space to user setup check as done",
                            "Assign Personal to Alex",
                          );
                          itEnablesContinueButton();

                          describe("when the Continue button is clicked", () => {
                            beforeEach(() => {
                              fireEvent.click(
                                screen.getByRole("button", { name: "Continue" }),
                              );
                            });

                            itRendersStepHeading("Create your first list");
                            itRendersStepEyebrow("Workspace · First list");
                            itRendersStepLead(
                              "renders the first list lead",
                              /Lists hold the tasks you want to act on/,
                              "Lists hold the tasks you want to act on. Give your first one a name, or skip and create one later.",
                            );
                            itRendersDonePhase("Instance", "100%");
                            itRendersActivePhase("Workspace", "100%", "2");
                            itRendersUpcomingPhase("Ready", "0%", "3");
                            itRendersField("List name", {
                              id: "onb-list",
                              placeholder: "Errands",
                              value: "Errands",
                            });
                            itRendersFieldIcon("List name", "list-checks");

                            it("renders a suggestion hint", () => {
                              // Assert
                              const hint = screen.getByText(
                                "Or start with one of these",
                              );
                              expect(hint).toHaveClass("suggest-hint");
                            });

                            it("renders an Errands suggestion tag as selected", () => {
                              // Assert
                              const tag = screen.getByRole("button", {
                                name: "Errands",
                              });
                              expect(tag).toHaveClass(
                                "tag",
                                "is-medium",
                                "is-primary",
                              );
                              expect(tag.closest(".tags")).toHaveClass(
                                "tags",
                                "mt-2",
                              );
                            });

                            it("renders a Shopping suggestion tag", () => {
                              // Assert
                              const tag = screen.getByRole("button", {
                                name: "Shopping",
                              });
                              expect(tag).toHaveClass("tag", "is-medium");
                              expect(tag).not.toHaveClass("is-primary");
                            });

                            it("renders an Ideas suggestion tag", () => {
                              // Assert
                              const tag = screen.getByRole("button", {
                                name: "Ideas",
                              });
                              expect(tag).toHaveClass("tag", "is-medium");
                              expect(tag).not.toHaveClass("is-primary");
                            });

                            it("renders a Reading suggestion tag", () => {
                              // Assert
                              const tag = screen.getByRole("button", {
                                name: "Reading",
                              });
                              expect(tag).toHaveClass("tag", "is-medium");
                              expect(tag).not.toHaveClass("is-primary");
                            });

                            it("renders a Goals suggestion tag", () => {
                              // Assert
                              const tag = screen.getByRole("button", {
                                name: "Goals",
                              });
                              expect(tag).toHaveClass("tag", "is-medium");
                              expect(tag).not.toHaveClass("is-primary");
                            });

                            it("renders a Skip for now button", () => {
                              // Assert
                              const skip = screen.getByRole("button", {
                                name: "Skip for now",
                              });
                              expect(skip).toHaveClass("skip-link");
                            });

                            itRendersBackButton();

                            it("renders a Create list button", () => {
                              // Assert
                              const button = screen.getByRole("button", {
                                name: "Create list",
                              });
                              expect(button).toHaveClass(
                                "button",
                                "is-primary",
                                "footer-grow",
                              );
                              expect(button).toBeEnabled();
                              expect(
                                button.querySelector("svg.lucide-arrow-right"),
                              ).toBeInTheDocument();
                              expect(
                                button.closest(".onboarding-footer"),
                              ).toBeInTheDocument();
                            });

                            const suggestionNames = [
                              "Errands",
                              "Shopping",
                              "Ideas",
                              "Reading",
                              "Goals",
                            ];
                            for (const name of suggestionNames) {
                              describe(`when the ${name} suggestion tag is clicked`, () => {
                                beforeEach(() => {
                                  fireEvent.change(
                                    screen.getByLabelText("List name"),
                                    { target: { value: "Custom" } },
                                  );
                                  fireEvent.click(
                                    screen.getByRole("button", { name }),
                                  );
                                });

                                it(`sets the List name field to ${name}`, () => {
                                  // Assert
                                  expect(
                                    screen.getByLabelText("List name"),
                                  ).toHaveValue(name);
                                });

                                it(`marks the ${name} suggestion tag as selected`, () => {
                                  // Assert
                                  expect(
                                    screen.getByRole("button", { name }),
                                  ).toHaveClass("tag", "is-medium", "is-primary");
                                });

                                it("deselects the other suggestion tags", () => {
                                  // Assert
                                  for (const other of suggestionNames) {
                                    if (other === name) {
                                      continue;
                                    }
                                    expect(
                                      screen.getByRole("button", { name: other }),
                                    ).not.toHaveClass("is-primary");
                                  }
                                });
                              });
                            }

                            describe("when a custom list name is entered", () => {
                              beforeEach(() => {
                                fireEvent.change(
                                  screen.getByLabelText("List name"),
                                  { target: { value: "Custom" } },
                                );
                              });

                              it("sets the List name field to Custom", () => {
                                // Assert
                                expect(
                                  screen.getByLabelText("List name"),
                                ).toHaveValue("Custom");
                              });

                              it("deselects all suggestion tags", () => {
                                // Assert
                                for (const name of suggestionNames) {
                                  expect(
                                    screen.getByRole("button", { name }),
                                  ).not.toHaveClass("is-primary");
                                }
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
          });
        });
      });
    });

    describe("when the Standalone Web choice card is clicked", () => {
      beforeEach(() => {
        fireEvent.click(screen.getByRole("button", { name: /Standalone Web/ }));
      });

      itSelectsChoiceCard("Standalone Web");
      itDeselectsChoiceCard("Local");
      itRendersPhaseFill("Instance", "67%");

      describe("when the Local choice card is clicked", () => {
        beforeEach(() => {
          fireEvent.click(screen.getByRole("button", { name: /Local/ }));
        });

        itSelectsChoiceCard("Local");
        itDeselectsChoiceCard("Standalone Web");
        itRendersPhaseFill("Instance", "50%");
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

  describe("when the instance has a HostingMode set", () => {
    beforeEach(async () => {
      cleanup();
      vi.mocked(getInstance).mockResolvedValue({
        ...createdInstance,
        HostingMode: "local",
      });
      renderOnboardingPage();
      await waitFor(() => {
        expect(
          screen.getByRole("heading", {
            name: "How should Listello be hosted?",
          }),
        ).toBeInTheDocument();
      });
    });

    itRendersStepHeading("How should Listello be hosted?");
  });

  describe("when the instance has a PersistenceLocation set", () => {
    beforeEach(async () => {
      cleanup();
      vi.mocked(getInstance).mockResolvedValue({
        ...createdInstance,
        HostingMode: "local",
        PersistenceLocation: instancePersistenceLocation,
      });
      renderOnboardingPage();
      await waitFor(() => {
        expect(
          screen.getByRole("heading", {
            name: "Choose a data directory",
          }),
        ).toBeInTheDocument();
      });
    });

    itRendersStepHeading("Choose a data directory");

    it("renders the Data directory field with the instance PersistenceLocation", () => {
      // Assert
      expect(screen.getByLabelText("Data directory")).toHaveValue(
        instancePersistenceLocation,
      );
    });
  });

  describe("when the instance PersistenceState is initialized", () => {
    beforeEach(async () => {
      cleanup();
      vi.mocked(getInstance).mockResolvedValue({
        ...createdInstance,
        HostingMode: "local",
        PersistenceLocation: instancePersistenceLocation,
        PersistenceState: "initialized",
      });
      renderOnboardingPage();
      await waitFor(() => {
        expect(
          screen.getByRole("heading", {
            name: "Name your space",
          }),
        ).toBeInTheDocument();
      });
    });

    itRendersStepHeading("Name your space");

    it("calls pairSpace with the space name when Continue is clicked", async () => {
      // Arrange
      vi.mocked(pairSpace).mockResolvedValue(createdInstance);

      // Act
      await act(async () => {
        fireEvent.click(screen.getByRole("button", { name: "Continue" }));
      });

      // Assert
      expect(pairSpace).toHaveBeenCalledWith({ name: "Personal" });
    });

    describe("when a custom space name is entered", () => {
      beforeEach(() => {
        fireEvent.change(screen.getByLabelText("Space name"), {
          target: { value: "Work" },
        });
      });

      describe("when the Continue button is clicked", () => {
        beforeEach(async () => {
          vi.mocked(pairSpace).mockResolvedValue({
            ...createdInstance,
            PersistenceState: "initialized",
            Space: { ID: "SP_1", Name: "Work" },
          });
          await act(async () => {
            fireEvent.click(screen.getByRole("button", { name: "Continue" }));
          });
        });

        describe("when the Back button is clicked", () => {
          beforeEach(() => {
            fireEvent.click(screen.getByRole("button", { name: "Back" }));
          });

          it("renders the Space name field with Work", () => {
            // Assert
            expect(screen.getByLabelText("Space name")).toHaveValue("Work");
          });

          describe("when the Continue button is clicked again", () => {
            beforeEach(async () => {
              await act(async () => {
                fireEvent.click(screen.getByRole("button", { name: "Continue" }));
              });
            });

            it("does not call pairSpace", () => {
              // Assert
              expect(pairSpace).toHaveBeenCalledTimes(1);
            });
          });

          describe("when the Space name is changed to Home", () => {
            beforeEach(() => {
              fireEvent.change(screen.getByLabelText("Space name"), {
                target: { value: "Home" },
              });
            });

            describe("when the Continue button is clicked again", () => {
              beforeEach(async () => {
                await act(async () => {
                  fireEvent.click(screen.getByRole("button", { name: "Continue" }));
                });
              });

              it("calls pairSpace with Home", () => {
                // Assert
                expect(pairSpace).toHaveBeenCalledWith({ name: "Home" });
              });
            });
          });
        });
      });
    });
  });

  describe("when the instance has a Space.ID set", () => {
    beforeEach(async () => {
      cleanup();
      vi.mocked(getInstance).mockResolvedValue({
        ...createdInstance,
        HostingMode: "local",
        PersistenceLocation: instancePersistenceLocation,
        PersistenceState: "initialized",
        Space: { ID: "SP_1", Name: "Personal" },
      });
      renderOnboardingPage();
      await waitFor(() => {
        expect(
          screen.getByRole("heading", {
            name: "What should we call you?",
          }),
        ).toBeInTheDocument();
      });
    });

    itRendersStepHeading("What should we call you?");
  });

  describe("when the instance has a Space named Work", () => {
    beforeEach(async () => {
      cleanup();
      vi.mocked(getInstance).mockResolvedValue({
        ...createdInstance,
        HostingMode: "local",
        PersistenceLocation: instancePersistenceLocation,
        PersistenceState: "initialized",
        Space: { ID: "SP_1", Name: "Work" },
      });
      renderOnboardingPage();
      await waitFor(() => {
        expect(
          screen.getByRole("heading", {
            name: "What should we call you?",
          }),
        ).toBeInTheDocument();
      });
    });

    describe("when the Back button is clicked", () => {
      beforeEach(() => {
        fireEvent.click(screen.getByRole("button", { name: "Back" }));
      });

      it("renders the Space name field with Work", () => {
        // Assert
        expect(screen.getByLabelText("Space name")).toHaveValue("Work");
      });
    });
  });
});

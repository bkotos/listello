import { cleanup, render, screen, waitFor } from "@testing-library/react";
import { createElement } from "react";
import type { ListelloInstanceResponse } from "api-types/listello-instance";
import { MemoryRouter } from "react-router-dom";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { AppProvider } from "./contexts/AppContext";
import { createQueryWrapper } from "./test/renderWithQueryClient";
import App from "./App";

vi.mock("./lib/api/instance-client", () => ({
  getInstance: vi.fn(),
  createInstance: vi.fn(),
  getDefaultPersistenceLocation: vi.fn(),
}));

vi.mock("./lib/api/list-client", () => ({
  getAllLists: vi.fn(),
  getList: vi.fn(),
  createList: vi.fn(),
}));

import { getDefaultPersistenceLocation, getInstance } from "./lib/api/instance-client";
import { getAllLists } from "./lib/api/list-client";

const existingInstance: ListelloInstanceResponse = {
  HostingMode: "",
  PersistenceLocation: "",
  PersistenceState: "",
  SetupState: "completed",
  Space: { ID: "", Name: "" },
  User: { ID: "", Name: "" },
};

afterEach(() => {
  cleanup();
  vi.clearAllMocks();
});

function renderApp(initialEntries = ["/"]) {
  const { QueryWrapper } = createQueryWrapper();

  return render(
    createElement(
      QueryWrapper,
      null,
      createElement(
        MemoryRouter,
        { initialEntries },
        createElement(AppProvider, null, createElement(App)),
      ),
    ),
  );
}

describe("App", () => {
  it("shows the onboarding page when the instance does not exist", async () => {
    // Arrange
    vi.mocked(getInstance).mockResolvedValue(null);
    vi.mocked(getAllLists).mockResolvedValue([]);
    renderApp();

    // Assert
    await waitFor(() => {
      expect(screen.getByRole("heading", { name: "Welcome to Listello" })).toBeInTheDocument();
    });
  });

  it("shows the inbox when the instance exists", async () => {
    // Arrange
    vi.mocked(getInstance).mockResolvedValue(existingInstance);
    vi.mocked(getAllLists).mockResolvedValue([]);
    renderApp();

    // Assert
    await waitFor(() => {
      expect(screen.getByRole("heading", { name: "Inbox" })).toBeInTheDocument();
    });
  });

  it.each(["incomplete", ""])(
    'shows the onboarding page when the instance setup is "%s"',
    async (setupState) => {
      // Arrange
      vi.mocked(getInstance).mockResolvedValue({
        ...existingInstance,
        SetupState: setupState,
      });
      vi.mocked(getAllLists).mockResolvedValue([]);
      renderApp();

      // Assert
      await waitFor(() => {
        expect(screen.getByRole("heading", { name: "Welcome to Listello" })).toBeInTheDocument();
      });
    },
  );

  describe("when visiting /onboarding and the instance does not exist", () => {
    beforeEach(async () => {
      vi.mocked(getInstance).mockResolvedValue(null);
      vi.mocked(getDefaultPersistenceLocation).mockResolvedValue({
        Location: "/Users/api/Library/Application Support/listello",
      });
      vi.mocked(getAllLists).mockResolvedValue([]);
      renderApp(["/onboarding"]);
      await waitFor(() => {
        expect(
          screen.getByRole("heading", { name: "Welcome to Listello" }),
        ).toBeInTheDocument();
      });
    });

    it("shows the Welcome heading", () => {
      // Assert
      expect(
        screen.getByRole("heading", { name: "Welcome to Listello" }),
      ).toBeInTheDocument();
    });

    describe("when the page is refreshed", () => {
      beforeEach(async () => {
        cleanup();
        renderApp(["/onboarding"]);
        await waitFor(() => {
          expect(
            screen.getByRole("heading", { name: "Welcome to Listello" }),
          ).toBeInTheDocument();
        });
      });

      it("shows the Welcome heading", () => {
        // Assert
        expect(
          screen.getByRole("heading", { name: "Welcome to Listello" }),
        ).toBeInTheDocument();
      });
    });
  });
});

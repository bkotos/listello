import { render } from "@testing-library/react";
import { createElement, type ReactElement } from "react";
import { MemoryRouter, Outlet, Route, Routes } from "react-router-dom";
import { vi } from "vitest";
import { createQueryWrapper } from "./renderWithQueryClient";

type RenderPageWithShellContextOptions = {
  path: string;
  initialEntry: string;
};

export function renderPageWithShellContext(
  page: ReactElement,
  { path, initialEntry }: RenderPageWithShellContextOptions,
) {
  const openSidebar = vi.fn();

  const { QueryWrapper } = createQueryWrapper();

  function Shell() {
    return createElement(Outlet, { context: { openSidebar } });
  }

  const view = render(
    createElement(
      QueryWrapper,
      null,
      createElement(
        MemoryRouter,
        { initialEntries: [initialEntry] },
        createElement(
          Routes,
          null,
          createElement(
            Route,
            { element: createElement(Shell) },
            createElement(Route, { path, element: page }),
          ),
        ),
      ),
    ),
  );

  return { ...view, openSidebar };
}

---
name: create-ui-component
description: >-
  Scaffolds and extends React UI in ui/src/components/ and ui/src/pages/ with
  Vitest + Testing Library using TDD. Use when adding or changing ItemRow,
  hover-reveal actions, dropdowns, mockup HTML, component specs, ListPage,
  or frontend UI following Listello architecture.
---

# Create UI Component

Guide for adding or changing React UI in `ui/src/components/` and `ui/src/pages/`. Read this skill before writing component or page specs.

Architecture context: see [README.md](../../../README.md). Layer order: [LAYER-ORDER.md](../LAYER-ORDER.md). TDD workflow: see [.cursor/rules/tdd.mdc](../../rules/tdd.mdc). Mockup: `mockup/task-management-system-bulma/`.

## Upstream dependencies

Presentational markup and interaction on **existing props** can proceed immediately.

**Stop and do not proceed** if the UI needs data or mutations that do not exist yet. See [LAYER-ORDER.md](../LAYER-ORDER.md).

| Need | How to verify | If missing |
|------|---------------|------------|
| New read/write from the API | Client function in `ui/src/lib/api/{resource}-client.ts` | [create-api-client](../create-api-client/SKILL.md) |
| Query/mutation hook | Hook in `ui/src/lib/api/{resource}-queries.ts` | [create-api-queries](../create-api-queries/SKILL.md) |
| DTO fields on the item/list | Type in `api-types` | [create-api-handler](../create-api-handler/SKILL.md) + `make api-types` |

Do not stub missing API in the component. Do not implement backend in this skill.

## Scope

**In scope:** `{Component}.tsx`, `{Component}.spec.ts`, page files that render the component, CSS in `ui/src/index.css` required by asserted classes (`hover-parent`, `hover-reveal`, `icon-btn`, …). Same-file child components after green.

**Out of scope** (mention as follow-ups only; do not implement unless asked):

- API clients (`{resource}-client.ts`) — use [create-api-client](../create-api-client/SKILL.md)
- React Query hooks — use [create-api-queries](../create-api-queries/SKILL.md)
- Go handlers, application, domain
- The Next.js mockup app (`mockup/`) — read it; do not edit it to “make tests pass”

## Architecture constraints

- Specs live next to the component: `ItemRow.tsx` + `ItemRow.spec.ts`.
- Specs use **Testing Library** (`render`, `screen`, `fireEvent`) and **`createElement`**, not JSX.
- One behavioral concern per `it`. When the user pastes mockup/Chrome HTML, split every visible control into its own test under a shared `describe` + `beforeEach`.
- Assert **roles, accessible names, classes, and lucide icon class names** — not `outerHTML` and not SVG path `d` attributes.
- `hover-reveal` is CSS opacity. Those nodes are in the DOM without `mouseEnter`. Do not simulate hover to assert structure.
- Do not add `@testing-library/user-event`. Use `fireEvent` (`click`, `mouseDown` on `document.body` for outside-click).
- Production JSX may match mockup structure and inline styles. Specs do not need to assert every style.
- Copy interaction patterns from existing UI (`AccountMenu` for dropdowns) instead of inventing new ones.

## Decision tree

1. **User pasted mockup/Chrome HTML?** → Treat it as the visual spec. Translate to semantic tests (see below). Implement only what that HTML (or stated click behavior) includes — omit extra mockup features (e.g. “Move to Inbox”) not in the dump.
2. **Static chrome (hover actions, layout)?** → `describe` + `beforeEach` that **renders**. One `it` per control.
3. **Interaction (open, close, toggle)?** → Nested `describe("when …")` whose `beforeEach` renders **and** performs the shared act (e.g. click Task options). Extra acts stay in the `it`.
4. **New widget vs extend existing?** → Prefer extending `ItemRow` / `ListPage` / `Sidebar`. New file only when it is a new route-level page or a clearly separate widget.
5. **Subtree gained its own state (open menu, popover)?** → After green, extract a same-file component if the user asks or the parent is noisy.
6. **One TDD loop per behavior.** Hover chrome, then open menu, then click-outside, then click-again-to-close are **separate** red → stop → green cycles. Do not spec the next interaction in the same red.

## Scaffold checklist

```
Task progress:
- [ ] Verify upstream if the UI needs new API data (stop if not — see LAYER-ORDER.md)
- [ ] Read the existing component spec and mockup source (if HTML was pasted)
- [ ] Decide: new describe vs new file; static chrome vs interaction
- [ ] Write failing spec(s) — describe + beforeEach; one UI piece per it
- [ ] Run tests — confirm failure is missing behavior (not a broken harness)
- [ ] STOP — summarize each new spec and its failure for user review
- [ ] (After approval) Implement the bare minimum production code + CSS the spec requires
- [ ] Re-run tests — confirm green
- [ ] (After green, if asked) Extract same-file child component; stay green
```

## Naming conventions

| Artifact | Convention |
|----------|------------|
| Component file | `ui/src/components/{Name}.tsx` (pages in `ui/src/pages/`) |
| Component spec | `{Name}.spec.ts` beside the component |
| Exported component | `export function {Name}` |
| Same-file child | unexported `function {Name}` in the same file |
| Nested describe (static) | `"hover actions"`, `"completed styling"` |
| Nested describe (interaction) | `"when the {control} is clicked"` |
| Test title | `"renders a {control}"` / `"closes when …"` |
| Accessible name | Match `aria-label` / visible text from the mockup (`"Set date"`, `"Task options"`) |
| Callback props | `onComplete`, `onUncomplete` — `vi.fn()` in specs |

## Translating mockup HTML

Do **not** snapshot the pasted HTML. Map it like this:

| In the HTML | In the spec |
|-------------|-------------|
| `role="button"` + `aria-label="Set date"` | `screen.getByRole("button", { name: "Set date" })` |
| `class="icon-btn hover-reveal"` | `toHaveClass("icon-btn", "hover-reveal")` |
| `aria-expanded="true"` / `"false"` | `toHaveAttribute("aria-expanded", …)` |
| `<svg class="lucide lucide-calendar">` | `querySelector("svg.lucide-calendar")` |
| `dropdown is-right is-active` | `closest(".dropdown")` + `toHaveClass("is-right", "is-active")` |
| `<a role="menuitem">Rename</a>` | `getByRole("menuitem", { name: "Rename" })` |
| Hover-only visibility | Assert the nodes exist; add `.hover-reveal` CSS in green |

Put the **distinctive missing control** first in a test (e.g. `"Set date"`) so red fails for missing behavior, not a nearby class that already exists.

## Test templates

### Shared arrange — static chrome

```ts
describe("hover actions", () => {
  beforeEach(() => {
    render(
      createElement(ItemRow, {
        item: { ...baseItem, Title: "Reply to the venue about the offsite" },
        onComplete: vi.fn(),
        onUncomplete: vi.fn(),
      }),
    );
  });

  it("renders a Set date icon button", () => {
    // Assert
    const setDate = screen.getByRole("button", { name: "Set date" });
    expect(setDate).toHaveClass("icon-btn");
    expect(setDate).toHaveAttribute("title", "Set date");
    expect(setDate.querySelector("svg.lucide-calendar")).toBeInTheDocument();
  });
});
```

### Shared arrange + act — open then assert each piece

```ts
describe("when the Task options button is clicked", () => {
  beforeEach(() => {
    render(
      createElement(ItemRow, {
        item: { ...baseItem, Title: "Reply to the venue about the offsite" },
        onComplete: vi.fn(),
        onUncomplete: vi.fn(),
      }),
    );
    fireEvent.click(screen.getByRole("button", { name: "Task options" }));
  });

  it("opens the dropdown", () => {
    // Assert
    const trigger = screen.getByRole("button", { name: "Task options" });
    expect(trigger).toHaveAttribute("aria-expanded", "true");
    expect(trigger.closest(".dropdown")).toHaveClass("is-right", "is-active");
  });

  it("renders a Rename menuitem", () => {
    // Assert
    const rename = screen.getByRole("menuitem", { name: "Rename" });
    expect(rename).toHaveClass("dropdown-item", "is-flex", "is-align-items-center");
    expect(rename.querySelector("svg.lucide-pencil")).toBeInTheDocument();
  });
});
```

### Extra act in the `it` — close

```ts
it("closes when clicking outside", () => {
  // Act
  fireEvent.mouseDown(document.body);

  // Assert
  const trigger = screen.getByRole("button", { name: "Task options" });
  expect(trigger).toHaveAttribute("aria-expanded", "false");
  expect(trigger.closest(".dropdown")).not.toHaveClass("is-active");
  expect(screen.queryByRole("menu")).not.toBeInTheDocument();
});
```

Mirror `AccountMenu`: outside click is `fireEvent.mouseDown(document.body)`.

Keep `afterEach(cleanup)` at the file level. Reuse existing fixtures (`baseItem`) instead of inventing new DTO shapes.

## Code templates

### Hover-reveal CSS (`ui/src/index.css`)

Add when specs assert `hover-parent` / `hover-reveal` (from mockup `globals.css`):

```css
.hover-reveal {
  opacity: 0;
  transition: opacity 0.15s ease;
}

.hover-parent:hover .hover-reveal,
.hover-reveal:focus-within,
.hover-reveal.is-visible {
  opacity: 1;
}
```

### Dropdown trigger + menu

Match Bulma + `AccountMenu`: `dropdown is-right`, `dropdown-trigger`, `aria-haspopup="menu"`, `aria-expanded`, `icon-btn` + `hover-reveal` when closed / `is-active` when open. Render `dropdown-menu` only when open unless an existing component already always mounts it.

### Same-file extract (after green)

Move state that belongs to the subtree into an unexported function in the same file:

```tsx
<TaskOptions />

function TaskOptions() {
  const [menuOpen, setMenuOpen] = useState(false);
  // …
}
```

Do not extract during red. Do not create a new file unless the user asks.

## TDD workflow

Follow `.cursor/rules/tdd.mdc`:

1. **Red** — Write the failing spec first. Confirm the failure is missing behavior (missing role/name, not a harness error).
2. **Stop** — Summarize each new `it` and its failure. **Do not implement** until the user approves.
3. **Green** — Implement the **bare minimum** to pass. No Escape-to-close, menuitem `onClick`, extra menu items, or popovers the current spec does not assert.
4. **Refactor** — Only while green, and only if asked or needed to keep the change set small.

Do not write production UI in the same turn as a new failing spec.

If the user asks to split a lumped spec into `describe` / `beforeEach` / per-control tests, **restructure the spec, re-run red, and stop again**. That is still the red phase.

## Verification

```bash
cd ui && npm test -- src/components/{Component}.spec.ts

# If the component is used on a page
cd ui && npm test -- src/pages/{Page}.spec.ts
```

If browser tools are available and the change is user-visible, exercise the flow. If not, say so — specs are the substitute.

## Further reading

Annotated ItemRow TDD walkthrough: [examples.md](examples.md)

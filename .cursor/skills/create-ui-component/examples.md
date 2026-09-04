# UI Component Examples

Annotated references from the Listello UI TDD loop for `ItemRow`. Read when translating mockup HTML into specs or adding hover/dropdown chrome.

## 1. Mockup HTML → semantic specs (not a snapshot)

The user pastes Chrome HTML for hover chrome. Do not `expect(el.outerHTML).toBe(...)`.

**Source (abbreviated):** `task-row hover-parent` row, `Set date` / `Add comment` in `hover-reveal`, closed `Task options` ellipsis dropdown.

**Spec file:** `ui/src/components/ItemRow.spec.ts`

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

  it("renders the row as a hover-parent button", () => {
    const row = screen.getByText("Reply to the venue about the offsite").closest(".task-row");
    expect(row).toHaveClass("hover-parent", "is-flex", "p-3");
    expect(row).toHaveAttribute("role", "button");
    expect(row).toHaveAttribute("tabindex", "0");
  });

  it("renders a Set date icon button", () => {
    const setDate = screen.getByRole("button", { name: "Set date" });
    expect(setDate).toHaveClass("icon-btn");
    expect(setDate.querySelector("svg.lucide-calendar")).toBeInTheDocument();
  });

  it("renders an Add comment icon button", () => { /* lucide-message-square */ });

  it("groups date and comment actions in a hover-reveal row", () => {
    const setDate = screen.getByRole("button", { name: "Set date" });
    expect(setDate.parentElement).toHaveClass("is-inline-flex", "is-align-items-center", "hover-reveal");
  });

  it("renders a closed Task options dropdown", () => {
    const options = screen.getByRole("button", { name: "Task options" });
    expect(options).toHaveClass("icon-btn", "hover-reveal");
    expect(options).toHaveAttribute("aria-expanded", "false");
    expect(options.closest(".dropdown")).toHaveClass("is-right");
  });
});
```

Key points:

- One `it` per visible control. Shared `beforeEach` only renders.
- No `fireEvent.mouseEnter` — hover-reveal nodes are always in the DOM.
- Red failed with `Unable to find … name "Set date"` (missing behavior), not a class typo on the existing checkbox.

**Green (minimum):** add `hover-parent` / `role="button"` / `tabIndex={0}`, the meta row with Calendar + MessageSquare, closed ellipsis dropdown, and `.hover-reveal` CSS in `ui/src/index.css`. Do not wire date/comment click handlers until a later spec.

## 2. Open menu — shared click in `beforeEach`

**User HTML (abbreviated):** `dropdown is-right is-active`, trigger `aria-expanded="true"` + `icon-btn is-active`, menu with Rename (pencil), divider, Delete (trash).

The pasted dump had **no** “Move to Inbox”. Do not spec it.

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
    const trigger = screen.getByRole("button", { name: "Task options" });
    expect(trigger).toHaveAttribute("aria-expanded", "true");
    expect(trigger.closest(".dropdown")).toHaveClass("is-right", "is-active");
  });

  it("marks the trigger as active", () => {
    const trigger = screen.getByRole("button", { name: "Task options" });
    expect(trigger).toHaveClass("icon-btn", "is-active");
    expect(trigger).not.toHaveClass("hover-reveal");
  });

  it("renders a Rename menuitem", () => {
    const rename = screen.getByRole("menuitem", { name: "Rename" });
    expect(rename.querySelector("svg.lucide-pencil")).toBeInTheDocument();
  });

  it("renders a Delete menuitem", () => {
    const remove = screen.getByRole("menuitem", { name: "Delete" });
    expect(remove.querySelector("svg.lucide-trash-2")).toBeInTheDocument();
  });

  it("separates Rename and Delete with a divider", () => {
    const content = screen.getByRole("menu").querySelector(".dropdown-content");
    const children = [...(content?.children ?? [])];
    expect(children[0]).toHaveTextContent("Rename");
    expect(children[1]).toHaveClass("dropdown-divider");
    expect(children[2]).toHaveTextContent("Delete");
  });
});
```

Key points:

- Shared act (open) lives in `beforeEach`. Each `it` asserts one piece of the open DOM.
- Red: `aria-expanded` stayed `"false"`; no `menuitem` named `"Rename"`.
- Green: `useState` + click sets open; render menu only when open. No `onClick` on Rename/Delete until specced.

## 3. Close interactions — extra act in the `it`

Same `describe` as open (menu already opened in `beforeEach`). Each close behavior is its **own TDD cycle**.

```ts
it("closes when clicking outside", () => {
  fireEvent.mouseDown(document.body);
  const trigger = screen.getByRole("button", { name: "Task options" });
  expect(trigger).toHaveAttribute("aria-expanded", "false");
  expect(screen.queryByRole("menu")).not.toBeInTheDocument();
});

it("closes when the Task options button is clicked again", () => {
  fireEvent.click(screen.getByRole("button", { name: "Task options" }));
  expect(screen.getByRole("button", { name: "Task options" })).toHaveAttribute(
    "aria-expanded",
    "false",
  );
});
```

Key points:

- Outside click matches `AccountMenu`: `mousedown` on `document.body`, listener on `document` with a `ref` + `contains`.
- Toggle is `setMenuOpen((open) => !open)` — not `setMenuOpen(true)` only.
- Do not add Escape-to-close in the same green unless a spec asserts it.
- Clicking the trigger is inside the ref, so outside-click must not steal the toggle.

## 4. Production: same-file `TaskOptions`

**File:** `ui/src/components/ItemRow.tsx`

After the open-menu specs were green, extract the dropdown (and its `menuOpen` state) into an unexported `TaskOptions` in the same file. Parent `ItemRow` keeps complete/uncomplete. Specs still target the row; they do not import `TaskOptions`.

```tsx
export function ItemRow({ item, onComplete, onUncomplete }: ItemRowProps) {
  // complete toggle + title + hover actions …
  return (
    <div role="button" tabIndex={0} className="task-row hover-parent is-flex p-3">
      {/* … */}
      <TaskOptions />
    </div>
  );
}

function TaskOptions() {
  const [menuOpen, setMenuOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!menuOpen) return;
    const onClick = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setMenuOpen(false);
      }
    };
    document.addEventListener("mousedown", onClick);
    return () => document.removeEventListener("mousedown", onClick);
  }, [menuOpen]);

  return (
    <div ref={containerRef} className={`dropdown is-right${menuOpen ? " is-active" : ""}`}>
      {/* trigger + conditional dropdown-menu */}
    </div>
  );
}
```

Key points:

- Extract **after** green, in a separate turn if the user asks to refactor.
- Stay green — re-run `ItemRow.spec.ts`.
- Copy dropdown a11y from `AccountMenu` (`aria-haspopup`, `aria-expanded`, `role="menu"` / `menuitem`).

## 5. What not to copy from the mockup

`mockup/task-management-system-bulma/components/item-row.tsx` is richer than a single HTML dump: rename-in-place, date popover, comments, Move to Inbox, Escape, click-outside-to-close-popovers.

Only implement what the **current** spec asserts. Later TDD loops add those behaviors one at a time.

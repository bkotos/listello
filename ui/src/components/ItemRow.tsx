import { useEffect, useRef, useState, type ReactNode } from "react";
import { Calendar, Check, Ellipsis, MessageSquare, Pencil, Trash2 } from "lucide-react";
import type { ItemDto } from "api-types";

export type ItemRowProps = {
  item: ItemDto;
  onComplete: (itemId: string) => void;
  onUncomplete: (itemId: string) => void;
  onDelete: (itemId: string) => void;
  onModifyTitle: (itemId: string, title: string) => void;
};

function isComplete(item: ItemDto): boolean {
  return item.State === "complete";
}

export function ItemRow({ item, onComplete, onUncomplete, onDelete, onModifyTitle }: ItemRowProps) {
  const [optimisticallyComplete, setOptimisticallyComplete] = useState(false);
  const [optimisticallyDeleted, setOptimisticallyDeleted] = useState(false);
  const [optimisticTitle, setOptimisticTitle] = useState<string | null>(null);
  const [renaming, setRenaming] = useState(false);
  const completed = isComplete(item) || optimisticallyComplete;
  const title = optimisticTitle ?? item.Title;

  function handleToggle() {
    if (!completed) {
      setOptimisticallyComplete(true);
      onComplete(item.ID);
      return;
    }
    setOptimisticallyComplete(false);
    onUncomplete(item.ID);
  }

  if (optimisticallyDeleted) {
    return null;
  }

  if (renaming) {
    return (
      <RenamingItemRow
        title={title}
        completed={completed}
        onToggle={handleToggle}
        onCommit={(nextTitle) => {
          setOptimisticTitle(nextTitle);
          setRenaming(false);
          onModifyTitle(item.ID, nextTitle);
        }}
        onCancel={() => setRenaming(false)}
      />
    );
  }

  return (
    <ItemRowShell>
      <CompleteToggle completed={completed} onToggle={handleToggle} />
      <div style={{ minWidth: 0, flex: "1 1 0" }}>
        <div className="is-flex" style={{ gap: "0.5rem", alignItems: "flex-start" }}>
          <p className={`is-size-6 mb-0${completed ? " muted line-through" : ""}`}>
            {title}
          </p>
        </div>
        <HoverActions />
      </div>
      <TaskOptions
        itemId={item.ID}
        onRename={() => setRenaming(true)}
        onDelete={(itemId) => {
          setOptimisticallyDeleted(true);
          onDelete(itemId);
        }}
      />
    </ItemRowShell>
  );
}

type ItemRowShellProps = {
  children: ReactNode;
};

function ItemRowShell({ children }: ItemRowShellProps) {
  return (
    <div
      role="button"
      tabIndex={0}
      className="task-row hover-parent is-flex p-3"
      style={{ gap: "0.75rem", alignItems: "flex-start" }}
    >
      {children}
    </div>
  );
}

type CompleteToggleProps = {
  completed: boolean;
  onToggle: () => void;
};

function CompleteToggle({ completed, onToggle }: CompleteToggleProps) {
  return (
    <button
      type="button"
      aria-label={completed ? "Mark incomplete" : "Mark complete"}
      className={`check-toggle${completed ? " is-checked" : ""}`}
      style={{ height: "1.25rem", width: "1.25rem", marginTop: "0.125rem" }}
      onClick={onToggle}
    >
      {completed ? <Check size={12} strokeWidth={2} aria-hidden /> : null}
    </button>
  );
}

type RenamingItemRowProps = {
  title: string;
  completed: boolean;
  onToggle: () => void;
  onCommit: (title: string) => void;
  onCancel: () => void;
};

function RenamingItemRow({ title, completed, onToggle, onCommit, onCancel }: RenamingItemRowProps) {
  const [draft, setDraft] = useState(title);

  return (
    <ItemRowShell>
      <CompleteToggle completed={completed} onToggle={onToggle} />
      <div style={{ minWidth: 0, flex: "1 1 0" }}>
        <div className="is-flex" style={{ gap: "0.5rem", alignItems: "flex-start" }}>
          <input
            className="input is-small"
            value={draft}
            autoFocus
            onChange={(e) => setDraft(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") {
                onCommit(draft);
              } else if (e.key === "Escape") {
                onCancel();
              }
            }}
          />
        </div>
      </div>
    </ItemRowShell>
  );
}

function HoverActions() {
  return (
    <div
      className="is-flex is-flex-wrap-wrap is-align-items-center mt-2 is-size-7 muted"
      style={{ gap: "0.5rem" }}
    >
      <span
        className="is-inline-flex is-align-items-center hover-reveal"
        style={{ gap: "0.25rem" }}
      >
        <button
          type="button"
          aria-label="Set date"
          title="Set date"
          className="icon-btn"
          style={{ height: "1.5rem", padding: "0 0.375rem", gap: "0.25rem" }}
        >
          <Calendar size={14} />
        </button>
        <button
          type="button"
          aria-label="Add comment"
          title="Add comment"
          className="icon-btn"
          style={{ height: "1.5rem", padding: "0 0.375rem", gap: "0.25rem" }}
        >
          <MessageSquare size={14} />
        </button>
      </span>
    </div>
  );
}

type TaskOptionsProps = {
  itemId: string;
  onRename: () => void;
  onDelete: (itemId: string) => void;
};

function TaskOptions({ itemId, onRename, onDelete }: TaskOptionsProps) {
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
    return () => {
      document.removeEventListener("mousedown", onClick);
    };
  }, [menuOpen]);

  return (
    <div
      ref={containerRef}
      className={`dropdown is-right${menuOpen ? " is-active" : ""}`}
    >
      <div className="dropdown-trigger">
        <button
          type="button"
          aria-label="Task options"
          aria-haspopup="menu"
          aria-expanded={menuOpen}
          className={`icon-btn ${menuOpen ? "is-active" : "hover-reveal"}`}
          style={{ height: "1.75rem", width: "1.75rem" }}
          onClick={() => setMenuOpen((open) => !open)}
        >
          <Ellipsis size={16} />
        </button>
      </div>
      {menuOpen ? (
        <div className="dropdown-menu" role="menu">
          <div className="dropdown-content">
            <a
              role="menuitem"
              className="dropdown-item is-flex is-align-items-center"
              style={{ gap: "0.625rem" }}
              onClick={onRename}
            >
              <span style={{ display: "inline-flex", color: "hsl(0deg 0% 45%)" }}>
                <Pencil size={16} />
              </span>
              <span>Rename</span>
            </a>
            <hr className="dropdown-divider" />
            <a
              role="menuitem"
              className="dropdown-item is-flex is-align-items-center"
              style={{ gap: "0.625rem", color: "hsl(348deg 86% 43%)" }}
              onClick={() => onDelete(itemId)}
            >
              <span style={{ display: "inline-flex", color: "hsl(348deg 86% 43%)" }}>
                <Trash2 size={16} />
              </span>
              <span>Delete</span>
            </a>
          </div>
        </div>
      ) : null}
    </div>
  );
}

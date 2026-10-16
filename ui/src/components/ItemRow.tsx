import { useEffect, useRef, useState, type ReactNode } from "react";
import { Calendar, Check, CornerDownRight, Ellipsis, MessageSquare, Pencil, Trash2 } from "lucide-react";
import type { ItemDto } from "api-types";

export type ItemRowProps = {
  item: ItemDto;
  selected: boolean;
  onSelect: (itemId: string) => void;
  onComplete: (itemId: string) => void;
  onUncomplete: (itemId: string) => void;
  onDelete: (itemId: string) => void;
  onModifyTitle: (itemId: string, title: string) => void;
  onLinkAsChild?: (itemId: string, parentId: string) => void;
  linkableParents?: { ID: string; Title: string }[];
};

function isComplete(item: ItemDto): boolean {
  return item.State === "complete";
}

export function ItemRow({ item, selected, onSelect, onComplete, onUncomplete, onDelete, onModifyTitle, onLinkAsChild, linkableParents }: ItemRowProps) {
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
        selected={selected}
        indented={Boolean(item.ParentID)}
        onSelect={() => onSelect(item.ID)}
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
    <ItemRowShell selected={selected} indented={Boolean(item.ParentID)} onSelect={() => onSelect(item.ID)}>
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
        onLinkAsChild={onLinkAsChild}
        linkableParents={linkableParents}
      />
    </ItemRowShell>
  );
}

type ItemRowShellProps = {
  children: ReactNode;
  selected: boolean;
  indented?: boolean;
  onSelect: () => void;
};

function ItemRowShell({ children, selected, indented, onSelect }: ItemRowShellProps) {
  return (
    <div
      role="button"
      tabIndex={0}
      className={`task-row hover-parent is-flex p-3${selected ? " is-active" : ""}${indented ? " is-child" : ""}`}
      style={{ gap: "0.75rem", alignItems: "flex-start" }}
      onClick={onSelect}
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
      onClick={(event) => {
        event.stopPropagation();
        onToggle();
      }}
    >
      {completed ? <Check size={12} strokeWidth={2} aria-hidden /> : null}
    </button>
  );
}

type RenamingItemRowProps = {
  title: string;
  selected: boolean;
  indented?: boolean;
  onSelect: () => void;
  completed: boolean;
  onToggle: () => void;
  onCommit: (title: string) => void;
  onCancel: () => void;
};

function RenamingItemRow({ title, selected, indented, onSelect, completed, onToggle, onCommit, onCancel }: RenamingItemRowProps) {
  const [draft, setDraft] = useState(title);

  return (
    <ItemRowShell selected={selected} indented={indented} onSelect={onSelect}>
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
            onBlur={() => onCommit(draft)}
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
  onLinkAsChild?: (itemId: string, parentId: string) => void;
  linkableParents?: { ID: string; Title: string }[];
};

function TaskOptions({ itemId, onRename, onDelete, onLinkAsChild, linkableParents = [] }: TaskOptionsProps) {
  const [menuOpen, setMenuOpen] = useState(false);
  const [pickingParent, setPickingParent] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);
  const canLinkAsChild = Boolean(onLinkAsChild) && linkableParents.length > 0;

  useEffect(() => {
    if (!menuOpen) return;
    const onClick = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setMenuOpen(false);
        setPickingParent(false);
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
          onClick={(event) => {
            event.stopPropagation();
            setMenuOpen((open) => {
              if (open) {
                setPickingParent(false);
              }
              return !open;
            });
          }}
        >
          <Ellipsis size={16} />
        </button>
      </div>
      {menuOpen ? (
        <div className="dropdown-menu" role="menu" onClick={(event) => event.stopPropagation()}>
          <div className="dropdown-content">
            {pickingParent
              ? linkableParents.map((parent) => (
                  <a
                    key={parent.ID}
                    role="menuitem"
                    className="dropdown-item is-flex is-align-items-center"
                    style={{ gap: "0.625rem" }}
                    onClick={() => {
                      onLinkAsChild?.(itemId, parent.ID);
                      setMenuOpen(false);
                      setPickingParent(false);
                    }}
                  >
                    <span>{parent.Title}</span>
                  </a>
                ))
              : (
                <>
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
                  {canLinkAsChild ? (
                    <a
                      role="menuitem"
                      className="dropdown-item is-flex is-align-items-center"
                      style={{ gap: "0.625rem" }}
                      onClick={() => setPickingParent(true)}
                    >
                      <span style={{ display: "inline-flex", color: "hsl(0deg 0% 45%)" }}>
                        <CornerDownRight size={16} />
                      </span>
                      <span>Link as child</span>
                    </a>
                  ) : null}
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
                </>
              )}
          </div>
        </div>
      ) : null}
    </div>
  );
}

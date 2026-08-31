import { useEffect, useRef, useState, type ReactNode } from "react";
import {
  BadgeCheck,
  Check,
  Hash,
  Inbox,
  LogOut,
  MoreHorizontal,
  Plus,
  Settings,
  User,
} from "lucide-react";
import { NavLink } from "react-router-dom";

const PLACEHOLDER_LISTS = [
  { id: "work", name: "Work" },
  { id: "personal", name: "Personal" },
  { id: "reading", name: "Reading" },
];

type SidebarProps = {
  onNavigate?: () => void;
};

export function Sidebar({ onNavigate }: SidebarProps) {
  return (
    <div
      className="is-flex is-flex-direction-column p-4"
      style={{ height: "100%", gap: "1.5rem" }}
    >
      <div className="is-flex is-align-items-center px-2 pt-1" style={{ gap: "0.5rem" }}>
        <span
          className="is-flex is-align-items-center is-justify-content-center has-background-primary"
          style={{ height: "1.75rem", width: "1.75rem", borderRadius: "0.5rem" }}
        >
          <Check size={16} strokeWidth={3} color="white" />
        </span>
        <span className="is-size-5 has-text-weight-semibold">Listello</span>
      </div>

      <aside className="menu">
        <ul className="menu-list">
          <li>
            <NavLink
              to="/inbox"
              className={({ isActive }) =>
                `is-flex is-align-items-center${isActive ? " is-active" : ""}`
              }
              style={{ gap: "0.625rem" }}
              onClick={onNavigate}
            >
              <Inbox size={16} />
              <span style={{ flex: "1 1 0" }}>Inbox</span>
            </NavLink>
          </li>
        </ul>
      </aside>

      <div
        className="is-flex is-flex-direction-column"
        style={{ flex: "1 1 0", minHeight: 0, gap: "0.25rem" }}
      >
        <div className="is-flex is-align-items-center is-justify-content-space-between px-2 pb-1">
          <span className="menu-label is-family-code mb-0" style={{ letterSpacing: "0.12em" }}>
            Lists
          </span>
          <button
            type="button"
            aria-label="Add list"
            className="icon-btn"
            style={{ height: "1.5rem", width: "1.5rem" }}
          >
            <Plus size={16} />
          </button>
        </div>

        <aside className="menu app-scroll">
          <ul className="menu-list">
            {PLACEHOLDER_LISTS.map((list) => (
              <li key={list.id}>
                <NavLink
                  to={`/lists/${list.id}`}
                  className={({ isActive }) =>
                    `is-flex is-align-items-center${isActive ? " is-active" : ""}`
                  }
                  style={{ gap: "0.625rem" }}
                  onClick={onNavigate}
                >
                  <Hash size={16} />
                  <span
                    style={{
                      flex: "1 1 0",
                      overflow: "hidden",
                      textOverflow: "ellipsis",
                      whiteSpace: "nowrap",
                    }}
                  >
                    {list.name}
                  </span>
                </NavLink>
              </li>
            ))}
          </ul>
        </aside>
      </div>

      <AccountMenu />
    </div>
  );
}

function AccountMenu() {
  const [open, setOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const onClick = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") setOpen(false);
    };
    document.addEventListener("mousedown", onClick);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("mousedown", onClick);
      document.removeEventListener("keydown", onKey);
    };
  }, [open]);

  return (
    <div
      ref={containerRef}
      className={`dropdown is-up is-right ${open ? "is-active" : ""}`}
      style={{ display: "block" }}
    >
      <div className="dropdown-menu" style={{ minWidth: "100%" }} role="menu">
        <div className="dropdown-content">
          <MenuItem icon={<User size={16} />} label="Profile" onClick={() => setOpen(false)} />
          <MenuItem icon={<Settings size={16} />} label="Settings" onClick={() => setOpen(false)} />
          <MenuItem
            icon={<BadgeCheck size={16} />}
            label="Delegation policy"
            onClick={() => setOpen(false)}
          />
          <hr className="dropdown-divider" />
          <MenuItem
            icon={<LogOut size={16} />}
            label="Sign out"
            destructive
            onClick={() => setOpen(false)}
          />
        </div>
      </div>

      <div
        className="dropdown-trigger is-flex is-align-items-center p-3"
        style={{
          gap: "0.5rem",
          width: "100%",
          borderRadius: "0.75rem",
          backgroundColor:
            "hsl(var(--bulma-primary-h) var(--bulma-primary-s) var(--bulma-primary-l) / 0.1)",
        }}
      >
        <div style={{ minWidth: 0, flex: "1 1 0" }}>
          <p
            className="is-family-code is-uppercase has-text-primary"
            style={{ fontSize: "0.6875rem", letterSpacing: "0.12em" }}
          >
            Signed in
          </p>
          <p className="is-size-7 has-text-weight-medium mt-1">You · Owner</p>
        </div>
        <button
          type="button"
          onClick={() => setOpen((v) => !v)}
          aria-label="Account options"
          aria-haspopup="menu"
          aria-expanded={open}
          className={`icon-btn ${open ? "is-active" : ""}`}
          style={{ height: "2rem", width: "2rem" }}
        >
          <MoreHorizontal size={16} />
        </button>
      </div>
    </div>
  );
}

function MenuItem({
  icon,
  label,
  onClick,
  destructive,
}: {
  icon: ReactNode;
  label: string;
  onClick: () => void;
  destructive?: boolean;
}) {
  return (
    <a
      role="menuitem"
      onClick={onClick}
      className="dropdown-item is-flex is-align-items-center"
      style={{ gap: "0.625rem", color: destructive ? "hsl(348deg 86% 43%)" : undefined }}
    >
      <span
        className="muted"
        style={{
          color: destructive ? "hsl(348deg 86% 43%)" : undefined,
          display: "inline-flex",
        }}
      >
        {icon}
      </span>
      <span>{label}</span>
    </a>
  );
}

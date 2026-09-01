import { Check, Hash, Inbox, Plus } from "lucide-react";
import { NavLink } from "react-router-dom";
import { AccountMenu } from "./AccountMenu.tsx";

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

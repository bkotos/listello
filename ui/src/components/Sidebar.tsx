import { useQueryClient } from "@tanstack/react-query";
import { Check, Hash, Inbox, Plus, Trash2 } from "lucide-react";
import { useState } from "react";
import { NavLink } from "react-router-dom";
import { useAppContext } from "../contexts/useAppContext";
import { deleteList } from "../lib/api/list-client";
import { listQueryKeys, useCreateListMutation } from "../lib/api/list-queries";
import { AccountMenu } from "./AccountMenu";

type SidebarProps = {
  onNavigate?: () => void;
};

export function Sidebar({ onNavigate }: SidebarProps) {
  const { lists } = useAppContext();
  const queryClient = useQueryClient();
  const { mutate: createList } = useCreateListMutation();
  const [adding, setAdding] = useState(false);
  const [newList, setNewList] = useState("");

  const submitList = () => {
    const name = newList.trim();
    if (!name) {
      return;
    }

    createList(name);
    setNewList("");
    setAdding(false);
    onNavigate?.();
  };

  async function handleDelete(listId: string) {
    await deleteList(listId);
    await queryClient.invalidateQueries({ queryKey: listQueryKeys.all });
  }

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
            onClick={() => setAdding((value) => !value)}
          >
            <Plus size={16} />
          </button>
        </div>

        <aside className="menu app-scroll">
          <ul className="menu-list">
            {lists.map((list) => (
              <li key={list.ID} className="hover-parent" style={{ position: "relative" }}>
                <NavLink
                  to={`/lists/${list.ID}`}
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
                    {list.Name}
                  </span>
                </NavLink>
                <button
                  type="button"
                  aria-label={`Delete ${list.Name}`}
                  className="icon-btn is-danger-hover hover-reveal"
                  style={{
                    position: "absolute",
                    right: "0.5rem",
                    top: "50%",
                    transform: "translateY(-50%)",
                    height: "1.5rem",
                    width: "1.5rem",
                  }}
                  onClick={() => {
                    void handleDelete(list.ID);
                  }}
                >
                  <Trash2 size={14} />
                </button>
              </li>
            ))}
          </ul>

          {adding && (
            <input
              autoFocus
              value={newList}
              placeholder="List name"
              className="input is-small mt-2"
              onChange={(event) => setNewList(event.target.value)}
              onKeyDown={(event) => {
                if (event.key === "Enter") {
                  submitList();
                }
                if (event.key === "Escape") {
                  setAdding(false);
                  setNewList("");
                }
              }}
              onBlur={() => {
                setAdding(false);
                setNewList("");
              }}
            />
          )}
        </aside>
      </div>

      <AccountMenu />
    </div>
  );
}

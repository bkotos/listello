import type { ClipboardEvent, KeyboardEvent, ReactNode } from "react";
import type { ItemDto } from "api-types";
import {
  Calendar,
  ChevronRight,
  CornerDownRight,
  Flag,
  Plus,
  Send,
  Tag,
  Trash2,
  X,
} from "lucide-react";

export type ItemDetailProps = {
  item: ItemDto;
  listName: string;
  onClose: () => void;
  onModifyTitle: (itemId: string, title: string) => void;
};

export function ItemDetail({ item, listName, onClose, onModifyTitle }: ItemDetailProps) {
  return (
    <aside className="app-detail is-hidden-touch">
      <div
        className="is-flex is-flex-direction-column"
        style={{ height: "100%", backgroundColor: "hsl(0deg 0% 100%)" }}
      >
        <header
          className="is-flex is-align-items-center is-justify-content-space-between px-5 py-3"
          style={{ gap: "0.5rem", borderBottom: "1px solid var(--app-border)" }}
        >
          <div className="is-flex is-align-items-center is-size-7 muted" style={{ gap: "0.5rem" }}>
            <span className="tag is-primary is-light is-family-code is-uppercase">{listName}</span>
            <span className="is-family-code is-uppercase" style={{ letterSpacing: "0.12em" }}>
              unclarified
            </span>
          </div>
          <div className="is-flex is-align-items-center" style={{ gap: "0.25rem" }}>
            <button
              type="button"
              aria-label="Delete item"
              className="icon-btn is-danger-hover"
              style={{ height: "2rem", width: "2rem" }}
            >
              <Trash2 size={16} />
            </button>
            <button
              type="button"
              aria-label="Close"
              className="icon-btn"
              style={{ height: "2rem", width: "2rem" }}
              onClick={onClose}
            >
              <X size={16} />
            </button>
          </div>
        </header>

        <div className="app-scroll">
          <div className="is-flex is-flex-direction-column px-5 py-5" style={{ gap: "1.5rem" }}>
            <div className="is-flex" style={{ gap: "0.75rem", alignItems: "flex-start" }}>
              <button
                type="button"
                aria-label="Mark complete"
                className="check-toggle"
                style={{ height: "1.5rem", width: "1.5rem", marginTop: "0.25rem" }}
              />
              <textarea
                key={item.Title}
                rows={1}
                defaultValue={item.Title}
                className="textarea is-shadowless"
                onBlur={(e) => onModifyTitle(item.ID, e.currentTarget.value)}
                onKeyDown={preventNewlines}
                onPaste={stripPastedNewlines}
                style={{
                  flex: "1 1 0",
                  resize: "none",
                  border: "none",
                  padding: 0,
                  minHeight: "unset",
                  backgroundColor: "transparent",
                  fontSize: "1.125rem",
                  fontWeight: 600,
                  lineHeight: 1.3,
                  boxShadow: "none",
                }}
              />
            </div>

            <div className="box p-0" style={{ boxShadow: "none", border: "1px solid var(--app-border)" }}>
              <PropertyRow icon={<Calendar size={16} />} label="Due date" divider={false}>
                <div className="is-flex is-align-items-center" style={{ gap: "0.5rem" }}>
                  <input className="input is-small" type="date" style={{ width: "auto" }} />
                </div>
              </PropertyRow>
              <PropertyRow icon={<Flag size={16} />} label="Priority">
                <div className="is-flex is-flex-wrap-wrap is-align-items-center" style={{ gap: "0.375rem" }}>
                  <button type="button" className="button is-small is-primary is-light" style={{ gap: "0.375rem" }}>
                    <span
                      style={{
                        height: "0.5rem",
                        width: "0.5rem",
                        borderRadius: "9999px",
                        backgroundColor: "hsl(0deg 0% 60%)",
                      }}
                    />
                    No priority
                  </button>
                  <button type="button" className="button is-small is-white" style={{ gap: "0.375rem" }}>
                    <span
                      style={{
                        height: "0.5rem",
                        width: "0.5rem",
                        borderRadius: "9999px",
                        backgroundColor: "hsl(0deg 0% 45%)",
                      }}
                    />
                    Low
                  </button>
                  <button type="button" className="button is-small is-white" style={{ gap: "0.375rem" }}>
                    <span
                      style={{
                        height: "0.5rem",
                        width: "0.5rem",
                        borderRadius: "9999px",
                        backgroundColor: "hsl(38deg 92% 50%)",
                      }}
                    />
                    Medium
                  </button>
                  <button type="button" className="button is-small is-white" style={{ gap: "0.375rem" }}>
                    <span
                      style={{
                        height: "0.5rem",
                        width: "0.5rem",
                        borderRadius: "9999px",
                        backgroundColor: "hsl(348deg 86% 55%)",
                      }}
                    />
                    High
                  </button>
                </div>
              </PropertyRow>
              <PropertyRow icon={<ChevronRight size={16} />} label="List">
                <div className="select is-small">
                  <select>
                    <option>Inbox</option>
                  </select>
                </div>
              </PropertyRow>
              <PropertyRow icon={<Tag size={16} />} label="Tags">
                <div className="is-flex is-flex-wrap-wrap is-align-items-center" style={{ gap: "0.375rem" }}>
                  <button
                    type="button"
                    className="button is-small is-white"
                    style={{ gap: "0.25rem", borderStyle: "dashed", borderColor: "var(--app-border)" }}
                  >
                    <Plus size={12} /> add
                  </button>
                </div>
              </PropertyRow>
            </div>

            <Section title="Notes">
              <textarea rows={3} className="textarea" placeholder="Add notes, links, or context…" />
            </Section>

            <Section title="Subtasks">
              <div className="is-flex is-flex-direction-column" style={{ gap: "0.25rem" }}>
                <div className="is-flex is-align-items-center px-1 py-1 muted" style={{ gap: "0.5rem" }}>
                  <CornerDownRight size={16} style={{ flexShrink: 0 }} />
                  <input
                    className="input is-small is-shadowless"
                    placeholder="Add a subtask…"
                    style={{ border: "none", paddingLeft: 0, backgroundColor: "transparent" }}
                  />
                </div>
              </div>
            </Section>

            <Section title="Comments">
              <div className="is-flex is-flex-direction-column" style={{ gap: "0.75rem" }}>
                <div className="field has-addons mb-0">
                  <div className="control is-expanded">
                    <input className="input is-small" placeholder="Write a comment…" />
                  </div>
                  <div className="control">
                    <button
                      type="button"
                      aria-label="Send comment"
                      disabled
                      className="button is-small is-primary"
                    >
                      <Send size={16} />
                    </button>
                  </div>
                </div>
              </div>
            </Section>
          </div>
        </div>
      </div>
    </aside>
  );
}

function preventNewlines(event: KeyboardEvent<HTMLTextAreaElement>) {
  if (event.key === "Enter") {
    event.preventDefault();
  }
}

function stripPastedNewlines(event: ClipboardEvent<HTMLTextAreaElement>) {
  event.preventDefault();
  event.currentTarget.value = (event.clipboardData?.getData("text") ?? "").replace(/\n/g, "");
}

type PropertyRowProps = {
  icon: ReactNode;
  label: string;
  children: ReactNode;
  divider?: boolean;
};

function PropertyRow({ icon, label, children, divider = true }: PropertyRowProps) {
  return (
    <div
      className="is-flex px-3 py-3"
      style={{
        gap: "0.75rem",
        alignItems: "center",
        borderTop: divider ? "1px solid var(--app-border)" : undefined,
      }}
    >
      <div
        className="is-flex is-align-items-center is-size-7 muted"
        style={{ width: "7rem", flexShrink: 0, gap: "0.5rem" }}
      >
        <span style={{ display: "inline-flex" }}>{icon}</span>
        {label}
      </div>
      <div style={{ minWidth: 0, flex: "1 1 0" }}>{children}</div>
    </div>
  );
}

type SectionProps = {
  title: string;
  children: ReactNode;
};

function Section({ title, children }: SectionProps) {
  return (
    <section className="is-flex is-flex-direction-column" style={{ gap: "0.625rem" }}>
      <div className="is-flex is-align-items-center" style={{ gap: "0.5rem" }}>
        <h2
          className="is-family-code is-uppercase muted mb-0"
          style={{ fontSize: "0.6875rem", letterSpacing: "0.12em" }}
        >
          {title}
        </h2>
      </div>
      {children}
    </section>
  );
}

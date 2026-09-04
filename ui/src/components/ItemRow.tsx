import { useState } from "react";
import { Calendar, Check, Ellipsis, MessageSquare } from "lucide-react";
import type { ItemDto } from "api-types";

export type ItemRowProps = {
  item: ItemDto;
  onComplete: (itemId: string) => void;
  onUncomplete: (itemId: string) => void;
};

function isComplete(item: ItemDto): boolean {
  return item.State === "complete";
}

export function ItemRow({ item, onComplete, onUncomplete }: ItemRowProps) {
  const [optimisticallyComplete, setOptimisticallyComplete] = useState(false);
  const completed = isComplete(item) || optimisticallyComplete;

  return (
    <div
      role="button"
      tabIndex={0}
      className="task-row hover-parent is-flex p-3"
      style={{ gap: "0.75rem", alignItems: "flex-start" }}
    >
      <button
        type="button"
        aria-label={completed ? "Mark incomplete" : "Mark complete"}
        className={`check-toggle${completed ? " is-checked" : ""}`}
        style={{ height: "1.25rem", width: "1.25rem", marginTop: "0.125rem" }}
        onClick={() => {
          if (!completed) {
            setOptimisticallyComplete(true);
            onComplete(item.ID);
            return;
          }
          setOptimisticallyComplete(false);
          onUncomplete(item.ID);
        }}
      >
        {completed ? <Check size={12} strokeWidth={2} aria-hidden /> : null}
      </button>

      <div style={{ minWidth: 0, flex: "1 1 0" }}>
        <div className="is-flex" style={{ gap: "0.5rem", alignItems: "flex-start" }}>
          <p
            className={`is-size-6 mb-0${completed ? " muted line-through" : ""}`}
          >
            {item.Title}
          </p>
        </div>
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
      </div>

      <div className="dropdown is-right">
        <div className="dropdown-trigger">
          <button
            type="button"
            aria-label="Task options"
            aria-haspopup="menu"
            aria-expanded="false"
            className="icon-btn hover-reveal"
            style={{ height: "1.75rem", width: "1.75rem" }}
          >
            <Ellipsis size={16} />
          </button>
        </div>
      </div>
    </div>
  );
}

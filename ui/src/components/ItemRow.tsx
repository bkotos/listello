import { useState } from "react";
import { Check } from "lucide-react";
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
      className="task-row is-flex p-3"
      style={{ gap: "0.75rem", alignItems: "flex-start" }}
    >
      <button
        type="button"
        aria-label={completed ? "Mark incomplete" : "Mark complete"}
        className={`check-toggle${completed ? " is-checked" : ""}`}
        style={{ height: "1.25rem", width: "1.25rem", marginTop: "0.125rem" }}
        onClick={
          !completed
            ? () => {
                setOptimisticallyComplete(true);
                onComplete(item.ID);
              }
            : () => {
                onUncomplete(item.ID);
              }
        }
      >
        {completed ? <Check size={12} strokeWidth={2} aria-hidden /> : null}
      </button>

      <div style={{ minWidth: 0, flex: "1 1 0" }}>
        <p
          className={`is-size-6 mb-0${completed ? " muted line-through" : ""}`}
        >
          {item.Title}
        </p>
      </div>
    </div>
  );
}

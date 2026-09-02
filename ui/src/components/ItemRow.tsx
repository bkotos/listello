import type { ItemDto } from "api-types";

export type ItemRowProps = {
  item: ItemDto;
};

function CheckIcon() {
  return (
    <svg viewBox="0 0 12 12" width={12} height={12} fill="none" aria-hidden="true">
      <path
        d="M2.5 6.2 5 8.5 9.5 3.5"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}

function isComplete(item: ItemDto): boolean {
  return item.State === "complete";
}

export function ItemRow({ item }: ItemRowProps) {
  const completed = isComplete(item);

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
      >
        {completed ? <CheckIcon /> : null}
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

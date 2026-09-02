import type { ItemDto } from "api-types";

export type ItemRowProps = {
  item: ItemDto;
};

export function ItemRow({ item }: ItemRowProps) {
  return (
    <div
      className="task-row is-flex p-3"
      style={{ gap: "0.75rem", alignItems: "flex-start" }}
    >
      <button
        type="button"
        aria-label="Mark complete"
        className="check-toggle"
        style={{ height: "1.25rem", width: "1.25rem", marginTop: "0.125rem" }}
      />

      <div style={{ minWidth: 0, flex: "1 1 0" }}>
        <p className="is-size-6 mb-0">{item.Title}</p>
      </div>
    </div>
  );
}

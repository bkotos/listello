import { useState, type KeyboardEvent, type ReactNode } from "react";
import { PanelLeft, Plus } from "lucide-react";

export type ContentAreaProps = {
  title: string;
  count?: number;
  icon?: ReactNode;
  capturePlaceholder: string;
  onOpenSidebar?: () => void;
  onCaptureSubmit?: (title: string) => void;
  children: ReactNode;
};

export function ContentArea({
  title,
  count,
  icon,
  capturePlaceholder,
  onOpenSidebar,
  onCaptureSubmit,
  children,
}: ContentAreaProps) {
  const [captureValue, setCaptureValue] = useState("");

  function handleCaptureKeyDown(event: KeyboardEvent<HTMLInputElement>) {
    if (event.key !== "Enter") {
      return;
    }

    const title = captureValue.trim();
    if (!title || !onCaptureSubmit) {
      return;
    }

    onCaptureSubmit(title);
    setCaptureValue("");
  }

  return (
    <div className="is-flex is-flex-direction-column" style={{ height: "100%" }}>
      <header
        className="is-flex is-align-items-center px-5 py-4"
        style={{ gap: "0.75rem", borderBottom: "1px solid var(--app-border)" }}
      >
        {onOpenSidebar && (
          <button
            type="button"
            onClick={onOpenSidebar}
            aria-label="Open menu"
            className="icon-btn is-hidden-desktop"
            style={{ height: "2.25rem", width: "2.25rem" }}
          >
            <PanelLeft size={20} />
          </button>
        )}
        <div className="is-flex is-align-items-center" style={{ gap: "0.5rem" }}>
          {icon}
          <h1 className="is-size-4 has-text-weight-semibold mb-0">{title}</h1>
          {count !== undefined && (
            <span
              className="is-family-code muted"
              aria-label={`${count} ${count === 1 ? "item" : "items"}`}
            >
              {count}
            </span>
          )}
        </div>
      </header>

      <div className="px-5 py-3" style={{ borderBottom: "1px solid var(--app-border)" }}>
        <div className="field mb-0">
          <div className="control has-icons-left">
            <input
              className="input is-medium"
              placeholder={capturePlaceholder}
              value={captureValue}
              onChange={(event) => setCaptureValue(event.target.value)}
              onKeyDown={handleCaptureKeyDown}
            />
            <span className="icon is-left">
              <Plus size={18} />
            </span>
          </div>
        </div>
      </div>

      <div className="app-scroll px-5 py-4">{children}</div>
    </div>
  );
}

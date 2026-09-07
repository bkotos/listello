import { ArrowLeft, ArrowRight, CircleCheck, Database, Globe, HardDrive } from "lucide-react";

export enum HostingMode {
  Local = "local",
  StandaloneWeb = "standalone-web",
}

type HostingStepProps = {
  hostingMode: HostingMode;
  onSelectHostingMode: (mode: HostingMode) => void;
};

export function HostingStep({
  hostingMode,
  onSelectHostingMode,
}: HostingStepProps) {
  const isLocal = hostingMode === HostingMode.Local;
  const isStandaloneWeb = hostingMode === HostingMode.StandaloneWeb;

  return (
    <div>
      <p className="step-eyebrow">Instance · Hosting</p>
      <h1 className="step-title text-balance">How should Listello be hosted?</h1>
      <p className="step-lead text-pretty">
        Your hosting mode decides where this instance runs and how your data is stored. You
        can migrate later.
      </p>
      <div className="step-content choice-list">
        <button
          type="button"
          className={`choice-card${isLocal ? " is-selected" : ""}`}
          aria-pressed={isLocal}
          onClick={() => onSelectHostingMode(HostingMode.Local)}
        >
          <span className="choice-icon">
            <HardDrive size={20} />
          </span>
          <span>
            <span className="choice-title is-block">Local</span>
            <span className="choice-desc">
              Runs directly on this machine, with local filesystem-backed persistence.
            </span>
            <span className="choice-meta">
              <Database size={12} />
              Uses SQLite
            </span>
          </span>
          {isLocal && (
            <span className="choice-check">
              <CircleCheck size={20} />
            </span>
          )}
        </button>
        <button
          type="button"
          className={`choice-card${isStandaloneWeb ? " is-selected" : ""}`}
          aria-pressed={isStandaloneWeb}
          onClick={() => onSelectHostingMode(HostingMode.StandaloneWeb)}
        >
          <span className="choice-icon">
            <Globe size={20} />
          </span>
          <span>
            <span className="choice-title is-block">Standalone Web</span>
            <span className="choice-desc">
              Runs independently in the browser as a PWA or standalone app.
            </span>
            <span className="choice-meta">
              <Database size={12} />
              Uses IndexedDB / OPFS
            </span>
          </span>
          {isStandaloneWeb && (
            <span className="choice-check">
              <CircleCheck size={20} />
            </span>
          )}
        </button>
      </div>
    </div>
  );
}

type HostingFooterProps = {
  onContinue: () => void;
};

export function HostingFooter({ onContinue }: HostingFooterProps) {
  return (
    <>
      <button type="button" className="button is-light">
        <span className="icon">
          <ArrowLeft size={18} />
        </span>
        <span>Back</span>
      </button>
      <button type="button" className="button is-primary footer-grow" onClick={onContinue}>
        <span>Continue</span>
        <span className="icon">
          <ArrowRight size={18} />
        </span>
      </button>
    </>
  );
}

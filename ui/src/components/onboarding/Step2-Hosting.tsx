import { ArrowLeft, ArrowRight, CircleCheck, Database, Globe, HardDrive } from "lucide-react";

export function HostingStep() {
  return (
    <div>
      <p className="step-eyebrow">Instance · Hosting</p>
      <h1 className="step-title text-balance">How should Listello be hosted?</h1>
      <p className="step-lead text-pretty">
        Your hosting mode decides where this instance runs and how your data is stored. You
        can migrate later.
      </p>
      <div className="step-content choice-list">
        <button type="button" className="choice-card is-selected" aria-pressed="true">
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
          <span className="choice-check">
            <CircleCheck size={20} />
          </span>
        </button>
        <button type="button" className="choice-card" aria-pressed="false">
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
        </button>
      </div>
    </div>
  );
}

export function HostingFooter() {
  return (
    <>
      <button type="button" className="button is-light">
        <span className="icon">
          <ArrowLeft size={18} />
        </span>
        <span>Back</span>
      </button>
      <button type="button" className="button is-primary footer-grow">
        <span>Continue</span>
        <span className="icon">
          <ArrowRight size={18} />
        </span>
      </button>
    </>
  );
}

import { ArrowLeft, ArrowRight, LoaderCircle } from "lucide-react";

export function InitializeStep() {
  return (
    <div>
      <p className="step-eyebrow">Instance · Initialize</p>
      <h1 className="step-title text-balance">Setting up persistence</h1>
      <p className="step-lead text-pretty">
        We're initializing the <strong>SQLite</strong> store for your{" "}
        <strong>local</strong> instance.
      </p>
      <div className="step-content">
        <div>
          <div className="setup-check">
            <span className="setup-check-status">
              <LoaderCircle size={18} className="spin" />
            </span>
            <span className="setup-check-label">Create instance</span>
          </div>
          <div className="setup-check is-pending">
            <span className="setup-check-status">
              <span className="check-toggle" style={{ width: 16, height: 16 }} />
            </span>
            <span className="setup-check-label">Prepare SQLite store</span>
          </div>
          <div className="setup-check is-pending">
            <span className="setup-check-status">
              <span className="check-toggle" style={{ width: 16, height: 16 }} />
            </span>
            <span className="setup-check-label">Initialize persistence</span>
          </div>
        </div>
      </div>
    </div>
  );
}

export function InitializeFooter() {
  return (
    <>
      <button type="button" className="button is-light">
        <span className="icon">
          <ArrowLeft size={18} />
        </span>
        <span>Back</span>
      </button>
      <button type="button" className="button is-primary footer-grow" disabled>
        <span>Continue</span>
        <span className="icon">
          <ArrowRight size={18} />
        </span>
      </button>
    </>
  );
}

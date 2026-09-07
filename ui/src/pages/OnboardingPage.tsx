import { ArrowRight, Check, Layers, ListChecks, Server, Sparkles } from "lucide-react";

function OnboardingPage() {
  return (
    <div className="onboarding-page">
      <header className="onboarding-topbar">
        <div className="onboarding-rail">
          <div className="brand-row">
            <span className="brand-mark">
              <span className="brand-dot">
                <Check size={18} strokeWidth={3} />
              </span>
              Listello
            </span>
            <span className="setup-badge">First-time setup</span>
          </div>

          <div className="phase-progress" aria-hidden="true">
            <div className="phase-seg is-active">
              <span className="phase-seg-head">
                <span className="phase-seg-index">1</span>
                <span className="phase-seg-label">Instance</span>
              </span>
              <span className="phase-seg-track">
                <span className="phase-seg-fill" style={{ width: "25%" }} />
              </span>
            </div>
            <div className="phase-seg is-upcoming">
              <span className="phase-seg-head">
                <span className="phase-seg-index">2</span>
                <span className="phase-seg-label">Workspace</span>
              </span>
              <span className="phase-seg-track">
                <span className="phase-seg-fill" style={{ width: "0%" }} />
              </span>
            </div>
            <div className="phase-seg is-upcoming">
              <span className="phase-seg-head">
                <span className="phase-seg-index">3</span>
                <span className="phase-seg-label">Ready</span>
              </span>
              <span className="phase-seg-track">
                <span className="phase-seg-fill" style={{ width: "0%" }} />
              </span>
            </div>
          </div>
        </div>
      </header>

      <div className="onboarding-body">
        <div className="onboarding-inner">
          <div>
            <span className="big-icon">
              <Sparkles size={26} />
            </span>
            <h1 className="step-title text-balance">Welcome to Listello</h1>
            <p className="step-lead text-pretty">
              Let's create your instance and set up a calm, GTD-style workspace. It
              only takes a minute, and you can change everything later.
            </p>
            <div className="step-content">
              <div className="setup-check is-done">
                <span className="setup-check-status">
                  <Server size={18} />
                </span>
                <span className="setup-check-label">Choose how Listello is hosted</span>
              </div>
              <div className="setup-check is-done">
                <span className="setup-check-status">
                  <Layers size={18} />
                </span>
                <span className="setup-check-label">Create your space and profile</span>
              </div>
              <div className="setup-check is-done">
                <span className="setup-check-status">
                  <ListChecks size={18} />
                </span>
                <span className="setup-check-label">Start your first list</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <footer className="onboarding-footer">
        <button type="button" className="button is-primary footer-grow">
          <span>Create instance</span>
          <span className="icon">
            <ArrowRight size={18} />
          </span>
        </button>
      </footer>
    </div>
  );
}

export default OnboardingPage;

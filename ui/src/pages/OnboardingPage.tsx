import { useState } from "react";
import {
  ArrowLeft,
  ArrowRight,
  Check,
  CircleCheck,
  Database,
  Globe,
  HardDrive,
  Layers,
  ListChecks,
  Server,
  Sparkles,
} from "lucide-react";
import { createInstance } from "../lib/api/instance-client";

function OnboardingPage() {
  const [instanceCreated, setInstanceCreated] = useState(false);

  async function handleCreateInstance() {
    await createInstance();
    setInstanceCreated(true);
  }

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
                <span
                  className="phase-seg-fill"
                  style={{ width: instanceCreated ? "50%" : "25%" }}
                />
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
          {instanceCreated ? (
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
          ) : (
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
          )}
        </div>
      </div>

      <footer className="onboarding-footer">
        {instanceCreated ? (
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
        ) : (
          <button type="button" className="button is-primary footer-grow" onClick={handleCreateInstance}>
            <span>Create instance</span>
            <span className="icon">
              <ArrowRight size={18} />
            </span>
          </button>
        )}
      </footer>
    </div>
  );
}

export default OnboardingPage;

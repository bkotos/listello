import { useState } from "react";
import { Check } from "lucide-react";
import { createInstance } from "../lib/api/instance-client";
import { HostingFooter, HostingStep } from "../components/onboarding/Step2-Hosting";
import { WelcomeFooter, WelcomeStep } from "../components/onboarding/Step1-Welcome";

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
          {instanceCreated ? <HostingStep /> : <WelcomeStep />}
        </div>
      </div>

      <footer className="onboarding-footer">
        {instanceCreated ? (
          <HostingFooter />
        ) : (
          <WelcomeFooter onCreateInstance={handleCreateInstance} />
        )}
      </footer>
    </div>
  );
}

export default OnboardingPage;

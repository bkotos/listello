import { useState } from "react";
import { Check } from "lucide-react";
import { createInstance } from "../lib/api/instance-client";
import { HostingFooter, HostingStep } from "../components/onboarding/Step2-Hosting";
import { WelcomeFooter, WelcomeStep } from "../components/onboarding/Step1-Welcome";

enum OnboardingStep {
  Welcome = "step1-welcome",
  Hosting = "step2-hosting",
}

function OnboardingPage() {
  const [step, setStep] = useState(OnboardingStep.Welcome);
  const [standaloneWebSelected, setStandaloneWebSelected] = useState(false);

  async function handleCreateInstance() {
    await createInstance();
    setStep(OnboardingStep.Hosting);
  }

  const isStep1Welcome = step === OnboardingStep.Welcome;
  const isStep2Hosting = step === OnboardingStep.Hosting;

  let instanceFill = "25%";
  if (isStep2Hosting) {
    instanceFill = standaloneWebSelected ? "67%" : "50%";
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
                  style={{
                    width: instanceFill,
                  }}
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
          {isStep1Welcome && <WelcomeStep />}
          {isStep2Hosting && (
            <HostingStep
              standaloneWebSelected={standaloneWebSelected}
              onSelectStandaloneWeb={() => setStandaloneWebSelected(true)}
            />
          )}
        </div>
      </div>

      <footer className="onboarding-footer">
        {isStep1Welcome && (
          <WelcomeFooter onCreateInstance={handleCreateInstance} />
        )}
        {isStep2Hosting && <HostingFooter />}
      </footer>
    </div>
  );
}

export default OnboardingPage;

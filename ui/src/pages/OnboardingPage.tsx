import { useState } from "react";
import { Check } from "lucide-react";
import { createInstance } from "../lib/api/instance-client";
import { HostingFooter, HostingMode, HostingStep } from "../components/onboarding/Step2-Hosting";
import {
  DataDirectoryFooter,
  DataDirectoryStep,
} from "../components/onboarding/Step3-DataDirectory";
import { WelcomeFooter, WelcomeStep } from "../components/onboarding/Step1-Welcome";

enum OnboardingStep {
  Welcome = "step1-welcome",
  Hosting = "step2-hosting",
  DataDirectory = "step3-data-directory",
}

function instancePhaseFill(step: OnboardingStep, hostingMode: HostingMode): string {
  if (step === OnboardingStep.DataDirectory) {
    return "75%";
  }
  if (step === OnboardingStep.Hosting) {
    return hostingMode === HostingMode.StandaloneWeb ? "67%" : "50%";
  }
  return "25%";
}

function OnboardingPage() {
  const [step, setStep] = useState(OnboardingStep.Welcome);
  const [hostingMode, setHostingMode] = useState(HostingMode.Local);

  async function handleCreateInstance() {
    await createInstance();
    setStep(OnboardingStep.Hosting);
  }

  function handleHostingContinue() {
    if (hostingMode === HostingMode.Local) {
      setStep(OnboardingStep.DataDirectory);
    }
  }

  const isStep1Welcome = step === OnboardingStep.Welcome;
  const isStep2Hosting = step === OnboardingStep.Hosting;
  const isStep3DataDirectory = step === OnboardingStep.DataDirectory;
  const instanceFill = instancePhaseFill(step, hostingMode);

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
              hostingMode={hostingMode}
              onSelectHostingMode={setHostingMode}
            />
          )}
          {isStep3DataDirectory && <DataDirectoryStep />}
        </div>
      </div>

      <footer className="onboarding-footer">
        {isStep1Welcome && (
          <WelcomeFooter onCreateInstance={handleCreateInstance} />
        )}
        {isStep2Hosting && <HostingFooter onContinue={handleHostingContinue} />}
        {isStep3DataDirectory && (
          <DataDirectoryFooter onBack={() => setStep(OnboardingStep.Hosting)} />
        )}
      </footer>
    </div>
  );
}

export default OnboardingPage;

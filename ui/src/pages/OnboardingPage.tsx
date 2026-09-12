import { useCallback, useState } from "react";
import { Check } from "lucide-react";
import { createInstance, initializePersistence, selectHostingMode, selectPersistenceLocation } from "../lib/api/instance-client";
import { useDefaultPersistenceLocationQuery } from "../lib/api/instance-queries";
import { createUser } from "../lib/api/user-client";
import { createSpace } from "../lib/api/space-client";
import { HostingFooter, HostingMode, HostingStep } from "../components/onboarding/Step2-Hosting";
import {
  DataDirectoryFooter,
  DataDirectoryStep,
} from "../components/onboarding/Step3-DataDirectory";
import { InitializeFooter, InitializeStep } from "../components/onboarding/Step4-Initialize";
import { NameSpaceFooter, NameSpaceStep } from "../components/onboarding/Step5-NameSpace";
import { YourNameFooter, YourNameStep } from "../components/onboarding/Step6-YourName";
import { SystemSetupFooter, SystemSetupStep } from "../components/onboarding/Step7-SystemSetup";
import { WelcomeFooter, WelcomeStep } from "../components/onboarding/Step1-Welcome";

enum OnboardingStep {
  Step1Welcome = "step1-welcome",
  Step2Hosting = "step2-hosting",
  Step3DataDirectory = "step3-data-directory",
  Step4Initialize = "step4-initialize",
  Step5NameSpace = "step5-name-space",
  Step6YourName = "step6-your-name",
  Step7SystemSetup = "step7-system-setup",
}

function instancePhaseFill(step: OnboardingStep, hostingMode: HostingMode): string {
  if (
    step === OnboardingStep.Step4Initialize ||
    step === OnboardingStep.Step5NameSpace ||
    step === OnboardingStep.Step6YourName ||
    step === OnboardingStep.Step7SystemSetup
  ) {
    return "100%";
  }
  if (step === OnboardingStep.Step3DataDirectory) {
    return "75%";
  }
  if (step === OnboardingStep.Step2Hosting) {
    return hostingMode === HostingMode.StandaloneWeb ? "67%" : "50%";
  }
  return "25%";
}

function workspacePhaseFill(step: OnboardingStep): string {
  if (step === OnboardingStep.Step7SystemSetup) {
    return "75%";
  }
  if (step === OnboardingStep.Step6YourName) {
    return "50%";
  }
  if (step === OnboardingStep.Step5NameSpace) {
    return "25%";
  }
  return "0%";
}

function OnboardingPage() {
  const { data: defaultPersistenceLocation } = useDefaultPersistenceLocationQuery();
  const [step, setStep] = useState(OnboardingStep.Step1Welcome);
  const [hostingMode, setHostingMode] = useState(HostingMode.Local);
  const [dataDirectoryOverride, setDataDirectoryOverride] = useState<string | null>(null);
  const [initializeComplete, setInitializeComplete] = useState(false);
  const [spaceName, setSpaceName] = useState("Personal");
  const [userName, setUserName] = useState("");
  const [systemSetupComplete, setSystemSetupComplete] = useState(false);

  const dataDirectory =
    dataDirectoryOverride ?? defaultPersistenceLocation?.Location ?? "";

  const handleInitializeComplete = useCallback(() => {
    setInitializeComplete(true);
  }, []);

  const handleSystemSetupComplete = useCallback(() => {
    setSystemSetupComplete(true);
  }, []);

  function handleInitializeBack() {
    setInitializeComplete(false);
    setStep(OnboardingStep.Step3DataDirectory);
  }

  async function handleCreateInstance() {
    await createInstance();
    setStep(OnboardingStep.Step2Hosting);
  }

  async function handleHostingContinue() {
    if (hostingMode === HostingMode.Local) {
      await selectHostingMode({ mode: hostingMode });
      setStep(OnboardingStep.Step3DataDirectory);
    }
  }

  async function handleDataDirectoryContinue() {
    setStep(OnboardingStep.Step4Initialize);
    await selectPersistenceLocation({ location: dataDirectory });
    await initializePersistence();
  }

  async function handleCreateUser() {
    await createUser({ name: userName });
    setStep(OnboardingStep.Step7SystemSetup);
  }

  async function handleCreateSpace() {
    await createSpace({ name: spaceName });
    setStep(OnboardingStep.Step6YourName);
  }

  const isStep1Welcome = step === OnboardingStep.Step1Welcome;
  const isStep2Hosting = step === OnboardingStep.Step2Hosting;
  const isStep3DataDirectory = step === OnboardingStep.Step3DataDirectory;
  const isStep4Initialize = step === OnboardingStep.Step4Initialize;
  const isStep5NameSpace = step === OnboardingStep.Step5NameSpace;
  const isStep6YourName = step === OnboardingStep.Step6YourName;
  const isStep7SystemSetup = step === OnboardingStep.Step7SystemSetup;
  const isWorkspacePhase = isStep5NameSpace || isStep6YourName || isStep7SystemSetup;
  const instanceFill = instancePhaseFill(step, hostingMode);
  const workspaceFill = workspacePhaseFill(step);

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
            <div className={`phase-seg ${isWorkspacePhase ? "is-done" : "is-active"}`}>
              <span className="phase-seg-head">
                <span className="phase-seg-index">
                  {isWorkspacePhase ? <Check size={12} strokeWidth={3} /> : "1"}
                </span>
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
            <div className={`phase-seg ${isWorkspacePhase ? "is-active" : "is-upcoming"}`}>
              <span className="phase-seg-head">
                <span className="phase-seg-index">2</span>
                <span className="phase-seg-label">Workspace</span>
              </span>
              <span className="phase-seg-track">
                <span
                  className="phase-seg-fill"
                  style={{ width: workspaceFill }}
                />
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
          {isStep3DataDirectory && (
            <DataDirectoryStep
              value={dataDirectory}
              onChange={setDataDirectoryOverride}
            />
          )}
          {isStep4Initialize && (
            <InitializeStep onComplete={handleInitializeComplete} />
          )}
          {isStep5NameSpace && <NameSpaceStep value={spaceName} onChange={setSpaceName} />}
          {isStep6YourName && <YourNameStep value={userName} onChange={setUserName} />}
          {isStep7SystemSetup && <SystemSetupStep spaceName={spaceName} userName={userName} onComplete={handleSystemSetupComplete} />}
        </div>
      </div>

      <footer className="onboarding-footer">
        {isStep1Welcome && (
          <WelcomeFooter onCreateInstance={handleCreateInstance} />
        )}
        {isStep2Hosting && (
          <HostingFooter
            onBack={() => setStep(OnboardingStep.Step1Welcome)}
            onContinue={handleHostingContinue}
          />
        )}
        {isStep3DataDirectory && (
          <DataDirectoryFooter
            onBack={() => setStep(OnboardingStep.Step2Hosting)}
            onContinue={handleDataDirectoryContinue}
          />
        )}
        {isStep4Initialize && (
          <InitializeFooter
            continueEnabled={initializeComplete}
            onBack={handleInitializeBack}
            onContinue={() => setStep(OnboardingStep.Step5NameSpace)}
          />
        )}
        {isStep5NameSpace && (
          <NameSpaceFooter
            onContinue={handleCreateSpace}
          />
        )}
        {isStep6YourName && <YourNameFooter continueEnabled={userName.length > 0} onContinue={handleCreateUser} />}
        {isStep7SystemSetup && (
          <SystemSetupFooter
            continueEnabled={systemSetupComplete}
            onContinue={() => {}}
          />
        )}
      </footer>
    </div>
  );
}

export default OnboardingPage;

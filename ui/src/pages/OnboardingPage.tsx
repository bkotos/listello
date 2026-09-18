import { useCallback, useEffect, useState } from "react";
import type { ListelloInstanceResponse } from "api-types/listello-instance";
import { Check } from "lucide-react";
import { createInstance, initializePersistence, pairSpace, pairUser, selectHostingMode, selectPersistenceLocation } from "../lib/api/instance-client";
import { createList } from "../lib/api/list-client";
import { useDefaultPersistenceLocationQuery, useInstanceQuery } from "../lib/api/instance-queries";
import { HostingFooter, HostingMode, HostingStep } from "../components/onboarding/Step2-Hosting";
import {
  DataDirectoryFooter,
  DataDirectoryStep,
} from "../components/onboarding/Step3-DataDirectory";
import { InitializeFooter, InitializeStep } from "../components/onboarding/Step4-Initialize";
import { NameSpaceFooter, NameSpaceStep } from "../components/onboarding/Step5-NameSpace";
import { YourNameFooter, YourNameStep } from "../components/onboarding/Step6-YourName";
import { SystemSetupFooter, SystemSetupStep } from "../components/onboarding/Step7-SystemSetup";
import { FirstListFooter, FirstListStep } from "../components/onboarding/Step8-FirstList";
import { CompleteFooter, CompleteStep } from "../components/onboarding/Step9-Complete";
import { WelcomeFooter, WelcomeStep } from "../components/onboarding/Step1-Welcome";

enum OnboardingStep {
  Step1Welcome = "step1-welcome",
  Step2Hosting = "step2-hosting",
  Step3DataDirectory = "step3-data-directory",
  Step4Initialize = "step4-initialize",
  Step5NameSpace = "step5-name-space",
  Step6YourName = "step6-your-name",
  Step7SystemSetup = "step7-system-setup",
  Step8FirstList = "step8-first-list",
  Step9Complete = "step9-complete",
}

function instancePhaseFill(step: OnboardingStep, hostingMode: HostingMode): string {
  if (
    step === OnboardingStep.Step4Initialize ||
    step === OnboardingStep.Step5NameSpace ||
    step === OnboardingStep.Step6YourName ||
    step === OnboardingStep.Step7SystemSetup ||
    step === OnboardingStep.Step8FirstList ||
    step === OnboardingStep.Step9Complete
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
  if (
    step === OnboardingStep.Step8FirstList ||
    step === OnboardingStep.Step9Complete
  ) {
    return "100%";
  }
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

function getDataDirectory(
  dataDirectoryOverride: string | null,
  instance: ListelloInstanceResponse | null | undefined,
  defaultPersistenceLocation: string | undefined,
): string {
  return (
    dataDirectoryOverride ??
    (instance?.PersistenceLocation || defaultPersistenceLocation) ??
    ""
  );
}

function onboardingStepFromInstance(
  instance: ListelloInstanceResponse | null | undefined,
): OnboardingStep | undefined {
  if (instance?.Space?.ID) {
    return OnboardingStep.Step6YourName;
  }
  if (instance?.PersistenceState === "initialized") {
    return OnboardingStep.Step5NameSpace;
  }
  if (instance?.PersistenceLocation) {
    return OnboardingStep.Step3DataDirectory;
  }
  if (instance?.HostingMode) {
    return OnboardingStep.Step2Hosting;
  }
}

function OnboardingPage() {
  const { data: instance } = useInstanceQuery();
  const { data: defaultPersistenceLocation } = useDefaultPersistenceLocationQuery();
  const [step, setStep] = useState(OnboardingStep.Step1Welcome);
  const [hostingMode, setHostingMode] = useState(HostingMode.Local);
  const [dataDirectoryOverride, setDataDirectoryOverride] = useState<string | null>(null);
  const [initializeComplete, setInitializeComplete] = useState(false);
  const [spaceName, setSpaceName] = useState("Personal");
  const [pairedSpaceName, setPairedSpaceName] = useState<string | null>(null);
  const [userName, setUserName] = useState("");
  const [listName, setListName] = useState("Errands");
  const [systemSetupComplete, setSystemSetupComplete] = useState(false);

  useEffect(() => {
    const nextStep = onboardingStepFromInstance(instance);
    if (nextStep) {
      setStep(nextStep);
    }
    if (instance?.Space?.Name) {
      setSpaceName(instance.Space.Name);
    }
  }, [instance]);

  const dataDirectory = getDataDirectory(
    dataDirectoryOverride,
    instance,
    defaultPersistenceLocation?.Location,
  );

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
    await pairUser({ name: userName });
    setStep(OnboardingStep.Step7SystemSetup);
  }

  async function handleCreateSpace() {
    if (pairedSpaceName !== spaceName) {
      await pairSpace({ name: spaceName });
      setPairedSpaceName(spaceName);
    }
    setStep(OnboardingStep.Step6YourName);
  }

  async function handleCreateList() {
    await createList(listName);
    setStep(OnboardingStep.Step9Complete);
  }

  const isStep1Welcome = step === OnboardingStep.Step1Welcome;
  const isStep2Hosting = step === OnboardingStep.Step2Hosting;
  const isStep3DataDirectory = step === OnboardingStep.Step3DataDirectory;
  const isStep4Initialize = step === OnboardingStep.Step4Initialize;
  const isStep5NameSpace = step === OnboardingStep.Step5NameSpace;
  const isStep6YourName = step === OnboardingStep.Step6YourName;
  const isStep7SystemSetup = step === OnboardingStep.Step7SystemSetup;
  const isStep8FirstList = step === OnboardingStep.Step8FirstList;
  const isStep9Complete = step === OnboardingStep.Step9Complete;
  const isWorkspacePhase =
    isStep5NameSpace || isStep6YourName || isStep7SystemSetup || isStep8FirstList;
  const isReadyPhase = isStep9Complete;
  const isInstanceDone = isWorkspacePhase || isReadyPhase;
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
            <div className={`phase-seg ${isInstanceDone ? "is-done" : "is-active"}`}>
              <span className="phase-seg-head">
                <span className="phase-seg-index">
                  {isInstanceDone ? <Check size={12} strokeWidth={3} /> : "1"}
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
            <div
              className={`phase-seg ${isReadyPhase ? "is-done" : isWorkspacePhase ? "is-active" : "is-upcoming"}`}
            >
              <span className="phase-seg-head">
                <span className="phase-seg-index">
                  {isReadyPhase ? <Check size={12} strokeWidth={3} /> : "2"}
                </span>
                <span className="phase-seg-label">Workspace</span>
              </span>
              <span className="phase-seg-track">
                <span
                  className="phase-seg-fill"
                  style={{ width: workspaceFill }}
                />
              </span>
            </div>
            <div className={`phase-seg ${isReadyPhase ? "is-active" : "is-upcoming"}`}>
              <span className="phase-seg-head">
                <span className="phase-seg-index">3</span>
                <span className="phase-seg-label">Ready</span>
              </span>
              <span className="phase-seg-track">
                <span
                  className="phase-seg-fill"
                  style={{ width: isReadyPhase ? "100%" : "0%" }}
                />
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
          {isStep7SystemSetup && (
            <SystemSetupStep
              spaceName={spaceName}
              userName={userName}
              onComplete={handleSystemSetupComplete}
            />
          )}
          {isStep8FirstList && (
            <FirstListStep value={listName} onChange={setListName} />
          )}
          {isStep9Complete && (
            <CompleteStep
              userName={userName}
              dataDirectory={dataDirectory}
              spaceName={spaceName}
              firstListName={listName}
            />
          )}
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
            onBack={() => setStep(OnboardingStep.Step3DataDirectory)}
            onContinue={handleCreateSpace}
          />
        )}
        {isStep6YourName && (
          <YourNameFooter
            continueEnabled={userName.length > 0}
            onBack={() => setStep(OnboardingStep.Step5NameSpace)}
            onContinue={handleCreateUser}
          />
        )}
        {isStep7SystemSetup && (
          <SystemSetupFooter
            continueEnabled={systemSetupComplete}
            onBack={() => setStep(OnboardingStep.Step6YourName)}
            onContinue={() => setStep(OnboardingStep.Step8FirstList)}
          />
        )}
        {isStep8FirstList && (
          <FirstListFooter onContinue={handleCreateList} />
        )}
        {isStep9Complete && <CompleteFooter />}
      </footer>
    </div>
  );
}

export default OnboardingPage;

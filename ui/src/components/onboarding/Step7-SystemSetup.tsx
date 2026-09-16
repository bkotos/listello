import { useEffect, useState } from "react";
import { ArrowLeft, ArrowRight, CircleCheck, LoaderCircle } from "lucide-react";

type SystemSetupStepProps = {
  spaceName: string;
  userName: string;
  onComplete: () => void;
};

export function SystemSetupStep({ spaceName, userName, onComplete }: SystemSetupStepProps) {
  const [doneCount, setDoneCount] = useState(0);
  const checks = [
    "Create Inbox",
    `Assign ${spaceName} to ${userName}`,
  ];

  useEffect(() => {
    const first = setTimeout(() => {
      setDoneCount(1);
    }, 500);
    const second = setTimeout(() => {
      setDoneCount(2);
      onComplete();
    }, 1000);
    return () => {
      clearTimeout(first);
      clearTimeout(second);
    };
  }, [onComplete]);

  return (
    <div>
      <p className="step-eyebrow">Workspace · Automatic</p>
      <h1 className="step-title text-balance">Getting things ready</h1>
      <p className="step-lead text-pretty">
        Listello is wiring up the essentials for <strong>{spaceName}</strong>.
      </p>
      <div className="step-content">
        <div>
          {checks.map((label, i) => {
            const done = i < doneCount;
            const active = i === doneCount;
            return (
              <div
                key={label}
                className={`setup-check${done ? " is-done" : active ? "" : " is-pending"}`}
              >
                <span className="setup-check-status">
                  {done ? (
                    <CircleCheck size={18} />
                  ) : active ? (
                    <LoaderCircle size={18} className="spin" />
                  ) : (
                    <span className="check-toggle" style={{ width: 16, height: 16 }} />
                  )}
                </span>
                <span className="setup-check-label">{label}</span>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}

type SystemSetupFooterProps = {
  continueEnabled: boolean;
  onContinue: () => void;
};

export function SystemSetupFooter({ continueEnabled, onContinue }: SystemSetupFooterProps) {
  return (
    <>
      <button type="button" className="button is-light">
        <span className="icon">
          <ArrowLeft size={18} />
        </span>
        <span>Back</span>
      </button>
      <button type="button" className="button is-primary footer-grow" disabled={!continueEnabled} onClick={onContinue}>
        <span>Continue</span>
        <span className="icon">
          <ArrowRight size={18} />
        </span>
      </button>
    </>
  );
}

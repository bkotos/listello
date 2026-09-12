import { ArrowRight } from "lucide-react";
import { useEffect } from "react";

type SystemSetupStepProps = {
  spaceName: string;
  userName: string;
  onComplete: () => void;
};

export function SystemSetupStep({ spaceName, userName, onComplete }: SystemSetupStepProps) {
  useEffect(() => {
    const timer = setTimeout(() => {
      onComplete();
    }, 500);
    return () => clearTimeout(timer);
  }, [onComplete]);

  return (
    <div>
      <p className="step-eyebrow">Workspace · Automatic</p>
      <h1 className="step-title text-balance">Getting things ready</h1>
      <p className="step-lead text-pretty">
        Listello is wiring up the essentials for <strong>{spaceName}</strong>.
      </p>
      <div className="step-content">
        <div className="setup-check">
          <span className="setup-check-status">
            <span className="check-toggle" style={{ width: 16, height: 16 }} />
          </span>
          <span className="setup-check-label">Create Inbox</span>
        </div>
        <div className="setup-check">
          <span className="setup-check-status">
            <span className="check-toggle" style={{ width: 16, height: 16 }} />
          </span>
          <span className="setup-check-label">Assign {spaceName} to {userName}</span>
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
    <button type="button" className="button is-primary footer-grow" disabled={!continueEnabled} onClick={onContinue}>
      <span>Continue</span>
      <span className="icon">
        <ArrowRight size={18} />
      </span>
    </button>
  );
}

import { useEffect, useState } from "react";
import { ArrowLeft, ArrowRight, CircleCheck, LoaderCircle } from "lucide-react";

const checks = [
  "Create instance",
  "Prepare SQLite store",
  "Initialize persistence",
];

type InitializeStepProps = {
  onComplete: () => void;
};

export function InitializeStep({ onComplete }: InitializeStepProps) {
  const [doneCount, setDoneCount] = useState(0);

  useEffect(() => {
    const first = setTimeout(() => {
      setDoneCount(1);
    }, 500);
    const second = setTimeout(() => {
      setDoneCount(2);
    }, 1000);
    const third = setTimeout(() => {
      setDoneCount(3);
      onComplete();
    }, 1500);
    return () => {
      clearTimeout(first);
      clearTimeout(second);
      clearTimeout(third);
    };
  }, [onComplete]);

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

type InitializeFooterProps = {
  continueEnabled: boolean;
  onBack: () => void;
};

export function InitializeFooter({ continueEnabled, onBack }: InitializeFooterProps) {
  return (
    <>
      <button type="button" className="button is-light" onClick={onBack}>
        <span className="icon">
          <ArrowLeft size={18} />
        </span>
        <span>Back</span>
      </button>
      <button
        type="button"
        className="button is-primary footer-grow"
        disabled={!continueEnabled}
      >
        <span>Continue</span>
        <span className="icon">
          <ArrowRight size={18} />
        </span>
      </button>
    </>
  );
}

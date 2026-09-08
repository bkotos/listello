import { useEffect, useState } from "react";
import { ArrowLeft, ArrowRight, CircleCheck, LoaderCircle } from "lucide-react";

const checks = [
  "Create instance",
  "Prepare SQLite store",
  "Initialize persistence",
];

export function InitializeStep() {
  const [doneCount, setDoneCount] = useState(0);

  useEffect(() => {
    const timer = setTimeout(() => {
      setDoneCount(1);
    }, 2000);
    return () => clearTimeout(timer);
  }, []);

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

export function InitializeFooter() {
  return (
    <>
      <button type="button" className="button is-light">
        <span className="icon">
          <ArrowLeft size={18} />
        </span>
        <span>Back</span>
      </button>
      <button type="button" className="button is-primary footer-grow" disabled>
        <span>Continue</span>
        <span className="icon">
          <ArrowRight size={18} />
        </span>
      </button>
    </>
  );
}

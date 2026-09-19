import { Inbox, Rocket } from "lucide-react";

type CompleteStepProps = {
  userName: string;
  dataDirectory: string;
  spaceName: string;
  firstListName: string;
};

export function CompleteStep({
  userName,
  dataDirectory,
  spaceName,
  firstListName,
}: CompleteStepProps) {
  return (
    <div>
      <span className="big-icon">
        <Rocket size={26} />
      </span>
      <h1 className="step-title text-balance">You're all set, {userName}</h1>
      <p className="step-lead text-pretty">
        Your instance is ready. Here is what we set up — you can change any of it later.
      </p>
      <div className="step-content box">
        <div className="summary-row">
          <span className="summary-key">Hosting</span>
          <span className="summary-val">Local</span>
        </div>
        <div className="summary-row">
          <span className="summary-key">Persistence</span>
          <span className="summary-val">SQLite</span>
        </div>
        <div className="summary-row">
          <span className="summary-key">Location</span>
          <span className="summary-val is-family-code">{dataDirectory}</span>
        </div>
        <div className="summary-row">
          <span className="summary-key">Space</span>
          <span className="summary-val">{spaceName}</span>
        </div>
        <div className="summary-row">
          <span className="summary-key">First list</span>
          <span className="summary-val">{firstListName}</span>
        </div>
        <div className="summary-row">
          <span className="summary-key">Inbox</span>
          <span className="summary-val">
            <span className="icon-text">
              <span className="icon has-text-primary">
                <Inbox size={16} />
              </span>
              <span>Ready</span>
            </span>
          </span>
        </div>
      </div>
    </div>
  );
}

type CompleteFooterProps = {
  onContinue: () => void;
};

export function CompleteFooter({ onContinue }: CompleteFooterProps) {
  return (
    <button type="button" className="button is-primary footer-grow" onClick={onContinue}>
      <span>Enter Listello</span>
      <span className="icon">
        <Rocket size={18} />
      </span>
    </button>
  );
}

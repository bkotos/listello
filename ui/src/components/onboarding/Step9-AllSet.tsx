import { Inbox, Rocket } from "lucide-react";

type AllSetStepProps = {
  userName: string;
  hostingLabel: string;
  persistence: string;
  location: string;
  spaceName: string;
  firstListName: string;
};

export function AllSetStep({
  userName,
  hostingLabel,
  persistence,
  location,
  spaceName,
  firstListName,
}: AllSetStepProps) {
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
          <span className="summary-val">{hostingLabel}</span>
        </div>
        <div className="summary-row">
          <span className="summary-key">Persistence</span>
          <span className="summary-val">{persistence}</span>
        </div>
        <div className="summary-row">
          <span className="summary-key">Location</span>
          <span className="summary-val is-family-code">{location}</span>
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

export function AllSetFooter() {
  return (
    <button type="button" className="button is-primary footer-grow">
      <span>Enter Listello</span>
      <span className="icon">
        <Rocket size={18} />
      </span>
    </button>
  );
}

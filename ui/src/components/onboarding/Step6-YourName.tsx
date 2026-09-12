import { ArrowLeft, ArrowRight, User } from "lucide-react";

export function YourNameStep() {
  return (
    <div>
      <p className="step-eyebrow">Workspace · You</p>
      <h1 className="step-title text-balance">What should we call you?</h1>
      <p className="step-lead text-pretty">
        Your name shows up on comments and activity. It's just for you — no account needed.
      </p>
      <div className="step-content">
        <div className="field">
          <label className="label" htmlFor="onb-user">
            Your name
          </label>
          <div className="control has-icons-left">
            <input
              id="onb-user"
              className="input"
              type="text"
              placeholder="e.g. Alex"
            />
            <span className="icon is-small is-left">
              <User size={16} />
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}

export function YourNameFooter() {
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

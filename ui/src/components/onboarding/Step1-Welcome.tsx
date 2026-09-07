import { ArrowRight, Layers, ListChecks, Server, Sparkles } from "lucide-react";

export function WelcomeStep() {
  return (
    <div>
      <span className="big-icon">
        <Sparkles size={26} />
      </span>
      <h1 className="step-title text-balance">Welcome to Listello</h1>
      <p className="step-lead text-pretty">
        Let's create your instance and set up a calm, GTD-style workspace. It
        only takes a minute, and you can change everything later.
      </p>
      <div className="step-content">
        <div className="setup-check is-done">
          <span className="setup-check-status">
            <Server size={18} />
          </span>
          <span className="setup-check-label">Choose how Listello is hosted</span>
        </div>
        <div className="setup-check is-done">
          <span className="setup-check-status">
            <Layers size={18} />
          </span>
          <span className="setup-check-label">Create your space and profile</span>
        </div>
        <div className="setup-check is-done">
          <span className="setup-check-status">
            <ListChecks size={18} />
          </span>
          <span className="setup-check-label">Start your first list</span>
        </div>
      </div>
    </div>
  );
}

type WelcomeFooterProps = {
  onCreateInstance: () => void;
};

export function WelcomeFooter({ onCreateInstance }: WelcomeFooterProps) {
  return (
    <button type="button" className="button is-primary footer-grow" onClick={onCreateInstance}>
      <span>Create instance</span>
      <span className="icon">
        <ArrowRight size={18} />
      </span>
    </button>
  );
}

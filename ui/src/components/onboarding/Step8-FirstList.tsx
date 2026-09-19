import { ArrowLeft, ArrowRight, ListChecks } from "lucide-react";

const suggestions = ["Errands", "Shopping", "Ideas", "Reading", "Goals"];

type FirstListStepProps = {
  value: string;
  onChange: (value: string) => void;
};

export function FirstListStep({ value, onChange }: FirstListStepProps) {
  return (
    <div>
      <p className="step-eyebrow">Workspace · First list</p>
      <h1 className="step-title text-balance">Create your first list</h1>
      <p className="step-lead text-pretty">
        Lists hold the tasks you want to act on. Give your first one a name, or skip and
        create one later.
      </p>
      <div className="step-content">
        <div className="field">
          <label className="label" htmlFor="onb-list">
            List name
          </label>
          <div className="control has-icons-left">
            <input
              id="onb-list"
              className="input"
              type="text"
              placeholder="Errands"
              value={value}
              onChange={(e) => onChange(e.target.value)}
            />
            <span className="icon is-small is-left">
              <ListChecks size={16} />
            </span>
          </div>
          <p className="suggest-hint">Or start with one of these</p>
          <div className="tags mt-2">
            {suggestions.map((name) => (
              <button
                key={name}
                type="button"
                className={`tag is-medium${name === value ? " is-primary" : ""}`}
                onClick={() => onChange(name)}
              >
                {name}
              </button>
            ))}
          </div>
        </div>
        <button type="button" className="skip-link">
          Skip for now
        </button>
      </div>
    </div>
  );
}

type FirstListFooterProps = {
  onBack: () => void;
  onContinue: () => void;
};

export function FirstListFooter({ onBack, onContinue }: FirstListFooterProps) {
  return (
    <>
      <button type="button" className="button is-light" onClick={onBack}>
        <span className="icon">
          <ArrowLeft size={18} />
        </span>
        <span>Back</span>
      </button>
      <button type="button" className="button is-primary footer-grow" onClick={onContinue}>
        <span>Create list</span>
        <span className="icon">
          <ArrowRight size={18} />
        </span>
      </button>
    </>
  );
}

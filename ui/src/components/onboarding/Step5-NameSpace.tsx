import { ArrowLeft, ArrowRight, Layers } from "lucide-react";

type NameSpaceStepProps = {
  value: string;
  onChange: (value: string) => void;
};

export function NameSpaceStep({ value, onChange }: NameSpaceStepProps) {
  return (
    <div>
      <p className="step-eyebrow">Workspace · Space</p>
      <h1 className="step-title text-balance">Name your space</h1>
      <p className="step-lead text-pretty">
        A space groups your lists together. Most people start with a single personal space.
      </p>
      <div className="step-content">
        <div className="field">
          <label className="label" htmlFor="onb-space">
            Space name
          </label>
          <div className="control has-icons-left">
            <input
              id="onb-space"
              className="input"
              type="text"
              placeholder="Personal"
              value={value}
              onChange={(e) => onChange(e.target.value)}
            />
            <span className="icon is-small is-left">
              <Layers size={16} />
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}

type NameSpaceFooterProps = {
  onContinue: () => void;
};

export function NameSpaceFooter({ onContinue }: NameSpaceFooterProps) {
  return (
    <>
      <button type="button" className="button is-light">
        <span className="icon">
          <ArrowLeft size={18} />
        </span>
        <span>Back</span>
      </button>
      <button type="button" className="button is-primary footer-grow" onClick={onContinue}>
        <span>Continue</span>
        <span className="icon">
          <ArrowRight size={18} />
        </span>
      </button>
    </>
  );
}

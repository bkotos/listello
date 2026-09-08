import { ArrowLeft, ArrowRight, Folder } from "lucide-react";

export const defaultDataDirectory =
  "/Users/jdoe/Library/Application Support/listello";

export function DataDirectoryStep() {
  return (
    <div>
      <p className="step-eyebrow">Instance · Local</p>
      <h1 className="step-title text-balance">Choose a data directory</h1>
      <p className="step-lead text-pretty">
        This is where Listello keeps its data for this instance. You can move it later.
      </p>
      <div className="step-content">
        <div className="field">
          <label className="label" htmlFor="onb-location">
            Data directory
          </label>
          <div className="control has-icons-left">
            <input
              id="onb-location"
              className="input"
              type="text"
              placeholder="~/listello"
              defaultValue={defaultDataDirectory}
            />
            <span className="icon is-small is-left">
              <Folder size={16} />
            </span>
          </div>
          <p className="help">Logs, the SQLite database, and config are stored here.</p>
        </div>
      </div>
    </div>
  );
}

type DataDirectoryFooterProps = {
  onBack: () => void;
  onContinue: () => void;
};

export function DataDirectoryFooter({ onBack, onContinue }: DataDirectoryFooterProps) {
  return (
    <>
      <button type="button" className="button is-light" onClick={onBack}>
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

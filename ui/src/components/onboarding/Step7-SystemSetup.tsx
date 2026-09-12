import { ArrowRight } from "lucide-react";

type SystemSetupStepProps = {
  spaceName: string;
  userName: string;
};

export function SystemSetupStep({ spaceName, userName }: SystemSetupStepProps) {
  return <div>Not implemented</div>;
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

import { Inbox } from "lucide-react";
import { ContentArea } from "../components/ContentArea";
import { useShellContext } from "../lib/useShellContext";

function InboxPage() {
  const { openSidebar } = useShellContext();

  return (
    <ContentArea
      title="Inbox"
      icon={<Inbox size={20} className="muted" />}
      capturePlaceholder="Capture something on your mind…"
      onOpenSidebar={openSidebar}
    >
      <p className="muted">No captured items yet.</p>
    </ContentArea>
  );
}

export default InboxPage;

import { useParams } from "react-router-dom";
import { ContentArea } from "../components/ContentArea.tsx";
import { useShellContext } from "../lib/useShellContext.ts";

const LIST_NAMES: Record<string, string> = {
  work: "Work",
  personal: "Personal",
  reading: "Reading",
};

function ListPage() {
  const { listId } = useParams();
  const { openSidebar } = useShellContext();
  const title = (listId && LIST_NAMES[listId]) ?? "List";

  return (
    <ContentArea
      title={title}
      capturePlaceholder={`Add to ${title}…`}
      onOpenSidebar={openSidebar}
    >
      <p className="muted">No items yet.</p>
    </ContentArea>
  );
}

export default ListPage;

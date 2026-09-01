import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { ContentArea } from "../components/ContentArea.tsx";
import { getList } from "../lib/api/list-client.ts";
import { useShellContext } from "../lib/useShellContext.ts";

function ListPage() {
  const { listId } = useParams();
  const { openSidebar } = useShellContext();
  const [title, setTitle] = useState("List");

  useEffect(() => {
    if (!listId) {
      return;
    }

    let cancelled = false;

    void getList(listId)
      .then((list) => {
        if (!cancelled) {
          setTitle(list.Name);
        }
      })
      .catch(() => {
        if (!cancelled) {
          setTitle("List");
        }
      });

    return () => {
      cancelled = true;
    };
  }, [listId]);

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

import { useEffect } from "react";
import { useParams } from "react-router-dom";
import { ContentArea } from "../components/ContentArea.tsx";
import { useAllItemsQuery } from "../lib/api/item-queries.ts";
import { useListQuery } from "../lib/api/list-queries.ts";
import { useShellContext } from "../lib/useShellContext.ts";

function ListPage() {
  const { listId } = useParams();
  const { openSidebar } = useShellContext();
  const { data: list } = useListQuery(listId);
  const { data: items } = useAllItemsQuery(listId);
  const title = list?.Name ?? "List";

  useEffect(() => {
    if (items) {
      console.log(items);
    }
  }, [items]);

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

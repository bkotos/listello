import { useParams } from "react-router-dom";
import { ContentArea } from "../components/ContentArea";
import { ItemRow } from "../components/ItemRow";
import { useAllItemsQuery } from "../lib/api/item-queries";
import { useListQuery } from "../lib/api/list-queries";
import { useShellContext } from "../lib/useShellContext";

function ListPage() {
  const { listId } = useParams();
  const { openSidebar } = useShellContext();
  const { data: list } = useListQuery(listId);
  const { data: items } = useAllItemsQuery(listId);
  const title = list?.Name ?? "List";

  return (
    <ContentArea
      title={title}
      count={items?.length}
      capturePlaceholder={`Add to ${title}…`}
      onOpenSidebar={openSidebar}
    >
      {items && items.length > 0 ? (
        <div
          className="is-flex is-flex-direction-column"
          style={{ maxWidth: "42rem", marginInline: "auto", gap: "1rem" }}
        >
          <div className="is-flex is-flex-direction-column" style={{ gap: "0.375rem" }}>
            {items.map((item) => (
              <ItemRow key={item.ID} item={item} />
            ))}
          </div>
        </div>
      ) : (
        <p className="muted">No items yet.</p>
      )}
    </ContentArea>
  );
}

export default ListPage;

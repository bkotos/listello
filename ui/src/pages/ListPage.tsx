import { useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { ContentArea } from "../components/ContentArea";
import { ItemDetail } from "../components/ItemDetail";
import { ItemRow } from "../components/ItemRow";
import { defineItem, completeItem, uncompleteItem, deleteItem, modifyItemTitle } from "../lib/api/item-client";
import { itemQueryKeys, useAllItemsQuery } from "../lib/api/item-queries";
import { useListQuery } from "../lib/api/list-queries";
import { useShellContext } from "../lib/useShellContext";

function ListPage() {
  const { listId } = useParams();
  const queryClient = useQueryClient();
  const { openSidebar } = useShellContext();
  const { data: list } = useListQuery(listId);
  const { data: items } = useAllItemsQuery(listId);
  const [selectedItemId, setSelectedItemId] = useState<string | null>(null);
  const title = list?.Name ?? "List";
  const selectedItem = items?.find((item) => item.ID === selectedItemId) ?? null;

  async function handleCaptureSubmit(itemTitle: string) {
    if (!listId) {
      return;
    }

    await defineItem(listId, { title: itemTitle });
    await queryClient.invalidateQueries({ queryKey: itemQueryKeys.byList(listId) });
  }

  async function handleComplete(itemId: string) {
    await completeItem(itemId);
    if (listId) {
      await queryClient.invalidateQueries({ queryKey: itemQueryKeys.byList(listId) });
    }
  }

  async function handleUncomplete(itemId: string) {
    await uncompleteItem(itemId);
    if (listId) {
      await queryClient.invalidateQueries({ queryKey: itemQueryKeys.byList(listId) });
    }
  }

  async function handleDelete(itemId: string) {
    await deleteItem(itemId);
    if (listId) {
      await queryClient.invalidateQueries({ queryKey: itemQueryKeys.byList(listId) });
    }
  }

  async function handleModifyTitle(itemId: string, itemTitle: string) {
    await modifyItemTitle(itemId, { title: itemTitle });
    if (listId) {
      await queryClient.invalidateQueries({ queryKey: itemQueryKeys.byList(listId) });
    }
  }

  return (
    <div className="is-flex" style={{ height: "100%", minWidth: 0 }}>
      <div className="is-flex is-flex-direction-column" style={{ flex: "1 1 0", minWidth: 0 }}>
        <ContentArea
          title={title}
          count={items?.length}
          capturePlaceholder={`Add to ${title}…`}
          onOpenSidebar={openSidebar}
          onCaptureSubmit={handleCaptureSubmit}
        >
          {items && items.length > 0 ? (
            <div
              className="is-flex is-flex-direction-column"
              style={{ maxWidth: "42rem", marginInline: "auto", gap: "1rem" }}
            >
              <div className="is-flex is-flex-direction-column" style={{ gap: "0.375rem" }}>
                {items.map((item) => (
                  <ItemRow
                    key={item.ID}
                    item={item}
                    selected={item.ID === selectedItemId}
                    onSelect={setSelectedItemId}
                    onComplete={handleComplete}
                    onUncomplete={handleUncomplete}
                    onDelete={handleDelete}
                    onModifyTitle={handleModifyTitle}
                  />
                ))}
              </div>
            </div>
          ) : (
            <p className="muted">No items yet.</p>
          )}
        </ContentArea>
      </div>
      {selectedItem ? (
        <ItemDetail item={selectedItem} listName={title} onClose={() => setSelectedItemId(null)} />
      ) : null}
    </div>
  );
}

export default ListPage;

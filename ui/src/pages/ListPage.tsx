import { useQueries, useQueryClient } from "@tanstack/react-query";
import { useNavigate, useParams } from "react-router-dom";
import { ContentArea } from "../components/ContentArea";
import { ItemDetail } from "../components/ItemDetail";
import { ItemRow } from "../components/ItemRow";
import { defineItem, completeItem, uncompleteItem, deleteItem, modifyItemTitle, commentItem, addSubtask, completeSubtask, uncompleteSubtask, deleteSubtask, linkAsChild, getAllItems, getItem } from "../lib/api/item-client";
import { itemQueryKeys, useAllItemsQuery } from "../lib/api/item-queries";
import { useAllListsQuery, useListQuery } from "../lib/api/list-queries";
import { useShellContext } from "../lib/useShellContext";

function ListPage() {
  const { listId, itemId } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { openSidebar } = useShellContext();
  const { data: list } = useListQuery(listId);
  const { data: lists } = useAllListsQuery();
  const { data: items } = useAllItemsQuery(listId);
  const missingParentIds = Array.from(
    new Set(
      (items ?? [])
        .map((item) => item.ParentID)
        .filter(
          (parentId) =>
            Boolean(parentId) && !(items ?? []).some((candidate) => candidate.ID === parentId),
        ),
    ),
  );
  const parentQueries = useQueries({
    queries: missingParentIds.map((id) => ({
      queryKey: itemQueryKeys.detail(id),
      queryFn: ({ signal }: { signal: AbortSignal }) => getItem(id, { signal }),
    })),
  });
  const otherListQueries = useQueries({
    queries: (lists ?? [])
      .filter((candidate) => candidate.ID !== listId)
      .map((candidate) => ({
        queryKey: itemQueryKeys.byList(candidate.ID),
        queryFn: ({ signal }: { signal: AbortSignal }) => getAllItems(candidate.ID, { signal }),
      })),
  });
  const fetchedParents = parentQueries.flatMap((query) => (query.data ? [query.data] : []));
  const otherListItems = otherListQueries.flatMap((query) => query.data ?? []);
  const title = list?.Name ?? "List";
  const selectedItem = items?.find((item) => item.ID === itemId) ?? null;

  function itemById(id: string | undefined) {
    if (!id) {
      return undefined;
    }
    return (
      (items ?? []).find((item) => item.ID === id) ??
      fetchedParents.find((item) => item.ID === id) ??
      otherListItems.find((item) => item.ID === id)
    );
  }

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

  async function handleComment(commentedItemId: string, body: string) {
    await commentItem(commentedItemId, { body });
    if (listId) {
      await queryClient.invalidateQueries({ queryKey: itemQueryKeys.byList(listId) });
    }
  }

  async function handleAddSubtask(itemId: string, title: string) {
    await addSubtask(itemId, { title });
    if (listId) {
      await queryClient.invalidateQueries({ queryKey: itemQueryKeys.byList(listId) });
    }
  }

  async function handleCompleteSubtask(itemId: string, subtaskId: string) {
    await completeSubtask(itemId, subtaskId);
    if (listId) {
      await queryClient.invalidateQueries({ queryKey: itemQueryKeys.byList(listId) });
    }
  }

  async function handleUncompleteSubtask(itemId: string, subtaskId: string) {
    await uncompleteSubtask(itemId, subtaskId);
    if (listId) {
      await queryClient.invalidateQueries({ queryKey: itemQueryKeys.byList(listId) });
    }
  }

  async function handleDeleteSubtask(itemId: string, subtaskId: string) {
    await deleteSubtask(itemId, subtaskId);
    if (listId) {
      await queryClient.invalidateQueries({ queryKey: itemQueryKeys.byList(listId) });
    }
  }

  async function handleLinkAsChild(itemId: string, parentId: string) {
    await linkAsChild(itemId, { parentID: parentId });
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
                    selected={item.ID === itemId}
                    onSelect={(selectedId) => {
                      const target = itemById(selectedId);
                      navigate(`/lists/${target?.ListID ?? listId}/items/${selectedId}`);
                    }}
                    onComplete={handleComplete}
                    onUncomplete={handleUncomplete}
                    onDelete={handleDelete}
                    onModifyTitle={handleModifyTitle}
                    onLinkAsChild={handleLinkAsChild}
                    parentTitle={itemById(item.ParentID)?.Title}
                    linkableParents={[...(items ?? []), ...otherListItems]
                      .filter((candidate) => candidate.ID !== item.ID && !candidate.ParentID)
                      .map((candidate) => ({ ID: candidate.ID, Title: candidate.Title }))}
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
        <ItemDetail
          key={selectedItem.ID}
          item={selectedItem}
          listName={title}
          onClose={() => navigate(`/lists/${listId}`)}
          onModifyTitle={handleModifyTitle}
          onComment={handleComment}
          onAddSubtask={handleAddSubtask}
          onCompleteSubtask={handleCompleteSubtask}
          onUncompleteSubtask={handleUncompleteSubtask}
          onDeleteSubtask={handleDeleteSubtask}
        />
      ) : null}
    </div>
  );
}

export default ListPage;

import type { AddSubtaskRequest, CommentItemRequest, DefineItemRequest, ItemCommentDto, ItemDto, LinkAsChildRequest, ModifyItemTitleRequest, MoveItemRequest, PrioritizeSubtaskRequest } from "api-types";
import { request } from "./util";

export async function getItem(
  id: string,
  init?: RequestInit,
): Promise<ItemDto> {
  return request<ItemDto>(`/api/items/${id}`, init);
}

export async function getAllItems(
  listId: string,
  init?: RequestInit,
): Promise<ItemDto[]> {
  return request<ItemDto[]>(`/api/lists/${listId}/items`, init);
}

export async function defineItem(
  listId: string,
  body: DefineItemRequest,
): Promise<ItemDto> {
  return request<ItemDto>(`/api/lists/${listId}/items`, {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export async function completeItem(itemId: string): Promise<ItemDto> {
  return request<ItemDto>(`/api/items/${itemId}/complete`, {
    method: "POST",
  });
}

export async function uncompleteItem(itemId: string): Promise<ItemDto> {
  return request<ItemDto>(`/api/items/${itemId}/uncomplete`, {
    method: "POST",
  });
}

export async function deleteItem(itemId: string): Promise<void> {
  return request<void>(`/api/items/${itemId}`, {
    method: "DELETE",
  });
}

export async function modifyItemTitle(
  itemId: string,
  body: ModifyItemTitleRequest,
): Promise<ItemDto> {
  return request<ItemDto>(`/api/items/${itemId}/title`, {
    method: "PATCH",
    body: JSON.stringify(body),
  });
}

export async function moveItem(
  itemId: string,
  body: MoveItemRequest,
): Promise<ItemDto> {
  return request<ItemDto>(`/api/items/${itemId}/move`, {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export async function getItemComments(
  itemId: string,
  init?: RequestInit,
): Promise<ItemCommentDto[]> {
  return request<ItemCommentDto[]>(`/api/items/${itemId}/comments`, init);
}

export async function commentItem(
  itemId: string,
  body: CommentItemRequest,
): Promise<ItemDto> {
  return request<ItemDto>(`/api/items/${itemId}/comment`, {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export async function deleteComment(
  itemId: string,
  commentId: string,
): Promise<ItemDto> {
  return request<ItemDto>(`/api/items/${itemId}/comment/${commentId}`, {
    method: "DELETE",
  });
}

export async function addSubtask(
  itemId: string,
  body: AddSubtaskRequest,
): Promise<ItemDto> {
  return request<ItemDto>(`/api/items/${itemId}/subtask`, {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export async function completeSubtask(
  itemId: string,
  subtaskId: string,
): Promise<ItemDto> {
  return request<ItemDto>(`/api/items/${itemId}/subtask/${subtaskId}/complete`, {
    method: "POST",
  });
}

export async function uncompleteSubtask(
  itemId: string,
  subtaskId: string,
): Promise<ItemDto> {
  return request<ItemDto>(`/api/items/${itemId}/subtask/${subtaskId}/uncomplete`, {
    method: "POST",
  });
}

export async function deleteSubtask(
  itemId: string,
  subtaskId: string,
): Promise<ItemDto> {
  return request<ItemDto>(`/api/items/${itemId}/subtask/${subtaskId}`, {
    method: "DELETE",
  });
}

export async function linkAsChild(
  itemId: string,
  body: LinkAsChildRequest,
): Promise<ItemDto> {
  return request<ItemDto>(`/api/items/${itemId}/link-as-child`, {
    method: "POST",
    body: JSON.stringify(body),
  });
}

export async function prioritizeSubtask(
  itemId: string,
  subtaskId: string,
  body: PrioritizeSubtaskRequest,
): Promise<ItemDto> {
  return request<ItemDto>(`/api/items/${itemId}/subtask/${subtaskId}/prioritize`, {
    method: "POST",
    body: JSON.stringify(body),
  });
}

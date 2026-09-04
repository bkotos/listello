import type { DefineItemRequest, ItemDto, ModifyItemTitleRequest } from "api-types";
import { request } from "./util";

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

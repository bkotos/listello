import type { ListResponse } from "api-types";
import { request } from "./util";

export async function getAllLists(init?: RequestInit): Promise<ListResponse[]> {
  return request<ListResponse[]>("/api/lists", init);
}

export async function getList(id: string, init?: RequestInit): Promise<ListResponse> {
  return request<ListResponse>(`/api/lists/${id}`, init);
}

export async function createList(name: string): Promise<ListResponse> {
  return request<ListResponse>("/api/lists", {
    method: "POST",
    body: JSON.stringify({ name }),
  });
}

export async function createFirstList(name: string): Promise<ListResponse> {
  return request<ListResponse>("/api/lists/first", {
    method: "POST",
    body: JSON.stringify({ name }),
  });
}

export async function deleteList(listId: string): Promise<void> {
  return request<void>(`/api/lists/${listId}`, {
    method: "DELETE",
  });
}

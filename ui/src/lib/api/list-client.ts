import type { ListResponse } from "api-types";
import { request } from "./util.ts";

export async function getAllLists(): Promise<ListResponse[]> {
  return request<ListResponse[]>("/api/lists");
}

export async function createList(name: string): Promise<ListResponse> {
  return request<ListResponse>("/api/lists", {
    method: "POST",
    body: JSON.stringify({ name }),
  });
}

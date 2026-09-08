import type { ListelloInstanceResponse, SelectPersistenceLocationRequest } from "api-types/listello-instance";
import { request } from "./util";

export async function createInstance(): Promise<ListelloInstanceResponse> {
  return request<ListelloInstanceResponse>("/api/instance", {
    method: "POST",
  });
}

export async function getInstance(init?: RequestInit): Promise<ListelloInstanceResponse | null> {
  return request<ListelloInstanceResponse | null>("/api/instance", init);
}

export async function selectPersistenceLocation(
  body: SelectPersistenceLocationRequest,
): Promise<ListelloInstanceResponse> {
  return request<ListelloInstanceResponse>("/api/instance/persistence-location", {
    method: "POST",
    body: JSON.stringify(body),
  });
}

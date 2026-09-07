import type { ListelloInstanceResponse } from "api-types/listello-instance";
import { request } from "./util";

export async function createInstance(): Promise<ListelloInstanceResponse> {
  return request<ListelloInstanceResponse>("/api/instance", {
    method: "POST",
  });
}

export async function getInstance(init?: RequestInit): Promise<ListelloInstanceResponse | null> {
  return request<ListelloInstanceResponse | null>("/api/instance", init);
}

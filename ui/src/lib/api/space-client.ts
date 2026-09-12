import type { SpaceResponse, CreateSpaceRequest } from "api-types";
import { request } from "./util";

export async function createSpace(body: CreateSpaceRequest): Promise<SpaceResponse> {
  return request<SpaceResponse>("/api/spaces", {
    method: "POST",
    body: JSON.stringify(body),
  });
}

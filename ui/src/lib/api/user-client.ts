import type { UserResponse, CreateUserRequest } from "api-types";
import { request } from "./util";

export async function createUser(body: CreateUserRequest): Promise<UserResponse> {
  return request<UserResponse>("/api/users", {
    method: "POST",
    body: JSON.stringify(body),
  });
}

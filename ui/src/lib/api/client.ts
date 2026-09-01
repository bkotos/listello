import type { ListResponse } from "api-types";

export class ApiError extends Error {
  readonly status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(path, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...init?.headers,
    },
  });

  if (!response.ok) {
    let message = response.statusText;
    try {
      const body = (await response.json()) as { error?: string };
      if (body.error) {
        message = body.error;
      }
    } catch {
      // Response body was not JSON.
    }
    throw new ApiError(response.status, message);
  }

  return response.json() as Promise<T>;
}

export async function getAllLists(): Promise<ListResponse[]> {
  return request<ListResponse[]>("/api/lists");
}

export async function createList(name: string): Promise<ListResponse> {
  return request<ListResponse>("/api/lists", {
    method: "POST",
    body: JSON.stringify({ name }),
  });
}

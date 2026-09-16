import type { ListelloInstanceResponse } from "api-types/listello-instance";

export function needsOnboarding(
  instance: ListelloInstanceResponse | null | undefined,
): boolean {
  return instance !== undefined && (instance === null || instance.SetupState !== "completed");
}

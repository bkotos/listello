import { describe, expect, it } from "vitest";
import { pluralize } from "./pluralize.ts";

describe("pluralize", () => {
  it("uses the singular form for a count of one", () => {
    expect(pluralize(1, "list")).toBe("1 list");
  });

  it("uses the default plural form for other counts", () => {
    expect(pluralize(0, "item")).toBe("0 items");
    expect(pluralize(3, "item")).toBe("3 items");
  });

  it("uses a custom plural form when provided", () => {
    expect(pluralize(2, "person", "people")).toBe("2 people");
  });
});

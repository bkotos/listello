// @vitest-environment node

import { describe, expect, it } from "vitest";
import { parsePort } from "./parsePort.ts";

describe("parsePort", () => {
  it("returns the fallback when value is undefined", () => {
    expect(parsePort(undefined, 3001)).toBe(3001);
  });

  it("parses a valid port string", () => {
    expect(parsePort("8080", 3001)).toBe(8080);
  });

  it("throws for invalid ports", () => {
    expect(() => parsePort("0", 3001)).toThrow("invalid port: 0");
    expect(() => parsePort("70000", 3001)).toThrow("invalid port: 70000");
  });
});

import { describe, expect, it } from "vitest";
import { runEngineCli } from "./engine-cli.js";

describe("engine CLI", () => {
  it("prints ok by default", async () => {
    const res = await runEngineCli();
    expect(res.exitCode).toBe(0);
    expect(res.stderr).toBe("");
    expect(res.stdout.trim()).toBe("engine: ok");
  });

  it("prints a version with --version", async () => {
    const res = await runEngineCli(["--version"]);
    expect(res.exitCode).toBe(0);
    expect(res.stderr).toBe("");
    expect(res.stdout.trim().length).toBeGreaterThan(0);
  });
});


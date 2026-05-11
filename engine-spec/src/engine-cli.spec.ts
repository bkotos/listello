import { readdirSync, readFileSync, mkdtempSync, rmSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { afterEach, beforeEach, describe, expect, it } from "vitest";
import { runEngineCli } from "./engine-cli.js";

/** Crockford base32: 26 chars, case-insensitive (see ideas/file-format.md). */
const ulidListelloFilenameRe = /^[0-9A-HJKMNP-TV-Z]{26}\.listello\.md$/i;

describe("engine CLI", () => {
  it("prints ok by default", async () => {
    // act
    const res = await runEngineCli();

    // assert
    expect(res.exitCode).toBe(0);
    expect(res.stderr).toBe("");
    expect(res.stdout.trim()).toBe("engine: ok");
  });

  it("prints a version with --version", async () => {
    // arrange
    const args = ["--version"];

    // act
    const res = await runEngineCli(args);

    // assert
    expect(res.exitCode).toBe(0);
    expect(res.stderr).toBe("");
    expect(res.stdout.trim().length).toBeGreaterThan(0);
  });
});

describe("listello create (CLI contract)", () => {
  let workDir: string;

  beforeEach(() => {
    workDir = mkdtempSync(join(tmpdir(), "listello-engine-create-"));
  });

  afterEach(() => {
    rmSync(workDir, { recursive: true, force: true });
  });

  async function runCreateBuyDiapers() {
    return runEngineCli(["create", "--title", "Buy diapers"], { cwd: workDir });
  }

  it("exits 0 with no stderr", async () => {
    // act
    const res = await runCreateBuyDiapers();

    // assert
    expect(res.exitCode).toBe(0);
    expect(res.stderr).toBe("");
  });

  it("creates exactly one {ULID}.listello.md in the working directory", async () => {
    // act
    await runCreateBuyDiapers();

    // assert
    const files = readdirSync(workDir);
    expect(files).toHaveLength(1);
    expect(files[0]).toMatch(ulidListelloFilenameRe);
  });

  it("writes YAML frontmatter with ---, title from --title, and id matching the filename stem", async () => {
    // act
    await runCreateBuyDiapers();

    // assert
    const [filename] = readdirSync(workDir);
    const content = readFileSync(join(workDir, filename), "utf8");
    expect(content.startsWith("---\n")).toBe(true);
    expect(content).toContain("title: Buy diapers");
    const ulid = filename.replace(/\.listello\.md$/i, "");
    expect(content).toContain(`id: ${ulid}`);
  });
});


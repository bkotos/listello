import fs from "node:fs/promises";
import { afterEach, beforeAll, beforeEach, describe, expect, it } from "vitest";
import { $ } from "zx";
import {
  createWorkdir,
  dbPathFor,
  parseCreatedListId,
  removeWorkdir,
  runListello,
} from "./support/listello-cli";

describe("listello cli", () => {
  let workdir: string;
  let dbPath: string;

  beforeAll(() => {
    $.verbose = false;
  });

  beforeEach(async () => {
    workdir = await createWorkdir();
    dbPath = dbPathFor(workdir);
  });

  afterEach(async () => {
    await removeWorkdir(workdir);
  });

  describe("list create", () => {
    it("creates a list in the given sqlite file", async () => {
      // Arrange
      const args = ["list", "create", "Groceries"];

      // Act
      const result = await runListello(dbPath, args);

      // Assert
      expect(result.exitCode).toBe(0);
      expect(result.stderr).toBe("");
      expect(result.stdout).toMatch(/Created list "Groceries" \(LS_/);
      await expect(fs.stat(dbPath)).resolves.toBeDefined();
    });
  });

  describe("item define", () => {
    it("defines an item on a list", async () => {
      // Arrange
      const createResult = await runListello(dbPath, [
        "list",
        "create",
        "Groceries",
      ]);
      const listId = parseCreatedListId(createResult.stdout);

      // Act
      const result = await runListello(dbPath, [
        "item",
        "define",
        listId,
        "Buy milk",
      ]);

      // Assert
      expect(result.exitCode).toBe(0);
      expect(result.stderr).toBe("");
      expect(result.stdout).toMatch(
        new RegExp(`^Defined item "Buy milk" \\(IT_[^)]+\\) on list ${listId}$`),
      );
    });
  });

  describe("item list", () => {
    it("lists items on a list", async () => {
      // Arrange
      const createResult = await runListello(dbPath, [
        "list",
        "create",
        "Groceries",
      ]);
      const listId = parseCreatedListId(createResult.stdout);
      await runListello(dbPath, ["item", "define", listId, "Buy milk"]);

      // Act
      const result = await runListello(dbPath, ["item", "list", listId]);

      // Assert
      expect(result.exitCode).toBe(0);
      expect(result.stderr).toBe("");
      expect(result.stdout).toMatch(/\[ \] IT_.+  Buy milk/);
    });
  });
});

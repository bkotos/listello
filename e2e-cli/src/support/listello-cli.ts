import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { $, fs } from "zx";

const repoRoot = path.resolve(
  fileURLToPath(new URL(".", import.meta.url)),
  "../../..",
);
const apiDir = path.join(repoRoot, "api");

export async function createWorkdir(): Promise<string> {
  return fs.mkdtemp(path.join(os.tmpdir(), "listello-e2e-"));
}

export function dbPathFor(workdir: string): string {
  return path.join(workdir, "test.db");
}

export async function runListello(
  dbPath: string,
  args: string[],
): Promise<{ stdout: string; stderr: string; exitCode: number }> {
  const result = await $({
    cwd: apiDir,
    quiet: true,
    nothrow: true,
  })`go run ./cmd/cli --db ${dbPath} ${args}`;

  return {
    stdout: result.stdout.trim(),
    stderr: result.stderr.trim(),
    exitCode: result.exitCode,
  };
}

export async function removeWorkdir(workdir: string): Promise<void> {
  await fs.rm(workdir, { recursive: true, force: true });
}

export function parseCreatedListId(output: string): string {
  const match = output.match(/Created list ".+" \((LS_[^)]+)\)/);
  if (!match) {
    throw new Error(`could not parse list id from output: ${output}`);
  }
  return match[1];
}

export function parseDefinedItemId(output: string): string {
  const match = output.match(/Defined item ".+" \((IT_[^)]+)\)/);
  if (!match) {
    throw new Error(`could not parse item id from output: ${output}`);
  }
  return match[1];
}

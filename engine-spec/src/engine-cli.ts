import { spawn } from "node:child_process";
import { fileURLToPath } from "node:url";
import path from "node:path";

export type EngineRunResult = {
  exitCode: number;
  stdout: string;
  stderr: string;
};

function repoRootFromHere(): string {
  // engine-spec/src/engine-cli.ts -> engine-spec/src -> engine-spec -> repo root
  const here = path.dirname(fileURLToPath(import.meta.url));
  return path.resolve(here, "..", "..");
}

export async function runEngineCli(
  args: string[] = [],
  opts?: { cwd?: string; env?: NodeJS.ProcessEnv; timeoutMs?: number },
): Promise<EngineRunResult> {
  const repoRoot = repoRootFromHere();
  const engineRoot = path.join(repoRoot, "engine");
  /** Process working directory (e.g. where `create` writes); Go module root is always `engineRoot` via `-C`. */
  const spawnCwd = opts?.cwd ?? engineRoot;
  const timeoutMs = opts?.timeoutMs ?? 30_000;

  return new Promise((resolve, reject) => {
    const child = spawn(
      "go",
      ["run", "-C", engineRoot, "./cmd/engine", "--", ...args],
      {
        cwd: spawnCwd,
        env: { ...process.env, ...opts?.env },
        stdio: ["ignore", "pipe", "pipe"],
      },
    );

    let stdout = "";
    let stderr = "";

    const timer = setTimeout(() => {
      child.kill("SIGKILL");
      reject(new Error(`engine CLI timed out after ${timeoutMs}ms`));
    }, timeoutMs);

    child.stdout?.setEncoding("utf8");
    child.stderr?.setEncoding("utf8");
    child.stdout?.on("data", (d) => (stdout += d));
    child.stderr?.on("data", (d) => (stderr += d));

    child.on("error", (err) => {
      clearTimeout(timer);
      reject(err);
    });

    child.on("close", (code) => {
      clearTimeout(timer);
      resolve({
        exitCode: code ?? 0,
        stdout,
        stderr,
      });
    });
  });
}


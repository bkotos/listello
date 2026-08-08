import { readFile } from "node:fs/promises";
import { createRequire } from "node:module";
import path from "node:path";
import { fileURLToPath } from "node:url";

import type { Listello } from "./types.js";

export type { Item, ItemAPI, List, ListAPI, Listello } from "./types.js";

const require = createRequire(import.meta.url);
const here = path.dirname(fileURLToPath(import.meta.url));

let loading: Promise<Listello> | undefined;

/**
 * Load the Go js/wasm module and return the listello API.
 * Safe to call multiple times; the module is only instantiated once.
 */
export async function loadListello(): Promise<Listello> {
  if (globalThis.listello) {
    return globalThis.listello;
  }
  if (!loading) {
    loading = instantiate();
  }
  return loading;
}

async function waitForListello(timeoutMs = 1000): Promise<Listello> {
  const deadline = Date.now() + timeoutMs;
  while (Date.now() < deadline) {
    if (globalThis.listello) {
      return globalThis.listello;
    }
    await new Promise<void>((resolve) => setImmediate(resolve));
  }
  throw new Error("listello API was not exported from listello-js.wasm");
}

async function instantiate(): Promise<Listello> {
  require(path.join(here, "wasm_exec.js"));

  const go = new Go();
  const wasmPath = path.join(here, "listello-js.wasm");
  const bytes = await readFile(wasmPath);
  const { instance } = await WebAssembly.instantiate(bytes, go.importObject);

  // go.run resolves when the Go program exits; listello-js blocks in select{}.
  void go.run(instance);

  return waitForListello();
}

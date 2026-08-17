# Listello monorepo architecture

Target layout for shipping Listello as a Go core, a Node/WASM SDK, and a UI — without forcing the UI (or other JS consumers) to talk to Go directly.

## Packages

```text
listello/
  lib-go/       # Go module: domain, application, adapters, CLIs
  lib-node/     # npm package: loader, TypeScript types, WASM artifact
  ui/           # UI app; depends on lib-node
```

Optional root `package.json` with npm/pnpm workspaces:

```json
{
  "private": true,
  "workspaces": ["lib-node", "ui"]
}
```

Go remains its own module (`lib-go/go.mod`). JS workspaces do not replace the Go module.

## Responsibilities

| Package | Owns | Does not own |
|---------|------|----------------|
| `lib-go` | Domain, application, ports/adapters, native CLI, `js/wasm` bridge entrypoint (`cmd/listello-js`), WASI CLI build | npm metadata, UI |
| `lib-node` | `package.json`, JS/TS loader (`wasm_exec` glue), `.d.ts` API surface, published `dist/` including the `.wasm` | Go sources, business rules |
| `ui` | Presentation and UX | Persistence, domain logic, WASM build |

Dependency direction:

```text
ui  →  lib-node  →  (embeds) listello-js.wasm
                      ↑
              built from lib-go
```

## Build flow

1. **Compile the JS-callable WASM from Go**, emitting into the Node package:

   ```bash
   # from lib-go (or root Make that cds there)
   GOOS=js GOARCH=wasm go build -o ../lib-node/dist/listello-js.wasm ./cmd/listello-js
   cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" ../lib-node/dist/
   ```

2. **Node package** wraps that artifact: load/instantiate WASM, expose `globalThis.listello` (or a default export that resolves after load).

3. **UI** depends on `@…/listello` (workspace protocol locally, published version in CI/release).

It is intentional that the WASM binary is built **into** `lib-node` (e.g. `dist/`), not kept as the consumer-facing home under `lib-go/bin`. Go owns the *build*; the Node package owns the *distribution artifact*.

Native and WASI CLIs can still emit under `lib-go/bin/` for local Go workflows (`make build`, `make build-wasm` / Wasmtime). Those are separate from the npm-shipped `js/wasm` module.

```text
lib-go
  cmd/listello        → native + wasip1 CLI
  cmd/listello-js     → GOOS=js GOARCH=wasm (Node SDK)

lib-node/dist
  listello-js.wasm    ← js/wasm build output
  wasm_exec.js
  index.js / index.d.ts
```

## Node SDK surface

The `js/wasm` bridge exports a nested JS API (not the WASI CLI). Example shape:

```ts
listello.list.createList(name: string): List
listello.item.defineItem(listId: string, title: string): Item
```

- Return values are plain JS objects derived from Go structs (e.g. `{ ID, Name }` unless the bridge maps to camelCase).
- TypeScript `.d.ts` files describe that JS shape; they do not carry Go types across the boundary.
- Errors need an explicit convention at the boundary (throw vs result object); `(T, error)` cannot cross as-is.

Go sketch:

```go
js.Global().Set("listello", map[string]any{
	"list": map[string]any{
		"createList": js.FuncOf(createList),
	},
	"item": map[string]any{
		"defineItem": js.FuncOf(defineItem),
	},
})
select {} // keep the runtime alive
```

`internal/` packages stay under `lib-go`. The bridge lives in the same Go module (`cmd/listello-js`), so it may import `internal/...`. External npm consumers never import Go packages.

## Persistence note

- **Native / WASI CLI:** file-backed SQLite (e.g. ncruces) and file event log, as today.
- **Node `js/wasm` SDK:** may need adapters suited to that target (same application ports; possibly different adapter implementations). The monorepo split exists partly so those composition roots can diverge without affecting the UI.

## Publishing

- Publish `lib-node` from `dist/` (`main`, `types`, `files` point at built outputs).
- Consumers of the npm package should not need a Go toolchain.
- Choose one policy for the `.wasm` blob: commit under `dist/` for simple installs, or produce it in CI on publish. Either is valid; document the choice in the Node package README when implemented.

## Migration from the current repo

Today the repo root *is* the Go module. Moving to this layout means:

1. Relocate Go sources into `lib-go/` and update `go.mod` / import paths as needed.
2. Add `lib-node` with loader + types; point the `js/wasm` `-o` path at its `dist/`.
3. Add `ui` as a workspace package depending on the Node library.
4. Keep root docs/Make as a thin orchestrator if useful (`make build-js-wasm` → output into `lib-node/dist`).

Until that move happens, this document is the target architecture — not a description of the current tree.

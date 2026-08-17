# @bkotos/listello

Node SDK for Listello. Wraps the Go core compiled as `GOOS=js GOARCH=wasm`.

## Setup

```bash
# from repo root
npm install

# build WASM into dist/ then compile TypeScript
npm run build:all -w @bkotos/listello
```

Requires Go (to build the wasm) and Node ≥ 18.

## Usage

```js
import { loadListello } from "@bkotos/listello";

const listello = await loadListello();
const list = listello.list.createList("Next actions");
console.log(list.ID, list.Name);
```

## Layout

| Path | Role |
|------|------|
| `src/` | TypeScript loader + types |
| `dist/` | Published outputs: `index.js`, `.d.ts`, `listello-js.wasm`, `wasm_exec.js` |

The `.wasm` is produced by `make -C ../lib-go build-js-wasm`.

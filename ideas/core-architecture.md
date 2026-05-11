# Listello Core Architecture Notes

## High-Level Direction

Listello should have a portable core domain engine written in Go.

The Go engine should be responsible for:

- Parsing markdown task files
- Parsing YAML frontmatter
- Validating schemas and domain rules
- Performing all CRUD/domain operations
- Rendering markdown back out
- Producing normalized JSON representations
- Serving as the single authoritative implementation of Listello behavior

The frontend/UI should NOT directly manipulate markdown files.

Instead:

```text
UI / CLI / Agent / MCP / Browser
    ↓
Go Core Engine
    ↓
Markdown + Structured JSON
```

---

# Why This Architecture Makes Sense

This is conceptually similar to projects like esbuild:

```text
Go core engine
    ↓
CLI / npm package / Wasm module
```

The key idea is:

- business/domain logic implemented once
- reusable everywhere
- deterministic behavior across environments
- portable test suite
- future-proof implementation boundary

Potential targets:

```text
Native CLI
Browser Wasm
Node.js npm package
MCP server
Desktop app
VS Code extension
Future mobile app
```

---

# Important Architectural Principle

## The Go Core Is The Source of Truth

All writes should go through the Go core.

The UI should never directly:

- edit YAML
- append markdown comments
- mutate task state
- manipulate file formatting

Instead, the UI sends intents/commands.

Example:

```text
User clicks "Complete Task"
    ↓
Frontend sends command
    ↓
Go core applies domain rules
    ↓
Go core returns updated markdown + JSON
```

---

# Core Responsibilities

The Go core should own:

## Parsing

```text
Markdown → structured domain model
```

Including:

- YAML frontmatter
- comments/activity
- metadata
- task relationships
- validation

---

## Rendering

```text
Domain model → markdown
```

Including:

- preserving human-authored formatting where possible
- stable formatting
- deterministic output

---

## CRUD / Domain Operations

Examples:

- create task
- update task
- complete task
- add comment
- archive task
- validate task
- format task

---

# Recommended Go API Shape

Example domain-oriented API:

```go
CreateTask(input CreateTaskInput) (*TaskDocument, error)

UpdateTask(doc TaskDocument, patch TaskPatch) (*TaskDocument, error)

AddComment(doc TaskDocument, input AddCommentInput) (*TaskDocument, error)

CompleteTask(doc TaskDocument, input CompleteTaskInput) (*TaskDocument, error)

ValidateTask(doc TaskDocument) []ValidationError

RenderMarkdown(doc TaskDocument) (string, error)

ParseMarkdown(markdown string) (*TaskDocument, error)
```

Important:

The API should be domain-oriented, NOT markdown-edit oriented.

---

# CLI As Stable Public Contract

The CLI should expose stable commands like:

```bash
listello parse task.listello.md --json

listello validate task.listello.md --json

listello create --title "Buy diapers"

listello update task.listello.md --set status=waiting

listello comment task.listello.md --author Brian --body "Called them today"

listello complete task.listello.md

listello format task.listello.md
```

The CLI output becomes the stable contract.

---

# Testing Strategy

Use TWO layers of testing.

---

## 1. Go Unit Tests

Purpose:

- fast TDD feedback loop
- parser internals
- domain rules
- validation logic

These are implementation-level tests.

Example:

```go
func TestParseFrontmatter(t *testing.T)
```

---

## 2. Language-Agnostic Contract Tests

These test the public behavior of the engine.

Prefer:

- Node.js
- Vitest
- shell tests
- fixture-based tests

These tests should invoke the CLI or Wasm API.

Example:

```ts
const output = execFileSync("./bin/listello", [
  "parse",
  "fixtures/basic-task.listello.md",
  "--json"
]);

expect(JSON.parse(output)).toEqual(expected);
```

The important thing:

Tests validate:

```text
markdown input → normalized output
```

NOT:

```text
Go implementation details
```

This enables future portability.

---

# Why Contract Tests Matter

If the engine is ever rewritten:

```text
Go → Rust
Go → Zig
Go → other runtime
```

The same fixtures and behavioral expectations can remain unchanged.

This protects the Listello format itself.

---

# Recommended Fixture Structure

```text
contract-tests/
  fixtures/
    parse-basic/
      input.listello.md
      expected.json

    update-status/
      input.listello.md
      command.json
      expected.listello.md
      expected.json

    add-comment/
      input.listello.md
      command.json
      expected.listello.md
      expected.json
```

---

# Round-Trip Safety

Very important.

The engine should aim to preserve:

- formatting
- spacing
- comments
- ordering
- human-written text

unless explicitly formatting.

Meaning:

```text
markdown
  ↓ parse
domain model
  ↓ mutate
updated domain model
  ↓ render
updated markdown
```

should avoid unnecessary rewrites.

Example:

Updating task status should not completely reformat the file.

---

# npm + Wasm Strategy

The frontend should likely remain TypeScript/React.

Recommended split:

```text
TypeScript
  ↓
UI/state/browser integration

Go
  ↓
parsing/domain engine
```

Potential npm packages:

```text
@listello/core
@listello/core-wasm
@listello/cli
```

---

# Recommended Long-Term Philosophy

Listello markdown becomes a portable document standard.

That means:

- AI agents
- editors
- CLIs
- MCP servers
- sync tools
- import/export pipelines
- VS Code extensions

can all share the same canonical parser + domain engine.

This is the strongest justification for using Go + Wasm + npm distribution.

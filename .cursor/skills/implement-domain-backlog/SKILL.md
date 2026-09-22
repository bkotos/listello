---
name: implement-domain-backlog
description: >-
  Walks docs/domain-model-backlog.md from the top and implements the next
  unchecked checkbox (one layer, one PR). Use when the user types
  /implement-domain-backlog, @implement-domain-backlog, or asks to knock off
  the next domain-model backlog item, continue the event-storming backlog, or
  implement the next command/event slice. Marks that checkbox [x] in the
  backlog file as part of the same PR, then stops for review and merge.
---

# Implement domain model backlog

Orchestrator for [docs/domain-model-backlog.md](../../../docs/domain-model-backlog.md). **Read that file first** (How to work a slice, one-PR-per-skill, UI mockup). Then do **exactly one** unchecked checkbox.

This skill does not replace the layer skills. After you know which checkbox it is, **read and follow** the matching layer skill (or the domain/bootstrap notes below).

## When to use

- User invokes `/implement-domain-backlog` or `@implement-domain-backlog`
- User asks to implement / knock off / continue the domain model backlog
- User asks to start at the top of the event-storming gap list

If they name a specific command and layer (e.g. "Delete List, domain only"), do that checkbox instead of the first unchecked one — still one checkbox, still one PR.

## Hard rules

1. **One checkbox → one skill → one PR.** Never implement domain-through-UI (or two layers) in one run.
2. **TDD** (`.cursor/rules/tdd.mdc`): red spec → **stop for review** → green only after the user approves.
3. **Do not start the next checkbox** until this PR is reviewed and **merged**. After merge they invoke `/implement-domain-backlog` again (or say "continue" in a follow-up once merge is on `master`).
4. **Mark the checkbox done in this PR.** Change that line from `- [ ]` to `- [x]` in `docs/domain-model-backlog.md`. Do this in green, after the layer work passes — not during red, not in a later PR.
5. **Upstream gates** ([LAYER-ORDER.md](../LAYER-ORDER.md)): if the layer below is missing or still a stub, **stop**. Do not write downstream tests or code. Do not check the box.

## Procedure

### 1. Pick the item

Read `docs/domain-model-backlog.md` from the top.

The work items are the `- [ ]` / `- [x]` lines under command headings (List, Subtask, Comment, …). Ignore the legend table.

Take the **first `- [ ]`**. Record:

| Field | From |
|---|---|
| Command | Nearest `### Command → Event` heading |
| Aggregate section | Nearest `##` (List, Subtask, Comment, Delegation Policy, Delegation, Listello Instance) |
| Layer label | The checkbox text (`Domain`, `Application service`, `UI`, …) |
| Context package | Instance section → `api/internal/listello-instance-context/`. Everything else → `api/internal/personal-productivity-context/` |

Tell the user which checkbox you are doing before writing code.

### 2. Map the layer

| Checkbox starts with | Do this |
|---|---|
| `Domain` | Gherkin + godog in that context's `domain/` package. No layer skill. Spec the command/event from [docs/domain-model.md](../../../docs/domain-model.md). Call the real domain API from steps (not a no-op). |
| `Application service` | [create-application-service](../create-application-service/SKILL.md) |
| `Adapter repository` | [create-adapter-repository](../create-adapter-repository/SKILL.md) |
| `Bootstrap` | Wire adapter → service in the composition root (`api/internal/bootstrap/`, `api/cmd/…`). No layer skill. |
| `API handler` | [create-api-handler](../create-api-handler/SKILL.md) |
| `CLI` | [create-cli-command](../create-cli-command/SKILL.md) |
| `API client` | [create-api-client](../create-api-client/SKILL.md) |
| `React Query hooks` | [create-api-queries](../create-api-queries/SKILL.md) — **new reads only**. If this is a write and the page already has a query key to invalidate, do not add a mutation hook. Mark the checkbox `[x]`, note N/A in the PR, and **immediately continue to the next `- [ ]` in this same invocation** (usually UI). That skip is the only time two checkboxes share a PR. |
| `UI` | [create-ui-component](../create-ui-component/SKILL.md). **Before specs:** boot the mockup (`cd mockup/task-management-system-bulma && npm run dev`), open it in Chrome (`http://localhost:3000`), inspect the screens for this command, use that as the visual spec. Do not edit the mockup. |

Hosting-mode lines that say "only if … needs new behavior": if no new behavior is required, mark `[x]` with no code and continue to the next `- [ ]` in this same invocation (same skip exception).

### 3. Red, then stop

Follow the layer skill's TDD loop. Run the relevant tests. Confirm the failure is missing behavior.

**Stop.** Summarize: which backlog checkbox, what spec, why it failed. Do not implement and do not check the box until the user approves.

### 4. Green (after approval)

Bare minimum to pass. Then:

1. In `docs/domain-model-backlog.md`, change **only** this command's checkbox from `- [ ]` to `- [x]` (the first matching unchecked line under that `###` heading).
2. Commit the layer work and the backlog check together (or two commits on the same branch).
3. Open **one** PR. Title/body name the command, the event, and the layer. Mention that the backlog checkbox is marked done.

### 5. Stop for merge

Do not start the next layer. Tell the user:

- Merge this PR.
- Run `/implement-domain-backlog` again to pick up the next `- [ ]` from the top.

If they say "continue" in this same chat: `git fetch origin master`, confirm the previous PR is on `master`, then start the next checkbox on a **new** branch as a **new** PR. If it is not merged, wait.

## Out of scope

- Checking off boxes you did not implement (except the React Query / hosting "only if" skip above)
- Editing the Next.js mockup to make `ui/` tests pass
- Implementing extra event-storming commands that are not the current checkbox
- Fixing domain-model misalignments unless they block this checkbox

## Verification

Use the layer skill's verification commands. Also confirm `docs/domain-model-backlog.md` shows `[x]` for this checkbox and still `[ ]` for later ones.

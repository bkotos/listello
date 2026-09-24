# Domain model implementation backlog

Tracks command/event pairs from [domain-model.md](domain-model.md) that are **not yet in the code domain**. Check items off as each vertical slice lands.

To implement the next command/event pair: **`/implement-domain-backlog`** (skill: [implement-domain-backlog](../.cursor/skills/implement-domain-backlog/SKILL.md)). One invoke walks every remaining layer under that `###` heading. Each checkbox is still its own PR; the agent marks it `[x]` in this file in that PR, waits for merge, then continues to the next layer. Invoke again when that command is fully checked off.

Source: event storming → `docs/domain-model.md`. This list is the gap vs `api/internal/personal-productivity-context/domain` and `api/internal/listello-instance-context/domain`. It does not track misalignments (extra pairing commands, `standalone-web` vs the storming hosting names, item-level `ItemPriorityChanged`, and so on).

## How to work a slice

Follow TDD (`.cursor/rules/tdd.mdc`): red spec → stop for review → green. Do not start a downstream layer until the upstream layer is implemented (not a stub). See [.cursor/skills/LAYER-ORDER.md](../.cursor/skills/LAYER-ORDER.md).

Adapter and CLI/handler are siblings: both depend on the application service, not on each other. Handler/CLI tests use mocks; the adapter is for persistence and production wiring.

| Checkbox | Layer | Skill / notes |
|---|---|---|
| Domain | Gherkin + godog | No skill file. Writes live in the context `domain/` package. |
| Application service | Use case + repository port | `create-application-service` |
| Adapter repository | SQLite (and Postgres where those tests apply) | `create-adapter-repository` |
| Bootstrap | Wire adapter → service in composition root | No skill yet |
| API handler | HTTP + view DTOs + `make api-types` + Bruno | `create-api-handler` |
| CLI | Cobra command + e2e | `create-cli-command` |
| API client | `ui/src/lib/api/{resource}-client.ts` | `create-api-client` |
| React Query hooks | `{resource}-queries.ts` | `create-api-queries` — **new reads only**. Skip for writes when the page already has a query key to invalidate. |
| UI | Page/component wiring | `create-ui-component` — boot the mockup in Chrome first (see below) |

Personal productivity commands use `api/internal/personal-productivity-context/`. Instance commands use `api/internal/listello-instance-context/` (same layer order; skills are written against the productivity paths).

The first command on a new aggregate creates the service and repository. Later commands on that aggregate extend them.

### One pull request per skill

Each **checkbox** is its own pull request. One **slash invoke** covers one **command/event pair** (one `###` heading) and walks the remaining checkboxes under it.

Do **not** implement a whole command (domain through UI) in one giant PR.

1. Do **one** layer (one skill).
2. Open a pull request for that layer only.
3. **Stop.** Wait for review and for the PR to be merged.
4. Then, in the **same** `/implement-domain-backlog` run, start the next checkbox under that command, on a new branch, as a new PR.

Do not stack the next skill onto an unmerged PR. Do not open the next PR until the previous one is merged. Do not start a different `###` command until the user invokes the skill again.

TDD red → stop → green still happens **inside** that one skill PR. That is not a reason to split red and green into two PRs, and it is not a reason to batch several skills into one.

### UI: boot the mockup in Chrome

When the next checkbox is **UI** (`create-ui-component`), do not guess the layout from memory or from the domain names alone.

1. Start the Next.js mockup:

   ```bash
   cd mockup/task-management-system-bulma
   npm run dev
   ```

2. Open it in **Chrome** (Next.js defaults to `http://localhost:3000`).
3. Navigate to the screens and controls that correspond to this command (list, item detail, comments, subtasks, delegation, onboarding/hosting, and so on).
4. Use what you see — structure, labels, icons, hover/menus, empty states, and interactions — as the visual spec for the `ui/` implementation.

Read the mockup source under `mockup/task-management-system-bulma/` as well. **Do not edit the mockup** to make `ui/` tests pass. If the mockup has no screen for this command yet, say so and implement the smallest UI that still matches Listello patterns.

---

## List

### Delete List → List Deleted

- [x] Domain
- [x] Application service
- [x] Adapter repository
- [x] Bootstrap
- [x] API handler
- [x] CLI
- [x] API client
- [x] React Query hooks
- [x] UI

---

## Subtask

Child entity of Item in the domain model. No `Subtask` type exists yet. (`Item.ChangePriority` already raises `ItemPriorityChanged` on **items**; that is not this work.)

### Add Subtask → Subtask Added to Item

- [ ] Domain
- [ ] Application service
- [ ] Adapter repository
- [ ] Bootstrap
- [ ] API handler
- [ ] CLI
- [ ] API client
- [ ] React Query hooks
- [ ] UI

### Complete Subtask → Subtask Completed on Item

- [ ] Domain
- [ ] Application service
- [ ] Adapter repository
- [ ] Bootstrap
- [ ] API handler
- [ ] CLI
- [ ] API client
- [ ] React Query hooks
- [ ] UI

### Uncomplete Subtask → Subtask Uncompleted

- [ ] Domain
- [ ] Application service
- [ ] Adapter repository
- [ ] Bootstrap
- [ ] API handler
- [ ] CLI
- [ ] API client
- [ ] React Query hooks
- [ ] UI

### Delete Subtask → Subtask Deleted on Item

- [ ] Domain
- [ ] Application service
- [ ] Adapter repository
- [ ] Bootstrap
- [ ] API handler
- [ ] CLI
- [ ] API client
- [ ] React Query hooks
- [ ] UI

### Prioritize → Subtask Priority Changed

Subtask-entity priority, not item-level `ChangePriority` / `ItemPriorityChanged`.

- [ ] Domain
- [ ] Application service
- [ ] Adapter repository
- [ ] Bootstrap
- [ ] API handler
- [ ] CLI
- [ ] API client
- [ ] React Query hooks
- [ ] UI

---

## Comment

Child entity of Item. No `Comment` type exists yet.

### Comment → Item Commented On

- [x] Domain
- [x] Application service
- [x] Adapter repository
- [x] Bootstrap
- [ ] API handler
- [ ] CLI
- [ ] API client
- [ ] React Query hooks
- [ ] UI

### Delete Comment → Item Comment Deleted

- [ ] Domain
- [ ] Application service
- [ ] Adapter repository
- [ ] Bootstrap
- [ ] API handler
- [ ] CLI
- [ ] API client
- [ ] React Query hooks
- [ ] UI

---

## Delegation Policy

No `DelegationPolicy` type exists yet.

### Allow Delegation → Delegation Permission Granted to Other Person

- [ ] Domain
- [ ] Application service
- [ ] Adapter repository
- [ ] Bootstrap
- [ ] API handler
- [ ] CLI
- [ ] API client
- [ ] React Query hooks
- [ ] UI

---

## Delegation

No `Delegation` type exists yet.

### Delegate Item → Item Delegated

- [ ] Domain
- [ ] Application service
- [ ] Adapter repository
- [ ] Bootstrap
- [ ] API handler
- [ ] CLI
- [ ] API client
- [ ] React Query hooks
- [ ] UI

### Request Clarification → Clarification Requested

- [ ] Domain
- [ ] Application service
- [ ] Adapter repository
- [ ] Bootstrap
- [ ] API handler
- [ ] CLI
- [ ] API client
- [ ] React Query hooks
- [ ] UI

### Provide Clarification → Clarification Provided

- [ ] Domain
- [ ] Application service
- [ ] Adapter repository
- [ ] Bootstrap
- [ ] API handler
- [ ] CLI
- [ ] API client
- [ ] React Query hooks
- [ ] UI

### Deny Delegation → Delegation Permission Revoked from Other Person

In [domain-model.md](domain-model.md) this sits on the Delegation aggregate. The command map puts Allow/Deny on Delegation Policy. Implement against the domain model unless we decide otherwise.

- [ ] Domain
- [ ] Application service
- [ ] Adapter repository
- [ ] Bootstrap
- [ ] API handler
- [ ] CLI
- [ ] API client
- [ ] React Query hooks
- [ ] UI

---

## Listello Instance

These live in `api/internal/listello-instance-context/`. Generic **Initialize Persistence → Persistence Initialized** already exists; this is the local-specific pair from the storming board.

### Initialize Local Persistence → Local Persistence Initialized

- [ ] Domain
- [ ] Application service
- [ ] Adapter repository
- [ ] Bootstrap
- [ ] API handler
- [ ] CLI
- [ ] API client
- [ ] React Query hooks
- [ ] UI

### Hosting modes `self-hosted` and `embedded`

Select Hosting Mode is already implemented through the stack. The ER lists `local, self-hosted, embedded`. The domain currently allows `local` and `standalone-web` only.

- [ ] Domain — accept `self-hosted` and `embedded`
- [ ] Application service — only if Select Hosting Mode needs new behavior
- [ ] Adapter repository — persist the new values
- [ ] Bootstrap
- [ ] API handler
- [ ] CLI
- [ ] API client
- [ ] React Query hooks
- [ ] UI — offer the storming modes where hosting is chosen

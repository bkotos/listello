# Domain model implementation backlog

Tracks command/event pairs from [domain-model.md](domain-model.md) that are **not yet in the code domain**. Check items off as each vertical slice lands.

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
| UI | Page/component wiring | `create-ui-component` |

Personal productivity commands use `api/internal/personal-productivity-context/`. Instance commands use `api/internal/listello-instance-context/` (same layer order; skills are written against the productivity paths).

The first command on a new aggregate creates the service and repository. Later commands on that aggregate extend them.

---

## List

### Delete List → List Deleted

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

- [ ] Domain
- [ ] Application service
- [ ] Adapter repository
- [ ] Bootstrap
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

---
name: create-cli-command
description: >-
  Scaffolds and wires Cobra CLI commands in api/cmd/cli/ that call application
  services. Use when adding or extending CLI commands, cobra subcommands, root
  wiring, or listello list/item commands following Listello hexagonal architecture.
---

# Create CLI Command

Guide for adding Cobra commands in `api/cmd/cli/`. Read this skill before changing CLI commands or root wiring.

Architecture context: see [README.md](../../../README.md) (CLI calls application services directly). Layer order: [LAYER-ORDER.md](../LAYER-ORDER.md). CLI design: [listello-cli-design.md](../../../listello-cli-design.md). TDD workflow: see [.cursor/rules/tdd.mdc](../../rules/tdd.mdc).

## Upstream dependencies

**Stop and do not proceed** until the application service method exists and is implemented. See [LAYER-ORDER.md](../LAYER-ORDER.md).

Before starting, verify:

| Check | How |
|-------|-----|
| Service method exists | `{Method}(...)` on `{Aggregate}Service` interface in `internal/application/` |
| Method is implemented | Not `return ..., fmt.Errorf("not implemented")` |
| Application tests pass | `go test ./internal/application/...` green for that method |

If any check fails → **stop**. Tell the user to use `create-application-service` first and get green tests. Do not write command tests or Cobra wiring.

Adapter repository is **not** required for this skill (command tests use mock service interfaces).

## Scope

**In scope:** Cobra command factory, command test (`cobra.Execute` + `SetArgs`), wiring on parent command or `root.go`, stdout/stderr assertions.

**Out of scope** (mention as follow-ups only; do not implement unless asked):

- Application service methods (`api/internal/application/`) — use [create-application-service](../create-application-service/SKILL.md)
- Adapter repositories (`api/internal/adapter/`) — use [create-adapter-repository](../create-adapter-repository/SKILL.md)
- Bootstrap / `main.go` service construction (`api/internal/bootstrap/`, `api/cmd/cli/main.go`)
- HTTP handlers (`api/cmd/server/`) — use [create-api-handler](../create-api-handler/SKILL.md)
- View DTOs / `api-types` (CLI uses args, flags, and human-readable output — not JSON DTOs unless `--json` is in scope)

## Architecture constraints

- Commands depend on application **service interfaces** — never call domain or adapters directly.
- Use Cobra `RunE` (return errors; do not `os.Exit` in commands).
- Shape: `listello <noun> <verb>` — parent group per aggregate (`list`, `item`), leaf command per action (`create`, `define`).
- Mutating commands print a one-line confirmation to `cmd.OutOrStdout()`.
- Return application/domain errors from `RunE`; root prints `error: %v` to stderr and exits non-zero.
- Do not print domain events in CLI output (events go to the event log via the application layer).
- Factory accepts the service **interface**: `func NewListCreate(listService application.ListService) *cobra.Command`.

## Do not duplicate domain logic

Domain rules live in `internal/domain` (Gherkin + godog). CLI commands are thin adapters — do not re-implement or re-test domain behavior here.

**In command code:**

- Do not validate args or flags that the domain/application layer already validates.
- Parse positional args/flags; pass values to the service; return errors unchanged from `RunE`.
- Do not add command-only business rules.

**In command tests:**

- Mock the **service interface** (`appmocks.NewMock{Service}Service`); assert `EXPECT().{Method}(...)` was called with correct arguments.
- Split into separate tests for **calls service** and **prints confirmation** — one concern per test.
- Do **not** add tests for domain validation failures (e.g. define on inbox, create list named "Inbox") — those belong in `internal/domain`.
- Do **not** re-assert domain field semantics on returned aggregates — only verify the command uses the service result in its confirmation output.

## Preconditions

These duplicate the upstream gate — all must pass:

1. **Application service method exists and is implemented** — the command calls a real service API, not a stub.
2. **Command name and args decided** — see [listello-cli-design.md](../../../listello-cli-design.md) for the command tree.

If the service method is missing or still `not implemented` → **stop** and use `create-application-service` first.

## Decision tree

1. **New aggregate group?** → Create `{resource}.go` with `New{Resource}()` parent command; register in `root.go`.
2. **New action on existing group?** → Create `{resource}_{action}.go`; `parent.AddCommand(New{Resource}{Action}(svc))`.
3. **New service dependency?** → Add parameter to `newRoot` in `root.go`; mention `main.go` bootstrap follow-up.
4. **Positional args?** → Use `cobra.ExactArgs`, `MinimumNArgs`, etc.; read via `args[0]`.
5. **Flags?** → `cmd.Flags().StringVar(&v, "name", "", "help")` — only add flags the spec requires.

## Scaffold checklist

```
Task progress:
- [ ] Verify upstream: service method implemented, application tests green (stop if not — see LAYER-ORDER.md)
- [ ] Read application service method signature
- [ ] Decide command path (e.g. list create), args, flags, confirmation message
- [ ] Create or extend parent group command ({resource}.go)
- [ ] Write failing command test (cobra Execute + SetArgs)
- [ ] Run tests — confirm failure is missing behavior
- [ ] STOP — summarize spec and failure for user review
- [ ] (After approval) Implement leaf command (RunE calls service, prints confirmation)
- [ ] Wire command on parent (AddCommand) or root.go
- [ ] Update newRoot signature if new service dependency
- [ ] Re-run tests — confirm green
```

## Naming conventions

| Artifact | Convention |
|----------|------------|
| Parent group file | `{resource}.go` (e.g. `list.go`, `item.go`) |
| Leaf command file | `{resource}_{action}.go` (e.g. `list_create.go`, `item_define.go`) |
| Test file | `{resource}_{action}_test.go` |
| Test package | `commands` (same package) |
| Parent factory | `New{Resource}(svc) *cobra.Command` |
| Leaf factory | `New{Resource}{Action}(svc) *cobra.Command` |
| Test name | `Test{Resource}{Action}_{Behavior}` |
| Test helper | `new{Resource}TestRoot(svc) *cobra.Command` in test file |
| Service mocks | `appmocks "github.com/bkotos/listello/internal/application/mocks"` (mockery-generated) |
| Service param | `{aggregate}Service` — type is `application.{Aggregate}Service` interface |

## Code templates

### Parent group command

```go
package commands

import (
	"github.com/spf13/cobra"

	application "github.com/bkotos/listello/internal/application"
)

func New{Resource}({aggregate}Service application.{Service}) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "{resource}",
		Short: "Manage {resources}",
	}
	cmd.AddCommand(New{Resource}{Action}({aggregate}Service))
	return cmd
}
```

### Leaf command (positional args)

```go
package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	application "github.com/bkotos/listello/internal/application"
)

func New{Resource}{Action}({aggregate}Service application.{Service}) *cobra.Command {
	return &cobra.Command{
		Use:   "{action} <arg>",
		Short: "{Action} a {resource}",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := {aggregate}Service.{Method}(args[0] /* map other args/flags */)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created {resource} %q (%s)\n", result.Name, result.ID)
			return nil
		},
	}
}
```

### Leaf command (multiple args)

```go
Use:   "{action} <list-id> <title>",
Args:  cobra.ExactArgs(2),
RunE: func(cmd *cobra.Command, args []string) error {
	item, err := {aggregate}Service.{Method}(args[0], args[1])
	// ...
},
```

### Root wiring (`root.go`)

```go
func newRoot(listService application.ListService, itemService application.ItemService) *cobra.Command {
	root := &cobra.Command{
		Use:           "listello",
		Short:         "Listello command-line interface",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(commands.NewList(listService))
	root.AddCommand(commands.NewItem(itemService))
	return root
}
```

## Test templates

Use mock service interfaces from `internal/application/mocks`. Capture stdout/stderr with `bytes.Buffer`. Split **calls service** and **prints confirmation** into separate tests.

### Calls application

```go
func Test{Resource}{Action}_CallsApplication(t *testing.T) {
	{aggregate}Service := appmocks.NewMock{Service}Service(t)
	{aggregate}Service.EXPECT().{Method}(/* args */).Return(domain.{Aggregate}{ID: "IT_1"}, nil)

	root := new{Resource}TestRoot({aggregate}Service)
	root.SetOut(&bytes.Buffer{})
	root.SetArgs([]string{"{resource}", "{action}", /* args */})

	err := root.Execute()
	require.NoError(t, err)
}
```

### Prints confirmation

```go
func Test{Resource}{Action}_PrintsConfirmation(t *testing.T) {
	expected := domain.{Aggregate}{ID: "IT_1", Title: "Buy milk"}
	{aggregate}Service := appmocks.NewMock{Service}Service(t)
	{aggregate}Service.EXPECT().{Method}(/* args */).Return(expected, nil)

	stdout := &bytes.Buffer{}
	root := new{Resource}TestRoot({aggregate}Service)
	root.SetOut(stdout)
	root.SetArgs([]string{"{resource}", "{action}", /* args */})

	require.NoError(t, root.Execute())
	assert.Contains(t, stdout.String(), `Expected confirmation line`)
}
```

One behavioral concern per test. Mock the service interface — do **not** stub repository ports. Do **not** add command tests for domain validation failures.

## Output conventions

| Situation | Behavior |
|-----------|----------|
| Success (mutate) | One-line confirmation to stdout via `fmt.Fprintf(cmd.OutOrStdout(), ...)` |
| Application/domain error | Return `error` from `RunE` (root prints to stderr) |
| Success (read/query) | Print human-readable summary to stdout (format per spec) |

Match existing commands; do not add `--json` or extra flags unless the spec requires them.

## TDD workflow

Follow `.cursor/rules/tdd.mdc`:

1. **Red** — Write failing command test first. Command may not exist yet (compile error) or print wrong output. Confirm failure is missing behavior.
2. **Stop** — Summarize spec and failure. **Do not implement** until the user approves.
3. **Green** — Implement command + wiring. Re-run tests.

Do not write command implementation in the same turn as a new failing spec.

## Verification

```bash
# CLI command tests only
cd api && go test ./cmd/cli/commands/...

# Full API tests
make -C api test

# Manual smoke test (after bootstrap wiring)
make -C api run ARGS='list create "Next actions"'
```

## Further reading

Annotated walkthrough of `list create`: [examples.md](examples.md)

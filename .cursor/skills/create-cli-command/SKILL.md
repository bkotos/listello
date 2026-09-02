---
name: create-cli-command
description: >-
  Scaffolds and wires Cobra CLI commands in api/cmd/cli/ that call application
  services. Use when adding or extending CLI commands, cobra subcommands, root
  wiring, or listello list/item commands following Listello hexagonal architecture.
---

# Create CLI Command

Guide for adding Cobra commands in `api/cmd/cli/`. Read this skill before changing CLI commands or root wiring.

Architecture context: see [README.md](../../../README.md) (CLI calls application services directly). CLI design conventions: [listello-cli-design.md](../../../listello-cli-design.md). TDD workflow: see [.cursor/rules/tdd.mdc](../../rules/tdd.mdc).

## Scope

**In scope:** Cobra command factory, command test (`cobra.Execute` + `SetArgs`), wiring on parent command or `root.go`, stdout/stderr assertions.

**Out of scope** (mention as follow-ups only; do not implement unless asked):

- Application service methods (`api/internal/application/`) — use [create-application-service](../create-application-service/SKILL.md)
- Adapter repositories (`api/internal/adapter/`) — use [create-adapter-repository](../create-adapter-repository/SKILL.md)
- Bootstrap / `main.go` service construction (`api/internal/bootstrap/`, `api/cmd/cli/main.go`)
- HTTP handlers (`api/cmd/server/`) — use [create-api-handler](../create-api-handler/SKILL.md)
- View DTOs / `api-types` (CLI uses args, flags, and human-readable output — not JSON DTOs unless `--json` is in scope)

## Architecture constraints

- Commands depend on application services — never call domain or adapters directly.
- Use Cobra `RunE` (return errors; do not `os.Exit` in commands).
- Shape: `listello <noun> <verb>` — parent group per aggregate (`list`, `item`), leaf command per action (`create`, `define`).
- Mutating commands print a one-line confirmation to `cmd.OutOrStdout()`.
- Return application/domain errors from `RunE`; root prints `error: %v` to stderr and exits non-zero.
- Do not print domain events in CLI output (events go to the event log via the application layer).
- Factory accepts the application service: `func NewListCreate(lists *application.ListService) *cobra.Command`.

## Preconditions

1. **Application service method exists** — the command calls an existing service API.
2. **Command name and args decided** — see [listello-cli-design.md](../../../listello-cli-design.md) for the command tree.

If the service method is missing, stop and use `create-application-service` first.

## Decision tree

1. **New aggregate group?** → Create `{resource}.go` with `New{Resource}()` parent command; register in `root.go`.
2. **New action on existing group?** → Create `{resource}_{action}.go`; `parent.AddCommand(New{Resource}{Action}(svc))`.
3. **New service dependency?** → Add parameter to `newRoot` in `root.go`; mention `main.go` bootstrap follow-up.
4. **Positional args?** → Use `cobra.ExactArgs`, `MinimumNArgs`, etc.; read via `args[0]`.
5. **Flags?** → `cmd.Flags().StringVar(&v, "name", "", "help")` — only add flags the spec requires.

## Scaffold checklist

```
Task progress:
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
| Test helper | `newTestRoot(svc) *cobra.Command` in test file |
| Stubs | `stubListRepository`, `stubEventPublisher` in `stubs_test.go` |

## Code templates

### Parent group command

```go
package commands

import (
	"github.com/spf13/cobra"

	application "github.com/bkotos/listello/internal/application"
)

func New{Resource}(svc *application.{Service}) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "{resource}",
		Short: "Manage {resources}",
	}
	cmd.AddCommand(New{Resource}{Action}(svc))
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

func New{Resource}{Action}(svc *application.{Service}) *cobra.Command {
	return &cobra.Command{
		Use:   "{action} <arg>",
		Short: "{Action} a {resource}",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := svc.{Method}(args[0] /* map other args/flags */)
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
	item, err := svc.{Method}(args[0], args[1])
	// ...
},
```

### Root wiring (`root.go`)

```go
func newRoot(lists *application.ListService /*, items *application.ItemService */) *cobra.Command {
	root := &cobra.Command{
		Use:           "listello",
		Short:         "Listello command-line interface",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(commands.NewList(lists))
	// root.AddCommand(commands.NewItem(items))
	return root
}
```

## Test templates

Use `newTestRoot` with stub repositories. Capture stdout/stderr with `bytes.Buffer`.

### Success — calls service and prints confirmation

```go
func Test{Resource}{Action}_CallsApplicationAndPrintsConfirmation(t *testing.T) {
	// Arrange
	var saved domain.{Aggregate}
	svc := application.New{Service}(
		&stub{Aggregate}Repository{
			saveFn: func(/* args */) error {
				saved = /* capture */
				return nil
			},
		},
		&stubEventPublisher{},
	)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	root := newTestRoot(svc)
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetArgs([]string{"{resource}", "{action}", "arg value"})

	// Act
	err := root.Execute()

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "expected value", saved.Name)
	assert.Contains(t, stdout.String(), `Created {resource} "arg value" (`)
	assert.Contains(t, stdout.String(), saved.ID)
	assert.Empty(t, stderr.String())
}
```

### Domain error — returns error, no stdout

```go
func Test{Resource}{Action}_PrintsDomainError(t *testing.T) {
	svc := application.New{Service}(&stub{Aggregate}Repository{}, &stubEventPublisher{})
	// ... setup root, SetArgs with invalid input
	err := root.Execute()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected domain message")
	assert.Empty(t, stdout.String())
}
```

One behavioral concern per test. Add stubs to `stubs_test.go` when a new repository port is needed.

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

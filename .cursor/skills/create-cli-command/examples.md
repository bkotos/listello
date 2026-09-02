# CLI Command Examples

Annotated references from the Listello codebase. Read when implementing a new Cobra command.

## 1. Parent group: `list`

**File:** `api/cmd/cli/commands/list.go`

```go
func NewList(lists *application.ListService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Manage lists",
	}
	cmd.AddCommand(NewListCreate(lists))
	return cmd
}
```

Key points:

- One parent command per aggregate (`list`, `item`, …).
- Parent has no `RunE` — only groups subcommands.
- Register leaf commands with `AddCommand`.

## 2. Leaf command: `list create`

**File:** `api/cmd/cli/commands/list_create.go`

```go
func NewListCreate(lists *application.ListService) *cobra.Command {
	return &cobra.Command{
		Use:   "create <name>",
		Short: "Create a list",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := lists.CreateList(args[0])
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created list %q (%s)\n", list.Name, list.ID)
			return nil
		},
	}
}
```

Key points:

- `Args: cobra.ExactArgs(1)` — one positional arg (`<name>`).
- Call application service; propagate errors unchanged.
- Print confirmation to `cmd.OutOrStdout()`, not `fmt.Println`.
- Do not mention domain events in output.

## 3. Root wiring

**File:** `api/cmd/cli/root.go`

```go
func newRoot(lists *application.ListService) *cobra.Command {
	root := &cobra.Command{
		Use:           "listello",
		Short:         "Listello command-line interface",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(commands.NewList(lists))
	return root
}
```

**File:** `api/cmd/cli/main.go` (out of skill scope — follow-up)

```go
listService := bootstrap.NewListService(db, eventLog)
run(newRoot(listService))
```

Adding `ItemService` commands requires `newRoot` to accept `items *application.ItemService` and `main.go` to construct it.

## 4. Test: success path

**File:** `api/cmd/cli/commands/list_create_test.go`

```go
func newTestRoot(lists *application.ListService) *cobra.Command {
	root := &cobra.Command{
		Use:           "listello",
		Short:         "Listello command-line interface",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.AddCommand(NewList(lists))
	return root
}

func TestListCreate_CallsApplicationAndPrintsConfirmation(t *testing.T) {
	var saved domain.List
	svc := application.NewListService(
		&stubListRepository{
			saveFn: func(list domain.List) error {
				saved = list
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
	root.SetArgs([]string{"list", "create", "Next actions"})

	err := root.Execute()

	require.NoError(t, err)
	assert.Equal(t, "Next actions", saved.Name)
	assert.Contains(t, stdout.String(), `Created list "Next actions" (`)
	assert.Contains(t, stdout.String(), saved.ID)
	assert.Empty(t, stderr.String())
}
```

Key points:

- `SetArgs` uses the full command path: `list`, `create`, then args.
- Stub captures what the service persisted.
- Assert stdout confirmation and empty stderr.
- Real `application.NewListService` with stubbed ports — not mockery.

## 5. Test: domain error

```go
func TestListCreate_PrintsDomainError(t *testing.T) {
	svc := application.NewListService(&stubListRepository{}, &stubEventPublisher{})
	// ... setup buffers and root
	root.SetArgs([]string{"list", "create", "Inbox"})

	err := root.Execute()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot create a list named Inbox")
	assert.Empty(t, stdout.String())
}
```

Key points:

- Domain validation fails through the real service (no stub override needed).
- Error returned from `Execute()` — root's `run()` would print it to stderr in production.
- No stdout on failure.

## 6. Test stubs

**File:** `api/cmd/cli/commands/stubs_test.go`

Hand-rolled stubs implement application ports (same pattern as HTTP handler stubs):

```go
type stubListRepository struct {
	saveFn    func(list domain.List) error
	getAllFn  func() ([]domain.List, error)
	getByIDFn func(id string) (domain.List, error)
}
```

Add `stubItemRepository` when wiring item commands.

## 7. CLI design reference

**File:** `listello-cli-design.md`

Command tree shape:

```text
listello
├── list
│   ├── create
│   └── show
├── item
│   ├── capture
│   ├── define
│   └── complete
```

Conventions:

- `listello <noun> <verb>` — reads like domain language.
- IDs: `LS_…` for lists, `IT_…` for items.
- Mutating commands: one-line confirmation in plain language.
- Errors: domain message verbatim, non-zero exit.

## 8. Next command example: `item define`

Given `ItemService.DefineItem(listID, title string)`:

1. **Parent:** `item.go` with `NewItem(items *application.ItemService)` if not exists.
2. **Leaf:** `item_define.go`:
   - `Use: "define <list-id> <title>"`
   - `Args: cobra.ExactArgs(2)`
   - `items.DefineItem(args[0], args[1])`
   - `fmt.Fprintf(cmd.OutOrStdout(), "Defined item %q (%s) on list %s\n", item.Title, item.ID, listID)`
3. **Wire:** `NewItem` → `AddCommand(NewItemDefine(items))`; register `NewItem` in `root.go`.
4. **Test:** `TestItemDefine_CallsApplicationAndPrintsConfirmation` with `SetArgs([]string{"item", "define", "LS_1", "Buy milk"})`.
5. **Follow-up:** `main.go` bootstrap for `ItemService` (out of skill scope unless asked).

# CLI Command Examples

Annotated references from the Listello codebase. Read when implementing a new Cobra command.

## 1. Parent group: `item`

**File:** `api/cmd/cli/commands/item.go`

```go
func NewItem(itemService application.ItemService) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "item",
		Short: "Manage items",
	}
	cmd.AddCommand(NewItemDefine(itemService))
	return cmd
}
```

Key points:

- One parent command per aggregate (`list`, `item`, …).
- Parent has no `RunE` — only groups subcommands.
- Factory accepts `application.ItemService` **interface**.

## 2. Leaf command: `item define`

**File:** `api/cmd/cli/commands/item_define.go`

```go
func NewItemDefine(itemService application.ItemService) *cobra.Command {
	return &cobra.Command{
		Use:   "define <list-id> <title>",
		Short: "Define an item on a list",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			listID := args[0]
			title := args[1]

			item, err := itemService.DefineItem(listID, title)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Defined item %q (%s) on list %s\n", item.Title, item.ID, listID)
			return nil
		},
	}
}
```

Key points:

- `Args: cobra.ExactArgs(2)` — list ID and title as positional args.
- Call application service; propagate errors unchanged (no command-level validation of title, etc.).
- Print confirmation to `cmd.OutOrStdout()`, not `fmt.Println`.
- Do not mention domain events in output.

## 3. Root wiring

**File:** `api/cmd/cli/root.go`

```go
func newRoot(listService application.ListService, itemService application.ItemService) *cobra.Command {
	root := &cobra.Command{ ... }
	root.AddCommand(commands.NewList(listService))
	root.AddCommand(commands.NewItem(itemService))
	return root
}
```

**File:** `api/cmd/cli/main.go`

```go
listService := bootstrap.NewListService(db, eventLog)
itemService := bootstrap.NewItemService(db, eventLog)
run(newRoot(listService, itemService))
```

## 4. Tests: `item define` (preferred pattern)

**File:** `api/cmd/cli/commands/item_define_test.go`

Split **calls service** and **prints confirmation** into separate tests. Mock the service interface — do not stub repositories. Do **not** test domain validation failures here (covered in `internal/domain`).

### Calls application

```go
func TestItemDefine_CallsApplication(t *testing.T) {
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().DefineItem("LS_1", "Buy milk").Return(domain.Item{ID: "IT_1"}, nil)

	root := newItemTestRoot(itemService)
	root.SetOut(&bytes.Buffer{})
	root.SetArgs([]string{"item", "define", "LS_1", "Buy milk"})

	require.NoError(t, root.Execute())
}
```

### Prints confirmation

```go
func TestItemDefine_PrintsConfirmation(t *testing.T) {
	expected := domain.Item{ID: "IT_1", ListID: "LS_1", Title: "Buy milk"}
	itemService := appmocks.NewMockItemService(t)
	itemService.EXPECT().DefineItem("LS_1", "Buy milk").Return(expected, nil)

	stdout := &bytes.Buffer{}
	root := newItemTestRoot(itemService)
	root.SetOut(stdout)
	root.SetArgs([]string{"item", "define", "LS_1", "Buy milk"})

	require.NoError(t, root.Execute())
	assert.Contains(t, stdout.String(), `Defined item "Buy milk" (IT_1) on list LS_1`)
}
```

## 5. Legacy: `list create` tests

**File:** `api/cmd/cli/commands/list_create_test.go`

Older commands may still use `application.NewListService` with stub repositories in one combined test. New commands should use mock service interfaces and split concerns as in §4.

Do **not** add new domain-error command tests (e.g. `TestListCreate_PrintsDomainError`) — domain validation belongs in godog features.

## 6. CLI design reference

**File:** `listello-cli-design.md`

```text
listello
├── list
│   └── create
├── item
│   └── define
```

Conventions:

- `listello <noun> <verb>` — reads like domain language.
- IDs: `LS_…` for lists, `IT_…` for items.
- Mutating commands: one-line confirmation in plain language.
- Errors: domain message verbatim via `RunE` return (root prints to stderr).

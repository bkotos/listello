# Listello CLI Design

A proposed command-line interface for Listello, derived from the
[event storming command map](listello-event-storming-command-map.md). Each
aggregate becomes a command group, and each command on the board becomes a
sub-command. The design favors a `listello <noun> <verb>` shape so the CLI
reads like the ubiquitous language of the domain.

Conventions used throughout:

- List IDs look like `LS_9f3a2c`, item IDs like `IT_4b81d0`, comment IDs like
  `CM_77e1af` — matching the domain's `newID` prefixes.
- Anywhere a `<list>` or `<item>` argument appears, the CLI accepts either an
  ID or a unique name/title prefix (resolved before invoking the application
  service).
- Mutating commands print a short one-line confirmation in plain user
  language. Domain events are a technical design detail — they are persisted
  to the event log but never surfaced in CLI output.
- Errors from the domain are printed verbatim to stderr and exit non-zero.

---

## Command Tree

```text
listello
├── list
│   ├── create        List created
│   ├── delete        List deleted
│   └── show          (view: List view)
├── item
│   ├── capture       Item captured
│   ├── define        Item defined
│   ├── title         Item title changed
│   ├── describe      Item description changed
│   ├── due           Due date added to item
│   ├── undue         Due date removed from item
│   ├── tag           Tag added to item
│   ├── untag         Tag removed from item
│   ├── move          Item moved to other list
│   ├── complete      Item completed
│   ├── uncomplete    Item uncompleted
│   ├── delete        Item deleted
│   └── show          (view: Item view)
├── subtask
│   ├── add           Subtask added to item
│   ├── link          Item linked as child of item
│   ├── prioritize    Subtask priority changed
│   ├── complete      Subtask completed on item
│   ├── uncomplete    Subtask uncompleted
│   └── delete        Subtask deleted on item
├── comment
│   ├── add           Item commented on
│   └── delete        Item comment deleted
├── delegation
│   ├── allow         Delegation permission granted to other person
│   ├── deny          Delegation permission revoked from other person
│   ├── assign        Item delegated
│   ├── ask           Clarification requested
│   ├── answer        Clarification provided
│   └── show          (views: Delegated items / Assigned to me / Needs clarification)
└── inbox             (shortcut for `item capture` / inbox view)
```

Global flags, available on every command:

```text
--db <path>       SQLite database path        (default: listello.db)
--events <path>   Domain event log path       (default: domain_events.log)
--json            Emit machine-readable JSON instead of human output
--quiet           Suppress confirmation output; exit code only
```

---

## `listello list`

### `listello list create`

> Event storming: **User → Create list → List created** (implemented end to end)

```text
SYNOPSIS
    listello list create <name>
```

```console
$ listello list create "Next actions"
Created list "Next actions" (LS_9f3a2c)```

Failure cases surface the domain errors directly:

```console
$ listello list create Inbox
error: cannot create a list named Inbox
```

### `listello list delete`

> Event storming: **User → Delete list → List deleted**

```text
SYNOPSIS
    listello list delete <list> [--force]
```

`--force` skips the confirmation prompt when the list still has outstanding
items.

```console
$ listello list delete "Someday maybe"
List "Someday maybe" has 3 outstanding items. Delete anyway? [y/N] y
Deleted list "Someday maybe" (LS_1c44b7)```

### `listello list show`

> Interface: **List view** (query, no domain event)

```text
SYNOPSIS
    listello list show [<list>] [--all]
```

Without arguments, shows all lists. With a list, shows its items. `--all`
includes completed items.

```console
$ listello list show
  LS_0inbox  Inbox          4 items
  LS_9f3a2c  Next actions   7 items (2 due soon)
  LS_b2d901  Groceries      12 items

$ listello list show "Next actions"
Next actions (LS_9f3a2c)
  [ ] IT_4b81d0  Renew passport            due 2026-08-15  #errands
  [ ] IT_77aa03  Draft Q3 planning doc                     #work !high
  [x] IT_20c9fe  Book dentist appointment
```

---

## `listello item`

### `listello item capture`

> Event storming: **User → Capture inbox item → Item captured** (domain implemented)

Captures a raw thought straight to the Inbox — the only list capture is
allowed on. Also available as the top-level shortcut `listello inbox <title>`.

```text
SYNOPSIS
    listello item capture <title>
    listello inbox <title>                (shortcut)
```

```console
$ listello item capture "Look into standing desk options"
Captured "Look into standing desk options" (IT_8d02e1) to Inbox```

```console
$ listello item capture --list "Next actions" "This will fail"
error: can only capture items to the inbox
```

### `listello item define`

> Event storming: **User → Define item → Item defined** (domain implemented)

Defines a well-formed item directly on a non-inbox list.

```text
SYNOPSIS
    listello item define <title> --list <list> [--description <text>]
                         [--due <date>] [--tag <tag>]...
```

```console
$ listello item define "Renew passport" --list "Next actions" \
    --due 2026-08-15 --tag errands
Defined "Renew passport" (IT_4b81d0) on "Next actions"```

### `listello item title`

> Event storming: **Owner → Modify item title → Item title changed** (domain implemented)

```text
SYNOPSIS
    listello item title <item> <new-title>
```

```console
$ listello item title IT_4b81d0 "Renew passport before trip"
Renamed IT_4b81d0 to "Renew passport before trip"```

### `listello item describe`

> Event storming: **Owner → Modify item description → Item description changed** (domain implemented)

```text
SYNOPSIS
    listello item describe <item> <description>
    listello item describe <item> --edit        (opens $EDITOR)
```

```console
$ listello item describe IT_4b81d0 "Bring old passport and two photos to the post office."
Updated description of "Renew passport before trip" (IT_4b81d0)```

### `listello item due` / `listello item undue`

> Event storming: **Requester or Owner → Modify due date → Due date added to item**
> and **Remove due date → Due date removed from item** (domain implemented)

```text
SYNOPSIS
    listello item due <item> <date>
    listello item undue <item>
```

Dates accept ISO 8601 (`2026-08-15` or full timestamps).

```console
$ listello item due IT_4b81d0 2026-08-15
Due date for "Renew passport before trip" set to 2026-08-15
$ listello item undue IT_4b81d0
Removed due date from "Renew passport before trip"```

### `listello item tag` / `listello item untag`

> Event storming: **Owner → Tag item → Tag added to item**
> and **Untag item → Tag removed from item** (domain implemented)

```text
SYNOPSIS
    listello item tag <item> <tag>...
    listello item untag <item> <tag>...
```

```console
$ listello item tag IT_4b81d0 errands travel
Tagged "Renew passport before trip" with #errands #travel
$ listello item untag IT_4b81d0 travel
Removed #travel from "Renew passport before trip"```

### `listello item move`

> Event storming: **Owner → Move item → Item moved to other list** (domain implemented)

The natural home of inbox processing: capture to Inbox, then move to a real
list once defined.

```text
SYNOPSIS
    listello item move <item> --to <list>
```

```console
$ listello item move IT_8d02e1 --to "Next actions"
Moved "Look into standing desk options" from Inbox to "Next actions"```

### `listello item complete` / `listello item uncomplete`

> Event storming: **Requester or Owner → Complete item → Item completed**
> and **Uncomplete item → Item uncompleted** (domain implemented)

```text
SYNOPSIS
    listello item complete <item>...
    listello item uncomplete <item>...
```

```console
$ listello item complete IT_4b81d0
✔ Completed "Renew passport before trip"
$ listello item uncomplete IT_4b81d0
Reopened "Renew passport before trip"```

### `listello item delete`

> Event storming: **Requester or Owner → Delete item → Item deleted** (domain implemented)

```text
SYNOPSIS
    listello item delete <item> [--force]
```

```console
$ listello item delete IT_20c9fe
Delete "Book dentist appointment"? [y/N] y
Deleted "Book dentist appointment" (IT_20c9fe)```

### `listello item show`

> Interface: **Item view** (query, no domain event)

```text
SYNOPSIS
    listello item show <item>
```

```console
$ listello item show IT_4b81d0
Renew passport before trip (IT_4b81d0)
  list:        Next actions (LS_9f3a2c)
  state:       outstanding
  due:         2026-08-15
  tags:        #errands
  description: Bring old passport and two photos to the post office.

  subtasks (1/3 done)
    [x] IT_a1b2c3  Gather documents            !high
    [ ] IT_d4e5f6  Get passport photos taken   !medium
    [ ] IT_g7h8i9  Fill out DS-82 form         !low

  comments
    CM_77e1af  bkotos  2026-08-01  "Post office closes at 4pm on Saturdays"
```

---

## `listello subtask`

Subtasks are items linked as children of a parent item, so this group is a
convenience layer over item commands scoped to a parent.

### `listello subtask add`

> Event storming: **Owner → Add subtask → Subtask added to item**

```text
SYNOPSIS
    listello subtask add <parent-item> <title> [--priority low|medium|high]
```

```console
$ listello subtask add IT_4b81d0 "Get passport photos taken" --priority medium
Added subtask "Get passport photos taken" (IT_d4e5f6) to "Renew passport before trip"```

### `listello subtask link`

> Event storming: **Owner → Link as child → Item linked as child of item** (domain implemented)

Turns an existing item into a subtask of another. Only one level of nesting is
allowed — the domain rejects linking under an item that already has a parent.

```text
SYNOPSIS
    listello subtask link <item> --parent <parent-item>
```

```console
$ listello subtask link IT_a1b2c3 --parent IT_4b81d0
Linked "Gather documents" as a subtask of "Renew passport before trip"```

```console
$ listello subtask link IT_x1 --parent IT_d4e5f6
error: cannot link as child of an item that already has a parent
```

### `listello subtask prioritize`

> Event storming: **Owner → Prioritize subtask → Subtask priority changed** (domain implemented)

```text
SYNOPSIS
    listello subtask prioritize <item> <low|medium|high|none>
```

```console
$ listello subtask prioritize IT_d4e5f6 high
Priority of "Get passport photos taken" changed to high```

### `listello subtask complete` / `uncomplete` / `delete`

> Event storming: **Owner → Complete subtask / Uncomplete subtask / Delete subtask**

Same behavior and output as the corresponding `item` commands, shown here for
completeness since the board models them as distinct commands:

```console
$ listello subtask complete IT_a1b2c3
✔ Completed subtask "Gather documents" (2/3 done on "Renew passport before trip")
$ listello subtask uncomplete IT_a1b2c3
Reopened subtask "Gather documents"
$ listello subtask delete IT_g7h8i9
Deleted subtask "Fill out DS-82 form" from "Renew passport before trip"```

---

## `listello comment`

### `listello comment add`

> Event storming: **User → Comment → Item commented on**

```text
SYNOPSIS
    listello comment add <item> <text>
```

```console
$ listello comment add IT_4b81d0 "Post office closes at 4pm on Saturdays"
Commented on "Renew passport before trip" (CM_77e1af)```

### `listello comment delete`

> Event storming: **User → Delete comment → Item comment deleted**

```text
SYNOPSIS
    listello comment delete <comment>
```

```console
$ listello comment delete CM_77e1af
Deleted comment CM_77e1af from "Renew passport before trip"```

---

## `listello delegation`

Covers both the Delegation Policy aggregate (`allow`/`deny`) and the
Delegation aggregate (`assign`/`ask`/`answer`).

### `listello delegation allow`

> Event storming: **Owner → Allow delegation → Delegation permission granted to other person**

```text
SYNOPSIS
    listello delegation allow <person>
```

```console
$ listello delegation allow ana@example.com
ana@example.com may now be delegated items```

### `listello delegation deny`

> Event storming: **Owner → Deny delegation → Delegation permission revoked from other person**

```text
SYNOPSIS
    listello delegation deny <person>
```

```console
$ listello delegation deny ana@example.com
ana@example.com may no longer be delegated items```

### `listello delegation assign`

> Event storming: **Owner → Delegate item → Item delegated**

```text
SYNOPSIS
    listello delegation assign <item> --to <person>
```

```console
$ listello delegation assign IT_77aa03 --to ana@example.com
Delegated "Draft Q3 planning doc" to ana@example.com```

Delegating to someone without permission fails against the policy aggregate:

```console
$ listello delegation assign IT_77aa03 --to sam@example.com
error: sam@example.com does not have delegation permission
```

### `listello delegation ask`

> Event storming: **Owner → Request clarification → Clarification requested**

```text
SYNOPSIS
    listello delegation ask <item> <question>
```

```console
$ listello delegation ask IT_77aa03 "Which teams should the plan cover?"
Requested clarification on "Draft Q3 planning doc" from ana@example.com```

### `listello delegation answer`

> Event storming: **Assignee → Provide clarification → Clarification provided**

```text
SYNOPSIS
    listello delegation answer <item> <answer>
```

```console
$ listello delegation answer IT_77aa03 "Platform and Growth only"
Provided clarification on "Draft Q3 planning doc"```

### `listello delegation show`

> Interfaces: **Delegated items view**, **Assigned to me view**, **Needs clarification view** (queries, no domain events)

```text
SYNOPSIS
    listello delegation show [--delegated | --assigned | --needs-clarification]
```

```console
$ listello delegation show --delegated
Delegated by me
  IT_77aa03  Draft Q3 planning doc   → ana@example.com   needs clarification

$ listello delegation show --assigned
Assigned to me
  IT_c0ffee  Review release notes    ← sam@example.com   due 2026-08-05

$ listello delegation show --needs-clarification
Needs clarification
  IT_77aa03  Draft Q3 planning doc   "Which teams should the plan cover?"
```

---

## Coverage Map

How each event storming command maps onto the CLI:

| Aggregate | Board command | CLI command |
|---|---|---|
| List | Create list | `listello list create` |
| List | Delete list | `listello list delete` |
| Item | Capture inbox item | `listello item capture` / `listello inbox` |
| Item | Define item | `listello item define` |
| Item | Modify item title | `listello item title` |
| Item | Modify item description | `listello item describe` |
| Item | Modify due date | `listello item due` |
| Item | Remove due date | `listello item undue` |
| Item | Tag item | `listello item tag` |
| Item | Untag item | `listello item untag` |
| Item | Move item | `listello item move` |
| Item | Complete item | `listello item complete` |
| Item | Uncomplete item | `listello item uncomplete` |
| Item | Delete item | `listello item delete` |
| Item | Add subtask | `listello subtask add` |
| Item | Link as child | `listello subtask link` |
| Item | Prioritize subtask | `listello subtask prioritize` |
| Item | Complete subtask | `listello subtask complete` |
| Item | Uncomplete subtask | `listello subtask uncomplete` |
| Item | Delete subtask | `listello subtask delete` |
| Item | Comment | `listello comment add` |
| Item | Delete comment | `listello comment delete` |
| Delegation Policy | Allow delegation | `listello delegation allow` |
| Delegation Policy | Deny delegation | `listello delegation deny` |
| Delegation | Delegate item | `listello delegation assign` |
| Delegation | Request clarification | `listello delegation ask` |
| Delegation | Provide clarification | `listello delegation answer` |

The `show` sub-commands (`list show`, `item show`, `delegation show`) have no
board command — they realize the interfaces listed per aggregate (List view,
Item view, and the three delegation views).

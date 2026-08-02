# Listello Event Storming Command Map

This document organizes the Listello event storming board by aggregate and maps:

- the actor invoking each command;
- the command being invoked;
- the domain event triggered by that command; and
- the interfaces associated with each aggregate.

> Some command names have been normalized slightly for clarity. For example, the board's `Subtask` command is represented as `Add subtask`, and `Prioritize` is represented as `Prioritize subtask`, based on their corresponding domain events.

---

# Delegation Policy Aggregate

## Interfaces

- Delegation policies view
- Create delegation policy view

## Actor, Command, and Domain Event Map

| Actor | Command | Domain Event | Implemented |
|---|---|---|---|
| Owner | Allow delegation | Delegation permission granted to other person | |
| Owner | Deny delegation | Delegation permission revoked from other person | |

---

# Delegation Aggregate

## Interfaces

- Delegated items view
- Assigned to me view
- Needs clarification view

## Actor, Command, and Domain Event Map

| Actor | Command | Domain Event | Implemented |
|---|---|---|---|
| Owner | Delegate item | Item delegated | |
| Owner | Request clarification | Clarification requested | |
| Assignee | Provide clarification | Clarification provided | |

---

# List Aggregate

## Interfaces

- List view
- Create list view

## Actor, Command, and Domain Event Map

| Actor | Command | Domain Event | Implemented |
|---|---|---|---|
| User | Create list | List created | ✅ |
| User | Delete list | List deleted | |

---

# Item Aggregate

## Interfaces

- Item view
- Item list
- Create item view
- Comment view

## Item Creation

| Actor | Command | Domain Event | Implemented |
|---|---|---|---|
| User | Capture inbox item | Item captured | ✅ |
| User | Define item | Item defined | ✅ |

## Item Details

| Actor | Command | Domain Event | Implemented |
|---|---|---|---|
| Owner | Modify item description | Item description changed | ✅ |
| Owner | Modify item title | Item title changed | ✅ |
| Requester or Owner | Modify due date | Due date added to item | ✅ |
| Requester or Owner | Remove due date | Due date removed from item | |
| Owner | Tag item | Tag added to item | ✅ |
| Owner | Untag item | Tag removed from item | |
| Owner | Link as child | Item linked as child of item | |
| Owner | Prioritize subtask | Subtask priority changed | ✅ |

## List Membership

| Actor | Command | Domain Event | Implemented |
|---|---|---|---|
| Owner | Move item | Item moved to other list | ✅ |

## Subtasks

| Actor | Command | Domain Event | Implemented |
|---|---|---|---|
| Owner | Add subtask | Subtask added to item | |
| Owner | Complete subtask | Subtask completed on item | |
| Owner | Uncomplete subtask | Subtask uncompleted | |
| Owner | Delete subtask | Subtask deleted on item | |

## Comments

| Actor | Command | Domain Event | Implemented |
|---|---|---|---|
| User | Comment | Item commented on | |
| User | Delete comment | Item comment deleted | |

## Item Lifecycle

| Actor | Command | Domain Event | Implemented |
|---|---|---|---|
| Requester or Owner | Complete item | Item completed | ✅ |
| Requester or Owner | Uncomplete item | Item uncompleted | |
| Requester or Owner | Delete item | Item deleted | |

---

# Condensed Command Map

```text
Delegation Policy
Owner
  ├─ Allow delegation
  │    └─ Delegation permission granted to other person
  └─ Deny delegation
       └─ Delegation permission revoked from other person

Delegation
Owner
  ├─ Delegate item
  │    └─ Item delegated
  └─ Request clarification
       └─ Clarification requested

Assignee
  └─ Provide clarification
       └─ Clarification provided

List
User
  ├─ Create list
  │    └─ List created
  └─ Delete list
       └─ List deleted

Item
User
  ├─ Capture inbox item
  │    └─ Item captured
  ├─ Define item
  │    └─ Item defined
  ├─ Comment
  │    └─ Item commented on
  └─ Delete comment
       └─ Item comment deleted

Owner
  ├─ Modify item description
  │    └─ Item description changed
  ├─ Modify item title
  │    └─ Item title changed
  ├─ Tag item
  │    └─ Tag added to item
  ├─ Untag item
  │    └─ Tag removed from item
  ├─ Link as child
  │    └─ Item linked as child of item
  ├─ Prioritize subtask
  │    └─ Subtask priority changed
  ├─ Move item
  │    └─ Item moved to other list
  ├─ Add subtask
  │    └─ Subtask added to item
  ├─ Complete subtask
  │    └─ Subtask completed on item
  ├─ Uncomplete subtask
  │    └─ Subtask uncompleted
  └─ Delete subtask
       └─ Subtask deleted on item

Requester or Owner
  ├─ Modify due date
  │    └─ Due date added to item
  ├─ Remove due date
  │    └─ Due date removed from item
  ├─ Complete item
  │    └─ Item completed
  ├─ Uncomplete item
  │    └─ Item uncompleted
  └─ Delete item
       └─ Item deleted
```

---

# Aggregate Inventory

| Aggregate | Interfaces |
|---|---|
| Delegation Policy | Delegation policies view; Create delegation policy view |
| Delegation | Delegated items view; Assigned to me view; Needs clarification view |
| List | List view; Create list view |
| Item | Item view; Item list; Create item view; Comment view |

## Actors Identified

- User
- Owner
- Requester
- Assignee

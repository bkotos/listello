// Domain model derived from the event storming board.
// Aggregates: List, Item, Subtask (part of Item), Delegation (stubbed in UI).

export type Priority = 'none' | 'low' | 'medium' | 'high'

export type ItemState = 'inbox' | 'defined'

export interface Subtask {
  id: string
  title: string
  completed: boolean
}

export interface Comment {
  id: string
  author: string
  body: string
  createdAt: number
}

// Domain events recorded for the activity feed (green "views" on the board).
export interface DomainEvent {
  id: string
  itemId: string
  name: string // e.g. "Item captured", "Tag added to item"
  detail?: string
  actor: string
  at: number
}

export interface Item {
  id: string
  listId: string
  title: string
  description: string
  state: ItemState
  completed: boolean
  priority: Priority
  dueDate: string | null // ISO yyyy-mm-dd
  tags: string[]
  subtasks: Subtask[]
  comments: Comment[]
  createdAt: number
}

export interface List {
  id: string
  name: string
}

export interface AppState {
  lists: List[]
  items: Item[]
  tags: string[]
  events: DomainEvent[]
  selectedListId: string // "inbox" is a virtual list
  selectedItemId: string | null
}

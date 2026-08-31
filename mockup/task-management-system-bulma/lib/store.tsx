'use client'

import {
  createContext,
  useContext,
  useMemo,
  useReducer,
  type ReactNode,
  type Dispatch,
} from 'react'
import type { AppState, DomainEvent, Item, Priority } from './types'
import { seedState } from './seed'

export const INBOX_ID = 'inbox'
export const CURRENT_USER = 'You'

const uid = () => Math.random().toString(36).slice(2, 10)

// Actions map 1:1 to the commands on the event storming board.
type Action =
  | { type: 'CREATE_LIST'; name: string }
  | { type: 'DELETE_LIST'; listId: string }
  | { type: 'SELECT_LIST'; listId: string }
  | { type: 'SELECT_ITEM'; itemId: string | null }
  | { type: 'CAPTURE_INBOX_ITEM'; title: string; listId?: string }
  | { type: 'DEFINE_ITEM'; itemId: string }
  | { type: 'MODIFY_TITLE'; itemId: string; title: string }
  | { type: 'MODIFY_DESCRIPTION'; itemId: string; description: string }
  | { type: 'SET_DUE_DATE'; itemId: string; dueDate: string }
  | { type: 'REMOVE_DUE_DATE'; itemId: string }
  | { type: 'ADD_TAG'; itemId: string; tag: string }
  | { type: 'REMOVE_TAG'; itemId: string; tag: string }
  | { type: 'PRIORITIZE'; itemId: string; priority: Priority }
  | { type: 'MOVE_ITEM'; itemId: string; listId: string }
  | { type: 'COMPLETE_ITEM'; itemId: string }
  | { type: 'UNCOMPLETE_ITEM'; itemId: string }
  | { type: 'DELETE_ITEM'; itemId: string }
  | { type: 'ADD_SUBTASK'; itemId: string; title: string }
  | { type: 'COMPLETE_SUBTASK'; itemId: string; subtaskId: string }
  | { type: 'UNCOMPLETE_SUBTASK'; itemId: string; subtaskId: string }
  | { type: 'DELETE_SUBTASK'; itemId: string; subtaskId: string }
  | { type: 'ADD_COMMENT'; itemId: string; body: string }
  | { type: 'DELETE_COMMENT'; itemId: string; commentId: string }

function recordEvent(
  state: AppState,
  itemId: string,
  name: string,
  detail?: string,
  actor = CURRENT_USER,
): DomainEvent[] {
  const event: DomainEvent = { id: uid(), itemId, name, detail, actor, at: Date.now() }
  return [event, ...state.events].slice(0, 200)
}

function mapItem(state: AppState, itemId: string, fn: (item: Item) => Item): Item[] {
  return state.items.map((it) => (it.id === itemId ? fn(it) : it))
}

function reducer(state: AppState, action: Action): AppState {
  switch (action.type) {
    case 'CREATE_LIST': {
      const name = action.name.trim()
      if (!name) return state
      const id = uid()
      return { ...state, lists: [...state.lists, { id, name }], selectedListId: id }
    }
    case 'DELETE_LIST': {
      const items = state.items.filter((it) => it.listId !== action.listId)
      const lists = state.lists.filter((l) => l.id !== action.listId)
      return {
        ...state,
        lists,
        items,
        selectedListId: state.selectedListId === action.listId ? INBOX_ID : state.selectedListId,
        selectedItemId: null,
      }
    }
    case 'SELECT_LIST':
      return { ...state, selectedListId: action.listId, selectedItemId: null }
    case 'SELECT_ITEM':
      return { ...state, selectedItemId: action.itemId }

    case 'CAPTURE_INBOX_ITEM': {
      const title = action.title.trim()
      if (!title) return state
      const id = uid()
      const toList = action.listId && action.listId !== INBOX_ID ? action.listId : INBOX_ID
      const item: Item = {
        id,
        listId: toList,
        title,
        description: '',
        state: toList === INBOX_ID ? 'inbox' : 'defined',
        completed: false,
        priority: 'none',
        dueDate: null,
        tags: [],
        subtasks: [],
        comments: [],
        createdAt: Date.now(),
      }
      return {
        ...state,
        items: [item, ...state.items],
        events: recordEvent(state, id, 'Item captured', title),
      }
    }
    case 'DEFINE_ITEM':
      return {
        ...state,
        items: mapItem(state, action.itemId, (it) => ({ ...it, state: 'defined' })),
        events: recordEvent(state, action.itemId, 'Defined created'),
      }
    case 'MODIFY_TITLE':
      return {
        ...state,
        items: mapItem(state, action.itemId, (it) => ({ ...it, title: action.title })),
        events: recordEvent(state, action.itemId, 'Item title changed', action.title),
      }
    case 'MODIFY_DESCRIPTION':
      return {
        ...state,
        items: mapItem(state, action.itemId, (it) => ({ ...it, description: action.description })),
        events: recordEvent(state, action.itemId, 'Item description changed'),
      }
    case 'SET_DUE_DATE':
      return {
        ...state,
        items: mapItem(state, action.itemId, (it) => ({ ...it, dueDate: action.dueDate })),
        events: recordEvent(state, action.itemId, 'Due date added to item', action.dueDate),
      }
    case 'REMOVE_DUE_DATE':
      return {
        ...state,
        items: mapItem(state, action.itemId, (it) => ({ ...it, dueDate: null })),
        events: recordEvent(state, action.itemId, 'Due date removed from item'),
      }
    case 'ADD_TAG': {
      const tag = action.tag.trim()
      if (!tag) return state
      return {
        ...state,
        tags: state.tags.includes(tag) ? state.tags : [...state.tags, tag],
        items: mapItem(state, action.itemId, (it) =>
          it.tags.includes(tag) ? it : { ...it, tags: [...it.tags, tag] },
        ),
        events: recordEvent(state, action.itemId, 'Tag added to item', tag),
      }
    }
    case 'REMOVE_TAG':
      return {
        ...state,
        items: mapItem(state, action.itemId, (it) => ({
          ...it,
          tags: it.tags.filter((t) => t !== action.tag),
        })),
        events: recordEvent(state, action.itemId, 'Tag removed from item', action.tag),
      }
    case 'PRIORITIZE':
      return {
        ...state,
        items: mapItem(state, action.itemId, (it) => ({ ...it, priority: action.priority })),
        events: recordEvent(state, action.itemId, 'Item priority changed', action.priority),
      }
    case 'MOVE_ITEM': {
      const target = state.lists.find((l) => l.id === action.listId)
      return {
        ...state,
        items: mapItem(state, action.itemId, (it) => ({
          ...it,
          listId: action.listId,
          state: 'defined',
        })),
        events: recordEvent(
          state,
          action.itemId,
          'Item moved to other list',
          target ? target.name : 'Inbox',
        ),
      }
    }
    case 'COMPLETE_ITEM':
      return {
        ...state,
        items: mapItem(state, action.itemId, (it) => ({ ...it, completed: true })),
        events: recordEvent(state, action.itemId, 'Item completed'),
      }
    case 'UNCOMPLETE_ITEM':
      return {
        ...state,
        items: mapItem(state, action.itemId, (it) => ({ ...it, completed: false })),
        events: recordEvent(state, action.itemId, 'Item uncompleted'),
      }
    case 'DELETE_ITEM':
      return {
        ...state,
        items: state.items.filter((it) => it.id !== action.itemId),
        selectedItemId: state.selectedItemId === action.itemId ? null : state.selectedItemId,
        events: recordEvent(state, action.itemId, 'Item deleted'),
      }

    case 'ADD_SUBTASK': {
      const title = action.title.trim()
      if (!title) return state
      return {
        ...state,
        items: mapItem(state, action.itemId, (it) => ({
          ...it,
          subtasks: [...it.subtasks, { id: uid(), title, completed: false }],
        })),
        events: recordEvent(state, action.itemId, 'Subtask added to item', title),
      }
    }
    case 'COMPLETE_SUBTASK':
      return {
        ...state,
        items: mapItem(state, action.itemId, (it) => ({
          ...it,
          subtasks: it.subtasks.map((s) =>
            s.id === action.subtaskId ? { ...s, completed: true } : s,
          ),
        })),
        events: recordEvent(state, action.itemId, 'Subtask completed on item'),
      }
    case 'UNCOMPLETE_SUBTASK':
      return {
        ...state,
        items: mapItem(state, action.itemId, (it) => ({
          ...it,
          subtasks: it.subtasks.map((s) =>
            s.id === action.subtaskId ? { ...s, completed: false } : s,
          ),
        })),
        events: recordEvent(state, action.itemId, 'Subtask uncompleted'),
      }
    case 'DELETE_SUBTASK':
      return {
        ...state,
        items: mapItem(state, action.itemId, (it) => ({
          ...it,
          subtasks: it.subtasks.filter((s) => s.id !== action.subtaskId),
        })),
        events: recordEvent(state, action.itemId, 'Subtask deleted on item'),
      }

    case 'ADD_COMMENT': {
      const body = action.body.trim()
      if (!body) return state
      return {
        ...state,
        items: mapItem(state, action.itemId, (it) => ({
          ...it,
          comments: [
            ...it.comments,
            { id: uid(), author: CURRENT_USER, body, createdAt: Date.now() },
          ],
        })),
        events: recordEvent(state, action.itemId, 'Item commented on'),
      }
    }
    case 'DELETE_COMMENT':
      return {
        ...state,
        items: mapItem(state, action.itemId, (it) => ({
          ...it,
          comments: it.comments.filter((c) => c.id !== action.commentId),
        })),
        events: recordEvent(state, action.itemId, 'Item comment deleted'),
      }

    default:
      return state
  }
}

const StoreContext = createContext<{ state: AppState; dispatch: Dispatch<Action> } | null>(null)

export function StoreProvider({ children }: { children: ReactNode }) {
  const [state, dispatch] = useReducer(reducer, undefined, seedState)
  const value = useMemo(() => ({ state, dispatch }), [state])
  return <StoreContext.Provider value={value}>{children}</StoreContext.Provider>
}

export function useStore() {
  const ctx = useContext(StoreContext)
  if (!ctx) throw new Error('useStore must be used within StoreProvider')
  return ctx
}

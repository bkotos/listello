'use client'

import { useEffect, useRef, useState } from 'react'
import {
  Calendar,
  MessageSquare,
  ListChecks,
  MoreHorizontal,
  Pencil,
  Trash2,
  Inbox,
  Send,
} from 'lucide-react'
import type { Item } from '@/lib/types'
import { useStore, INBOX_ID } from '@/lib/store'
import { formatDue, priorityMeta } from '@/lib/format'

type Popover = 'menu' | 'date' | 'comment' | null

const CheckIcon = () => (
  <svg viewBox="0 0 12 12" width={12} height={12} fill="none" aria-hidden>
    <path
      d="M2.5 6.2 5 8.5 9.5 3.5"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
)

export function ItemRow({ item }: { item: Item }) {
  const { dispatch, state } = useStore()
  const active = state.selectedItemId === item.id
  const due = formatDue(item.dueDate)
  const pr = priorityMeta(item.priority)
  const doneSubs = item.subtasks.filter((s) => s.completed).length

  const [popover, setPopover] = useState<Popover>(null)
  const [editing, setEditing] = useState(false)
  const [titleDraft, setTitleDraft] = useState(item.title)
  const [commentDraft, setCommentDraft] = useState('')
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!popover) return
    const onClick = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setPopover(null)
      }
    }
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setPopover(null)
    }
    document.addEventListener('mousedown', onClick)
    document.addEventListener('keydown', onKey)
    return () => {
      document.removeEventListener('mousedown', onClick)
      document.removeEventListener('keydown', onKey)
    }
  }, [popover])

  const select = () => dispatch({ type: 'SELECT_ITEM', itemId: item.id })

  const startRename = () => {
    setTitleDraft(item.title)
    setEditing(true)
    setPopover(null)
  }

  const commitRename = () => {
    const next = titleDraft.trim()
    if (next && next !== item.title) {
      dispatch({ type: 'MODIFY_TITLE', itemId: item.id, title: next })
    }
    setEditing(false)
  }

  const submitComment = () => {
    const body = commentDraft.trim()
    if (!body) return
    dispatch({ type: 'ADD_COMMENT', itemId: item.id, body })
    setCommentDraft('')
    setPopover(null)
  }

  const hasMeta =
    item.tags.length > 0 || item.dueDate || item.subtasks.length > 0 || item.comments.length > 0
  const showQuickActions = !item.dueDate || item.comments.length === 0

  return (
    <div
      ref={containerRef}
      role="button"
      tabIndex={0}
      onClick={select}
      onKeyDown={(e) => {
        if (editing) return
        if (e.key === 'Enter' || e.key === ' ') {
          e.preventDefault()
          select()
        }
      }}
      className={`task-row hover-parent is-flex p-3 ${active ? 'is-active' : ''}`}
      style={{ gap: '0.75rem', alignItems: 'flex-start' }}
    >
      <button
        type="button"
        aria-label={item.completed ? 'Mark incomplete' : 'Mark complete'}
        onClick={(e) => {
          e.stopPropagation()
          dispatch({
            type: item.completed ? 'UNCOMPLETE_ITEM' : 'COMPLETE_ITEM',
            itemId: item.id,
          })
        }}
        className={`check-toggle ${item.completed ? 'is-checked' : ''}`}
        style={{ height: '1.25rem', width: '1.25rem', marginTop: '0.125rem' }}
      >
        {item.completed && <CheckIcon />}
      </button>

      <div style={{ minWidth: 0, flex: '1 1 0' }}>
        <div className="is-flex" style={{ gap: '0.5rem', alignItems: 'flex-start' }}>
          {item.priority !== 'none' && !editing && (
            <span
              aria-hidden
              style={{
                marginTop: '0.4rem',
                height: '0.5rem',
                width: '0.5rem',
                flexShrink: 0,
                borderRadius: '9999px',
                backgroundColor: pr.dot,
              }}
            />
          )}
          {editing ? (
            <input
              autoFocus
              className="input is-small"
              value={titleDraft}
              onClick={(e) => e.stopPropagation()}
              onChange={(e) => setTitleDraft(e.target.value)}
              onBlur={commitRename}
              onKeyDown={(e) => {
                e.stopPropagation()
                if (e.key === 'Enter' && !e.nativeEvent.isComposing && e.keyCode !== 229) {
                  e.preventDefault()
                  commitRename()
                } else if (e.key === 'Escape') {
                  setEditing(false)
                }
              }}
            />
          ) : (
            <p className={`is-size-6 ${item.completed ? 'muted line-through' : ''}`}>{item.title}</p>
          )}
        </div>

        {(hasMeta || showQuickActions) && !editing && (
          <div
            className="is-flex is-flex-wrap-wrap is-align-items-center mt-2 is-size-7 muted"
            style={{ gap: '0.5rem' }}
          >
            {item.dueDate && (
              <button
                type="button"
                aria-label="Change date"
                onClick={(e) => {
                  e.stopPropagation()
                  setPopover(popover === 'date' ? null : 'date')
                }}
                className={`icon-btn ${popover === 'date' ? 'is-active' : ''}`}
                style={{
                  gap: '0.25rem',
                  padding: '0.125rem 0.375rem',
                  color:
                    due.tone === 'past'
                      ? 'hsl(348deg 86% 43%)'
                      : due.tone === 'today'
                        ? 'hsl(var(--bulma-primary-h) var(--bulma-primary-s) var(--bulma-primary-l))'
                        : undefined,
                }}
              >
                <Calendar size={14} />
                {due.label}
              </button>
            )}
            {item.subtasks.length > 0 && (
              <span className="is-inline-flex is-align-items-center" style={{ gap: '0.25rem' }}>
                <ListChecks size={14} />
                {doneSubs}/{item.subtasks.length}
              </span>
            )}
            {item.comments.length > 0 && (
              <button
                type="button"
                aria-label="Add comment"
                onClick={(e) => {
                  e.stopPropagation()
                  setCommentDraft('')
                  setPopover(popover === 'comment' ? null : 'comment')
                }}
                className={`icon-btn ${popover === 'comment' ? 'is-active' : ''}`}
                style={{ gap: '0.25rem', padding: '0.125rem 0.375rem' }}
              >
                <MessageSquare size={14} />
                {item.comments.length}
              </button>
            )}
            {item.tags.map((tag) => (
              <span key={tag} className="tag is-primary is-light is-family-code is-uppercase">
                {tag}
              </span>
            ))}

            {showQuickActions && (
              <span
                className={`is-inline-flex is-align-items-center hover-reveal ${
                  popover === 'date' || popover === 'comment' ? 'is-visible' : ''
                }`}
                style={{ gap: '0.25rem' }}
              >
                {!item.dueDate && (
                  <HoverAction
                    label="Set date"
                    active={popover === 'date'}
                    onClick={(e) => {
                      e.stopPropagation()
                      setPopover(popover === 'date' ? null : 'date')
                    }}
                  >
                    <Calendar size={14} />
                  </HoverAction>
                )}
                {item.comments.length === 0 && (
                  <HoverAction
                    label="Add comment"
                    active={popover === 'comment'}
                    onClick={(e) => {
                      e.stopPropagation()
                      setCommentDraft('')
                      setPopover(popover === 'comment' ? null : 'comment')
                    }}
                  >
                    <MessageSquare size={14} />
                  </HoverAction>
                )}
              </span>
            )}
          </div>
        )}

        {popover === 'date' && !editing && (
          <div
            onClick={(e) => e.stopPropagation()}
            className="box is-inline-flex is-align-items-center mt-2 p-2"
            style={{ gap: '0.5rem' }}
          >
            <input
              type="date"
              autoFocus
              className="input is-small"
              style={{ width: 'auto' }}
              value={item.dueDate ?? ''}
              onChange={(e) => {
                const v = e.target.value
                if (v) dispatch({ type: 'SET_DUE_DATE', itemId: item.id, dueDate: v })
                else dispatch({ type: 'REMOVE_DUE_DATE', itemId: item.id })
              }}
            />
            {item.dueDate && (
              <button
                type="button"
                className="button is-small is-danger is-light"
                onClick={() => {
                  dispatch({ type: 'REMOVE_DUE_DATE', itemId: item.id })
                  setPopover(null)
                }}
              >
                Clear
              </button>
            )}
          </div>
        )}

        {popover === 'comment' && !editing && (
          <div onClick={(e) => e.stopPropagation()} className="field has-addons mt-2 mb-0">
            <div className="control is-expanded">
              <input
                autoFocus
                className="input is-small"
                value={commentDraft}
                placeholder="Leave a comment…"
                onChange={(e) => setCommentDraft(e.target.value)}
                onKeyDown={(e) => {
                  e.stopPropagation()
                  if (e.key === 'Enter' && !e.nativeEvent.isComposing && e.keyCode !== 229) {
                    e.preventDefault()
                    submitComment()
                  }
                }}
              />
            </div>
            <div className="control">
              <button
                type="button"
                aria-label="Post comment"
                onClick={submitComment}
                disabled={!commentDraft.trim()}
                className="button is-small is-primary"
              >
                <Send size={12} />
              </button>
            </div>
          </div>
        )}
      </div>

      {!editing && (
        <div className={`dropdown is-right ${popover === 'menu' ? 'is-active' : ''}`}>
          <div className="dropdown-trigger">
            <button
              type="button"
              aria-label="Task options"
              aria-haspopup="menu"
              aria-expanded={popover === 'menu'}
              onClick={(e) => {
                e.stopPropagation()
                setPopover(popover === 'menu' ? null : 'menu')
              }}
              className={`icon-btn ${popover === 'menu' ? 'is-active' : 'hover-reveal'}`}
              style={{ height: '1.75rem', width: '1.75rem' }}
            >
              <MoreHorizontal size={16} />
            </button>
          </div>

          {popover === 'menu' && (
            <div className="dropdown-menu" role="menu" onClick={(e) => e.stopPropagation()}>
              <div className="dropdown-content">
                <MenuItem icon={<Pencil size={16} />} label="Rename" onClick={startRename} />
                {item.listId !== INBOX_ID && (
                  <MenuItem
                    icon={<Inbox size={16} />}
                    label="Move to Inbox"
                    onClick={() => {
                      dispatch({ type: 'MOVE_ITEM', itemId: item.id, listId: INBOX_ID })
                      setPopover(null)
                    }}
                  />
                )}
                <hr className="dropdown-divider" />
                <MenuItem
                  icon={<Trash2 size={16} />}
                  label="Delete"
                  destructive
                  onClick={() => {
                    dispatch({ type: 'DELETE_ITEM', itemId: item.id })
                    setPopover(null)
                  }}
                />
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  )
}

function HoverAction({
  label,
  active,
  onClick,
  children,
}: {
  label: string
  active?: boolean
  onClick: (e: React.MouseEvent) => void
  children: React.ReactNode
}) {
  return (
    <button
      type="button"
      aria-label={label}
      title={label}
      onClick={onClick}
      className={`icon-btn ${active ? 'is-active' : ''}`}
      style={{ height: '1.5rem', padding: '0 0.375rem', gap: '0.25rem' }}
    >
      {children}
    </button>
  )
}

function MenuItem({
  icon,
  label,
  onClick,
  destructive,
}: {
  icon: React.ReactNode
  label: string
  onClick: () => void
  destructive?: boolean
}) {
  return (
    <a
      role="menuitem"
      onClick={onClick}
      className="dropdown-item is-flex is-align-items-center"
      style={{ gap: '0.625rem', color: destructive ? 'hsl(348deg 86% 43%)' : undefined }}
    >
      <span style={{ display: 'inline-flex', color: destructive ? 'hsl(348deg 86% 43%)' : 'hsl(0deg 0% 45%)' }}>
        {icon}
      </span>
      <span>{label}</span>
    </a>
  )
}

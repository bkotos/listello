'use client'

import { useState } from 'react'
import {
  X,
  Trash2,
  Calendar,
  Flag,
  Tag as TagIcon,
  CornerDownRight,
  Send,
  Plus,
  ChevronRight,
} from 'lucide-react'
import type { Item, Priority } from '@/lib/types'
import { useStore, INBOX_ID } from '@/lib/store'
import { PRIORITIES, priorityMeta, formatDue, relativeTime } from '@/lib/format'

const CheckIcon = ({ size = 12 }: { size?: number }) => (
  <svg viewBox="0 0 12 12" width={size} height={size} fill="none" aria-hidden>
    <path
      d="M2.5 6.2 5 8.5 9.5 3.5"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
)

export function DetailPanel({ item, onClose }: { item: Item; onClose: () => void }) {
  const { state, dispatch } = useStore()
  const [subtask, setSubtask] = useState('')
  const [comment, setComment] = useState('')
  const [tagInput, setTagInput] = useState('')
  const [showTagInput, setShowTagInput] = useState(false)

  const list = state.lists.find((l) => l.id === item.listId)
  const listLabel = item.listId === INBOX_ID ? 'Inbox' : (list?.name ?? 'Inbox')
  const due = formatDue(item.dueDate)
  const itemEvents = state.events.filter((e) => e.itemId === item.id).slice(0, 8)

  return (
    <div className="is-flex is-flex-direction-column" style={{ height: '100%', backgroundColor: 'hsl(0deg 0% 100%)' }}>
      <header
        className="is-flex is-align-items-center is-justify-content-space-between px-5 py-3"
        style={{ gap: '0.5rem', borderBottom: '1px solid var(--app-border)' }}
      >
        <div className="is-flex is-align-items-center is-size-7 muted" style={{ gap: '0.5rem' }}>
          <span className="tag is-primary is-light is-family-code is-uppercase">{listLabel}</span>
          {item.state === 'inbox' && (
            <span className="is-family-code is-uppercase" style={{ letterSpacing: '0.12em' }}>
              unclarified
            </span>
          )}
        </div>
        <div className="is-flex is-align-items-center" style={{ gap: '0.25rem' }}>
          <button
            type="button"
            aria-label="Delete item"
            onClick={() => dispatch({ type: 'DELETE_ITEM', itemId: item.id })}
            className="icon-btn is-danger-hover"
            style={{ height: '2rem', width: '2rem' }}
          >
            <Trash2 size={16} />
          </button>
          <button
            type="button"
            aria-label="Close"
            onClick={onClose}
            className="icon-btn"
            style={{ height: '2rem', width: '2rem' }}
          >
            <X size={16} />
          </button>
        </div>
      </header>

      <div className="app-scroll">
        <div className="is-flex is-flex-direction-column px-5 py-5" style={{ gap: '1.5rem' }}>
          {/* Title + complete */}
          <div className="is-flex" style={{ gap: '0.75rem', alignItems: 'flex-start' }}>
            <button
              type="button"
              aria-label={item.completed ? 'Mark incomplete' : 'Mark complete'}
              onClick={() =>
                dispatch({
                  type: item.completed ? 'UNCOMPLETE_ITEM' : 'COMPLETE_ITEM',
                  itemId: item.id,
                })
              }
              className={`check-toggle ${item.completed ? 'is-checked' : ''}`}
              style={{ height: '1.5rem', width: '1.5rem', marginTop: '0.25rem' }}
            >
              {item.completed && <CheckIcon size={14} />}
            </button>
            <TitleField item={item} />
          </div>

          {/* Property list */}
          <div className="box p-0" style={{ boxShadow: 'none', border: '1px solid var(--app-border)' }}>
            <PropertyRow icon={<Calendar size={16} />} label="Due date" divider={false}>
              <div className="is-flex is-align-items-center" style={{ gap: '0.5rem' }}>
                <input
                  type="date"
                  className="input is-small"
                  style={{ width: 'auto' }}
                  value={item.dueDate ?? ''}
                  onChange={(e) =>
                    e.target.value
                      ? dispatch({ type: 'SET_DUE_DATE', itemId: item.id, dueDate: e.target.value })
                      : dispatch({ type: 'REMOVE_DUE_DATE', itemId: item.id })
                  }
                />
                {item.dueDate && (
                  <button
                    type="button"
                    onClick={() => dispatch({ type: 'REMOVE_DUE_DATE', itemId: item.id })}
                    className="button is-small is-ghost"
                    style={{ color: due.tone === 'past' ? 'hsl(348deg 86% 43%)' : undefined }}
                  >
                    clear
                  </button>
                )}
              </div>
            </PropertyRow>

            <PropertyRow icon={<Flag size={16} />} label="Priority">
              <div className="is-flex is-flex-wrap-wrap is-align-items-center" style={{ gap: '0.375rem' }}>
                {PRIORITIES.map((p) => {
                  const activeP = item.priority === p.value
                  return (
                    <button
                      key={p.value}
                      type="button"
                      onClick={() =>
                        dispatch({
                          type: 'PRIORITIZE',
                          itemId: item.id,
                          priority: p.value as Priority,
                        })
                      }
                      className={`button is-small ${activeP ? 'is-primary is-light' : 'is-white'}`}
                      style={{ gap: '0.375rem' }}
                    >
                      <span
                        style={{ height: '0.5rem', width: '0.5rem', borderRadius: '9999px', backgroundColor: p.dot }}
                      />
                      {p.label}
                    </button>
                  )
                })}
              </div>
            </PropertyRow>

            <PropertyRow icon={<ChevronRight size={16} />} label="List">
              <div className="select is-small">
                <select
                  value={item.listId}
                  onChange={(e) =>
                    dispatch({ type: 'MOVE_ITEM', itemId: item.id, listId: e.target.value })
                  }
                >
                  <option value={INBOX_ID}>Inbox</option>
                  {state.lists.map((l) => (
                    <option key={l.id} value={l.id}>
                      {l.name}
                    </option>
                  ))}
                </select>
              </div>
            </PropertyRow>

            <PropertyRow icon={<TagIcon size={16} />} label="Tags">
              <div className="is-flex is-flex-wrap-wrap is-align-items-center" style={{ gap: '0.375rem' }}>
                {item.tags.map((tag) => (
                  <span key={tag} className="tags has-addons mb-0" style={{ display: 'inline-flex' }}>
                    <span className="tag is-primary is-light is-family-code is-uppercase">{tag}</span>
                    <a
                      className="tag is-delete is-light"
                      aria-label={`Remove ${tag}`}
                      onClick={() => dispatch({ type: 'REMOVE_TAG', itemId: item.id, tag })}
                    />
                  </span>
                ))}
                {showTagInput ? (
                  <input
                    autoFocus
                    list="tag-suggestions"
                    className="input is-small"
                    style={{ width: '6rem' }}
                    value={tagInput}
                    onChange={(e) => setTagInput(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.nativeEvent.isComposing || e.keyCode === 229) return
                      if (e.key === 'Enter' && tagInput.trim()) {
                        dispatch({ type: 'ADD_TAG', itemId: item.id, tag: tagInput.trim() })
                        setTagInput('')
                        setShowTagInput(false)
                      }
                      if (e.key === 'Escape') {
                        setTagInput('')
                        setShowTagInput(false)
                      }
                    }}
                    onBlur={() => {
                      if (tagInput.trim())
                        dispatch({ type: 'ADD_TAG', itemId: item.id, tag: tagInput.trim() })
                      setTagInput('')
                      setShowTagInput(false)
                    }}
                    placeholder="tag"
                  />
                ) : (
                  <button
                    type="button"
                    onClick={() => setShowTagInput(true)}
                    className="button is-small is-white"
                    style={{ gap: '0.25rem', borderStyle: 'dashed', borderColor: 'var(--app-border)' }}
                  >
                    <Plus size={12} /> add
                  </button>
                )}
                <datalist id="tag-suggestions">
                  {state.tags.map((t) => (
                    <option key={t} value={t} />
                  ))}
                </datalist>
              </div>
            </PropertyRow>
          </div>

          {/* Description */}
          <Section title="Notes">
            <DescriptionField item={item} />
          </Section>

          {/* Subtasks */}
          <Section
            title="Subtasks"
            meta={
              item.subtasks.length > 0
                ? `${item.subtasks.filter((s) => s.completed).length}/${item.subtasks.length}`
                : undefined
            }
          >
            <div className="is-flex is-flex-direction-column" style={{ gap: '0.25rem' }}>
              {item.subtasks.map((s) => (
                <div
                  key={s.id}
                  className="hover-parent is-flex is-align-items-center px-1 py-1"
                  style={{ gap: '0.625rem' }}
                >
                  <button
                    type="button"
                    aria-label={s.completed ? 'Mark subtask incomplete' : 'Mark subtask complete'}
                    onClick={() =>
                      dispatch({
                        type: s.completed ? 'UNCOMPLETE_SUBTASK' : 'COMPLETE_SUBTASK',
                        itemId: item.id,
                        subtaskId: s.id,
                      })
                    }
                    className={`check-toggle ${s.completed ? 'is-checked' : ''}`}
                    style={{ height: '1.125rem', width: '1.125rem', borderRadius: '0.25rem' }}
                  >
                    {s.completed && <CheckIcon />}
                  </button>
                  <span className={`is-size-6 ${s.completed ? 'muted line-through' : ''}`} style={{ flex: '1 1 0' }}>
                    {s.title}
                  </span>
                  <button
                    type="button"
                    aria-label="Delete subtask"
                    onClick={() => dispatch({ type: 'DELETE_SUBTASK', itemId: item.id, subtaskId: s.id })}
                    className="icon-btn is-danger-hover hover-reveal"
                    style={{ height: '1.5rem', width: '1.5rem' }}
                  >
                    <Trash2 size={14} />
                  </button>
                </div>
              ))}
              <div className="is-flex is-align-items-center px-1 py-1 muted" style={{ gap: '0.5rem' }}>
                <CornerDownRight size={16} style={{ flexShrink: 0 }} />
                <input
                  className="input is-small is-shadowless"
                  style={{ border: 'none', paddingLeft: 0, backgroundColor: 'transparent' }}
                  value={subtask}
                  onChange={(e) => setSubtask(e.target.value)}
                  onKeyDown={(e) => {
                    if (e.nativeEvent.isComposing || e.keyCode === 229) return
                    if (e.key === 'Enter' && subtask.trim()) {
                      dispatch({ type: 'ADD_SUBTASK', itemId: item.id, title: subtask })
                      setSubtask('')
                    }
                  }}
                  placeholder="Add a subtask…"
                />
              </div>
            </div>
          </Section>

          {/* Comments */}
          <Section title="Comments" meta={item.comments.length || undefined}>
            <div className="is-flex is-flex-direction-column" style={{ gap: '0.75rem' }}>
              {item.comments.map((c) => (
                <div key={c.id} className="hover-parent is-flex" style={{ gap: '0.625rem' }}>
                  <div
                    className="is-flex is-align-items-center is-justify-content-center is-family-code has-text-primary has-text-weight-semibold"
                    style={{
                      height: '1.75rem',
                      width: '1.75rem',
                      flexShrink: 0,
                      borderRadius: '9999px',
                      fontSize: '0.6875rem',
                      backgroundColor:
                        'hsl(var(--bulma-primary-h) var(--bulma-primary-s) var(--bulma-primary-l) / 0.15)',
                    }}
                  >
                    {c.author.slice(0, 1)}
                  </div>
                  <div style={{ minWidth: 0, flex: '1 1 0' }}>
                    <div className="is-flex is-align-items-center" style={{ gap: '0.5rem' }}>
                      <span className="is-size-7 has-text-weight-medium">{c.author}</span>
                      <span className="is-size-7 muted">{relativeTime(c.createdAt)}</span>
                      <button
                        type="button"
                        aria-label="Delete comment"
                        onClick={() =>
                          dispatch({ type: 'DELETE_COMMENT', itemId: item.id, commentId: c.id })
                        }
                        className="icon-btn is-danger-hover hover-reveal"
                        style={{ marginLeft: 'auto', height: '1.5rem', width: '1.5rem' }}
                      >
                        <Trash2 size={14} />
                      </button>
                    </div>
                    <p className="is-size-7 muted">{c.body}</p>
                  </div>
                </div>
              ))}
              <div className="field has-addons mb-0">
                <div className="control is-expanded">
                  <input
                    className="input is-small"
                    value={comment}
                    onChange={(e) => setComment(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.nativeEvent.isComposing || e.keyCode === 229) return
                      if (e.key === 'Enter' && comment.trim()) {
                        dispatch({ type: 'ADD_COMMENT', itemId: item.id, body: comment })
                        setComment('')
                      }
                    }}
                    placeholder="Write a comment…"
                  />
                </div>
                <div className="control">
                  <button
                    type="button"
                    aria-label="Send comment"
                    disabled={!comment.trim()}
                    onClick={() => {
                      if (!comment.trim()) return
                      dispatch({ type: 'ADD_COMMENT', itemId: item.id, body: comment })
                      setComment('')
                    }}
                    className="button is-small is-primary"
                  >
                    <Send size={16} />
                  </button>
                </div>
              </div>
            </div>
          </Section>

          {/* Activity */}
          {itemEvents.length > 0 && (
            <Section title="Activity">
              <ol className="is-flex is-flex-direction-column" style={{ gap: '0.5rem' }}>
                {itemEvents.map((e) => (
                  <li key={e.id} className="is-flex is-align-items-center is-size-7 muted" style={{ gap: '0.5rem' }}>
                    <span
                      className="has-background-primary"
                      style={{ height: '0.375rem', width: '0.375rem', borderRadius: '9999px', flexShrink: 0 }}
                    />
                    <span className="has-text-weight-medium has-text-dark">{e.name}</span>
                    {e.detail && <span style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>· {e.detail}</span>}
                    <span className="is-family-code" style={{ marginLeft: 'auto', flexShrink: 0 }}>
                      {relativeTime(e.at)}
                    </span>
                  </li>
                ))}
              </ol>
            </Section>
          )}
        </div>
      </div>
    </div>
  )
}

function TitleField({ item }: { item: Item }) {
  const { dispatch } = useStore()
  const [value, setValue] = useState(item.title)
  const [editing, setEditing] = useState(false)
  const current = editing ? value : item.title

  return (
    <textarea
      rows={1}
      value={current}
      onFocus={() => {
        setValue(item.title)
        setEditing(true)
      }}
      onChange={(e) => setValue(e.target.value)}
      onBlur={() => {
        if (value.trim() && value !== item.title) {
          dispatch({ type: 'MODIFY_TITLE', itemId: item.id, title: value.trim() })
        }
        setEditing(false)
      }}
      onKeyDown={(e) => {
        if (e.nativeEvent.isComposing || e.keyCode === 229) return
        if (e.key === 'Enter') {
          e.preventDefault()
          ;(e.target as HTMLTextAreaElement).blur()
        }
      }}
      className={`textarea is-shadowless ${item.completed ? 'muted line-through' : ''}`}
      style={{
        flex: '1 1 0',
        resize: 'none',
        border: 'none',
        padding: 0,
        minHeight: 'unset',
        backgroundColor: 'transparent',
        fontSize: '1.125rem',
        fontWeight: 600,
        lineHeight: 1.3,
        boxShadow: 'none',
      }}
    />
  )
}

function DescriptionField({ item }: { item: Item }) {
  const { dispatch } = useStore()
  const [value, setValue] = useState(item.description)
  const [editing, setEditing] = useState(false)
  const current = editing ? value : item.description

  return (
    <textarea
      rows={3}
      className="textarea"
      value={current}
      onFocus={() => {
        setValue(item.description)
        setEditing(true)
      }}
      onChange={(e) => setValue(e.target.value)}
      onBlur={() => {
        if (value !== item.description) {
          dispatch({ type: 'MODIFY_DESCRIPTION', itemId: item.id, description: value })
        }
        setEditing(false)
      }}
      placeholder="Add notes, links, or context…"
    />
  )
}

function PropertyRow({
  icon,
  label,
  children,
  divider = true,
}: {
  icon: React.ReactNode
  label: string
  children: React.ReactNode
  divider?: boolean
}) {
  return (
    <div
      className="is-flex px-3 py-3"
      style={{
        gap: '0.75rem',
        alignItems: 'center',
        borderTop: divider ? '1px solid var(--app-border)' : undefined,
      }}
    >
      <div className="is-flex is-align-items-center is-size-7 muted" style={{ width: '7rem', flexShrink: 0, gap: '0.5rem' }}>
        <span style={{ display: 'inline-flex' }}>{icon}</span>
        {label}
      </div>
      <div style={{ minWidth: 0, flex: '1 1 0' }}>{children}</div>
    </div>
  )
}

function Section({
  title,
  meta,
  children,
}: {
  title: string
  meta?: string | number
  children: React.ReactNode
}) {
  return (
    <section className="is-flex is-flex-direction-column" style={{ gap: '0.625rem' }}>
      <div className="is-flex is-align-items-center" style={{ gap: '0.5rem' }}>
        <h2 className="is-family-code is-uppercase muted mb-0" style={{ fontSize: '0.6875rem', letterSpacing: '0.12em' }}>
          {title}
        </h2>
        {meta !== undefined && (
          <span className="is-family-code is-size-7 muted">· {meta}</span>
        )}
      </div>
      {children}
    </section>
  )
}

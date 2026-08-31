'use client'

import { useState } from 'react'
import { Inbox, Plus, PanelLeft } from 'lucide-react'
import { useStore, INBOX_ID } from '@/lib/store'
import { ItemRow } from './item-row'

export function TaskList({ onOpenSidebar }: { onOpenSidebar: () => void }) {
  const { state, dispatch } = useStore()
  const [capture, setCapture] = useState('')

  const isInbox = state.selectedListId === INBOX_ID
  const list = state.lists.find((l) => l.id === state.selectedListId)
  const title = isInbox ? 'Inbox' : (list?.name ?? 'List')

  const items = state.items.filter((it) => it.listId === state.selectedListId)
  const active = items.filter((it) => !it.completed)
  const completed = items.filter((it) => it.completed)

  const submitCapture = () => {
    const value = capture.trim()
    if (!value) return
    dispatch({
      type: 'CAPTURE_INBOX_ITEM',
      title: value,
      listId: isInbox ? undefined : state.selectedListId,
    })
    setCapture('')
  }

  return (
    <div className="is-flex is-flex-direction-column" style={{ height: '100%' }}>
      <header
        className="is-flex is-align-items-center px-5 py-4"
        style={{ gap: '0.75rem', borderBottom: '1px solid var(--app-border)' }}
      >
        <button
          type="button"
          onClick={onOpenSidebar}
          aria-label="Open menu"
          className="icon-btn is-hidden-desktop"
          style={{ height: '2.25rem', width: '2.25rem' }}
        >
          <PanelLeft size={20} />
        </button>
        <div className="is-flex is-align-items-center" style={{ gap: '0.5rem' }}>
          {isInbox && <Inbox size={20} className="muted" />}
          <h1 className="is-size-4 has-text-weight-semibold mb-0">{title}</h1>
          <span className="is-family-code muted">{active.length}</span>
        </div>
      </header>

      <div className="px-5 py-3" style={{ borderBottom: '1px solid var(--app-border)' }}>
        <div className="field mb-0">
          <div className="control has-icons-left">
            <input
              className="input is-medium"
              value={capture}
              onChange={(e) => setCapture(e.target.value)}
              onKeyDown={(e) => {
                if (e.nativeEvent.isComposing || e.keyCode === 229) return
                if (e.key === 'Enter') submitCapture()
              }}
              placeholder={isInbox ? 'Capture something on your mind…' : `Add to ${title}…`}
            />
            <span className="icon is-left">
              <Plus size={18} />
            </span>
          </div>
        </div>
      </div>

      <div className="app-scroll px-5 py-4">
        {items.length === 0 ? (
          <EmptyState isInbox={isInbox} title={title} />
        ) : (
          <div
            className="is-flex is-flex-direction-column"
            style={{ maxWidth: '42rem', marginInline: 'auto', gap: '1rem' }}
          >
            <div className="is-flex is-flex-direction-column" style={{ gap: '0.375rem' }}>
              {active.map((item) => (
                <ItemRow key={item.id} item={item} />
              ))}
            </div>

            {completed.length > 0 && (
              <div className="is-flex is-flex-direction-column" style={{ gap: '0.375rem' }}>
                <p
                  className="is-family-code is-uppercase muted px-1 pt-2"
                  style={{ fontSize: '0.6875rem', letterSpacing: '0.12em' }}
                >
                  Completed · {completed.length}
                </p>
                {completed.map((item) => (
                  <ItemRow key={item.id} item={item} />
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  )
}

function EmptyState({ isInbox, title }: { isInbox: boolean; title: string }) {
  return (
    <div
      className="is-flex is-flex-direction-column is-align-items-center is-justify-content-center has-text-centered"
      style={{ maxWidth: '24rem', marginInline: 'auto', gap: '0.5rem', paddingBlock: '5rem' }}
    >
      <span
        className="is-flex is-align-items-center is-justify-content-center has-text-primary"
        style={{
          height: '3rem',
          width: '3rem',
          borderRadius: '9999px',
          backgroundColor: 'hsl(var(--bulma-primary-h) var(--bulma-primary-s) var(--bulma-primary-l) / 0.12)',
        }}
      >
        <Inbox size={24} />
      </span>
      <p className="has-text-weight-medium">
        {isInbox ? 'Your inbox is empty' : `Nothing in ${title} yet`}
      </p>
      <p className="is-size-7 muted">
        {isInbox
          ? 'Capture anything on your mind above. Clarify and organize it later.'
          : 'Add an item above to start filling this list.'}
      </p>
    </div>
  )
}

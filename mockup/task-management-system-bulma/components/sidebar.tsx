'use client'

import { useEffect, useRef, useState } from 'react'
import {
  Inbox,
  Plus,
  Trash2,
  Hash,
  Check,
  MoreHorizontal,
  Settings,
  User,
  LogOut,
  BadgeCheck,
} from 'lucide-react'
import { useStore, INBOX_ID } from '@/lib/store'

export function Sidebar({ onNavigate }: { onNavigate?: () => void }) {
  const { state, dispatch } = useStore()
  const [newList, setNewList] = useState('')
  const [adding, setAdding] = useState(false)

  const inboxCount = state.items.filter((it) => it.listId === INBOX_ID && !it.completed).length

  const countFor = (listId: string) =>
    state.items.filter((it) => it.listId === listId && !it.completed).length

  const select = (listId: string) => {
    dispatch({ type: 'SELECT_LIST', listId })
    onNavigate?.()
  }

  const submitList = () => {
    if (!newList.trim()) return
    dispatch({ type: 'CREATE_LIST', name: newList })
    setNewList('')
    setAdding(false)
    onNavigate?.()
  }

  return (
    <div
      className="is-flex is-flex-direction-column p-4"
      style={{ height: '100%', gap: '1.5rem' }}
    >
      <div className="is-flex is-align-items-center px-2 pt-1" style={{ gap: '0.5rem' }}>
        <span
          className="is-flex is-align-items-center is-justify-content-center has-background-primary"
          style={{ height: '1.75rem', width: '1.75rem', borderRadius: '0.5rem' }}
        >
          <Check size={16} strokeWidth={3} color="white" />
        </span>
        <span className="is-size-5 has-text-weight-semibold">Listello</span>
      </div>

      <aside className="menu">
        <ul className="menu-list">
          <li>
            <NavButton
              active={state.selectedListId === INBOX_ID}
              onClick={() => select(INBOX_ID)}
              icon={<Inbox size={16} />}
              label="Inbox"
              count={inboxCount}
            />
          </li>
        </ul>
      </aside>

      <div className="is-flex is-flex-direction-column" style={{ flex: '1 1 0', minHeight: 0, gap: '0.25rem' }}>
        <div className="is-flex is-align-items-center is-justify-content-space-between px-2 pb-1">
          <span className="menu-label is-family-code mb-0" style={{ letterSpacing: '0.12em' }}>
            Lists
          </span>
          <button
            type="button"
            onClick={() => setAdding((v) => !v)}
            aria-label="Add list"
            className="icon-btn"
            style={{ height: '1.5rem', width: '1.5rem' }}
          >
            <Plus size={16} />
          </button>
        </div>

        <aside className="menu app-scroll">
          <ul className="menu-list">
            {state.lists.map((list) => (
              <li key={list.id} className="hover-parent" style={{ position: 'relative' }}>
                <NavButton
                  active={state.selectedListId === list.id}
                  onClick={() => select(list.id)}
                  icon={<Hash size={16} />}
                  label={list.name}
                  count={countFor(list.id)}
                />
                <button
                  type="button"
                  aria-label={`Delete ${list.name}`}
                  onClick={(e) => {
                    e.stopPropagation()
                    dispatch({ type: 'DELETE_LIST', listId: list.id })
                  }}
                  className="icon-btn is-danger-hover hover-reveal"
                  style={{
                    position: 'absolute',
                    right: '0.5rem',
                    top: '50%',
                    transform: 'translateY(-50%)',
                    height: '1.5rem',
                    width: '1.5rem',
                  }}
                >
                  <Trash2 size={14} />
                </button>
              </li>
            ))}
          </ul>

          {adding && (
            <input
              autoFocus
              value={newList}
              onChange={(e) => setNewList(e.target.value)}
              onKeyDown={(e) => {
                if (e.nativeEvent.isComposing || e.keyCode === 229) return
                if (e.key === 'Enter') submitList()
                if (e.key === 'Escape') {
                  setAdding(false)
                  setNewList('')
                }
              }}
              onBlur={submitList}
              placeholder="List name"
              className="input is-small mt-2"
            />
          )}
        </aside>
      </div>

      <AccountMenu />
    </div>
  )
}

function AccountMenu() {
  const [open, setOpen] = useState(false)
  const [flash, setFlash] = useState<string | null>(null)
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!open) return
    const onClick = (e: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setOpen(false)
      }
    }
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') setOpen(false)
    }
    document.addEventListener('mousedown', onClick)
    document.addEventListener('keydown', onKey)
    return () => {
      document.removeEventListener('mousedown', onClick)
      document.removeEventListener('keydown', onKey)
    }
  }, [open])

  const runAction = (label: string) => {
    setOpen(false)
    setFlash(label)
    window.setTimeout(() => setFlash(null), 1800)
  }

  return (
    <div ref={containerRef} className={`dropdown is-up is-right ${open ? 'is-active' : ''}`} style={{ display: 'block' }}>
      <div className="dropdown-menu" style={{ minWidth: '100%' }} role="menu">
        <div className="dropdown-content">
          <MenuItem icon={<User size={16} />} label="Profile" onClick={() => runAction('Profile')} />
          <MenuItem
            icon={<Settings size={16} />}
            label="Settings"
            onClick={() => runAction('Settings')}
          />
          <MenuItem
            icon={<BadgeCheck size={16} />}
            label="Delegation policy"
            onClick={() => runAction('Delegation policy')}
          />
          <hr className="dropdown-divider" />
          <MenuItem
            icon={<LogOut size={16} />}
            label="Sign out"
            destructive
            onClick={() => runAction('Signed out')}
          />
        </div>
      </div>

      <div
        className="dropdown-trigger is-flex is-align-items-center p-3"
        style={{
          gap: '0.5rem',
          width: '100%',
          borderRadius: '0.75rem',
          backgroundColor: 'hsl(var(--bulma-primary-h) var(--bulma-primary-s) var(--bulma-primary-l) / 0.1)',
        }}
      >
        <div style={{ minWidth: 0, flex: '1 1 0' }}>
          <p
            className="is-family-code is-uppercase has-text-primary"
            style={{ fontSize: '0.6875rem', letterSpacing: '0.12em' }}
          >
            {flash ?? 'Signed in'}
          </p>
          <p className="is-size-7 has-text-weight-medium mt-1">You · Owner</p>
        </div>
        <button
          type="button"
          onClick={() => setOpen((v) => !v)}
          aria-label="Account options"
          aria-haspopup="menu"
          aria-expanded={open}
          className={`icon-btn ${open ? 'is-active' : ''}`}
          style={{ height: '2rem', width: '2rem' }}
        >
          <MoreHorizontal size={16} />
        </button>
      </div>
    </div>
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
      <span className="muted" style={{ color: destructive ? 'hsl(348deg 86% 43%)' : undefined, display: 'inline-flex' }}>
        {icon}
      </span>
      <span>{label}</span>
    </a>
  )
}

function NavButton({
  active,
  onClick,
  icon,
  label,
  count,
}: {
  active: boolean
  onClick: () => void
  icon: React.ReactNode
  label: string
  count: number
}) {
  return (
    <a
      onClick={onClick}
      className={`is-flex is-align-items-center ${active ? 'is-active' : ''}`}
      style={{ gap: '0.625rem' }}
    >
      <span style={{ display: 'inline-flex' }}>{icon}</span>
      <span style={{ flex: '1 1 0', textAlign: 'left', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
        {label}
      </span>
      {count > 0 && (
        <span className="is-family-code is-size-7" style={{ opacity: active ? 0.7 : 0.55 }}>
          {count}
        </span>
      )}
    </a>
  )
}

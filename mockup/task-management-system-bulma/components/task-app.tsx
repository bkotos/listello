'use client'

import { useState } from 'react'
import { StoreProvider, useStore } from '@/lib/store'
import { Sidebar } from './sidebar'
import { TaskList } from './task-list'
import { DetailPanel } from './detail-panel'

function Shell() {
  const { state, dispatch } = useStore()
  const [mobileNav, setMobileNav] = useState(false)
  const selected = state.items.find((it) => it.id === state.selectedItemId) ?? null

  return (
    <div className="app-shell">
      {/* Sidebar — persistent on desktop */}
      <aside className="app-sidebar is-hidden-touch">
        <Sidebar />
      </aside>

      {/* Mobile sidebar drawer */}
      {mobileNav && (
        <div className="is-hidden-desktop">
          <div className="app-overlay" onClick={() => setMobileNav(false)} aria-hidden />
          <div className="app-drawer is-left">
            <Sidebar onNavigate={() => setMobileNav(false)} />
          </div>
        </div>
      )}

      {/* Main list */}
      <main className="app-main">
        <TaskList onOpenSidebar={() => setMobileNav(true)} />
      </main>

      {/* Detail — side panel on desktop */}
      {selected && (
        <aside className="app-detail is-hidden-touch">
          <DetailPanel
            item={selected}
            onClose={() => dispatch({ type: 'SELECT_ITEM', itemId: null })}
          />
        </aside>
      )}

      {/* Detail — overlay on mobile/tablet */}
      {selected && (
        <div className="is-hidden-desktop">
          <div
            className="app-overlay"
            style={{ zIndex: 50 }}
            onClick={() => dispatch({ type: 'SELECT_ITEM', itemId: null })}
            aria-hidden
          />
          <div className="app-drawer is-right" style={{ zIndex: 51 }}>
            <DetailPanel
              item={selected}
              onClose={() => dispatch({ type: 'SELECT_ITEM', itemId: null })}
            />
          </div>
        </div>
      )}
    </div>
  )
}

export function TaskApp() {
  return (
    <StoreProvider>
      <Shell />
    </StoreProvider>
  )
}

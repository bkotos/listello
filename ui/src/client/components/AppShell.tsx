import { useState } from "react";
import { Outlet } from "react-router-dom";
import { Sidebar } from "./Sidebar.tsx";

export function AppShell() {
  const [mobileNav, setMobileNav] = useState(false);

  return (
    <div className="app-shell">
      <aside className="app-sidebar is-hidden-touch">
        <Sidebar />
      </aside>

      {mobileNav && (
        <div className="is-hidden-desktop">
          <div className="app-overlay" onClick={() => setMobileNav(false)} aria-hidden />
          <div className="app-drawer is-left">
            <Sidebar onNavigate={() => setMobileNav(false)} />
          </div>
        </div>
      )}

      <main className="app-main">
        <Outlet context={{ openSidebar: () => setMobileNav(true) }} />
      </main>
    </div>
  );
}

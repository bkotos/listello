import { NavLink, Route, Routes } from "react-router-dom";
import HomePage from "./pages/HomePage.tsx";
import InboxPage from "./pages/InboxPage.tsx";

function App() {
  return (
    <>
      <header>
        <p className="brand">Listello</p>
        <nav>
          <NavLink to="/" end>
            Lists
          </NavLink>
          <NavLink to="/inbox">Inbox</NavLink>
        </nav>
      </header>
      <main>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/inbox" element={<InboxPage />} />
        </Routes>
      </main>
    </>
  );
}

export default App;

import { Navigate, Route, Routes } from "react-router-dom";
import { AppShell } from "./components/AppShell.tsx";
import InboxPage from "./pages/InboxPage.tsx";
import ListPage from "./pages/ListPage.tsx";

function App() {
  return (
    <Routes>
      <Route element={<AppShell />}>
        <Route index element={<Navigate to="/inbox" replace />} />
        <Route path="inbox" element={<InboxPage />} />
        <Route path="lists/:listId" element={<ListPage />} />
      </Route>
    </Routes>
  );
}

export default App;

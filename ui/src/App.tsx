import { Navigate, Route, Routes } from "react-router-dom";
import { AppShell } from "./components/AppShell";
import { useInstanceQuery } from "./lib/api/instance-queries";
import InboxPage from "./pages/InboxPage";
import ListPage from "./pages/ListPage";
import OnboardingPage from "./pages/OnboardingPage";

function App() {
  const { data: instance } = useInstanceQuery();

  if (instance === null) {
    return <OnboardingPage />;
  }

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

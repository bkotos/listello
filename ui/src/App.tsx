import { Navigate, Outlet, Route, Routes } from "react-router-dom";
import { AppShell } from "./components/AppShell";
import { useInstanceQuery } from "./lib/api/instance-queries";
import InboxPage from "./pages/InboxPage";
import ListPage from "./pages/ListPage";
import OnboardingPage from "./pages/OnboardingPage";

function App() {
  return (
    <Routes>
      <Route path="onboarding" element={<OnboardingPage />} />
      <Route element={<RequireInstance />}>
        <Route element={<AppShell />}>
          <Route index element={<Navigate to="/inbox" replace />} />
          <Route path="inbox" element={<InboxPage />} />
          <Route path="lists/:listId" element={<ListPage />} />
        </Route>
      </Route>
    </Routes>
  );
}

function RequireInstance() {
  const { data: instance } = useInstanceQuery();

  if (instance === null) {
    return <Navigate to="/onboarding" replace />;
  }

  return <Outlet />;
}

export default App;

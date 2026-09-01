import { useOutletContext } from "react-router-dom";

type ShellContext = {
  openSidebar: () => void;
};

export function useShellContext() {
  return useOutletContext<ShellContext>();
}

import { Outlet, ScrollRestoration, useMatch } from "react-router-dom";
import { TopNav } from "../components/layout/TopNav";
import { Footer } from "../components/layout/Footer";

export function AppLayout() {
  const isSpace = useMatch("/space/*") !== null;

  return (
    <div className={`min-h-screen flex flex-col ${isSpace ? "pt-16" : ""}`}>
      <TopNav />
      <main className="flex-1">
        <Outlet />
      </main>
      {!isSpace && <Footer />}
      <ScrollRestoration />
    </div>
  );
}

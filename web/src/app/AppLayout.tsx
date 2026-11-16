import { Outlet } from "react-router-dom";
import { TopNav } from "../components/layout/TopNav";
import { Footer } from "../components/layout/Footer";

export function AppLayout() {
  return (
    <div className="min-h-screen flex flex-col">
      <TopNav />
      <main className="flex-1">
        <Outlet />
      </main>
      <Footer />
    </div>
  );
}

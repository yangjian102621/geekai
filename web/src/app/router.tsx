import { lazy, Suspense } from "react";
import { createBrowserRouter } from "react-router-dom";
import { AppProviders } from "./providers";
import { AppLayout } from "./AppLayout";
import { HomePage } from "../pages/HomePage";
import { InvitePage } from "../pages/InvitePage";
import { PricingPage } from "../pages/PricingPage";
import { NotFoundPage } from "../pages/NotFoundPage";

const SpacePage = lazy(() => import("../pages/SpacePage"));

export const router = createBrowserRouter([
  {
    path: "/",
    element: (
      <AppProviders>
        <AppLayout />
      </AppProviders>
    ),
    children: [
      { index: true, element: <HomePage /> },
      { path: "invite", element: <InvitePage /> },
      { path: "pricing", element: <PricingPage /> },
      { path: "space/*", element: <Suspense fallback={<div className="h-full" aria-busy="true" />}><SpacePage /></Suspense> },
      { path: "*", element: <NotFoundPage /> },
    ],
  },
]);

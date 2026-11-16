import { createBrowserRouter } from "react-router-dom";
import { AppProviders } from "./providers";
import { AppLayout } from "./AppLayout";
import { HomePage } from "../pages/HomePage";
import { InvitePage } from "../pages/InvitePage";
import { PricingPage } from "../pages/PricingPage";
import { NotFoundPage } from "../pages/NotFoundPage";

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
      { path: "*", element: <NotFoundPage /> },
    ],
  },
]);

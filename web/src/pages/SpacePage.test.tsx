// @vitest-environment jsdom
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { cleanup, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter, Route, Routes } from "react-router-dom";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { api } from "../api/client";
import { useAuth } from "../features/auth/useAuth";
import type { AuthUser } from "../features/auth/useAuth";
import SpacePage from "./SpacePage";

vi.mock("../api/client", () => ({ api: { get: vi.fn(), post: vi.fn() } }));
vi.mock("../features/auth/useAuth", () => ({ useAuth: vi.fn() }));
vi.mock("react-i18next", () => ({
  useTranslation: () => ({ t: (key: string) => key, i18n: { language: "zh-CN" } }),
}));

const user: AuthUser = {
  id: 7,
  username: "alice",
  nickname: "Alice",
  email: "alice@example.com",
  power: 880,
  vip: true,
  expired_time: 0,
};

function renderSpace(path = "/space") {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } });
  return render(
    <QueryClientProvider client={client}>
      <MemoryRouter initialEntries={[path]}>
        <Routes>
          <Route path="/space/*" element={<SpacePage />} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  );
}

describe("Personal Space", () => {
  beforeEach(() => {
    vi.mocked(useAuth).mockReturnValue({ data: user, isLoading: false, isError: false } as ReturnType<typeof useAuth>);
    vi.mocked(api.get).mockImplementation(async (path) => {
      switch (path) {
        case "/user/model-token":
          return { data: { code: 0, data: { enabled: false } } };
        case "/powerLog/stats":
          return { data: { data: { total: 120, today: 20, balance: 880 } } };
        case "/invite/stats":
          return { data: { data: { invite_count: 2, reward_total: 400, today_invite: 0, invite_code: "ABC", invite_link: "https://site.example/invite?aff=ABC" } } };
        case "/powerLog/list":
          return {
            data: {
              data: {
                items: [{
                  id: 1,
                  user_id: 7,
                  type: 4,
                  type_str: "邀请",
                  amount: 200,
                  mark: 1,
                  balance: 880,
                  created_at: 1700000000,
                }],
                total: 1,
              },
            },
          };
        default:
          throw new Error(`unexpected request: ${path}`);
      }
    });
    vi.mocked(api.post).mockResolvedValue({
      data: {
        data: {
          items: [{
            id: 1,
            user_id: 7,
            type: 4,
            type_str: "邀请",
            amount: 200,
            mark: 1,
            balance: 880,
            created_at: 1700000000,
          }],
          total: 1,
        },
      },
    });
  });

  afterEach(() => {
    cleanup();
    vi.clearAllMocks();
  });

  it("asks signed-out visitors to sign in instead of showing personal data", () => {
    vi.mocked(useAuth).mockReturnValue({ data: null, isLoading: false, isError: false } as ReturnType<typeof useAuth>);
    renderSpace();
    expect(screen.getByText("space.login_required")).toBeTruthy();
    expect(screen.getByRole("button", { name: "nav.login" })).toBeTruthy();
    expect(api.get).not.toHaveBeenCalled();
  });

  it("loads credit statistics and history for the signed-in user", async () => {
    renderSpace();
    expect(screen.getByRole("heading", { level: 1, name: "space.points" })).toBeTruthy();
    expect(await screen.findByText("120")).toBeTruthy();
    expect(screen.getAllByText("880").length).toBeGreaterThan(0);
    expect(screen.getByText("400")).toBeTruthy();
    expect(api.get).toHaveBeenCalledWith("/powerLog/stats");
    await waitFor(() => expect(api.post).toHaveBeenCalledWith("/powerLog/list", { page: 1, page_size: 10, date: undefined }));
    expect(screen.getByRole("link", { name: /space.invite/ })).toBeTruthy();
  });

  it("renders account details on the profile tab", () => {
    renderSpace("/space/profile");
    expect(screen.getByRole("heading", { level: 1, name: "space.profile" })).toBeTruthy();
    expect(screen.getByRole("heading", { level: 2, name: "alice" })).toBeTruthy();
    expect(screen.getByText("space.username")).toBeTruthy();
    expect(screen.getByText("space.mobile")).toBeTruthy();
    expect(screen.getByText("space.balance")).toBeTruthy();
    expect(screen.getByRole("img", { name: "alice" }).getAttribute("src")).toBe("/assets/logo.png");
    expect(screen.queryByText("Alice")).toBeNull();
    expect(screen.queryByText("alice@example.com")).toBeNull();
    expect(screen.queryByText("ID: 7")).toBeNull();
    for (const field of ["space.nickname", "space.email", "space.vip", "space.expired_at"]) {
      expect(screen.queryByText(field)).toBeNull();
    }
  });
});

// @vitest-environment jsdom
import { cleanup, render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { useAuth } from "../features/auth/useAuth";
import { useInviteSummary } from "../features/invite/useInviteSummary";
import { InvitePage } from "./InvitePage";

vi.mock("react-i18next", () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock("../features/auth/useAuth", () => ({ useAuth: vi.fn() }));
vi.mock("qrcode.react", () => ({ QRCodeCanvas: () => null }));
vi.mock("../features/invite/useInviteSummary", () => ({
  useInviteSummary: vi.fn(),
  useInviteReward: () => ({ data: 200 }),
}));
vi.mock("../hooks/useLatestSetup", () => ({
  useDownloadLinks: () => ({ windows: "/update/setup.exe" }),
  recommendedHref: () => "/update/setup.exe",
}));
vi.mock("../components/DownloadMenu", () => ({ DownloadMenu: ({ trigger }: { trigger: React.ReactNode }) => trigger }));

describe("GeekAI invitation page", () => {
  beforeEach(() => {
    vi.mocked(useAuth).mockReturnValue({ data: { id: 1, username: "tester" }, isLoading: false } as ReturnType<typeof useAuth>);
    vi.mocked(useInviteSummary).mockReturnValue({ data: undefined, isLoading: false, isError: false, refetch: vi.fn() } as unknown as ReturnType<typeof useInviteSummary>);
  });
  afterEach(cleanup);

  function show(path = "/invite") {
    render(<MemoryRouter initialEntries={[path]}><InvitePage /></MemoryRouter>);
  }

  it("shows service failures instead of fabricating invitation statistics", () => {
    vi.mocked(useInviteSummary).mockReturnValue({ isError: true, refetch: vi.fn() } as unknown as ReturnType<typeof useInviteSummary>);
    show();
    expect(screen.getByText("invite.error")).toBeTruthy();
    expect(screen.getByRole("button", { name: "重试" })).toBeTruthy();
  });

  it("shows the configured reward and actual GeekAI invitation totals", () => {
    vi.mocked(useInviteSummary).mockReturnValue({
      data: { invite_code: "TEST", invite_url: "https://site.example/invite?aff=TEST", invited_count: 3, total_invite_points: 600, today_invites: 1 },
      isLoading: false,
    } as ReturnType<typeof useInviteSummary>);
    show();
    expect(screen.getByText(/获得 200 积分/)).toBeTruthy();
    expect(screen.getByText("600")).toBeTruthy();
    expect(screen.getByDisplayValue("https://site.example/invite?aff=TEST")).toBeTruthy();
    expect(screen.queryByText(/500|15,000|上限 30/)).toBeNull();
  });

  it("offers registration to invited guests", () => {
    vi.mocked(useAuth).mockReturnValue({ data: null, isLoading: false } as ReturnType<typeof useAuth>);
    show("/invite?aff=TEST");
    expect(screen.getByRole("button", { name: "注册并下载" })).toBeTruthy();
  });

  it("uses the normal installer link for signed-in invitees", () => {
    show("/invite?aff=TEST");
    expect(screen.getByRole("link", { name: "nav.download" }).getAttribute("href")).toBe("/update/setup.exe");
  });
});

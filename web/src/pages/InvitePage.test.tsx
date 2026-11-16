// @vitest-environment jsdom
import { cleanup, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { useAuth } from "../features/auth/useAuth";
import { useInviteReward, useInviteSummary } from "../features/invite/useInviteSummary";
import { InvitePage } from "./InvitePage";

vi.mock("react-i18next", () => ({ useTranslation: () => ({ t: (key: string) => key }) }));
vi.mock("../features/auth/useAuth", () => ({ useAuth: vi.fn() }));
vi.mock("../features/invite/useInviteSummary", () => ({
  useInviteSummary: vi.fn(),
  useInviteReward: vi.fn(),
}));
vi.mock("qrcode.react", () => ({ QRCodeCanvas: () => null }));

function show(path = "/invite") {
  render(
    <MemoryRouter initialEntries={[path]}>
      <InvitePage />
    </MemoryRouter>,
  );
}

function setMacUserAgent(architecture?: string) {
  Object.defineProperty(window.navigator, "userAgent", {
    value: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
    configurable: true,
  });
  Object.defineProperty(window.navigator, "userAgentData", {
    value: architecture ? { architecture } : undefined,
    configurable: true,
  });
}

describe("GeekAI invitation page", () => {
  beforeEach(() => {
    vi.mocked(useAuth).mockReturnValue({
      data: { id: 1, username: "tester" },
      isLoading: false,
      isError: false,
    } as ReturnType<typeof useAuth>);
    vi.mocked(useInviteSummary).mockReturnValue({
      data: undefined,
      isLoading: false,
      isError: false,
      refetch: vi.fn(),
    } as unknown as ReturnType<typeof useInviteSummary>);
    vi.mocked(useInviteReward).mockReturnValue({ data: 200 } as ReturnType<typeof useInviteReward>);
    vi.stubGlobal("fetch", vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        version: "1.6.15",
        url: "/update/ThinkBuddySetup.exe",
        platforms: {
          "mac-arm64": "/update/ThinkBuddy-mac-arm64.dmg",
          "mac-x64": "/update/ThinkBuddy-mac-x64.dmg",
        },
      }),
    }));
  });

  afterEach(() => {
    cleanup();
    vi.unstubAllGlobals();
    vi.clearAllMocks();
  });

  it("shows service failures instead of fabricating invitation statistics", () => {
    vi.mocked(useInviteSummary).mockReturnValue({
      isError: true,
      refetch: vi.fn(),
    } as unknown as ReturnType<typeof useInviteSummary>);
    show();
    expect(screen.getByText("invite.error")).toBeTruthy();
    expect(screen.getByRole("button", { name: "invite.retry" })).toBeTruthy();
  });

  it("shows the configured reward to guests instead of a fixed amount", () => {
    vi.mocked(useAuth).mockReturnValue({ data: null, isLoading: false } as ReturnType<typeof useAuth>);
    show();
    expect(screen.getByText("invite.reward.per_invite")).toBeTruthy();
    expect(screen.getByText("200")).toBeTruthy();
    expect(screen.getAllByRole("button", { name: "invite.guest.login_btn" }).length).toBeGreaterThan(0);
  });

  it("shows the configured reward and actual GeekAI invitation totals to signed-in users", () => {
    vi.mocked(useInviteSummary).mockReturnValue({
      data: {
        invite_code: "TEST",
        invite_url: "https://site.example/invite?aff=TEST",
        invited_count: 3,
        total_invite_points: 600,
        today_invites: 1,
      },
      isLoading: false,
    } as ReturnType<typeof useInviteSummary>);
    show();
    expect(screen.getByText(/600/)).toBeTruthy();
    expect(screen.getByDisplayValue("https://site.example/invite?aff=TEST")).toBeTruthy();
    expect(screen.getByText("invite.rewards.inviter_desc")).toBeTruthy();
  });

  it("offers sign-in and download to invited guests", () => {
    vi.mocked(useAuth).mockReturnValue({ data: null, isLoading: false } as ReturnType<typeof useAuth>);
    show("/invite?aff=TEST");
    expect(screen.getByRole("button", { name: "invite.invitee.login_btn" })).toBeTruthy();
  });

  it("uses the normal installer link for signed-in invitees", async () => {
    show("/invite?aff=TEST");
    const link = screen.getByRole("link", { name: "invite.invitee.download_btn" });
    await waitFor(() => expect(link.getAttribute("href")).toBe("/update/ThinkBuddySetup.exe"));
  });

  it("offers both macOS installers when the Mac architecture is unknown", async () => {
    setMacUserAgent(undefined);
    show("/invite?aff=TEST");
    const arm = screen.getByRole("link", { name: "download.mac_arm64" });
    const intel = screen.getByRole("link", { name: "download.mac_intel" });
    await waitFor(() => expect(arm.getAttribute("href")).toBe("/update/ThinkBuddy-mac-arm64.dmg"));
    expect(intel.getAttribute("href")).toBe("/update/ThinkBuddy-mac-x64.dmg");
  });
});

// @vitest-environment jsdom
import { renderHook, waitFor } from "@testing-library/react";
import { afterEach, describe, expect, it, vi } from "vitest";

async function loadHookModule() {
  vi.resetModules();
  return import("./useLatestSetup");
}

describe("download links", () => {
  afterEach(() => vi.unstubAllGlobals());

  it("reads Windows and macOS installers from the backend manifest", async () => {
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
    const { useDownloadLinks } = await loadHookModule();
    const { result } = renderHook(() => useDownloadLinks());
    await waitFor(() => expect(result.current.version).toBe("1.6.15"));
    expect(result.current.windows).toBe("/update/ThinkBuddySetup.exe");
    expect(result.current.macArm64).toBe("/update/ThinkBuddy-mac-arm64.dmg");
    expect(result.current.macX64).toBe("/update/ThinkBuddy-mac-x64.dmg");
    expect(result.current.unreachable).toBeUndefined();
  });

  it("falls back to the platform download endpoint when the manifest is unreachable", async () => {
    vi.stubGlobal("fetch", vi.fn().mockRejectedValue(new Error("offline")));
    const { useDownloadLinks } = await loadHookModule();
    const { result } = renderHook(() => useDownloadLinks());
    await waitFor(() => expect(result.current.unreachable).toBe(true));
    expect(result.current.windows).toBe("/api/v1/download");
    expect(result.current.macArm64).toBeNull();
    expect(result.current.macX64).toBeNull();
  });
});

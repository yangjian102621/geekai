import { useEffect, useState } from "react";

export type PlatformKey = "windows" | "mac-arm64" | "mac-x64";
export type MacArch = "arm64" | "x64";

export interface DownloadLinks {
  version: string;
  windows: string | null;
  macArm64: string | null; // null = 暂未提供
  macX64: string | null;
  // 接口不可用（网络失败）时为 true：链接未确认，UI 应显示“下载状态暂不可确认”并提供重试，
  // 不能把网络失败说成已确认有包。
  unreachable?: boolean;
}

export interface DownloadOption {
  key: PlatformKey;
  labelKey: string;
  noteKey: string;
  href: string | null;
  // 为空字符串表示不显示角标（未知 mac 架构不显示“当前设备”）
  badgeKey: string;
  available: boolean;
}

// 后端 /api/v1/latest 返回安装包的同源相对路径（如 /update/ThinkBuddySetup-dt1.6.14.exe），
// mac 包按清单 mac.<arch>.url 输出（缺分支/缺包不输出该 key，即“暂未提供”）。
// 接口不可用或缺少平台地址时显示不可用，不推导安装包地址。
interface LatestResponse {
  version?: string;
  url?: string;
  platforms?: Record<string, string>;
}

const FALLBACK: DownloadLinks = {
  version: "",
  windows: null,
  macArm64: null,
  macX64: null,
  unreachable: true,
};

let cache: DownloadLinks | null = null;
let inflight: Promise<DownloadLinks> | null = null;

function loadLinks(): Promise<DownloadLinks> {
  if (cache) return Promise.resolve(cache);
  if (inflight) return inflight;
  inflight = (async () => {
    try {
      const res = await fetch("/api/v1/latest", { cache: "no-cache" });
      if (!res.ok) return FALLBACK;
      const data = (await res.json()) as LatestResponse;
      // 接口成功但缺 key → 视为暂未提供（null），不伪造链接。
      const links: DownloadLinks = {
        version: data.version ?? "",
        windows: data.url || data.platforms?.windows || null,
        macArm64: data.platforms?.["mac-arm64"] || null,
        macX64: data.platforms?.["mac-x64"] || null,
      };
      cache = links;
      return links;
    } catch {
      // 网络失败时保留重试入口。
      return FALLBACK;
    } finally {
      inflight = null;
    }
  })();
  return inflight;
}

type LinksListener = (links: DownloadLinks) => void;
const listeners = new Set<LinksListener>();

function broadcast(links: DownloadLinks) {
  listeners.forEach((fn) => fn(links));
}

// 网络失败后手动重试：清掉缓存与进行中的请求，重新拉取并通知所有订阅方。
export function refreshDownloadLinks(): void {
  cache = null;
  inflight = null;
  void loadLinks().then((resolved) => {
    cache = resolved;
    broadcast(resolved);
  });
}

export function useDownloadLinks(): DownloadLinks {
  const [links, setLinks] = useState<DownloadLinks>(cache ?? FALLBACK);
  useEffect(() => {
    let active = true;
    const listener = (resolved: DownloadLinks) => {
      if (active) setLinks(resolved);
    };
    listeners.add(listener);
    loadLinks().then((resolved) => {
      if (active) {
        cache = resolved;
        setLinks(resolved);
      }
    });
    return () => {
      active = false;
      listeners.delete(listener);
    };
  }, []);
  return links;
}

// 决定默认下载哪个包：Mac 给 mac 安装包，其他平台给 Windows 包。
export function currentOs(): "mac" | "windows" | "other" {
  const ua = typeof navigator === "undefined" ? "" : navigator.userAgent;
  if (/Macintosh|Mac OS X/i.test(ua)) return "mac";
  if (/Windows/i.test(ua)) return "windows";
  return "other";
}

// 通过 userAgentData.architecture 识别 mac 架构；拿不到就返回 unknown，绝不默认 Intel 或 ARM。
export function currentMacArch(): MacArch | "unknown" {
  if (typeof navigator === "undefined") return "unknown";
  const arch = (navigator as any).userAgentData?.architecture;
  if (arch === "arm") return "arm64";
  if (arch === "x86") return "x64";
  return "unknown";
}

// mac 已知架构时返回对应 platform 参数；Windows 或未知架构返回 null。
export function recommendedPlatformParam(): "mac-arm64" | "mac-x64" | null {
  if (currentOs() !== "mac") return null;
  const arch = currentMacArch();
  return arch === "unknown" ? null : arch === "arm64" ? "mac-arm64" : "mac-x64";
}

export function downloadOptions(links: DownloadLinks): DownloadOption[] {
  const os = currentOs();
  const macArch = currentMacArch();
  // 未知 mac 架构时不显示“当前设备”。
  const macBadgeKnown = os === "mac" && macArch !== "unknown";
  return [
    {
      key: "windows",
      labelKey: "download.win",
      noteKey: "download.win_note",
      href: links.windows,
      badgeKey: os === "windows" ? "download.current" : "",
      available: !!links.windows,
    },
    {
      key: "mac-arm64",
      labelKey: "download.mac_arm64",
      noteKey: "download.mac_arm64_note",
      href: links.macArm64,
      badgeKey: macBadgeKnown && macArch === "arm64" ? "download.current" : "",
      available: !!links.macArm64,
    },
    {
      key: "mac-x64",
      labelKey: "download.mac_intel",
      noteKey: "download.mac_intel_note",
      href: links.macX64,
      badgeKey: macBadgeKnown && macArch === "x64" ? "download.current" : "",
      available: !!links.macX64,
    },
  ];
}

export function recommendedHref(links: DownloadLinks): string | null {
  const os = currentOs();
  if (os === "mac") {
    const arch = currentMacArch();
    if (arch === "unknown") return null; // 未知架构：不默认，让用户选
    const target = arch === "arm64" ? links.macArm64 : links.macX64;
    return target ?? null; // 架构已知但包暂未提供 → null，不伪造链接
  }
  return links.windows || null;
}

// 兼容旧调用点（顶栏、页脚、定价页）：直接返回当前设备推荐的安装包地址；未知 none 时返回 null。
export function useLatestSetup(): string | null {
  return recommendedHref(useDownloadLinks());
}

import type { ReactNode } from "react";
import { useTranslation } from "react-i18next";
import { currentOs, downloadOptions, type DownloadLinks, type DownloadOption } from "../hooks/useLatestSetup";
import { PlatformIcon } from "./PlatformIcon";

function ChevronIcon() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-slate-400">
      <path d="m9 18 6-6-6-6" />
    </svg>
  );
}

// 下载下拉面板，右上角导航和首页 hero 共用：macOS 一行（悬停展开二级：Apple 芯片 / Intel）+ Windows 一行。
// - anchor="bottom"：面板在按钮下方（导航右对齐），macOS 二级菜单向左展开。
// - anchor="right"：面板展开在按钮右侧（hero 用），macOS 二级菜单向右展开。
// 包暂未提供（available=false）时子项禁用并显示“暂未提供”，不伪造链接。
export function DownloadMenu({
  links,
  trigger,
  anchor = "bottom",
  align = "right",
  className = "",
}: {
  links: DownloadLinks;
  trigger: ReactNode;
  anchor?: "bottom" | "right";
  align?: "right" | "center";
  className?: string;
}) {
  const { t } = useTranslation();
  const [windows, macArm64, macX64] = downloadOptions(links);
  const panelPosition =
    anchor === "right"
      ? "top-0 left-full pl-2"
      : `top-full pt-2 ${align === "center" ? "left-1/2 -translate-x-1/2" : "right-0"}`;
  // 导航面板右对齐时二级菜单向左展开；hero 面板在按钮右侧时二级菜单向右展开。
  const macSubPosition = anchor === "right" ? "left-full top-0 pl-2" : "right-full top-0 pr-2";
  const isMac = currentOs() === "mac";

  const macItem = (option: DownloadOption) =>
    option.href ? (
      <a
        key={option.key}
        href={option.href}
        className="flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-slate-50 whitespace-nowrap"
      >
        <span className="flex-1 text-sm text-slate-700">{t(option.labelKey)}</span>
        {option.badgeKey && (
          <span className="text-[10px] font-semibold text-brand-600">{t(option.badgeKey)}</span>
        )}
      </a>
    ) : (
      <div
        key={option.key}
        className="flex items-center gap-3 px-3 py-2 rounded-xl opacity-55 cursor-not-allowed whitespace-nowrap"
        title={t("download.unavailable")}
      >
        <span className="flex-1 text-sm text-slate-500">{t(option.labelKey)}</span>
        <span className="text-[10px] font-semibold text-slate-400">{t("download.unavailable")}</span>
      </div>
    );

  return (
    <div className={`relative group ${className}`}>
      {trigger}
      <div className={`absolute hidden group-hover:block z-30 ${panelPosition}`}>
        <div className="w-64 rounded-2xl border border-slate-200 bg-white p-1.5 shadow-xl text-left">
          <div className="relative group/mac">
            <div className="flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-slate-50 cursor-default">
              <span className="text-slate-500"><PlatformIcon platform="mac-arm64" /></span>
              <span className="flex-1 text-sm text-slate-700">{t("download.mac")}</span>
              {isMac && (
                <span className="text-[10px] font-semibold text-brand-600">{t("download.current")}</span>
              )}
              <ChevronIcon />
            </div>
            <div className={`absolute hidden group-hover/mac:block z-40 ${macSubPosition}`}>
              <div className="w-48 rounded-2xl border border-slate-200 bg-white p-1.5 shadow-xl text-left">
                {macItem(macArm64)}
                {macItem(macX64)}
              </div>
            </div>
          </div>
          {windows.href ? (
            <a
              href={windows.href}
              className="flex items-center gap-3 px-3 py-2 rounded-xl hover:bg-slate-50"
            >
              <span className="text-slate-500"><PlatformIcon platform="windows" /></span>
              <span className="flex-1 text-sm text-slate-700">
                {t(windows.labelKey)}
                <span className="text-slate-400">（{t(windows.noteKey)}）</span>
              </span>
              {windows.badgeKey && (
                <span className="text-[10px] font-semibold text-brand-600">{t(windows.badgeKey)}</span>
              )}
            </a>
          ) : (
            <div
              className="flex items-center gap-3 px-3 py-2 rounded-xl opacity-55 cursor-not-allowed"
              title={t("download.unavailable")}
            >
              <span className="text-slate-500"><PlatformIcon platform="windows" /></span>
              <span className="flex-1 text-sm text-slate-500">{t(windows.labelKey)}</span>
              <span className="text-[10px] font-semibold text-slate-400">{t("download.unavailable")}</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

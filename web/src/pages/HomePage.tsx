import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useReveal } from "../hooks/useReveal";
import { currentMacArch, currentOs, downloadOptions, recommendedHref, refreshDownloadLinks, useDownloadLinks } from "../hooks/useLatestSetup";
import { DownloadMenu } from "../components/DownloadMenu";
import { PlatformIcon } from "../components/PlatformIcon";
import { ChevronDown } from "lucide-react";

export function HomePage() {
  const { t } = useTranslation();
  const ref = useReveal<HTMLDivElement>();
  const [lightbox, setLightbox] = useState<string | null>(null);
  const links = useDownloadLinks();
  const options = downloadOptions(links);
  const downloadUrl = recommendedHref(links);

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => { if (e.key === "Escape") setLightbox(null); };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, []);

  return (
    <div ref={ref}>
      {/* Hero */}
      <section className="hero-bg relative overflow-hidden pt-28 pb-24 lg:pt-36 lg:pb-32">
        <div className="hero-aurora" aria-hidden="true" />
        <div className="hero-grid-overlay absolute inset-0 pointer-events-none" />
        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-brand-50 border border-brand-500/20 text-brand-700 text-xs font-semibold mb-6" data-reveal>
            <span className="w-1.5 h-1.5 rounded-full bg-brand-500" />
            <span>{t("hero.badge")}</span>
          </div>
          <h1 className="text-[clamp(1.5rem,7vw,2.25rem)] sm:text-5xl lg:text-6xl font-bold tracking-tight text-slate-900 leading-tight mb-6" data-reveal>
            {t("hero.title")}
          </h1>
          <p className="max-w-2xl mx-auto text-base sm:text-lg text-slate-600 leading-relaxed mb-10" data-reveal>
            {t("hero.subtitle")}
          </p>
          {/* relative z-30：hero 每行都带 data-reveal（transform 会生成层叠上下文），
              否则下方的支持文案行会盖住下载下拉面板。 */}
          <div className="relative z-30 flex flex-col sm:flex-row items-center justify-center gap-3 mb-8" data-reveal>
            <DownloadMenu
              links={links}
              anchor="right"
              className="w-full sm:w-auto"
              trigger={
                <a
                  href={downloadUrl ?? "#download"}
                  className="btn-primary inline-flex items-center gap-2 px-6 py-3 rounded-xl text-sm font-semibold w-full sm:w-auto justify-center"
                >
                  <DownloadIcon />
                  <span>{t("hero.cta_download")}</span>
                  <ChevronDown size={14} />
                </a>
              }
            />
          </div>
          <p className="text-xs text-slate-400 mt-6" data-reveal>{t("hero.platforms")}</p>
        </div>
      </section>

      {/* 产品演示 */}
      <section className="py-20 lg:py-28 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14" data-reveal>
            <h2 className="text-3xl sm:text-4xl font-bold tracking-tight text-slate-900 mb-3">{t("demo.title")}</h2>
            <p className="max-w-2xl mx-auto text-slate-600">{t("demo.subtitle")}</p>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {[
              { img: "/assets/demo-reading.png", title: t("demo.card1.title"), desc: t("demo.card1.desc") },
              { img: "/assets/demo-knowledge-base.png", title: t("demo.card2.title"), desc: t("demo.card2.desc") },
              { img: "/assets/demo-smart-writing.png", title: t("demo.card3.title"), desc: t("demo.card3.desc") },
            ].map((c) => (
              <div key={c.title} className="bg-white border border-slate-200 rounded-2xl overflow-hidden shadow-sm" data-reveal>
                <div className="demo-card-visual h-44 relative" onClick={() => setLightbox(c.img)}>
                  <img src={c.img} alt={c.title} className="demo-card-img" />
                </div>
                <div className="p-6">
                  <h3 className="text-lg font-semibold text-slate-900 mb-2">{c.title}</h3>
                  <p className="text-sm text-slate-600 leading-relaxed">{c.desc}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* 核心能力矩阵 */}
      <section id="features" className="py-20 lg:py-28 bg-slate-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14" data-reveal>
            <h2 className="text-3xl sm:text-4xl font-bold tracking-tight text-slate-900 mb-3">{t("features.title")}</h2>
            <p className="text-slate-600">{t("features.subtitle")}</p>
            <div className="stats-bar" data-reveal>
              <div className="stat-item">
                <div className="stat-number">6<span className="stat-plus">+</span></div>
                <div className="stat-label">{t("stats.modules")}</div>
              </div>
              <div className="stat-divider" />
              <div className="stat-item">
                <div className="stat-number">101<span className="stat-plus">+</span></div>
                <div className="stat-label">{t("stats.apps")}</div>
              </div>
              <div className="stat-divider" />
              <div className="stat-item">
                <div className="stat-number">30<span className="stat-plus">+</span></div>
                <div className="stat-label">{t("stats.models")}</div>
              </div>
              <div className="stat-divider" />
              <div className="stat-item">
                <div className="stat-number">3</div>
                <div className="stat-label">{t("stats.platforms")}</div>
              </div>
            </div>
          </div>
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {features.map((f) => (
              <div key={f.titleKey} className="feature-card bg-white border border-slate-200 rounded-2xl p-7" data-reveal>
                <div className="feature-icon-wrap w-12 h-12 rounded-xl flex items-center justify-center mb-5" dangerouslySetInnerHTML={{ __html: f.icon }} />
                <h3 className="text-lg font-semibold text-slate-900 mb-2">{t(f.titleKey)}</h3>
                <p className="text-sm text-slate-600 leading-relaxed">{t(f.descKey)}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* 大字错位排版区 */}
      <section className="bigtext-section py-20 lg:py-28 bg-white overflow-hidden">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h2 className="bigtext-title" data-reveal>{t("bigtext.title")}</h2>
          <p className="bigtext-subtitle" data-reveal>{t("bigtext.subtitle")}</p>
          <div className="bigtext-words" data-reveal>
            {["阅读", "笔记", "研究", "写作", "问答", "视频", "课程"].map((w, i) => (
              <span key={w} className="big-word" style={{ "--d": `${i * 120}ms` } as React.CSSProperties}>{w}</span>
            ))}
          </div>
        </div>
      </section>

      {/* 版本时间线 */}
      <section id="versions" className="py-12 lg:py-16 bg-white">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-8" data-reveal>
            <h2 className="text-2xl sm:text-3xl font-bold tracking-tight text-slate-900 mb-2">{t("versions.title")}</h2>
            <p className="text-sm text-slate-600">{t("versions.subtitle")}</p>
          </div>
          <div className="relative pl-7" data-reveal>
            <div className="timeline-line absolute left-2 top-2 bottom-2 w-0.5" />
            {[
              { v: "v1.6.14", d: "2026.10.8", h: t("versions.v1614.highlight") },
              { v: "v1.6.13", d: "2026.10.3", h: t("versions.v1613.highlight") },
              { v: "v1.6.12", d: "2026.9.28", h: t("versions.v1612.highlight") },
              { v: "v1.6.11", d: "2026.9.24", h: t("versions.v1611.highlight") },
              { v: "v1.6.10", d: "2026.9.22", h: t("versions.v1610.highlight") },
            ].map((item) => (
              <div key={item.v} className="relative mb-5">
                <span className="timeline-dot absolute -left-7 top-1 w-3 h-3 rounded-full bg-brand-500" />
                <div className="flex flex-wrap items-baseline gap-2 mb-0.5">
                  <span className="text-base font-bold text-slate-900">{item.v}</span>
                  <span className="text-xs text-slate-500">{item.d}</span>
                </div>
                <p className="text-xs text-slate-600">{item.h}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* 集成生态 */}
      <section id="ecosystem" className="py-20 lg:py-28 bg-slate-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-14" data-reveal>
            <h2 className="text-3xl sm:text-4xl font-bold tracking-tight text-slate-900 mb-3">{t("ecosystem.title")}</h2>
            <p className="max-w-2xl mx-auto text-slate-600">{t("ecosystem.subtitle")}</p>
          </div>
          <div className="space-y-10">
            {[
              { cat: t("ecosystem.cat.im"), items: ["飞书 / Feishu", "微信 / WeChat", "Telegram"] },
              { cat: t("ecosystem.cat.tools"), items: ["MCP Server", "101+ CLI Apps", "第三方 Capability 插件"] },
              { cat: t("ecosystem.cat.kb"), items: ["LightRAG", "GraphRAG", "PageIndex", "MarginNote 4", "GitHub 仓库同步"] },
              { cat: t("ecosystem.cat.models"), items: ["OpenAI Codex", "Gemini", "通义", "豆包", "Zhipu", "…数十家"] },
            ].map((g) => (
              <div key={g.cat} data-reveal>
                <h3 className="text-sm font-semibold text-slate-500 uppercase tracking-wider mb-3">{g.cat}</h3>
                <div className="flex flex-wrap gap-2">
                  {g.items.map((it) => (
                    <span key={it} className="ecosystem-chip bg-white border border-slate-200 rounded-full px-4 py-2 text-sm font-medium text-slate-700">{it}</span>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* 下载区 */}
      <section id="download" className="pb-12 lg:pb-16 bg-slate-50">
        <div className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-6" data-reveal>
            <h2 className="text-2xl sm:text-3xl font-bold tracking-tight text-slate-900 mb-2">{t("download.title")}</h2>
            {links.version && <p className="text-xs text-slate-500">v{links.version}</p>}
          </div>
          {links.unreachable && (
            <div className="flex items-center justify-center gap-2 mb-3 text-xs text-slate-500" data-reveal>
              <span>⚠ {t("download.unreachable")}</span>
              <button
                type="button"
                onClick={refreshDownloadLinks}
                className="text-brand-600 font-semibold hover:underline"
              >
                {t("download.retry")}
              </button>
            </div>
          )}
          {!links.unreachable && currentOs() === "mac" && currentMacArch() === "unknown" && (
            <p className="text-xs text-slate-500 text-center mb-3" data-reveal>{t("download.mac_unknown_hint")}</p>
          )}
          <div className="bg-white border border-slate-200 rounded-2xl overflow-hidden divide-y divide-slate-100" data-reveal>
            {options.map((option) =>
              option.available ? (
                <a
                  key={option.key}
                  href={option.href ?? undefined}
                  className="flex items-center gap-4 px-5 py-4 hover:bg-slate-50 transition"
                >
                  <span className="text-slate-600"><PlatformIcon platform={option.key} size={22} /></span>
                  <span className="flex-1 text-left">
                    <span className="block text-sm font-semibold text-slate-900">{t(option.labelKey)}</span>
                    <span className="block text-xs text-slate-500">{t(option.noteKey)}</span>
                  </span>
                  {option.badgeKey && (
                    <span
                      className={`rounded-full border px-2 py-0.5 text-[10px] font-semibold ${
                        option.badgeKey === "download.current"
                          ? "border-slate-200 text-slate-500"
                          : "border-brand-500/30 bg-brand-50 text-brand-600"
                      }`}
                    >
                      {t(option.badgeKey)}
                    </span>
                  )}
                  <span className="text-slate-400"><DownloadIcon /></span>
                </a>
              ) : (
                <div
                  key={option.key}
                  className="flex items-center gap-4 px-5 py-4 opacity-55 cursor-not-allowed"
                  title={t("download.unavailable")}
                >
                  <span className="text-slate-500"><PlatformIcon platform={option.key} size={22} /></span>
                  <span className="flex-1 text-left">
                    <span className="block text-sm font-semibold text-slate-500">{t(option.labelKey)}</span>
                    <span className="block text-xs text-slate-400">{t(option.noteKey)}</span>
                  </span>
                  <span className="rounded-full border border-slate-200 px-2 py-0.5 text-[10px] font-semibold text-slate-400">
                    {t("download.unavailable")}
                  </span>
                  <span className="text-slate-300"><DownloadIcon /></span>
                </div>
              )
            )}
          </div>
          {/* 定价入口暂时隐藏 */}
          <p className="hidden text-xs text-slate-400 mt-1 text-center">
            {t("credits.more_pricing")}
            <Link to="/pricing" className="text-brand-600 font-medium hover:underline">定价页</Link>
          </p>
        </div>
      </section>

      {/* Lightbox */}
      {lightbox && (
        <div
          className={`lightbox-overlay active`}
          onClick={() => setLightbox(null)}
        >
          <img src={lightbox} alt="放大预览" onClick={(e) => e.stopPropagation()} />
        </div>
      )}
    </div>
  );
}

const features = [
  { titleKey: "features.reading.title", descKey: "features.reading.desc", icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-brand-700"><path d="M2 3h6a4 4 0 0 1 4 4v14a3 3 0 0 0-3-3H2z"/><path d="M22 3h-6a4 4 0 0 0-4 4v14a3 3 0 0 1 3-3h7z"/></svg>' },
  { titleKey: "features.knowledge.title", descKey: "features.knowledge.desc", icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-brand-700"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/></svg>' },
  { titleKey: "features.learning.title", descKey: "features.learning.desc", icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-brand-700"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"/><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"/></svg>' },
  { titleKey: "features.video.title", descKey: "features.video.desc", icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-brand-700"><polygon points="23 7 16 12 23 17 23 7"/><rect x="1" y="5" width="15" height="14" rx="2" ry="2"/></svg>' },
  { titleKey: "features.cowriter.title", descKey: "features.cowriter.desc", icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-brand-700"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>' },
  { titleKey: "features.agent.title", descKey: "features.agent.desc", icon: '<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="text-brand-700"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>' },
];

function DownloadIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
      <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
      <polyline points="7 10 12 15 17 10" />
      <line x1="12" y1="15" x2="12" y2="3" />
    </svg>
  );
}

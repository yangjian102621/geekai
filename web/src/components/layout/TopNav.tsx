import { useEffect, useState } from "react";
import { Link, NavLink } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Menu, X, LogOut, ChevronDown } from "lucide-react";
import { useAuth } from "../../features/auth/useAuth";
import { downloadOptions, recommendedHref, useDownloadLinks } from "../../hooks/useLatestSetup";
import { DownloadMenu } from "../DownloadMenu";
import { PlatformIcon } from "../PlatformIcon";
import { docsURL, openLoginPage, removeUserToken } from "../../config";
import { api } from "../../api/client";

export function TopNav() {
  const { t, i18n } = useTranslation();
  const [open, setOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);
  const { data: user } = useAuth();
  const links = useDownloadLinks();
  const options = downloadOptions(links);
  const downloadUrl = recommendedHref(links);

  useEffect(() => {
    const onScroll = () => setScrolled(window.scrollY > 8);
    window.addEventListener("scroll", onScroll, { passive: true });
    onScroll();
    return () => window.removeEventListener("scroll", onScroll);
  }, []);

  const linkClass = ({ isActive }: { isActive: boolean }) =>
    `nav-link text-sm font-medium transition-colors ${
      isActive ? "text-brand-600" : "text-slate-700 hover:text-slate-900"
    }`;

  const handleLogin = openLoginPage;

  const handleLogout = async () => {
    try {
      await api.get("/user/logout");
    } catch {
      return;
    }
    removeUserToken();
    window.location.reload();
  };

  const toggleLang = () => {
    const next = i18n.language?.startsWith("zh") ? "en" : "zh-CN";
    i18n.changeLanguage(next);
  };

  const langLabel = i18n.language?.startsWith("zh") ? "EN" : "中";

  return (
    <header
      data-topnav
      className={`topnav fixed top-0 inset-x-0 z-50 ${
        scrolled ? "topnav-scrolled" : "bg-white/60"
      }`}
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16 relative">
          <Link to="/" className="flex items-center gap-2 z-10">
            <img src="/assets/logo.png" alt="ThinkBuddy" className="h-8 w-8" />
            <span className="text-lg font-bold tracking-tight text-slate-900">ThinkBuddy</span>
          </Link>

          <nav className="hidden md:flex items-center gap-8 absolute left-1/2 -translate-x-1/2">
            <NavLink to="/" end className={linkClass}>{t("nav.home")}</NavLink>
            {/* 定价暂时隐藏 */}
            <NavLink to="/invite" className={linkClass}>{t("nav.invite")}</NavLink>
            {docsURL && <a href={docsURL} target="_blank" rel="noopener noreferrer" className="nav-link text-sm font-medium transition-colors text-slate-700 hover:text-slate-900">
              {t("nav.docs")}
            </a>}
          </nav>

          <div className="hidden md:flex items-center gap-3 z-10">
            {/* 语言切换 */}
            <button
              onClick={toggleLang}
              className="text-xs font-semibold text-slate-500 hover:text-brand-600 px-2 py-1 rounded-md transition"
              aria-label="toggle language"
            >
              {langLabel}
            </button>

            {/* 登录态 */}
            {user ? (
              <div className="flex items-center gap-1">
                <span className="auth-user-name" title={user.username}>{user.username}</span>
                <button onClick={handleLogout} className="auth-logout-btn" title={t("auth.logout")}>
                  <LogOut size={13} />
                </button>
              </div>
            ) : (
              <button onClick={handleLogin} className="auth-login-btn">{t("nav.login")}</button>
            )}

            <DownloadMenu
              links={links}
              trigger={
                <a
                  href={downloadUrl ?? "/#download"}
                  className="btn-primary inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-semibold"
                >
                  {t("nav.download")}
                  <ChevronDown size={14} />
                </a>
              }
            />
          </div>

          <button
            className="md:hidden p-2 -mr-2"
            onClick={() => setOpen((v) => !v)}
            aria-label="menu"
            aria-expanded={open}
          >
            {open ? <X size={24} /> : <Menu size={24} />}
          </button>
        </div>

        {open && (
          <div className="md:hidden border-t border-slate-200 py-4">
            <nav className="flex flex-col gap-1">
              <NavLink to="/" end onClick={() => setOpen(false)} className="block py-2 px-3 text-sm font-medium text-slate-700 rounded hover:bg-slate-50">
                {t("nav.home")}
              </NavLink>
              <NavLink to="/invite" onClick={() => setOpen(false)} className="block py-2 px-3 text-sm font-medium text-slate-700 rounded hover:bg-slate-50">
                {t("nav.invite")}
              </NavLink>
              {docsURL && <a href={docsURL} target="_blank" rel="noopener noreferrer" onClick={() => setOpen(false)} className="block py-2 px-3 text-sm font-medium text-slate-700 rounded hover:bg-slate-50">
                {t("nav.docs")}
              </a>}
              {user ? (
                <button onClick={handleLogout} className="text-left py-2 px-3 text-sm font-medium text-slate-700 rounded hover:bg-slate-50">
                  {t("auth.logout")} ({user.username})
                </button>
              ) : (
                <button onClick={() => { setOpen(false); handleLogin(); }} className="text-left py-2 px-3 text-sm font-medium text-slate-700 rounded hover:bg-slate-50">
                  {t("nav.login")}
                </button>
              )}
              {options.length > 1 ? (
                options.map((option) =>
                  option.href ? (
                    <a
                      key={option.key}
                      href={option.href}
                      onClick={() => setOpen(false)}
                      className="flex items-center gap-3 py-2 px-3 text-sm font-medium text-slate-700 rounded hover:bg-slate-50"
                    >
                      <PlatformIcon platform={option.key} size={16} />
                      <span className="flex-1">{t(option.labelKey)}</span>
                      {option.badgeKey && (
                        <span className="text-[10px] font-semibold text-brand-600">{t(option.badgeKey)}</span>
                      )}
                    </a>
                  ) : (
                    <div
                      key={option.key}
                      className="flex items-center gap-3 py-2 px-3 text-sm font-medium text-slate-400 rounded cursor-not-allowed"
                      title={t("download.unavailable")}
                    >
                      <PlatformIcon platform={option.key} size={16} />
                      <span className="flex-1">{t(option.labelKey)}</span>
                      <span className="text-[10px] font-semibold text-slate-400">{t("download.unavailable")}</span>
                    </div>
                  )
                )
              ) : (
                <a
                  href={downloadUrl ?? "/#download"}
                  className="btn-primary block text-center px-4 py-2 rounded-lg text-sm font-semibold mt-2"
                >
                  {t("nav.download")}
                </a>
              )}
            </nav>
          </div>
        )}
      </div>
    </header>
  );
}

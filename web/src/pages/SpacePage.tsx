import { Coins, UserRound } from "lucide-react";
import { useTranslation } from "react-i18next";
import { Link, Navigate, Route, Routes, useLocation } from "react-router-dom";
import { openLoginPage } from "../config";
import { useAuth } from "../features/auth/useAuth";
import { PointsPanel } from "../features/space/PointsPanel";
import { ProfilePanel } from "../features/space/ProfilePanel";

// 个人空间：积分明细和账号信息。管理功能由独立管理后台（/admin）提供，
// 这里只承载官网登录用户自己的数据。
export default function SpacePage() {
  const { t } = useTranslation();
  const { data: user, isLoading, isError } = useAuth();
  const location = useLocation();
  const selected = location.pathname.replace(/\/$/, "");
  const tabs = [
    { key: "/space", label: t("space.points"), icon: <Coins size={18} strokeWidth={1.75} aria-hidden="true" /> },
    { key: "/space/profile", label: t("space.profile"), icon: <UserRound size={18} strokeWidth={1.75} aria-hidden="true" /> },
  ];
  const current = tabs.find(tab => tab.key === selected) ?? tabs[0];
  const description = current.key === "/space/profile" ? t("space.profile_hint") : t("space.points_description");

  return (
    <section className="min-h-[calc(100dvh-4rem)] section-soft pt-8 pb-16 sm:pt-10">
      <div className="mx-auto w-full max-w-5xl px-4 sm:px-6 lg:px-8">
        <header className="mb-8">
          <h1 className="font-serif text-3xl font-normal tracking-[-0.03em] text-ink">{current.label}</h1>
          <p className="mt-2 text-muted">{description}</p>
        </header>

        {isLoading ? (
          <div className="flex flex-col items-center justify-center gap-3 py-24 text-muted" role="status">
            <span className="h-8 w-8 animate-spin rounded-full border-2 border-hairline border-t-brand-500" aria-hidden="true" />
            <span>{t("space.loading")}</span>
          </div>
        ) : !user ? (
          <div className="rounded-xl border border-hairline bg-canvas px-6 py-16 text-center">
            <UserRound className="mx-auto text-muted-soft/40" size={40} aria-hidden="true" />
            <h2 className="mt-4 text-lg font-semibold text-ink">{t("space.login_required")}</h2>
            <p className="mt-2 text-sm text-muted">{t("space.login_hint")}</p>
            <button type="button" onClick={openLoginPage} className="btn-primary mt-6 rounded-lg px-6 py-2.5 text-sm font-medium">
              {t("nav.login")}
            </button>
            {isError && <p role="alert" className="mt-4 text-sm text-rose-600">{t("space.load_error")}</p>}
          </div>
        ) : (
          <>
            <nav aria-label={t("space.navigate")} className="mb-6 flex gap-1 overflow-x-auto rounded-xl border border-hairline bg-canvas p-1">
              {tabs.map(tab => {
                const active = tab.key === current.key;
                return (
                  <Link
                    key={tab.key}
                    to={tab.key}
                    aria-current={active ? "page" : undefined}
                    className={`inline-flex flex-1 items-center justify-center gap-2 whitespace-nowrap rounded-lg px-4 py-2.5 text-sm font-medium transition-colors ${
                      active ? "bg-surface-card text-ink" : "text-muted hover:bg-surface-soft hover:text-ink"
                    }`}
                  >
                    {tab.icon}
                    {tab.label}
                  </Link>
                );
              })}
            </nav>

            <Routes>
              <Route index element={<PointsPanel userId={user.id} />} />
              <Route path="profile" element={<ProfilePanel user={user} />} />
              <Route path="*" element={<Navigate to="/space" replace />} />
            </Routes>
          </>
        )}
      </div>
    </section>
  );
}

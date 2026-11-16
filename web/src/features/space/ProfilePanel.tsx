import { BadgeCheck, CalendarClock, Coins, Mail, Phone, UserRound } from "lucide-react";
import { useTranslation } from "react-i18next";
import type { AuthUser } from "../auth/useAuth";
import { formatNumber } from "./api";

function formatExpiry(value: number | undefined, language: string, never: string): string {
  if (!value) return never;
  return new Date(value * 1000).toLocaleString(language);
}

export function ProfilePanel({ user }: { user: AuthUser }) {
  const { t, i18n } = useTranslation();
  const rows = [
    { key: "username", icon: <UserRound size={16} />, value: user.username },
    { key: "nickname", icon: <BadgeCheck size={16} />, value: user.nickname },
    { key: "email", icon: <Mail size={16} />, value: user.email },
    { key: "mobile", icon: <Phone size={16} />, value: user.mobile },
    { key: "balance", icon: <Coins size={16} />, value: formatNumber(user.power, i18n.language) },
    { key: "vip", icon: <BadgeCheck size={16} />, value: user.vip ? t("space.vip_yes") : t("space.vip_no") },
    {
      key: "expired_at",
      icon: <CalendarClock size={16} />,
      value: formatExpiry(user.expired_time, i18n.language, t("space.never")),
    },
  ];

  return (
    <div className="overflow-hidden rounded-xl border border-hairline bg-canvas">
      <div className="flex items-center gap-4 border-b border-hairline-soft px-6 py-5">
        {user.avatar ? (
          <img src={user.avatar} alt={user.username} className="h-14 w-14 rounded-xl object-cover" />
        ) : (
          <div className="flex h-14 w-14 items-center justify-center rounded-xl bg-brand-50 text-xl font-bold text-brand-600">
            {user.username.slice(0, 1).toUpperCase()}
          </div>
        )}
        <div>
          <h2 className="text-lg font-semibold text-ink">{user.nickname || user.username}</h2>
          <p className="text-sm text-muted">ID: {user.id}</p>
        </div>
      </div>
      <dl className="divide-y divide-hairline-soft px-6">
        {rows.map((row) => (
          <div key={row.key} className="grid grid-cols-[120px_minmax(0,1fr)] items-center gap-4 py-4 text-sm">
            <dt className="flex items-center gap-2 text-muted">
              {row.icon}
              {t(`space.${row.key}`)}
            </dt>
            <dd className="min-w-0 break-all text-ink">{row.value || "—"}</dd>
          </div>
        ))}
      </dl>
    </div>
  );
}

import { Coins, Phone, UserRound } from "lucide-react";
import { useTranslation } from "react-i18next";
import type { AuthUser } from "../auth/useAuth";
import { formatNumber } from "./api";

export function ProfilePanel({ user }: { user: AuthUser }) {
  const { t, i18n } = useTranslation();
  const rows = [
    { key: "username", icon: <UserRound size={16} />, value: user.username },
    { key: "mobile", icon: <Phone size={16} />, value: user.mobile },
    { key: "balance", icon: <Coins size={16} />, value: formatNumber(user.power, i18n.language) },
  ];

  return (
    <div className="overflow-hidden rounded-xl border border-hairline bg-canvas">
      <div className="flex items-center gap-4 border-b border-hairline-soft px-6 py-5">
        <img
          src={user.avatar && user.avatar !== "/images/avatar/user.png" ? user.avatar : "/assets/logo.png"}
          alt={user.username}
          className="h-14 w-14 rounded-xl object-cover"
          onError={(event) => {
            if (event.currentTarget.getAttribute("src") !== "/assets/logo.png") {
              event.currentTarget.src = "/assets/logo.png";
            }
          }}
        />
        <h2 className="min-w-0 break-all text-lg font-semibold text-ink">{user.username}</h2>
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

import { useQuery } from "@tanstack/react-query";
import { CalendarDays, Coins, Gift, TrendingDown, TrendingUp } from "lucide-react";
import { useState, type ReactNode } from "react";
import { useTranslation } from "react-i18next";
import { Link } from "react-router-dom";
import {
  fetchInviteStats,
  fetchPowerLogs,
  fetchPowerStats,
  formatNumber,
  formatTime,
  type PowerLog,
} from "./api";

function todayString(offsetDays = 0): string {
  const date = new Date();
  date.setDate(date.getDate() + offsetDays);
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  return `${date.getFullYear()}-${month}-${day}`;
}

export function PointsPanel({ userId }: { userId: number }) {
  const { t, i18n } = useTranslation();
  const [page, setPage] = useState(1);
  const pageSize = 10;
  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");

  const stats = useQuery({
    queryKey: ["space-power-stats", userId],
    queryFn: fetchPowerStats,
  });
  const invite = useQuery({
    queryKey: ["space-invite-stats", userId],
    queryFn: fetchInviteStats,
    retry: false,
  });
  const logs = useQuery({
    queryKey: ["space-power-logs", userId, page, from, to],
    queryFn: () =>
      fetchPowerLogs({
        page,
        page_size: pageSize,
        date: from && to ? [from, to] : undefined,
      }),
  });

  const total = logs.data?.total ?? 0;
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const hasError = stats.isError || logs.isError;
  const typeLabel = (log: PowerLog) =>
    t(`space.types.${log.type}`, { defaultValue: log.type_str || t("space.other") });

  return (
    <div className="space-y-6">
      {hasError && (
        <div role="alert" className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-amber-200 bg-amber-50 px-5 py-4 text-sm text-amber-900">
          <span>{t("space.load_error")}</span>
          <button
            type="button"
            className="rounded-lg bg-amber-600 px-4 py-1.5 font-semibold text-white hover:bg-amber-700"
            onClick={() => {
              void stats.refetch();
              void logs.refetch();
            }}
          >
            {t("space.retry")}
          </button>
        </div>
      )}

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <StatCard icon={<Coins size={20} />} label={t("space.balance")} value={formatNumber(stats.data?.balance, i18n.language)} loading={stats.isLoading} />
        <StatCard icon={<TrendingDown size={20} />} label={t("space.total_consumed")} value={formatNumber(stats.data?.total, i18n.language)} loading={stats.isLoading} />
        <StatCard icon={<TrendingUp size={20} />} label={t("space.today_consumed")} value={formatNumber(stats.data?.today, i18n.language)} loading={stats.isLoading} />
        <StatCard icon={<Gift size={20} />} label={t("space.invite_rewards")} value={formatNumber(invite.data?.reward_total, i18n.language)} loading={invite.isLoading} />
      </div>

      <div className="rounded-xl border border-hairline bg-canvas">
        <div className="flex flex-wrap items-center justify-between gap-3 border-b border-hairline-soft px-5 py-4">
          <div>
            <h2 className="text-base font-semibold text-ink">{t("space.points")}</h2>
            <p className="mt-1 text-sm text-muted">{t("space.points_description")}</p>
          </div>
          <Link to="/invite" className="inline-flex items-center gap-2 rounded-lg border border-hairline px-3 py-1.5 text-sm font-medium text-body-strong hover:border-brand-500 hover:text-brand-600">
            <Gift size={15} />
            {t("space.invite")}
            {invite.data ? <span className="text-xs text-muted-soft">{t("space.invite_count", { count: invite.data.invite_count })}</span> : null}
          </Link>
        </div>

        <div className="flex flex-wrap items-end gap-3 border-b border-hairline-soft px-5 py-4">
          <label className="flex flex-col gap-1 text-xs font-medium text-muted">
            <span>{t("space.start_date")}</span>
            <input
              type="date"
              value={from}
              max={to || todayString()}
              onChange={(event) => {
                setFrom(event.target.value);
                setPage(1);
              }}
              className="rounded-lg border border-hairline px-3 py-2 text-sm text-body-strong outline-none focus:border-brand-500"
            />
          </label>
          <label className="flex flex-col gap-1 text-xs font-medium text-muted">
            <span>{t("space.end_date")}</span>
            <input
              type="date"
              value={to}
              min={from || undefined}
              max={todayString()}
              onChange={(event) => {
                setTo(event.target.value);
                setPage(1);
              }}
              className="rounded-lg border border-hairline px-3 py-2 text-sm text-body-strong outline-none focus:border-brand-500"
            />
          </label>
          <button
            type="button"
            className="rounded-lg border border-hairline px-3 py-2 text-sm text-body hover:bg-surface-soft"
            onClick={() => {
              setFrom("");
              setTo("");
              setPage(1);
            }}
          >
            {t("space.clear")}
          </button>
          <div className="ml-auto flex items-center gap-2 text-xs text-muted-soft">
            <CalendarDays size={14} />
            {from || to ? t("space.date_range") : t("space.all")}
          </div>
        </div>

        <div className="overflow-x-auto">
          <table className="min-w-full divide-y divide-hairline-soft text-sm">
            <thead className="bg-surface-soft text-left text-xs font-semibold uppercase tracking-wide text-muted">
              <tr>
                <th className="px-5 py-3">{t("space.time")}</th>
                <th className="px-5 py-3">{t("space.type")}</th>
                <th className="px-5 py-3 text-right">{t("space.amount")}</th>
                <th className="px-5 py-3 text-right">{t("space.balance_after")}</th>
                <th className="px-5 py-3">{t("space.detail")}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-hairline-soft">
              {logs.isLoading ? (
                <tr>
                  <td colSpan={5} className="px-5 py-10 text-center text-muted">{t("space.loading")}</td>
                </tr>
              ) : logs.data?.items.length ? (
                logs.data.items.map((log) => (
                  <tr key={log.id} className="hover:bg-surface-soft">
                    <td className="whitespace-nowrap px-5 py-3 text-body">{formatTime(log.created_at, i18n.language)}</td>
                    <td className="whitespace-nowrap px-5 py-3">
                      <span className="rounded-full bg-surface-card px-2.5 py-1 text-xs font-medium text-body-strong">{typeLabel(log)}</span>
                    </td>
                    <td className={`whitespace-nowrap px-5 py-3 text-right font-semibold ${log.mark === 1 ? "text-emerald-600" : "text-rose-600"}`}>
                      {log.mark === 1 ? "+" : "-"}{formatNumber(Math.abs(log.amount), i18n.language)}
                    </td>
                    <td className="whitespace-nowrap px-5 py-3 text-right text-body">{formatNumber(log.balance, i18n.language)}</td>
                    <td className="min-w-[180px] px-5 py-3 text-muted">
                      {log.remark || log.model || "—"}
                    </td>
                  </tr>
                ))
              ) : (
                <tr>
                  <td colSpan={5} className="px-5 py-10 text-center text-muted">{t("space.no_records")}</td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        <div className="flex flex-wrap items-center justify-between gap-3 border-t border-hairline-soft px-5 py-4 text-sm text-muted">
          <span>{t("space.page", { page, total: totalPages })}</span>
          <div className="flex items-center gap-2">
            <button
              type="button"
              disabled={page <= 1}
              onClick={() => setPage((value) => Math.max(1, value - 1))}
              className="rounded-lg border border-hairline px-3 py-1.5 disabled:cursor-not-allowed disabled:opacity-40"
            >
              {t("space.previous")}
            </button>
            <button
              type="button"
              disabled={page >= totalPages}
              onClick={() => setPage((value) => Math.min(totalPages, value + 1))}
              className="rounded-lg border border-hairline px-3 py-1.5 disabled:cursor-not-allowed disabled:opacity-40"
            >
              {t("space.next")}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

function StatCard({ icon, label, value, loading }: { icon: ReactNode; label: string; value: string; loading: boolean }) {
  return (
    <div className="rounded-xl border border-hairline bg-canvas p-5">
      <div className="flex items-center gap-2 text-muted">
        <span className="flex h-9 w-9 items-center justify-center rounded-xl bg-brand-50 text-brand-600">{icon}</span>
        <span className="text-sm">{label}</span>
      </div>
      <div className="mt-4 font-serif text-[1.75rem] font-normal tracking-[-0.02em] text-ink">
        {loading ? <span className="inline-block h-7 w-24 animate-pulse rounded bg-surface-card" /> : value}
      </div>
    </div>
  );
}

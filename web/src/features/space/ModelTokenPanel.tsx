import { useQuery, useQueryClient } from "@tanstack/react-query";
import { Copy, KeyRound } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { fetchModelToken, revokeModelToken, rotateModelToken, type ModelTokenStatus } from "./api";

export function ModelTokenPanel({ userId }: { userId: number }) {
  const { t } = useTranslation();
  const client = useQueryClient();
  const queryKey = ["model-token", userId];
  const status = useQuery({ queryKey, queryFn: fetchModelToken, retry: false });
  const [secret, setSecret] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const [copied, setCopied] = useState(false);

  async function changeToken(revoke: boolean) {
    if (status.data?.enabled && !window.confirm(t(revoke ? "model_api.confirm_revoke" : "model_api.confirm_reset"))) return;
    setBusy(true);
    setError("");
    setSecret("");
    setCopied(false);
    try {
      if (revoke) {
        await revokeModelToken();
        client.setQueryData<ModelTokenStatus>(queryKey, { enabled: false });
      } else {
        const result = await rotateModelToken();
        // 明文只留在当前组件中，不写入持久存储或查询缓存。
        setSecret(result.token);
        client.setQueryData<ModelTokenStatus>(queryKey, { enabled: true, created_at: result.created_at });
      }
    } catch {
      setError(t("model_api.action_error"));
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="rounded-xl border border-hairline bg-canvas px-6 py-5" aria-labelledby="model-api-title">
      <h2 id="model-api-title" className="flex items-center gap-2 text-base font-semibold text-ink">
        <KeyRound size={18} aria-hidden="true" />{t("model_api.title")}
      </h2>
      <p className="mt-2 text-sm leading-6 text-muted">{t("model_api.description")}</p>
      <dl className="mt-4 space-y-3 text-sm">
        <div>
          <dt className="text-muted">{t("model_api.base_url")}</dt>
          <dd className="mt-1 select-all break-all font-mono text-ink">{window.location.origin}/v1</dd>
        </div>
        <div>
          <dt className="text-muted">{t("model_api.protocol")}</dt>
          <dd className="mt-1 text-ink">OpenAI Chat Completions</dd>
        </div>
      </dl>
      {status.isLoading ? <p role="status" className="mt-4 text-sm text-muted">{t("space.loading")}</p> : status.isError ? (
        <div role="alert" className="mt-4 flex items-center gap-3 text-sm text-rose-700">
          {t("space.load_error")}
          <button type="button" className="underline" onClick={() => void status.refetch()}>{t("space.retry")}</button>
        </div>
      ) : (
        <>
          <p className="mt-4 text-sm text-muted">{t(status.data?.enabled ? "model_api.enabled" : "model_api.disabled")}</p>
          {secret && (
            <div className="mt-3">
              <label htmlFor="model-api-token" className="text-sm font-medium text-ink">{t("model_api.token")}</label>
              <div className="mt-2 flex gap-2">
                <input id="model-api-token" readOnly value={secret} onFocus={(event) => event.currentTarget.select()}
                  className="min-w-0 flex-1 rounded-lg border border-hairline bg-surface-soft px-3 py-2 font-mono text-sm text-ink" />
                <button type="button" className="inline-flex shrink-0 items-center gap-1.5 rounded-lg border border-hairline px-3 py-2 text-sm text-ink hover:bg-surface-soft" onClick={async () => {
                  try { await navigator.clipboard.writeText(secret); setCopied(true); }
                  catch { setError(t("model_api.copy_error")); }
                }}>
                  <Copy size={15} aria-hidden="true" />{t(copied ? "model_api.copied" : "model_api.copy")}
                </button>
              </div>
              <p className="mt-2 text-xs leading-5 text-muted">{t("model_api.save_once")}</p>
            </div>
          )}
          <div className="mt-4 flex flex-wrap gap-3">
            <button type="button" disabled={busy} onClick={() => void changeToken(false)} className="btn-primary rounded-lg px-4 py-2 text-sm font-medium disabled:opacity-50">
              {t(busy ? "model_api.processing" : status.data?.enabled ? "model_api.reset" : "model_api.create")}
            </button>
            {status.data?.enabled && <button type="button" disabled={busy} onClick={() => void changeToken(true)} className="rounded-lg border border-hairline px-4 py-2 text-sm text-muted hover:bg-surface-soft disabled:opacity-50">{t("model_api.revoke")}</button>}
          </div>
        </>
      )}
      {error && <p role="alert" className="mt-3 text-sm text-rose-700">{error}</p>}
    </section>
  );
}

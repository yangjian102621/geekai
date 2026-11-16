import { useState } from "react";
import { useSearchParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { QRCodeCanvas } from "qrcode.react";
import { Check, Copy, Download, Link as LinkIcon, Users, Gift } from "lucide-react";
import { useAuth } from "../features/auth/useAuth";
import { useInviteReward, useInviteSummary } from "../features/invite/useInviteSummary";
import type { InviteSummary } from "../api/client";
import {
  currentMacArch,
  currentOs,
  downloadOptions,
  recommendedHref,
  useDownloadLinks,
  type DownloadLinks,
} from "../hooks/useLatestSetup";
import { openLoginPage } from "../config";

export function InvitePage() {
  const [searchParams] = useSearchParams();
  const aff = searchParams.get("aff");
  const { data: user, isLoading: authLoading, isError: authError } = useAuth();
  const { data: reward } = useInviteReward();
  const { data: summary, isLoading: summaryLoading, isError: summaryError, refetch } = useInviteSummary(!!user && !aff);
  const links = useDownloadLinks();

  if (authLoading) return <LoadingView />;
  if (aff) return <InviteeView isLoggedIn={!!user} authError={authError} links={links} />;
  if (!user) return <GuestView reward={reward ?? null} />;
  if (summaryLoading && !summary) return <LoadingView />;
  if (summaryError || !summary) return <ErrorView onRetry={() => void refetch()} />;
  return <LoggedInView summary={summary} reward={reward ?? null} />;
}

/* ---------- 加载 / 出错 ---------- */
function LoadingView() {
  const { t } = useTranslation();
  return (
    <section className="invite-hero py-32">
      <div className="invite-hero-bg" />
      <div className="relative max-w-3xl mx-auto px-4 text-center" role="status">
        <div className="invite-loading-spinner mx-auto" />
        <p className="mt-4 text-body">{t("invite.loading")}</p>
      </div>
    </section>
  );
}

function ErrorView({ onRetry }: { onRetry: () => void }) {
  const { t } = useTranslation();
  return (
    <section className="invite-hero py-32">
      <div className="invite-hero-bg" />
      <div className="relative max-w-3xl mx-auto px-4 text-center" role="alert">
        <p className="text-body-strong mb-5">{t("invite.error")}</p>
        <button type="button" className="btn-primary px-5 py-2 rounded-lg" onClick={onRetry}>
          {t("invite.retry")}
        </button>
      </div>
    </section>
  );
}

/* ---------- 好友邀请视图（?aff=） ---------- */
function InviteeView({
  isLoggedIn,
  authError,
  links,
}: {
  isLoggedIn: boolean;
  authError: boolean;
  links: DownloadLinks;
}) {
  const { t } = useTranslation();
  const macUnknown = currentOs() === "mac" && currentMacArch() === "unknown";
  const [, macArm64, macX64] = downloadOptions(links);
  const downloadUrl = recommendedHref(links) ?? "/#download";

  return (
    <section className="invite-hero pt-28 pb-20">
      <div className="invite-hero-bg" />
      <div className="relative max-w-3xl mx-auto px-4 text-center">
        <div className="invite-hero-badge">
          <span className="w-1.5 h-1.5 rounded-full bg-coral animate-pulse" />
          <span>{t("invite.invitee.badge")}</span>
        </div>
        <h1 className="invite-hero-title">{t("invite.invitee.title")}</h1>
        <p className="invite-hero-subtitle">{t("invite.invitee.desc")}</p>
        <div className="invite-hero-cta-row">
          {isLoggedIn ? (
            macUnknown ? (
              <>
                <a href={macArm64.href ?? "/#download"} className="invite-cta-primary">
                  <Download size={20} />
                  {t("download.mac_arm64")}
                </a>
                <a href={macX64.href ?? "/#download"} className="invite-cta-primary">
                  <Download size={20} />
                  {t("download.mac_intel")}
                </a>
              </>
            ) : (
              <a href={downloadUrl} className="invite-cta-primary">
                <Download size={20} />
                {t("invite.invitee.download_btn")}
              </a>
            )
          ) : (
            <button type="button" className="invite-cta-primary" onClick={openLoginPage}>
              <LinkIcon size={20} />
              {t("invite.invitee.login_btn")}
            </button>
          )}
        </div>
        <p className="text-sm text-body">
          {isLoggedIn ? t("invite.invitee.logged_in_hint") : t("invite.invitee.login_note")}
        </p>
        <p className="text-xs text-muted mt-2">{t("invite.invitee.download_note")}</p>
        {authError && <p role="alert" className="mt-6 text-body-strong">{t("invite.auth_error")}</p>}
      </div>
    </section>
  );
}

/* ---------- 未登录视图 ---------- */
function GuestView({ reward }: { reward: number | null }) {
  const { t } = useTranslation();
  return (
    <div>
      <section className="invite-hero pt-28 pb-20">
        <div className="invite-hero-bg" />
        <div className="relative max-w-4xl mx-auto px-4 text-center">
          <div className="invite-hero-badge">
            <span className="w-1.5 h-1.5 rounded-full bg-coral animate-pulse" />
            <span>{t("invite.guest.badge")}</span>
          </div>
          <h1 className="invite-hero-title">
            {t("invite.title")}
            <br />
            <span className="invite-hero-title-accent">{t("invite.subtitle")}</span>
          </h1>
          <p className="invite-hero-subtitle">
            {reward != null && reward > 0 ? t("invite.reward.per_invite", { points: reward.toLocaleString() }) : t("invite.reward.fallback")}
          </p>
          <div className="invite-hero-cta-row">
            <button type="button" className="invite-cta-primary" onClick={openLoginPage}>
              <LinkIcon size={20} />
              {t("invite.guest.login_btn")}
            </button>
          </div>
          {reward != null && reward > 0 && (
            <div className="invite-hero-stats">
              <div className="text-center">
                <div className="invite-hero-stat-num">{reward.toLocaleString()}</div>
                <div className="invite-hero-stat-label">{t("invite.reward.unit")}</div>
              </div>
            </div>
          )}
        </div>
      </section>

      <StepsSection />
      <RewardsSection reward={reward} />
      <FaqSection reward={reward} />
      <BottomCta />
    </div>
  );
}

/* ---------- 已登录视图 ---------- */
function LoggedInView({ summary, reward }: { summary: InviteSummary; reward: number | null }) {
  const { t } = useTranslation();
  return (
    <div>
      <section className="invite-hero pt-28 pb-20">
        <div className="invite-hero-bg" />
        <div className="relative max-w-4xl mx-auto px-4 text-center">
          <div className="invite-hero-badge">
            <span className="w-1.5 h-1.5 rounded-full bg-coral animate-pulse" />
            <span>{t("invite.guest.badge")}</span>
          </div>
          <h1 className="invite-hero-title">
            {t("invite.title")}
            <br />
            <span className="invite-hero-title-accent">{t("invite.subtitle")}</span>
          </h1>
          <div className="invite-hero-stats">
            <div className="text-center">
              <div className="invite-hero-stat-num">{summary.invited_count.toLocaleString()}</div>
              <div className="invite-hero-stat-label">{t("invite.stats.invited")}</div>
            </div>
            <div className="invite-hero-stat-divider" />
            <div className="text-center">
              <div className="invite-hero-stat-num">{summary.total_invite_points.toLocaleString()}</div>
              <div className="invite-hero-stat-label">{t("invite.stats.reward")}</div>
            </div>
            <div className="invite-hero-stat-divider" />
            <div className="text-center">
              <div className="invite-hero-stat-num">{summary.today_invites.toLocaleString()}</div>
              <div className="invite-hero-stat-label">{t("invite.stats.today")}</div>
            </div>
          </div>
        </div>
      </section>

      <section className="invite-share-section">
        <div className="max-w-4xl mx-auto px-4">
          <div className="invite-share-card">
            <div className="invite-share-qr-col">
              <div className="invite-qr-wrapper">
                <QRCodeCanvas value={summary.invite_url} size={180} level="M" />
              </div>
              <p className="text-xs text-body mt-3">{t("invite.qr_hint")}</p>
            </div>
            <div className="invite-share-link-col">
              <label htmlFor="invite-link" className="invite-link-label">{t("invite.share_title")}</label>
              <div className="invite-share-link-row">
                <input
                  id="invite-link"
                  readOnly
                  value={summary.invite_url}
                  className="invite-share-link-input"
                  onFocus={(event) => event.target.select()}
                />
                <CopyButton text={summary.invite_url} />
              </div>
              <p className="text-sm text-body mt-3">{t("invite.share_hint")}</p>
            </div>
          </div>
        </div>
      </section>

      <StepsSection />
      <RewardsSection reward={reward} />
      <FaqSection reward={reward} />
      <BottomCta url={summary.invite_url} />
    </div>
  );
}

/* ---------- 三步流程 ---------- */
function StepsSection() {
  const { t } = useTranslation();
  const steps = [
    { num: 1, title: t("invite.steps.step1_title"), desc: t("invite.steps.step1_desc") },
    { num: 2, title: t("invite.steps.step2_title"), desc: t("invite.steps.step2_desc") },
    { num: 3, title: t("invite.steps.step3_title"), desc: t("invite.steps.step3_desc") },
  ];
  return (
    <section className="invite-steps-section py-16 lg:py-20">
      <div className="max-w-5xl mx-auto px-4 text-center">
        <h2 className="invite-section-title">{t("invite.steps.title")}</h2>
        <p className="invite-section-subtitle mb-10">{t("invite.steps.subtitle")}</p>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {steps.map((step) => (
            <div key={step.num} className="invite-step-card">
              <div className="invite-step-num">{step.num}</div>
              <div className="invite-step-title">{step.title}</div>
              <div className="invite-step-desc">{step.desc}</div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

/* ---------- 奖励卡片 ---------- */
function RewardsSection({ reward }: { reward: number | null }) {
  const { t } = useTranslation();
  return (
    <section className="invite-rewards-section py-16 lg:py-20">
      <div className="max-w-4xl mx-auto px-4">
        <div className="text-center mb-10">
          <h2 className="invite-section-title">{t("invite.rewards.title")}</h2>
          <p className="invite-section-subtitle">{t("invite.rewards.subtitle")}</p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="invite-big-card invite-big-card-inviter">
            <div className="invite-big-card-top"><Users size={22} /><h2 className="invite-big-card-title">{t("invite.rewards.inviter_title")}</h2></div>
            <p className="invite-big-card-desc">
              {reward != null && reward > 0
                ? t("invite.rewards.inviter_desc", { points: reward.toLocaleString() })
                : t("invite.rewards.inviter_fallback")}
            </p>
          </div>
          <div className="invite-big-card invite-big-card-invitee">
            <div className="invite-big-card-top"><Gift size={22} /><h2 className="invite-big-card-title">{t("invite.rewards.usage_title")}</h2></div>
            <p className="invite-big-card-desc">{t("invite.rewards.usage_desc")}</p>
          </div>
        </div>
      </div>
    </section>
  );
}

/* ---------- FAQ ---------- */
function FaqSection({ reward }: { reward: number | null }) {
  const { t } = useTranslation();
  const faqs = [
    { q: t("invite.faq.q1"), a: t("invite.faq.a1") },
    {
      q: t("invite.faq.q2"),
      a: reward != null && reward > 0 ? t("invite.faq.a2_with_reward", { points: reward.toLocaleString() }) : t("invite.faq.a2"),
    },
    { q: t("invite.faq.q3"), a: t("invite.faq.a3") },
    { q: t("invite.faq.q4"), a: t("invite.faq.a4") },
  ];
  return (
    <section className="invite-faq-section py-16 lg:py-20">
      <div className="max-w-3xl mx-auto px-4">
        <h2 className="invite-section-title text-center">{t("invite.faq.title")}</h2>
        <div className="invite-faq-list mt-8">
          {faqs.map((faq, index) => (
            <details key={faq.q} className="invite-faq-item" open={index === 0}>
              <summary className="invite-faq-q">{faq.q}</summary>
              <p className="invite-faq-a">{faq.a}</p>
            </details>
          ))}
        </div>
      </div>
    </section>
  );
}

/* ---------- 底部 CTA ---------- */
function BottomCta({ url }: { url?: string }) {
  const { t } = useTranslation();
  const [copied, setCopied] = useState(false);
  const copy = async () => {
    if (!url) return;
    try {
      await writeClipboard(url);
      setCopied(true);
      window.setTimeout(() => setCopied(false), 2000);
    } catch {
      // 剪贴板被拒绝时保留原始页面，不打断用户操作
    }
  };
  return (
    <section className="invite-bottom-cta py-16 lg:py-20 text-center">
      <div className="max-w-3xl mx-auto px-4">
        <h2 className="invite-bottom-cta-title">{t("invite.bottom_cta_title")}</h2>
        <p className="invite-bottom-cta-sub">{t("invite.bottom_cta_sub")}</p>
        {url ? (
          <button type="button" className="invite-cta-primary" onClick={() => void copy()}>
            <LinkIcon size={20} />
            {copied ? t("invite.copied") : t("invite.bottom_cta_copy")}
          </button>
        ) : (
          <button type="button" className="invite-cta-primary" onClick={openLoginPage}>
            <LinkIcon size={20} />
            {t("invite.guest.login_btn")}
          </button>
        )}
      </div>
    </section>
  );
}

/* ---------- 复制按钮 ---------- */
function CopyButton({ text }: { text: string }) {
  const { t } = useTranslation();
  const [copied, setCopied] = useState(false);
  const [copyError, setCopyError] = useState(false);
  const copy = async () => {
    try {
      await writeClipboard(text);
      setCopied(true);
      setCopyError(false);
    } catch {
      setCopyError(true);
    }
  };
  return (
    <>
      <button type="button" onClick={() => void copy()} className="invite-copy-btn-big">
        {copied ? <Check size={16} /> : <Copy size={16} />}
        {copied ? t("invite.copied") : t("invite.copy")}
      </button>
      {copyError && <p role="alert" className="text-sm text-body-strong mt-3">{t("invite.copy_manual")}</p>}
    </>
  );
}

async function writeClipboard(text: string): Promise<void> {
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(text);
    return;
  }
  const input = document.createElement("input");
  input.value = text;
  document.body.appendChild(input);
  input.select();
  document.execCommand("copy");
  document.body.removeChild(input);
}

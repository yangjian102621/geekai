import { useState } from "react";
import { useSearchParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { QRCodeCanvas } from "qrcode.react";
import { Check, Copy, Link as LinkIcon, Users, Gift } from "lucide-react";
import { useAuth } from "../features/auth/useAuth";
import { useInviteReward, useInviteSummary } from "../features/invite/useInviteSummary";
import { useDownloadLinks, recommendedHref } from "../hooks/useLatestSetup";
import { DownloadMenu } from "../components/DownloadMenu";
import { openLoginPage } from "../config";

export function InvitePage() {
  const { t } = useTranslation();
  const [searchParams] = useSearchParams();
  const aff = searchParams.get("aff");
  const { data: user, isLoading: authLoading, isError: authError } = useAuth();
  const { data: reward } = useInviteReward();
  const { data: summary, isLoading, isError, refetch } = useInviteSummary(!!user && !aff);
  const links = useDownloadLinks();
  const [copied, setCopied] = useState(false);
  const [copyError, setCopyError] = useState(false);

  const copyLink = async () => {
    if (!summary) return;
    try {
      await navigator.clipboard.writeText(summary.invite_url);
      setCopied(true);
      setCopyError(false);
    } catch {
      setCopyError(true);
    }
  };

  if (authLoading || (!!user && !aff && isLoading)) {
    return (
      <section className="invite-hero py-32">
        <div className="invite-hero-bg" />
        <div className="relative max-w-3xl mx-auto px-4 text-center" role="status">
          <div className="invite-loading-spinner mx-auto" />
          <p className="mt-4 text-slate-600">{t("invite.loading")}</p>
        </div>
      </section>
    );
  }

  return (
    <div>
      <section className="invite-hero pt-28 pb-20">
        <div className="invite-hero-bg" />
        <div className="relative max-w-4xl mx-auto px-4 text-center">
          <h1 className="invite-hero-title">{aff ? "好友邀请你加入 ThinkBuddy" : t("invite.title")}</h1>
          <p className="invite-hero-subtitle">
            {aff
              ? "通过邀请链接注册账号，下载客户端，开启你的学习旅程。"
              : "分享 ThinkBuddy，邀请好友一起阅读、研究和创作。"}
          </p>
          {reward != null && (
            <p className="text-slate-700 mb-8">
              {reward > 0 ? `每成功邀请一位好友注册，邀请人可获得 ${reward.toLocaleString()} 积分。` : "邀请好友一起使用 ThinkBuddy。"}
            </p>
          )}
          {user && aff ? (
            <DownloadMenu
              links={links}
              align="center"
              className="inline-block"
              trigger={<a href={recommendedHref(links) ?? "/#download"} className="invite-cta-primary">{t("nav.download")}</a>}
            />
          ) : !user ? (
            <button className="invite-cta-primary" onClick={openLoginPage}>
              <LinkIcon size={20} />
              {aff ? "注册并下载" : t("invite.guest.login_btn")}
            </button>
          ) : null}
          {authError && <p role="alert" className="mt-6 text-slate-700">暂时无法确认登录状态，请稍后重试。</p>}
          {summary && !aff && (
            <div className="invite-hero-stats">
              <div className="text-center">
                <div className="invite-hero-stat-num">{summary.invited_count}</div>
                <div className="invite-hero-stat-label">累计邀请</div>
              </div>
              <div className="invite-hero-stat-divider" />
              <div className="text-center">
                <div className="invite-hero-stat-num">{summary.total_invite_points.toLocaleString()}</div>
                <div className="invite-hero-stat-label">奖励积分</div>
              </div>
              <div className="invite-hero-stat-divider" />
              <div className="text-center">
                <div className="invite-hero-stat-num">{summary.today_invites}</div>
                <div className="invite-hero-stat-label">今日邀请</div>
              </div>
            </div>
          )}
        </div>
      </section>

      {!!user && !aff && (isError || !summary) ? (
        <section className="py-16 px-4 text-center" role="alert">
          <p className="text-slate-700 mb-5">{t("invite.error")}</p>
          <button className="btn-primary px-5 py-2 rounded-lg" onClick={() => void refetch()}>重试</button>
        </section>
      ) : summary && !aff ? (
        <section className="invite-share-section">
          <div className="max-w-4xl mx-auto px-4">
            <div className="invite-share-card">
              <div className="invite-share-qr-col">
                <div className="invite-qr-wrapper"><QRCodeCanvas value={summary.invite_url} size={180} level="M" /></div>
                <p className="text-xs text-slate-600 mt-3">{t("invite.qr_hint")}</p>
              </div>
              <div className="invite-share-link-col">
                <label htmlFor="invite-link" className="invite-link-label">{t("invite.share_title")}</label>
                <div className="invite-share-link-row">
                  <input id="invite-link" readOnly value={summary.invite_url} className="invite-share-link-input" onFocus={(event) => event.target.select()} />
                  <button onClick={() => void copyLink()} className="invite-copy-btn-big">
                    {copied ? <Check size={16} /> : <Copy size={16} />}
                    {copied ? t("invite.copied") : t("invite.copy")}
                  </button>
                </div>
                {copyError && <p role="alert" className="text-sm text-slate-700 mt-3">请选中上方链接手动复制。</p>}
                <p className="text-sm text-slate-600 mt-3">好友通过此链接注册后，系统会记录邀请关系。</p>
              </div>
            </div>
          </div>
        </section>
      ) : null}

      <section className="invite-steps-section py-16 lg:py-20">
        <div className="max-w-5xl mx-auto px-4 text-center">
          <h2 className="invite-section-title">三步邀请好友</h2>
          <p className="invite-section-subtitle mb-10">一起开启高效学习</p>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {[
              { title: "获取邀请链接", desc: "登录账号后，复制你的专属邀请链接。" },
              { title: "分享给好友", desc: "好友打开链接，完成账号注册。" },
              { title: "获得积分", desc: "注册成功后，邀请奖励计入积分余额。" },
            ].map((step, index) => (
              <div key={step.title} className="invite-step-card">
                <div className="invite-step-num">{index + 1}</div>
                <div className="invite-step-title">{step.title}</div>
                <div className="invite-step-desc">{step.desc}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      <section className="invite-rewards-section py-16 lg:py-20">
        <div className="max-w-4xl mx-auto px-4 grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="invite-big-card invite-big-card-inviter">
            <div className="invite-big-card-top"><Users size={22} /><h2 className="invite-big-card-title">邀请好友</h2></div>
            <p className="invite-big-card-desc">好友通过专属链接注册，系统自动绑定邀请关系；奖励按当前活动配置发放。</p>
          </div>
          <div className="invite-big-card invite-big-card-invitee">
            <div className="invite-big-card-top"><Gift size={22} /><h2 className="invite-big-card-title">积分直接使用</h2></div>
            <p className="invite-big-card-desc">奖励进入统一积分账户，用于 AI 服务消费。</p>
          </div>
        </div>
      </section>
    </div>
  );
}

import { useState } from "react";
import { Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useReveal } from "../hooks/useReveal";
import { useLatestSetup } from "../hooks/useLatestSetup";

export function PricingPage() {
  const { t } = useTranslation();
  const ref = useReveal<HTMLDivElement>();
  const [yearly, setYearly] = useState(false);
  const downloadUrl = useLatestSetup();

  const proMonthly = 10;
  const proYearly = 96;

  return (
    <div ref={ref} className="bg-slate-50 min-h-screen">
      <section className="pt-28 pb-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <p className="inline-block px-3 py-1 rounded-full bg-brand-50 border border-brand-500/20 text-brand-700 text-xs font-semibold mb-4" data-reveal>
            {t("pricing.badge")}
          </p>
          <h1 className="text-3xl sm:text-5xl font-bold tracking-tight text-slate-900 mb-4" data-reveal>{t("pricing.title")}</h1>
          <p className="max-w-2xl mx-auto text-slate-600" data-reveal>{t("pricing.subtitle")}</p>

          <div className="inline-flex items-center gap-1 mt-8 p-1 bg-slate-200/60 rounded-full" data-reveal>
            <button
              onClick={() => setYearly(false)}
              className={`billing-toggle px-5 py-2 rounded-full text-sm font-semibold ${!yearly ? "billing-toggle-active" : ""}`}
            >
              {t("pricing.monthly")}
            </button>
            <button
              onClick={() => setYearly(true)}
              className={`billing-toggle px-5 py-2 rounded-full text-sm font-semibold ${yearly ? "billing-toggle-active" : ""}`}
            >
              {t("pricing.yearly")}
            </button>
          </div>
        </div>
      </section>

      {/* 个人套餐 */}
      <section className="pb-20">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Free */}
            <div className="pricing-card bg-white border border-slate-200 rounded-2xl p-8 flex flex-col" data-reveal>
              <h3 className="text-lg font-bold text-slate-900 mb-1">{t("pricing.free.name")}</h3>
              <p className="text-sm text-slate-500 mb-4">{t("pricing.free.tagline")}</p>
              <div className="mb-1"><span className="text-3xl font-bold text-slate-900">¥0</span><span className="text-sm text-slate-500">{t("pricing.per_month")}</span></div>
              <p className="text-xs text-brand-600 font-medium mb-6">{t("pricing.free.trial_note")}</p>
              <ul className="space-y-3 text-sm text-slate-600 mb-8">
                <li>{t("pricing.free.f1")}{t("pricing.free.f1u")}</li>
                <li>{t("pricing.free.f2")}</li>
                <li>{t("pricing.free.f3")}</li>
                <li>{t("pricing.free.f4")}</li>
                <li>{t("pricing.free.f5")}</li>
              </ul>
              <a href={downloadUrl ?? "/#download"} className="btn-secondary mt-auto text-center px-5 py-2.5 rounded-lg text-sm font-semibold">{t("pricing.free.cta")}</a>
            </div>

            {/* Pro */}
            <div className="pricing-card pricing-card-featured bg-white border-2 border-brand-500 rounded-2xl p-8 flex flex-col relative" data-reveal>
              <span className="absolute -top-3 left-1/2 -translate-x-1/2 px-3 py-1 rounded-full bg-brand-600 text-white text-xs font-semibold">{t("pricing.pro.badge")}</span>
              <h3 className="text-lg font-bold text-slate-900 mb-1">{t("pricing.pro.name")}</h3>
              <p className="text-sm text-slate-500 mb-4">{t("pricing.pro.tagline")}</p>
              <div className="mb-1">
                <span className="text-3xl font-bold text-slate-900">¥{yearly ? proYearly : proMonthly}</span>
                <span className="text-sm text-slate-500">{yearly ? "/年" : t("pricing.per_month")}</span>
              </div>
              <p className="text-xs text-brand-600 font-medium mb-6">{t("pricing.pro.trial_note")}</p>
              <ul className="space-y-3 text-sm text-slate-600 mb-8">
                <li>{t("pricing.pro.f1")}</li>
                <li className="text-brand-600 font-medium">{t("pricing.pro.f1b")}</li>
                <li>{t("pricing.pro.f2")}</li>
                <li>{t("pricing.pro.f3")}</li>
                <li>{t("pricing.pro.f4")}</li>
                <li>{t("pricing.pro.f5")}</li>
                <li>{t("pricing.pro.f6")}</li>
              </ul>
              <a href={downloadUrl ?? "/#download"} className="btn-primary mt-auto text-center px-5 py-2.5 rounded-lg text-sm font-semibold">{t("pricing.pro.cta")}</a>
            </div>
          </div>
        </div>
      </section>

      {/* FAQ */}
      <section className="py-16 bg-white">
        <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
          <h2 className="text-center text-2xl sm:text-3xl font-bold tracking-tight text-slate-900 mb-8" data-reveal>{t("pricing.faq.title")}</h2>
          <div className="space-y-3">
            {(["q1", "q2", "q3", "q4", "q5"] as const).map((q) => (
              <details key={q} className="faq-item bg-white border border-slate-200 rounded-xl p-5" data-reveal>
                <summary className="font-semibold text-slate-900 cursor-pointer">{t(`pricing.faq.${q}`)}</summary>
                <p className="mt-3 text-sm text-slate-600 leading-relaxed">{t(`pricing.faq.a${q.slice(1)}`)}</p>
              </details>
            ))}
          </div>
          <div className="text-center mt-8">
            <Link to="/" className="inline-flex items-center justify-center px-6 py-3 rounded-xl text-sm font-semibold text-white bg-brand-600 hover:bg-brand-700 transition">
              返回首页
            </Link>
          </div>
        </div>
      </section>
    </div>
  );
}

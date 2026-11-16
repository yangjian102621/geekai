import { useTranslation } from "react-i18next";
import { useLocation, Link } from "react-router-dom";
import { useLatestSetup } from "../../hooks/useLatestSetup";
import { docsURL } from "../../config";

export function Footer() {
  const { t } = useTranslation();
  const loc = useLocation();
  const isInvite = loc.pathname.startsWith("/invite");
  const downloadUrl = useLatestSetup();

  if (isInvite) {
    return (
      <footer className="bg-slate-900 text-slate-400 mt-auto">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 text-center">
          <Link to="/" className="text-xs text-slate-400 hover:text-slate-200 transition">
            ← 返回首页
          </Link>
          <p className="text-xs text-slate-500 mt-3">© 2026 ThinkBuddy. Powered by ThinkBuddy (hanyoai).</p>
        </div>
      </footer>
    );
  }

  return (
    <footer className="bg-slate-900 text-slate-300">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
        <div className="footer-grid">
          <div>
            <div className="flex items-center gap-2 mb-3">
              <img src="/assets/logo.png" alt="ThinkBuddy" className="h-7 w-7 rounded-lg" />
              <span className="text-lg font-bold text-white">ThinkBuddy</span>
            </div>
            <p className="text-sm text-slate-400 max-w-xs">{t("footer.tagline")}</p>
          </div>
          <div>
            <h4 className="text-xs font-semibold text-white uppercase tracking-wider mb-3">{t("footer.col_product")}</h4>
            <ul className="space-y-2 text-sm">
              <li><a href={downloadUrl ?? "/#download"} className="text-slate-400 hover:text-white transition">{t("footer.download")}</a></li>
              <li><a href="#versions" className="text-slate-400 hover:text-white transition">{t("footer.versions")}</a></li>
            </ul>
          </div>
          <div>
            <h4 className="text-xs font-semibold text-white uppercase tracking-wider mb-3">{t("footer.col_resources")}</h4>
            <ul className="space-y-2 text-sm">
              {docsURL && <li><a href={docsURL} target="_blank" rel="noopener noreferrer" className="text-slate-400 hover:text-white transition">{t("footer.docs")}</a></li>}
              <li><a href="#" className="text-slate-400 hover:text-white transition">{t("footer.faq")}</a></li>
            </ul>
          </div>
          <div>
            <h4 className="text-xs font-semibold text-white uppercase tracking-wider mb-3">{t("footer.col_legal")}</h4>
            <ul className="space-y-2 text-sm">
              <li><a href="#" className="text-slate-400 hover:text-white transition">{t("footer.privacy")}</a></li>
              <li><a href="#" className="text-slate-400 hover:text-white transition">{t("footer.terms")}</a></li>
            </ul>
          </div>
        </div>
        <div className="border-t border-white/10 mt-12 pt-6 text-center">
          <p className="text-xs text-slate-500">{t("footer.copyright")}</p>
        </div>
      </div>
    </footer>
  );
}

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
      <footer className="footer-shell text-on-dark-soft mt-auto">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 text-center">
          <Link to="/" className="text-xs text-on-dark-soft hover:text-on-dark transition">
            ← 返回首页
          </Link>
          <p className="text-xs text-on-dark-soft/70 mt-3">© 2026 ThinkBuddy. Powered by ThinkBuddy (hanyoai).</p>
        </div>
      </footer>
    );
  }

  return (
    <footer className="footer-shell">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-14">
        <div className="footer-grid">
          <div>
            <div className="flex items-center gap-2 mb-3">
              <img src="/assets/logo.png" alt="ThinkBuddy" className="footer-logo h-7 w-7 rounded-lg" />
              <span className="text-lg font-semibold text-on-dark">ThinkBuddy</span>
            </div>
            <p className="text-sm text-on-dark-soft max-w-xs">{t("footer.tagline")}</p>
          </div>
          <div>
            <h4 className="text-xs font-semibold text-on-dark uppercase tracking-wider mb-3">{t("footer.col_product")}</h4>
            <ul className="space-y-2 text-sm">
              <li><a href={downloadUrl ?? "/#download"} className="text-on-dark-soft hover:text-on-dark transition">{t("footer.download")}</a></li>
              <li><a href="#versions" className="text-on-dark-soft hover:text-on-dark transition">{t("footer.versions")}</a></li>
            </ul>
          </div>
          <div>
            <h4 className="text-xs font-semibold text-on-dark uppercase tracking-wider mb-3">{t("footer.col_resources")}</h4>
            <ul className="space-y-2 text-sm">
              {docsURL && <li><a href={docsURL} target="_blank" rel="noopener noreferrer" className="text-on-dark-soft hover:text-on-dark transition">{t("footer.docs")}</a></li>}
              <li><a href="#" className="text-on-dark-soft hover:text-on-dark transition">{t("footer.faq")}</a></li>
            </ul>
          </div>
          <div>
            <h4 className="text-xs font-semibold text-on-dark uppercase tracking-wider mb-3">{t("footer.col_legal")}</h4>
            <ul className="space-y-2 text-sm">
              <li><a href="#" className="text-on-dark-soft hover:text-on-dark transition">{t("footer.privacy")}</a></li>
              <li><a href="#" className="text-on-dark-soft hover:text-on-dark transition">{t("footer.terms")}</a></li>
            </ul>
          </div>
        </div>
        <div className="border-t border-white/10 mt-12 pt-6 text-center">
          <p className="text-xs text-on-dark-soft/70">{t("footer.copyright")}</p>
        </div>
      </div>
    </footer>
  );
}

import { Link } from "react-router-dom";

export function NotFoundPage() {
  return (
    <div className="flex flex-col items-center justify-center py-32 text-center">
      <h1 className="text-6xl font-bold text-slate-300">404</h1>
      <p className="mt-4 text-slate-600">页面不存在</p>
      <Link
        to="/"
        className="mt-6 inline-flex items-center px-5 py-2 rounded-lg text-sm font-semibold text-white bg-brand-600 hover:bg-brand-700 transition"
      >
        返回首页
      </Link>
    </div>
  );
}

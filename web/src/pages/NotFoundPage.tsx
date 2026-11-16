import { Link } from "react-router-dom";

export function NotFoundPage() {
  return (
    <div className="flex flex-col items-center justify-center py-32 text-center">
      <h1 className="font-serif text-7xl font-normal text-muted-soft/50">404</h1>
      <p className="mt-4 text-body">页面不存在</p>
      <Link
        to="/"
        className="mt-6 inline-flex items-center px-5 py-2 rounded-lg text-sm font-medium text-white bg-primary hover:bg-primary-hover transition"
      >
        返回首页
      </Link>
    </div>
  );
}

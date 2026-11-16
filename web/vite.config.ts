import { cpSync } from "node:fs";
import { resolve, extname } from "node:path";
import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), "");
  const apiTarget = env.VITE_API_HOST || "http://localhost:5678";
  const adminTarget = "http://localhost:8889";
  const apiProxy = Object.fromEntries(
    ["/api", "/static", "/healthz", "/update"].map((path) => [
      path,
      { target: apiTarget, changeOrigin: true, ws: true },
    ]),
  );

  return {
    plugins: [
      react(),
      tailwindcss(),
      {
        name: "geekai-admin",
        apply: "build",
        closeBundle() {
          cpSync(resolve("../admin-web/dist"), resolve("dist/admin"), { recursive: true });
          for (const directory of ["images", "medias"]) {
            cpSync(resolve("../admin-web/public", directory), resolve("dist", directory), { recursive: true });
          }
        },
      },
      {
        name: "geekai-admin-preview",
        configurePreviewServer(server) {
          server.middlewares.use((req, _res, next) => {
            const url = new URL(req.url || "/", "http://localhost");
            const account = /^\/(?:login|register)\/?$/.test(url.pathname);
            const admin = /^\/admin(?:\/|$)/.test(url.pathname) && !extname(url.pathname);
            if (account || admin) req.url = "/admin/index.html" + url.search;
            next();
          });
        },
      },
    ],
    server: {
      host: "0.0.0.0",
      port: 8888,
      strictPort: true,
      proxy: {
        ...apiProxy,
        "/admin": { target: adminTarget, changeOrigin: true, ws: true },
        "^/(login|register)(/|\\?|$)": {
          target: adminTarget,
          changeOrigin: true,
          rewrite: (path) => "/admin/index.html" + (path.includes("?") ? path.slice(path.indexOf("?")) : ""),
        },
        "^/(images|medias)/": {
          target: adminTarget,
          changeOrigin: true,
          rewrite: (path) => "/admin" + path,
        },
      },
    },
    preview: { host: "0.0.0.0", proxy: apiProxy },
  };
});

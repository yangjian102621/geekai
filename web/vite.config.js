import vue from '@vitejs/plugin-vue'
import path from 'path'
import AutoImport from 'unplugin-auto-import/vite'
import { defineConfig, loadEnv } from 'vite'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd(), '')
  const apiHost = env.VITE_API_HOST || 'http://localhost:5000'
  return {
    plugins: [
      vue(),
      AutoImport({
        imports: [
          'vue',
          'vue-router',
          {
            '@vueuse/core': ['useMouse', 'useFetch'],
          },
        ],
        dts: true, // 生成 TypeScript 声明文件
      }),
    ],
    // 门户站点与后台合并为同一个 SPA：站点路由挂在 /，后台挂在 /admin/*。
    // 静态资源统一从根路径输出（/assets/...），nginx 只需按路径回退到 index.html。
    base: '/',
    build: {
      outDir: 'dist', // 构建输出目录
      target: 'esnext', // 允许使用顶层 await 等最新特性
    },

    // 开发环境与依赖预构建同样对齐到 esnext，避免 esbuild 报错
    esbuild: {
      target: 'esnext',
      supported: {
        'top-level-await': true,
      },
    },
    optimizeDeps: {
      esbuildOptions: {
        target: 'esnext',
        supported: {
          'top-level-await': true,
        },
      },
    },

    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
        '~@': path.resolve(__dirname, './src'),
      },
    },

    css: {
      preprocessorOptions: {
        scss: {
          api: 'modern-compiler', // 使用新的 Sass API，避免 legacy JS API 警告
        },
      },
    },

    server: {
      host: true,
      port: 3000, // 前端开发端口
      strictPort: true,
      open: false, // 可选：启动服务器时自动打开浏览器
      ...(process.env.NODE_ENV === 'development'
        ? {
            allowedHosts: ['localhost', '127.0.0.1', 'sapi.geekai.me'],
            proxy: {
              '/api': {
                target: apiHost,
                changeOrigin: true,
                ws: true,
              },
              '/v1/events': {
                target: apiHost,
                changeOrigin: true,
              },
              '/static/upload/': {
                target: apiHost,
                changeOrigin: true,
              },
              '/static/': {
                target: apiHost,
                changeOrigin: true,
              },
              '/healthz': {
                target: apiHost,
                changeOrigin: true,
              },
              '/update/': {
                target: apiHost,
                changeOrigin: true,
              },
            },
          }
        : {}),
    },
  }
})

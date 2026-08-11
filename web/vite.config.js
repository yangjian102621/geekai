import vue from '@vitejs/plugin-vue'
import path from 'path'
import AutoImport from 'unplugin-auto-import/vite'
import { defineConfig, loadEnv } from 'vite'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd(), '')
  const apiHost = env.VITE_API_HOST || 'http://localhost:5678'
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
    base: env.VITE_BASE_URL,
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
      port: 8888, // 设置你想要的端口号
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
              '/static/upload/': {
                target: apiHost,
                changeOrigin: true,
              },
            },
          }
        : {}),
    },
  }
})

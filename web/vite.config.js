import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import path from 'path';

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd());

  console.log('VITE_API_HOST from .env:', env);

  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, 'src'),
      },
    },
    css: {
      preprocessorOptions: {
        stylus: {
          paths: [path.resolve(__dirname, 'src')],
        },
      },
    },
    build: {
      sourcemap: false,
      minify: 'terser',
      rollupOptions: {
        output: {
          manualChunks: {
            vendor: ['vue', 'vue-router', 'pinia'],
          },
        },
      },
    },
    server: {
      port: 8888,
      host: true,
      proxy: {
        '/api': {
          target: env.VITE_APP_API_HOST,
          changeOrigin: true,
        },
        '/static/upload/': {
          target: env.VITE_APP_API_HOST,
          changeOrigin: true,
        },
      },
    },
    esbuild: {
      // 模拟 transpileDependencies: true
    },
  };
});

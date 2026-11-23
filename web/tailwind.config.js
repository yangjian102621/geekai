/* 门户站点设计系统色板 —— 暖灰 + 靛蓝：
   暖调浅灰画布、靛蓝主 CTA、暖近黑文字、暖炭黑深色面。
   颜色以 hex 定义以支持 Tailwind 的透明度修饰符（如 bg-canvas/80）。 */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        /* 表面（暖灰阶梯） */
        canvas: '#f7f6f3',
        'surface-soft': '#f0ede8',
        'surface-card': '#e7e3dd',
        'cream-strong': '#ded9d2',
        hairline: '#e1dcd6',
        'hairline-soft': '#eae5df',
        /* 文字（暖近黑） */
        ink: '#1a1815',
        body: '#474039',
        'body-strong': '#29251f',
        muted: '#6f6961',
        'muted-soft': '#948d84',
        /* 深色产品面 */
        dark: '#1a1815',
        'dark-elevated': '#29251f',
        'dark-soft': '#221f1b',
        'on-dark': '#f7f6f3',
        'on-dark-soft': '#a49d94',
        /* 强调：靛蓝只用于主 CTA 与关键强调 */
        primary: '#3b5bdb',
        'primary-hover': '#2f4bc0',
        'primary-disabled': '#e1dcd6',
        coral: '#3b5bdb',
        'coral-soft': '#e8ecfb',
        teal: '#2f7d8c',
        'teal-soft': '#e4edf1',
        amber: '#8b7bd8',
        'amber-soft': '#efebe5',
        /* brand-* 映射到靛蓝阶 */
        'brand-50': '#eef1fc',
        'brand-100': '#dfe4fa',
        'brand-200': '#c3cdf5',
        'brand-300': '#9fb0ee',
        'brand-400': '#7d93e6',
        'brand-500': '#3b5bdb',
        'brand-600': '#2f4bc0',
        'brand-700': '#273da0',
        'brand-800': '#21337f',
        'brand-900': '#1b2a67',
        /* 语义色 */
        success: '#3f9d63',
        warning: '#c4922a',
        error: '#c04a4a',
        'on-primary': '#ffffff',
      },
      fontFamily: {
        /* DESIGN.md：The Future 统一显示/正文，回退 Arial */
        serif:
          '"The Future", Arial, "PingFang SC", "Microsoft YaHei", "Noto Sans SC", sans-serif',
        sans:
          '"The Future", Arial, "PingFang SC", "Microsoft YaHei", "Noto Sans SC", system-ui, sans-serif',
        /* 技术标签/等宽：PP Neue Montreal Mono，回退 Georgia */
        mono: '"PP Neue Montreal Mono", Georgia, "SF Mono", Menlo, Consolas, monospace',
      },
      /* 阴影克制：只在悬浮态用极淡中性阴影 */
      boxShadow: {
        soft: '0 1px 3px rgba(26, 24, 21, 0.08)',
        lift: '0 4px 12px rgba(26, 24, 21, 0.1)',
      },
      /* 容器宽度按规范收敛到 1200px */
      maxWidth: {
        shell: '1200px',
      },
    },
  },
  plugins: [],
}

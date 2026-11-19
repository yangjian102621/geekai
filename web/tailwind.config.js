/* 门户站点（原 React web）设计系统色板 —— ThinkBuddy 暖纸与大地色方向。
   颜色以 hex 定义以支持 Tailwind 的透明度修饰符（如 bg-canvas/80）。 */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        /* 表面（暖中性阶梯） */
        canvas: '#faf9f5',
        'surface-soft': '#f5f2ea',
        'surface-card': '#f0ebe0',
        'cream-strong': '#e8e1d3',
        hairline: '#e3ddd1',
        'hairline-soft': '#ece7dc',
        /* 文字（暖近黑，不着色） */
        ink: '#141413',
        body: '#3d3d3a',
        'body-strong': '#252523',
        muted: '#6b6862',
        'muted-soft': '#757168',
        /* 深墨表面 */
        dark: '#191816',
        'dark-elevated': '#252320',
        'dark-soft': '#1f1e1b',
        'on-dark': '#faf9f5',
        'on-dark-soft': '#a09d96',
        /* 强调色（稀缺使用） */
        primary: '#8a5a2b',
        'primary-hover': '#75491f',
        'primary-disabled': '#e6dccf',
        coral: '#b7874e',
        'coral-soft': '#f6ecdd',
        teal: '#2f7365',
        'teal-soft': '#e5f1ee',
        amber: '#8a6015',
        'amber-soft': '#f9eed9',
        /* brand-* 映射到青铜琥珀阶 */
        'brand-50': '#f8f1e7',
        'brand-100': '#f0e3d1',
        'brand-200': '#e2cbb0',
        'brand-300': '#cda77d',
        'brand-400': '#b3814e',
        'brand-500': '#8a5a2b',
        'brand-600': '#75491f',
        'brand-700': '#5f3b19',
        'brand-800': '#4a2e13',
        'brand-900': '#38220e',
      },
      fontFamily: {
        serif:
          'Georgia, "Times New Roman", "Noto Serif SC", "Source Han Serif SC", "Songti SC", "STSong", "SimSun", serif',
        sans:
          '-apple-system, BlinkMacSystemFont, "Inter", "Segoe UI", Roboto, "PingFang SC", "Microsoft YaHei", "Noto Sans SC", system-ui, sans-serif',
        mono: '"JetBrains Mono", ui-monospace, SFMono-Regular, Menlo, Consolas, monospace',
      },
    },
  },
  plugins: [],
}

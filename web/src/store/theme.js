// src/store/index.js
import { defineStore } from 'pinia'

export const useThemeStore = defineStore('theme', {
  state: () => ({
    theme: localStorage.getItem('theme') || 'light', // 默认从 localStorage 获取主题
  }),
  actions: {
    setTheme(theme) {
      this.theme = theme
      document.documentElement.setAttribute('data-theme', theme)

      // 同时设置 dark 类，以便 Tailwind 的 dark: 前缀能够工作
      if (theme === 'dark') {
        document.documentElement.classList.add('dark')
      } else {
        document.documentElement.classList.remove('dark')
      }

      localStorage.setItem('theme', theme) // 保存到 localStorage
    },
  },
})

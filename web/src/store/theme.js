// src/store/index.js
import { defineStore } from "pinia";

export const useThemeStore = defineStore("theme", {
  state: () => ({
    theme: localStorage.getItem("theme") || "light" // 默认从 localStorage 获取主题
  }),
  actions: {
    setTheme(theme) {
      this.theme = theme;
      document.documentElement.setAttribute("data-theme", theme);
      localStorage.setItem("theme", theme); // 保存到 localStorage
    }
  }
});

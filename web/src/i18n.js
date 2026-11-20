import { createI18n } from 'vue-i18n'
import zhCN from './locales/zh-CN.json'
import en from './locales/en.json'

// 与原 React 站点（i18next-browser-languagedetector）保持同一 localStorage key，
// 已保存的语言偏好在迁移后依然生效。
const STORAGE_KEY = 'i18nextLng'

function detectLanguage() {
  try {
    const stored = window.localStorage.getItem(STORAGE_KEY)
    if (stored && ['zh-CN', 'en'].includes(stored)) return stored
  } catch (e) {
    /* ignore */
  }
  return navigator.language?.toLowerCase().startsWith('zh') ? 'zh-CN' : 'en'
}

const i18n = createI18n({
  legacy: false,
  locale: detectLanguage(),
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    en: en,
  },
  silentTranslationWarn: true,
  missingWarn: false,
  fallbackWarn: false,
})

export function toggleLanguage(locale) {
  const next = locale || (i18n.global.locale.value.startsWith('zh') ? 'en' : 'zh-CN')
  i18n.global.locale.value = next
  document.documentElement.setAttribute('lang', next)
  try {
    window.localStorage.setItem(STORAGE_KEY, next)
  } catch (e) {
    /* ignore */
  }
  return next
}

export default i18n

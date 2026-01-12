import { createI18n } from 'vue-i18n'
import zhCN from './locales/zh-CN'
import zhTW from './locales/zh-TW'
import en from './locales/en'

// Detect browser language
function getDefaultLocale(): string {
  const browserLang = navigator.language || (navigator as any).userLanguage
  
  if (browserLang.startsWith('zh')) {
    // Distinguish between Simplified and Traditional Chinese
    if (browserLang.includes('TW') || browserLang.includes('HK') || browserLang.includes('Hant')) {
      return 'zh-TW'
    }
    return 'zh-CN'
  }
  
  return 'en'
}

// Get saved locale or detect from browser
const savedLocale = localStorage.getItem('locale')
const defaultLocale = savedLocale || getDefaultLocale()

export const i18n = createI18n({
  legacy: false,
  locale: defaultLocale,
  fallbackLocale: 'en',
  messages: {
    'zh-CN': zhCN,
    'zh-TW': zhTW,
    'en': en
  }
})

// Helper to change locale
export function setLocale(locale: string) {
  ;(i18n.global.locale as any).value = locale
  localStorage.setItem('locale', locale)
  document.documentElement.setAttribute('lang', locale)
}

export const availableLocales = [
  { code: 'zh-CN', name: '简体中文', flag: '🇨🇳' },
  { code: 'en', name: 'English', flag: '🇺🇸' }
]

export default i18n

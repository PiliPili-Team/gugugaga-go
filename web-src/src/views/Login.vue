<script setup lang="ts">
/**
 * Login.vue - Premium Login Page
 * Glass morphism design with multiple wallpaper sources
 */

import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import { api } from '@/services/api'
import { setLocale, availableLocales } from '@/i18n'
import { Eye, EyeOff, Loader2, HardDrive, AlertCircle, Check, RefreshCw, Sun, Moon, Monitor } from 'lucide-vue-next'

const { t, locale } = useI18n()
const appStore = useAppStore()

// Appearance mode
const currentAppearance = ref(localStorage.getItem('appearance') || 'dark')

// Get current appearance icon component
const appearanceIcon = {
  light: Sun,
  dark: Moon,
  system: Monitor
}

// Cycle through appearance modes: dark -> light -> system -> dark
function cycleAppearance() {
  const modes = ['dark', 'light', 'system']
  const currentIndex = modes.indexOf(currentAppearance.value)
  const nextIndex = (currentIndex + 1) % modes.length
  currentAppearance.value = modes[nextIndex]
  localStorage.setItem('appearance', currentAppearance.value)
  applyAppearance(currentAppearance.value)
}

// Apply appearance mode
function applyAppearance(mode: string) {
  const root = document.documentElement
  let actualTheme = mode
  
  if (mode === 'system') {
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    actualTheme = prefersDark ? 'dark' : 'light'
  }
  
  root.setAttribute('data-theme', actualTheme)
  
  // Update PWA status bar color
  updateStatusBarColor(actualTheme)
}

// Update PWA status bar color based on theme
function updateStatusBarColor(theme: string) {
  // Remove media-specific theme-color meta tags and use a single one
  const existingMetas = document.querySelectorAll('meta[name="theme-color"]')
  existingMetas.forEach(meta => meta.remove())
  
  // Create new theme-color meta tag
  const themeColorMeta = document.createElement('meta')
  themeColorMeta.setAttribute('name', 'theme-color')
  document.head.appendChild(themeColorMeta)
  
  // Update apple-mobile-web-app-status-bar-style
  // iOS status bar styles:
  // - 'default': white background + black text/icons (for light mode)
  // - 'black': black background + white text/icons
  // - 'black-translucent': transparent background + white text/icons (content extends under status bar)
  let statusBarMeta = document.querySelector('meta[name="apple-mobile-web-app-status-bar-style"]')
  if (!statusBarMeta) {
    statusBarMeta = document.createElement('meta')
    statusBarMeta.setAttribute('name', 'apple-mobile-web-app-status-bar-style')
    document.head.appendChild(statusBarMeta)
  }
  
  if (theme === 'light') {
    // Light mode: white background, black text/icons in status bar
    themeColorMeta.setAttribute('content', '#f8fafc') // Match light mode background
    statusBarMeta.setAttribute('content', 'default') // iOS: white bg + black text/icons
  } else {
    // Dark mode: dark background, white text/icons
    themeColorMeta.setAttribute('content', '#0f0c18') // Match dark mode background
    statusBarMeta.setAttribute('content', 'black-translucent') // iOS: transparent bg + white text/icons
  }
}

// Language switcher
const showLanguageMenu = ref(false)

function changeLanguage(code: string) {
  setLocale(code)
  showLanguageMenu.value = false
}

// Form state
const username = ref('')
const password = ref('')
const showPassword = ref(false)
const isLoading = ref(false)
const error = ref('')

// Background - always fresh on login
const backgroundUrl = ref('')
const backgroundCopyright = ref('')
const backgroundSource = ref('')
const isLoadingBg = ref(false)
const lastWallpaperUrl = ref('')
const isFetchingWallpaper = ref(false) // 防止重复请求

// 多种壁纸源
const wallpaperSources = [
  'bing',
  'tmdb',
  'unsplash',
  'picsum'
]
let currentSourceIndex = 0

// Fetch wallpaper from different sources
async function fetchWallpaperFromSource(source: string): Promise<{ url: string; copyright: string; source: string }> {
  switch (source) {
    case 'bing':
      return await api.fetchBingWallpaper()
    case 'tmdb':
      return await api.fetchTMDBWallpaper()
    case 'unsplash':
      // Unsplash random nature/landscape photos
      const unsplashId = Date.now()
      return {
        url: `https://source.unsplash.com/1920x1080/?nature,landscape&sig=${unsplashId}`,
        copyright: 'Unsplash',
        source: 'unsplash'
      }
    case 'picsum':
      // Lorem Picsum random photos
      const picsumId = Math.floor(Math.random() * 1000)
      return {
        url: `https://picsum.photos/1920/1080?random=${picsumId}`,
        copyright: 'Lorem Picsum',
        source: 'picsum'
      }
    default:
      return { url: '', copyright: '', source: '' }
  }
}

// Fetch random wallpaper with retry and ensure different from last
async function fetchBackground(forceRefresh = false) {
  // 防止重复请求
  if (isFetchingWallpaper.value) return
  isFetchingWallpaper.value = true
  isLoadingBg.value = true
  
  try {
    let attempts = 0
    const maxAttempts = wallpaperSources.length * 2
    
    while (attempts < maxAttempts) {
      // 轮换壁纸源
      if (forceRefresh) {
        currentSourceIndex = (currentSourceIndex + 1) % wallpaperSources.length
      } else {
        currentSourceIndex = Math.floor(Math.random() * wallpaperSources.length)
      }
      
      const source = wallpaperSources[currentSourceIndex]
      const response = await fetchWallpaperFromSource(source)
      
      // 确保和上一张不同
      if (response.url && response.url !== lastWallpaperUrl.value) {
        backgroundUrl.value = response.url
        backgroundCopyright.value = response.copyright || ''
        backgroundSource.value = response.source || source
        lastWallpaperUrl.value = response.url
        break
      }
      
      attempts++
    }
  } catch (e) {
    // Silently fail - use fallback background
  } finally {
    isLoadingBg.value = false
    isFetchingWallpaper.value = false
  }
}

// 强制刷新壁纸
function refreshWallpaper() {
  fetchBackground(true)
}

// Handle login
async function handleLogin() {
  if (!username.value || !password.value) return
  
  isLoading.value = true
  error.value = ''
  
  try {
    const response = await api.login({
      username: username.value,
      password: password.value
    })
    
    if (response.success) {
      appStore.login(username.value)
    } else {
      error.value = t('login.error')
    }
  } catch (e: any) {
    error.value = e.message || t('login.error')
  } finally {
    isLoading.value = false
  }
}

// Handle Enter key
function handleKeyDown(e: KeyboardEvent) {
  if (e.key === 'Enter') {
    handleLogin()
  }
}

// Listen for system theme changes
let systemThemeMediaQuery: MediaQueryList | null = null

function handleSystemThemeChange(e: MediaQueryListEvent) {
  if (currentAppearance.value === 'system') {
    const newTheme = e.matches ? 'dark' : 'light'
    document.documentElement.setAttribute('data-theme', newTheme)
    updateStatusBarColor(newTheme)
  }
}

onMounted(() => {
  // Apply saved appearance mode
  applyAppearance(currentAppearance.value)
  
  // Listen for system theme changes
  systemThemeMediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  systemThemeMediaQuery.addEventListener('change', handleSystemThemeChange)
  
  // Always fetch fresh wallpaper on login page load
  fetchBackground()
})

onUnmounted(() => {
  // Remove system theme listener
  if (systemThemeMediaQuery) {
    systemThemeMediaQuery.removeEventListener('change', handleSystemThemeChange)
  }
})
</script>

<template>
  <div class="login-page">
    <!-- Background -->
    <div 
      class="login-background"
      :class="{ 'bg-loaded': backgroundUrl && !isLoadingBg }"
      :style="backgroundUrl ? { backgroundImage: `url(${backgroundUrl})` } : {}"
    >
      <div class="bg-overlay"></div>
      <!-- Loading spinner while fetching wallpaper -->
      <div v-if="isLoadingBg" class="bg-loading">
        <Loader2 :size="32" class="animate-spin" />
      </div>
    </div>
    
    <!-- Wallpaper Source Badge & Refresh Button -->
    <div class="wallpaper-controls">
      <button 
        class="wallpaper-refresh-btn" 
        @click="refreshWallpaper"
        :disabled="isLoadingBg"
        :title="t('login.refreshWallpaper')"
      >
        <RefreshCw :size="16" :class="{ 'animate-spin': isLoadingBg }" />
      </button>
      <div v-if="backgroundSource && !isLoadingBg" class="wallpaper-badge">
        <span class="badge-source">{{ backgroundSource.toUpperCase() }}</span>
        <span v-if="backgroundCopyright" class="badge-copyright">{{ backgroundCopyright }}</span>
      </div>
    </div>

    <!-- Language Switcher (Top Right) -->
    <div class="login-lang-switcher-top">
      <button 
        class="lang-btn-top"
        @click="showLanguageMenu = !showLanguageMenu"
      >
        <span class="lang-current-top">
          {{ availableLocales.find(l => l.code === locale)?.flag || '' }}
          {{ availableLocales.find(l => l.code === locale)?.name || locale }}
        </span>
      </button>
      
      <Transition name="dropdown">
        <div v-if="showLanguageMenu" class="lang-menu-top glass-card">
          <button
            v-for="lang in availableLocales"
            :key="lang.code"
            class="lang-option-top"
            :class="{ active: locale === lang.code }"
            @click="changeLanguage(lang.code)"
          >
            <span class="lang-flag-top">{{ lang.flag }}</span>
            <span class="lang-name-top">{{ lang.name }}</span>
            <Check v-if="locale === lang.code" :size="14" />
          </button>
        </div>
      </Transition>
    </div>

    <!-- Login Card -->
    <div class="login-container">
      <div class="login-card glass-card">
        <!-- Appearance Switcher (Top Right of Card) -->
        <button 
          class="appearance-btn-card"
          @click="cycleAppearance"
          :title="t('settings.' + (currentAppearance === 'light' ? 'themeLight' : currentAppearance === 'dark' ? 'themeDark' : 'themeSystem'))"
        >
          <component :is="appearanceIcon[currentAppearance as keyof typeof appearanceIcon]" :size="20" />
        </button>

        <!-- Logo -->
        <div class="login-logo">
          <div class="logo-icon">
            <HardDrive :size="32" />
          </div>
          <h1>GD Watcher</h1>
        </div>

        <!-- Header -->
        <div class="login-header">
          <h2>{{ t('login.title') }}</h2>
          <p>{{ t('login.subtitle') }}</p>
        </div>

        <!-- Error -->
        <Transition name="shake">
          <div v-if="error" class="login-error">
            <AlertCircle :size="16" />
            <span>{{ error }}</span>
          </div>
        </Transition>

        <!-- Form -->
        <form class="login-form" @submit.prevent="handleLogin">
          <div class="form-group">
            <label>{{ t('login.username') }}</label>
            <input
              type="text"
              v-model="username"
              class="input"
              :placeholder="t('login.usernamePlaceholder')"
              autocomplete="username"
              @keydown="handleKeyDown"
            />
          </div>

          <div class="form-group">
            <label>{{ t('login.password') }}</label>
            <div class="password-input">
              <input
                :type="showPassword ? 'text' : 'password'"
                v-model="password"
                class="input"
                :placeholder="t('login.passwordPlaceholder')"
                autocomplete="current-password"
                @keydown="handleKeyDown"
              />
              <button 
                type="button" 
                class="password-toggle"
                @click="showPassword = !showPassword"
              >
                <EyeOff v-if="showPassword" :size="18" />
                <Eye v-else :size="18" />
              </button>
            </div>
          </div>

          <button 
            type="submit" 
            class="login-btn"
            :disabled="isLoading || !username || !password"
          >
            <Loader2 v-if="isLoading" :size="18" class="animate-spin" />
            <span>{{ isLoading ? t('login.logging') : t('login.submit') }}</span>
          </button>
        </form>

        <!-- Footer -->
        <div class="login-footer">
          <p>{{ t('login.footer') }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.login-page {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: var(--space-4);
}

/* ========== Background ========== */
.login-background {
  position: fixed;
  inset: 0;
  background: var(--color-bg-primary);
  background-size: cover;
  background-position: center;
  z-index: -1;
  transition: background-image var(--duration-slow) var(--ease-default);
}

.bg-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    135deg,
    rgba(10, 15, 26, 0.85) 0%,
    rgba(10, 15, 26, 0.75) 50%,
    rgba(10, 15, 26, 0.85) 100%
  );
  backdrop-filter: blur(2px);
}

.login-background.bg-loaded {
  animation: fadeInBg 0.8s ease-out;
}

@keyframes fadeInBg {
  from {
    opacity: 0.5;
    transform: scale(1.05);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.bg-loading {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-quaternary);
  z-index: 1;
}

/* ========== Wallpaper Controls ========== */
.wallpaper-controls {
  position: fixed;
  bottom: calc(var(--space-4) + env(safe-area-inset-bottom));
  left: var(--space-4);
  display: flex;
  align-items: flex-end;
  gap: var(--space-2);
  z-index: 10;
}

.wallpaper-refresh-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-lg);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.wallpaper-refresh-btn:hover:not(:disabled) {
  background: rgba(0, 0, 0, 0.7);
  color: var(--color-text-primary);
  border-color: rgba(255, 255, 255, 0.2);
}

.wallpaper-refresh-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.wallpaper-badge {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-3);
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(8px);
  border-radius: var(--radius-lg);
  max-width: 300px;
}

.badge-source {
  font-size: var(--text-xs);
  font-weight: var(--font-bold);
  color: var(--color-accent);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.badge-copyright {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  line-height: var(--leading-snug);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

/* ========== Container ========== */
.login-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-4);
  width: 100%;
  max-width: 420px;
}

/* ========== Card ========== */
.login-card {
  position: relative;
  width: 100%;
  padding: var(--space-8);
  background: rgba(15, 22, 41, 0.65);
  backdrop-filter: blur(32px) saturate(180%);
  -webkit-backdrop-filter: blur(32px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: var(--radius-2xl);
  box-shadow: 
    0 0 0 1px rgba(255, 255, 255, 0.08) inset,
    0 8px 32px rgba(0, 0, 0, 0.4),
    0 2px 8px rgba(0, 0, 0, 0.2);
}

.login-card::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.05) 0%,
    rgba(255, 255, 255, 0.02) 50%,
    rgba(255, 255, 255, 0.05) 100%
  );
  border-radius: var(--radius-2xl);
  pointer-events: none;
  z-index: -1;
}

/* ========== Appearance Switcher (Top Right of Card) ========== */
.appearance-btn-card {
  position: absolute;
  top: var(--space-4);
  right: var(--space-4);
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  padding: 0;
  color: rgba(255, 255, 255, 0.7);
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
  z-index: 10;
}

.appearance-btn-card:hover {
  color: rgba(255, 255, 255, 1);
  background: rgba(255, 255, 255, 0.1);
  transform: scale(1.1);
}

.appearance-btn-card:active {
  transform: scale(0.95);
}

.appearance-btn-card svg {
  flex-shrink: 0;
}

/* ========== Language Switcher (Top Right) ========== */
.login-lang-switcher-top {
  position: fixed;
  /* Account for safe area (notch) on PWA */
  top: calc(env(safe-area-inset-top) + var(--space-4));
  right: var(--space-4);
  z-index: 100;
}

.lang-btn-top {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: rgba(255, 255, 255, 0.9);
  background: rgba(15, 22, 41, 0.8);
  backdrop-filter: blur(12px) saturate(180%);
  -webkit-backdrop-filter: blur(12px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.lang-btn-top:hover {
  background: rgba(15, 22, 41, 0.95);
  border-color: rgba(255, 255, 255, 0.25);
  transform: translateY(-1px);
}

.lang-current-top {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  white-space: nowrap;
}

.lang-menu-top {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 140px;
  padding: var(--space-2);
  background: rgba(15, 22, 41, 0.75);
  backdrop-filter: blur(24px) saturate(180%);
  -webkit-backdrop-filter: blur(24px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: var(--radius-lg);
  z-index: 100;
  box-shadow: 
    0 8px 32px rgba(0, 0, 0, 0.4),
    0 0 0 1px rgba(255, 255, 255, 0.05) inset;
}

.lang-option-top {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  width: 100%;
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-sm);
  color: rgba(255, 255, 255, 0.8);
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.lang-option-top:hover {
  color: rgba(255, 255, 255, 1);
  background: rgba(255, 255, 255, 0.1);
}

.lang-option-top.active {
  color: var(--color-accent);
  background: rgba(59, 130, 246, 0.15);
}

.lang-option-top svg {
  margin-left: auto;
  color: var(--color-accent);
  flex-shrink: 0;
}

.lang-flag-top {
  font-size: var(--text-lg);
  flex-shrink: 0;
}

.lang-name-top {
  flex: 1;
  text-align: left;
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: all var(--duration-fast) var(--ease-default);
  transform-origin: top right;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: scale(0.95) translateY(-4px);
}

/* ========== Logo ========== */
.login-logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-3);
  margin-bottom: var(--space-6);
}

.logo-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  background: var(--gradient-primary);
  border-radius: var(--radius-xl);
  color: white;
  box-shadow: var(--shadow-glow);
}

.login-logo h1 {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* ========== Header ========== */
.login-header {
  text-align: center;
  margin-bottom: var(--space-6);
}

.login-header h2 {
  font-size: var(--text-xl);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
  margin-bottom: var(--space-2);
}

.login-header p {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
}

/* ========== Error ========== */
.login-error {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-4);
  background: var(--color-error-light);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: var(--radius-lg);
  color: var(--color-error);
  font-size: var(--text-sm);
  margin-bottom: var(--space-4);
}

/* ========== Form ========== */
.login-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.form-group label {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
}

.login-form .input {
  height: 48px;
  background: var(--color-glass);
  border-color: var(--color-glass-border);
}

.login-form .input:focus {
  background: var(--color-glass-hover);
}

.password-input {
  position: relative;
}

.password-input .input {
  padding-right: var(--space-12);
}

.password-toggle {
  position: absolute;
  right: var(--space-3);
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-2);
  color: var(--color-text-tertiary);
  border-radius: var(--radius-md);
  transition: all var(--duration-fast) var(--ease-default);
}

.password-toggle:hover {
  color: var(--color-text-primary);
  background: var(--color-glass);
}

/* ========== Button ========== */
.login-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  height: 48px;
  margin-top: var(--space-2);
  font-size: var(--text-base);
  font-weight: var(--font-semibold);
  color: white;
  background: var(--gradient-primary);
  border: none;
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: var(--shadow-glow);
}

.login-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* ========== Footer ========== */
.login-footer {
  margin-top: var(--space-6);
  padding-top: var(--space-4);
  border-top: 1px solid var(--color-glass-border);
  text-align: center;
}

.login-footer p {
  font-size: var(--text-xs);
  color: var(--color-text-quaternary);
}

/* ========== Copyright ========== */
.bg-copyright {
  font-size: var(--text-xs);
  color: rgba(255, 255, 255, 0.4);
  text-align: center;
  max-width: 300px;
}

/* ========== Animations ========== */
.shake-enter-active {
  animation: shake 0.5s ease-out;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  10%, 30%, 50%, 70%, 90% { transform: translateX(-4px); }
  20%, 40%, 60%, 80% { transform: translateX(4px); }
}

/* ========== Responsive ========== */
@media (max-width: 480px) {
  .login-card {
    padding: var(--space-6);
  }
  
  .logo-icon {
    width: 48px;
    height: 48px;
  }
  
  .login-logo h1 {
    font-size: var(--text-xl);
  }
}
</style>

<!-- Light Theme Styles (non-scoped) -->
<style>
/* Light Theme - Appearance Switcher */
[data-theme="light"] .appearance-btn-card {
  color: rgba(0, 0, 0, 0.6) !important;
}

[data-theme="light"] .appearance-btn-card:hover {
  color: rgba(0, 0, 0, 0.9) !important;
  background: rgba(0, 0, 0, 0.08) !important;
}

/* Light Theme - Language Switcher */
[data-theme="light"] .lang-btn-top {
  background: rgba(255, 255, 255, 0.85) !important;
  border-color: rgba(0, 0, 0, 0.1) !important;
  color: #1e293b !important;
  backdrop-filter: blur(12px) saturate(180%) !important;
  -webkit-backdrop-filter: blur(12px) saturate(180%) !important;
}

[data-theme="light"] .lang-btn-top:hover {
  background: rgba(255, 255, 255, 0.95) !important;
  border-color: rgba(0, 0, 0, 0.15) !important;
  color: #0f172a !important;
}

[data-theme="light"] .lang-menu-top {
  background: rgba(255, 255, 255, 0.9) !important;
  border-color: rgba(0, 0, 0, 0.1) !important;
  backdrop-filter: blur(24px) saturate(180%) !important;
  -webkit-backdrop-filter: blur(24px) saturate(180%) !important;
  box-shadow: 
    0 8px 32px rgba(0, 0, 0, 0.12),
    0 0 0 1px rgba(0, 0, 0, 0.05) inset !important;
}

[data-theme="light"] .lang-option-top {
  color: #475569 !important;
}

[data-theme="light"] .lang-option-top:hover {
  color: #1e293b !important;
  background: rgba(0, 0, 0, 0.04) !important;
}

[data-theme="light"] .lang-option-top.active {
  color: var(--color-accent) !important;
  background: rgba(59, 130, 246, 0.1) !important;
}

/* Light Theme - Login Card */
[data-theme="light"] .login-card {
  background: rgba(255, 255, 255, 0.85) !important;
  border-color: rgba(0, 0, 0, 0.1) !important;
  backdrop-filter: blur(32px) saturate(180%) !important;
  -webkit-backdrop-filter: blur(32px) saturate(180%) !important;
  box-shadow: 
    0 0 0 1px rgba(0, 0, 0, 0.05) inset,
    0 8px 32px rgba(0, 0, 0, 0.1),
    0 2px 8px rgba(0, 0, 0, 0.05) !important;
}

[data-theme="light"] .login-card::before {
  background: linear-gradient(
    135deg,
    rgba(0, 0, 0, 0.02) 0%,
    rgba(0, 0, 0, 0.01) 50%,
    rgba(0, 0, 0, 0.02) 100%
  ) !important;
}

/* Light Theme - Background Overlay */
[data-theme="light"] .bg-overlay {
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.3) 0%,
    rgba(255, 255, 255, 0.2) 50%,
    rgba(255, 255, 255, 0.3) 100%
  ) !important;
}

/* Light Theme - Wallpaper Controls */
[data-theme="light"] .wallpaper-refresh-btn {
  background: rgba(255, 255, 255, 0.7) !important;
  border-color: rgba(0, 0, 0, 0.1) !important;
  color: #475569 !important;
  backdrop-filter: blur(8px) !important;
}

[data-theme="light"] .wallpaper-refresh-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.9) !important;
  color: #1e293b !important;
  border-color: rgba(0, 0, 0, 0.15) !important;
}

[data-theme="light"] .wallpaper-badge {
  background: rgba(255, 255, 255, 0.7) !important;
  backdrop-filter: blur(8px) !important;
}
</style>

<script setup lang="ts">
/**
 * Dashboard.vue - Main Dashboard Layout
 * Premium Sidebar Layout with Glass Effects
 */

import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore, useConfigStore, useLogsStore } from '@/stores'
import { setLocale, availableLocales } from '@/i18n'
import { api } from '@/services/api'

// Icons
import {
  LayoutDashboard,
  Server,
  Puzzle,
  Route,
  Eye,
  Settings,
  Key,
  FileText,
  Zap,
  ChevronsLeft,
  Menu,
  LogOut,
  Check,
  UserCircle,
  ChevronRight,
  ChevronLeft,
  Palette,
  MoreHorizontal,
  Sun,
  Moon,
  Monitor
} from 'lucide-vue-next'

// Panels
import DashboardPanel from './panels/DashboardPanel.vue'
import BasicPanel from './panels/BasicPanel.vue'
import IntegrationsPanel from './panels/IntegrationsPanel.vue'
import MappingsPanel from './panels/MappingsPanel.vue'
import TargetPanel from './panels/TargetPanel.vue'
import AdvancedPanel from './panels/AdvancedPanel.vue'
import OAuthPanel from './panels/OAuthPanel.vue'

// Business Components
import LogViewer from '@/components/business/LogViewer.vue'
import QuickActions from '@/components/business/QuickActions.vue'
import GoogleDriveIcon from '@/components/base/GoogleDriveIcon.vue'
import UserIcon from '@/components/base/UserIcon.vue'
import ColorPicker from '@/components/base/ColorPicker.vue'

const { t, locale } = useI18n()
const appStore = useAppStore()
const configStore = useConfigStore()
const logsStore = useLogsStore()

// State
const isSidebarCollapsed = ref(localStorage.getItem('sidebar_collapsed') === 'true')
const isMobileSidebarOpen = ref(false)
const activeTab = ref(localStorage.getItem('active_tab') || 'dashboard')
const showLanguageMenu = ref(false)
const showUserMenu = ref(false)
const showThemeMenu = ref(false)
const showAppearanceMenu = ref(false)
const showCustomThemePicker = ref(false)
const currentTheme = ref(localStorage.getItem('theme') || 'purple')
const currentAppearance = ref(localStorage.getItem('appearance') || 'dark')

// Available themes (color themes)
const availableThemes = [
  { id: 'purple', name: '紫色主题', preview: 'linear-gradient(135deg, #8b5cf6 0%, #3b82f6 100%)' },
  { id: 'blue', name: '蓝色主题', preview: 'linear-gradient(135deg, #3b82f6 0%, #06b6d4 100%)' },
  { id: 'green', name: '绿色主题', preview: 'linear-gradient(135deg, #22c55e 0%, #10b981 100%)' },
  { id: 'orange', name: '橙色主题', preview: 'linear-gradient(135deg, #f59e0b 0%, #ef4444 100%)' },
  { id: 'pink', name: '粉色主题', preview: 'linear-gradient(135deg, #ec4899 0%, #8b5cf6 100%)' }
]

// Available appearance modes (light/dark)
const appearanceModes = [
  { id: 'system', icon: Monitor, labelKey: 'settings.themeSystem' },
  { id: 'light', icon: Sun, labelKey: 'settings.themeLight' },
  { id: 'dark', icon: Moon, labelKey: 'settings.themeDark' }
]

// Close all menus
function closeAllMenus() {
  showLanguageMenu.value = false
  showUserMenu.value = false
  showThemeMenu.value = false
  showAppearanceMenu.value = false
  showCustomThemePicker.value = false
}

// Change theme
function changeTheme(themeId: string) {
  currentTheme.value = themeId
  localStorage.setItem('theme', themeId)
  applyTheme(themeId)
  showThemeMenu.value = false
}

// Apply color theme
function applyTheme(themeId: string) {
  const root = document.documentElement
  const themes: Record<string, any> = {
    purple: {
      '--color-accent': '#3b82f6',
      '--color-accent-hover': '#60a5fa',
      '--color-secondary': '#8b5cf6',
      '--color-accent-rgb': '59, 130, 246',
      '--color-secondary-rgb': '139, 92, 246'
    },
    blue: {
      '--color-accent': '#3b82f6',
      '--color-accent-hover': '#60a5fa',
      '--color-secondary': '#06b6d4',
      '--color-accent-rgb': '59, 130, 246',
      '--color-secondary-rgb': '6, 182, 212'
    },
    green: {
      '--color-accent': '#22c55e',
      '--color-accent-hover': '#34d399',
      '--color-secondary': '#10b981',
      '--color-accent-rgb': '34, 197, 94',
      '--color-secondary-rgb': '16, 185, 129'
    },
    orange: {
      '--color-accent': '#f59e0b',
      '--color-accent-hover': '#fbbf24',
      '--color-secondary': '#ef4444',
      '--color-accent-rgb': '245, 158, 11',
      '--color-secondary-rgb': '239, 68, 68'
    },
    pink: {
      '--color-accent': '#ec4899',
      '--color-accent-hover': '#f472b6',
      '--color-secondary': '#8b5cf6',
      '--color-accent-rgb': '236, 72, 153',
      '--color-secondary-rgb': '139, 92, 246'
    }
  }
  
  const theme = themes[themeId] || themes.purple
  Object.keys(theme).forEach(key => {
    root.style.setProperty(key, theme[key])
  })
}

// Change appearance mode (light/dark/system)
function changeAppearance(mode: string) {
  currentAppearance.value = mode
  localStorage.setItem('appearance', mode)
  applyAppearance(mode)
  showAppearanceMenu.value = false
}

// Apply appearance mode
function applyAppearance(mode: string) {
  const root = document.documentElement
  let actualTheme = mode
  
  if (mode === 'system') {
    // Detect system preference
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

// Listen for system theme changes
let systemThemeMediaQuery: MediaQueryList | null = null

function handleSystemThemeChange(e: MediaQueryListEvent) {
  if (currentAppearance.value === 'system') {
    const newTheme = e.matches ? 'dark' : 'light'
    document.documentElement.setAttribute('data-theme', newTheme)
    updateStatusBarColor(newTheme)
  }
}

// Handle custom theme apply
function handleCustomThemeApply(colors: { primary: string; secondary: string; opacity: number }) {
  const root = document.documentElement
  root.style.setProperty('--color-accent', colors.primary)
  root.style.setProperty('--color-accent-hover', adjustBrightness(colors.primary, 20))
  root.style.setProperty('--color-secondary', colors.secondary)
  currentTheme.value = 'custom'
  localStorage.setItem('theme', 'custom')
  localStorage.setItem('customTheme', JSON.stringify(colors))
}

// Helper to adjust brightness
function adjustBrightness(color: string, percent: number): string {
  const num = parseInt(color.replace('#', ''), 16)
  const r = Math.min(255, Math.max(0, (num >> 16) + percent * 2.55))
  const g = Math.min(255, Math.max(0, ((num >> 8) & 0x00FF) + percent * 2.55))
  const b = Math.min(255, Math.max(0, (num & 0x0000FF) + percent * 2.55))
  return `#${((r << 16) | (g << 8) | b).toString(16).padStart(6, '0')}`
}

// Listen for alert events to close menus
onMounted(() => {
  window.addEventListener('show-alert', closeAllMenus)
  window.addEventListener('show-confirm', closeAllMenus)
  
  // Apply saved theme and appearance
  applyTheme(currentTheme.value)
  applyAppearance(currentAppearance.value)
  
  // Listen for system theme changes
  systemThemeMediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  systemThemeMediaQuery.addEventListener('change', handleSystemThemeChange)
})

onUnmounted(() => {
  window.removeEventListener('show-alert', closeAllMenus)
  window.removeEventListener('show-confirm', closeAllMenus)
  
  // Remove system theme listener
  if (systemThemeMediaQuery) {
    systemThemeMediaQuery.removeEventListener('change', handleSystemThemeChange)
  }
})


// Navigation items
const navItems = computed(() => [
  { id: 'dashboard', icon: LayoutDashboard, label: t('nav.dashboard') },
  { id: 'basic', icon: Server, label: t('nav.basic') },
  { id: 'integrations', icon: Puzzle, label: t('nav.integrations') },
  { id: 'mappings', icon: Route, label: t('nav.mappings') },
  { id: 'target', icon: Eye, label: t('nav.target') },
  { id: 'advanced', icon: Settings, label: t('nav.advanced') },
  { id: 'oauth', icon: Key, label: t('nav.oauth') }
])

const utilityItems = computed(() => [
  { id: 'logs', icon: FileText, label: t('nav.logs') },
  { id: 'actions', icon: Zap, label: t('nav.actions') }
])

// Mobile TabBar items - 4 main items for bottom navigation
const mobileTabMainItems = computed(() => [
  { id: 'dashboard', icon: LayoutDashboard, label: t('nav.dashboard') },
  { id: 'basic', icon: Server, label: t('nav.basic') },
  { id: 'logs', icon: FileText, label: t('nav.logs') },
  { id: 'actions', icon: Zap, label: t('nav.actions') }
])



// Mobile TabBar more items
const mobileTabMoreItems = computed(() => [
  { id: 'integrations', icon: Puzzle, label: t('nav.integrations') },
  { id: 'mappings', icon: Route, label: t('nav.mappings') },
  { id: 'target', icon: Eye, label: t('nav.target') },
  { id: 'advanced', icon: Settings, label: t('nav.advanced') },
  { id: 'oauth', icon: Key, label: t('nav.oauth') }
])

// More menu state
const showMoreMenu = ref(false)

function selectFromMoreMenu(id: string) {
  selectTab(id)
  showMoreMenu.value = false
}

// Check if current tab is in more menu
const isMoreTabActive = computed(() => {
  return mobileTabMoreItems.value.some(item => item.id === activeTab.value)
})

// Detect PWA mode
const isPWA = ref(false)
onMounted(() => {
  // Check if running in standalone mode (PWA)
  if (window.matchMedia('(display-mode: standalone)').matches || 
      (window.navigator as any).standalone === true) {
    isPWA.value = true
  }
})

// Get current page title for mobile header
const currentPageTitle = computed(() => {
  const allItems = [...navItems.value, ...utilityItems.value]
  const currentItem = allItems.find(item => item.id === activeTab.value)
  return currentItem?.label || t('nav.dashboard')
})

// ...

// Current panel component
const currentPanel = computed(() => {
  const panels: Record<string, any> = {
    dashboard: DashboardPanel,
    basic: BasicPanel,
    integrations: IntegrationsPanel,
    mappings: MappingsPanel,
    target: TargetPanel, // Changed from ignore
    advanced: AdvancedPanel,
    oauth: OAuthPanel,
    logs: LogViewer,
    actions: QuickActions
  }
  return panels[activeTab.value] || DashboardPanel
})

// Current language info
const currentLanguageFlag = computed(() => {
  const current = availableLocales.find(l => l.code === locale.value)
  return current?.flag || '🌐'
})

const currentLanguageName = computed(() => {
  const current = availableLocales.find(l => l.code === locale.value)
  return current?.name || locale.value
})

// Handlers
function toggleSidebar() {
  if (window.innerWidth < 1024) {
    isMobileSidebarOpen.value = !isMobileSidebarOpen.value
  } else {
    isSidebarCollapsed.value = !isSidebarCollapsed.value
    localStorage.setItem('sidebar_collapsed', String(isSidebarCollapsed.value))
  }
}

function toggleUserMenu() {
  showUserMenu.value = !showUserMenu.value
  // Close other submenus when toggling main menu
  if (!showUserMenu.value) {
    showLanguageMenu.value = false
    showThemeMenu.value = false
  }
}

function selectTab(id: string) {
  activeTab.value = id
  localStorage.setItem('active_tab', id)
  // 关闭更多菜单（如果打开）
  if (showMoreMenu.value) {
    showMoreMenu.value = false
  }
  if (window.innerWidth < 1024) {
    isMobileSidebarOpen.value = false
  }
}

function handleLogout() {
  // Close all menus first
  showUserMenu.value = false
  showLanguageMenu.value = false
  showThemeMenu.value = false
  
  // Use custom alert instead of system confirm
  showConfirmDialog({
    title: t('common.confirm'),
    message: t('common.confirmLogout') || '确定要退出登录吗？',
    onConfirm: () => {
      appStore.logout()
    }
  })
}

// Custom alert/confirm dialog
const showAlertDialog = ref(false)
const showConfirmDialogRef = ref(false)
const dialogTitle = ref('')
const dialogMessage = ref('')
const dialogOnConfirm = ref<(() => void) | null>(null)

function showAlert(title: string, message: string) {
  dialogTitle.value = title
  dialogMessage.value = message
  showAlertDialog.value = true
}

function showConfirmDialog(options: { title: string; message: string; onConfirm: () => void }) {
  dialogTitle.value = options.title
  dialogMessage.value = options.message
  dialogOnConfirm.value = options.onConfirm
  showConfirmDialogRef.value = true
}

function handleDialogConfirm() {
  if (dialogOnConfirm.value) {
    dialogOnConfirm.value()
  }
  showConfirmDialogRef.value = false
  dialogOnConfirm.value = null
}

function handleDialogCancel() {
  showConfirmDialogRef.value = false
  window.dispatchEvent(new CustomEvent('confirm-result', { detail: false }))
  dialogOnConfirm.value = null
}

function handleAlertClose() {
  showAlertDialog.value = false
}

function changeLanguage(code: string) {
  setLocale(code)
  showLanguageMenu.value = false
  showUserMenu.value = false
  showThemeMenu.value = false
}

// Close user menu when clicking outside
function handleClickOutside(event: MouseEvent) {
  const target = event.target as HTMLElement
  // Check if click is outside user menu container and menu itself
  if (!target.closest('.user-menu-container') && !target.closest('.user-menu')) {
    showUserMenu.value = false
    showLanguageMenu.value = false
    showThemeMenu.value = false
  }
}

// Lifecycle
let logsInterval: number | null = null

onMounted(async () => {
  // Listen for clicks outside user menu
  document.addEventListener('click', handleClickOutside)
  
  // Listen for global dialog events
  window.addEventListener('show-alert', ((e: CustomEvent) => {
    // Close all menus when alert is shown
    showUserMenu.value = false
    showLanguageMenu.value = false
    showThemeMenu.value = false
    showAlert(e.detail.title, e.detail.message)
    if (e.detail.resolve) {
      // Wait for dialog to close
      const checkClosed = setInterval(() => {
        if (!showAlertDialog.value) {
          clearInterval(checkClosed)
          e.detail.resolve()
        }
      }, 100)
    }
  }) as EventListener)
  
  window.addEventListener('show-confirm', ((e: CustomEvent) => {
    // Close all menus when confirm is shown
    showUserMenu.value = false
    showLanguageMenu.value = false
    showThemeMenu.value = false
    showConfirmDialog({
      title: e.detail.title,
      message: e.detail.message,
      onConfirm: () => {
        window.dispatchEvent(new CustomEvent('confirm-result', { detail: true }))
      }
    })
    
    // Handle cancel
    const handleCancel = () => {
      window.dispatchEvent(new CustomEvent('confirm-result', { detail: false }))
    }
    
    // Listen for cancel
    window.addEventListener('confirm-cancel', handleCancel as EventListener, { once: true })
  }) as EventListener)
  
  // Fetch initial data
  await configStore.fetchConfig()
  
  await logsStore.fetchLogs()
  
  // Start polling logs
  logsInterval = window.setInterval(() => {
    logsStore.fetchLogs()
  }, 3000)
})

onUnmounted(() => {
  if (logsInterval) {
    clearInterval(logsInterval)
  }
  document.removeEventListener('click', handleClickOutside)
})
</script>

<template>
  <div class="dashboard" :class="{ 'sidebar-collapsed': isSidebarCollapsed }">
    <!-- Mobile Overlay -->
    <Transition name="fade">
      <div 
        v-if="isMobileSidebarOpen" 
        class="mobile-overlay"
        @click="isMobileSidebarOpen = false"
      />
    </Transition>

    <!-- Sidebar -->
    <aside 
      class="sidebar glass-card"
      :class="{ 'mobile-open': isMobileSidebarOpen }"
    >
      <!-- Logo -->
      <div class="sidebar-logo">
        <div class="logo-icon">
          <GoogleDriveIcon :size="24" />
        </div>
        <Transition name="fade">
          <span v-if="!isSidebarCollapsed" class="logo-text">GD Watcher</span>
        </Transition>
        <!-- Web端收起/展开按钮 - 仅在非移动端显示 -->
        <button 
          class="sidebar-toggle-btn desktop-only"
          @click="toggleSidebar"
          :title="isSidebarCollapsed ? t('sidebar.expand') : t('sidebar.collapse')"
        >
          <ChevronLeft v-if="!isSidebarCollapsed" :size="16" />
          <ChevronRight v-else :size="16" />
        </button>
      </div>

      <!-- Navigation -->
      <nav class="sidebar-nav">
        <div class="nav-section">
          <span v-if="!isSidebarCollapsed" class="nav-label">{{ t('nav.dashboard') }}</span>
          <div class="nav-items">
            <button
              v-for="item in navItems"
              :key="item.id"
              class="nav-item"
              :class="{ active: activeTab === item.id }"
              @click="selectTab(item.id)"
              :title="isSidebarCollapsed ? item.label : undefined"
            >
              <component :is="item.icon" :size="20" />
              <Transition name="fade">
                <span v-if="!isSidebarCollapsed">{{ item.label }}</span>
              </Transition>
              <span v-if="activeTab === item.id && isSidebarCollapsed" class="active-dot" />
            </button>
          </div>
        </div>

        <div class="nav-section">
          <span v-if="!isSidebarCollapsed" class="nav-label">{{ t('nav.actions') }}</span>
          <div class="nav-items">
            <button
              v-for="item in utilityItems"
              :key="item.id"
              class="nav-item"
              :class="{ active: activeTab === item.id }"
              @click="selectTab(item.id)"
              :title="isSidebarCollapsed ? item.label : undefined"
            >
              <component :is="item.icon" :size="20" />
              <Transition name="fade">
                <span v-if="!isSidebarCollapsed">{{ item.label }}</span>
              </Transition>
              <span v-if="activeTab === item.id && isSidebarCollapsed" class="active-dot" />
            </button>
          </div>
        </div>
      </nav>


    </aside>

    <!-- Mobile Top Navigation Bar -->
    <header class="mobile-top-nav">
      <div class="mobile-nav-content">
        <!-- Left spacer for symmetry -->
        <div class="nav-spacer"></div>
        
        <!-- Centered Title -->
        <div class="mobile-nav-title">
          <span class="nav-title-text">{{ currentPageTitle }}</span>
        </div>
        
        <!-- Right Actions -->
        <div class="mobile-nav-actions">
          <div class="user-menu-container mobile-user-menu">
            <button 
              class="mobile-user-menu-btn"
              @click.stop="toggleUserMenu"
            >
              <UserIcon :size="22" />
            </button>
            
            <Transition name="dropdown">
              <div v-if="showUserMenu" class="user-menu glass-card" @click.stop @mousedown.stop @touchstart.stop>
                <!-- User Info -->
                <div class="user-info">
                  <div class="user-avatar-large">
                    <UserIcon :size="40" />
                  </div>
                  <div class="user-details">
                    <div class="user-name">{{ appStore.userName || '管理员' }}</div>
                  </div>
                </div>
                
                <div class="user-menu-divider"></div>
                
                <!-- Language Switcher -->
                <button 
                  class="user-menu-item"
                  @click="showLanguageMenu = !showLanguageMenu"
                >
                  <span class="menu-item-icon">{{ currentLanguageFlag }}</span>
                  <span class="menu-item-label">{{ currentLanguageName }}</span>
                  <ChevronRight :size="16" class="menu-item-arrow" />
                </button>
                
                <!-- Language Menu -->
                <Transition name="slide">
                  <div v-if="showLanguageMenu" class="language-submenu">
                    <button
                      v-for="loc in availableLocales"
                      :key="loc.code"
                      class="lang-menu-item"
                      :class="{ active: locale === loc.code }"
                      @click="changeLanguage(loc.code)"
                    >
                      <span class="lang-flag">{{ loc.flag }}</span>
                      <span class="lang-name">{{ loc.name }}</span>
                      <Check v-if="locale === loc.code" :size="16" class="lang-check" />
                    </button>
                  </div>
                </Transition>
                
                <!-- Theme Switcher -->
                <button 
                  class="user-menu-item"
                  @click="showThemeMenu = !showThemeMenu"
                >
                  <Palette :size="18" class="menu-item-icon" />
                  <span class="menu-item-label">{{ t('header.theme') }}</span>
                  <ChevronRight :size="16" class="menu-item-arrow" />
                </button>
                
                <!-- Theme Menu -->
                <Transition name="slide">
                  <div v-if="showThemeMenu" class="theme-submenu">
                    <button
                      v-for="theme in availableThemes"
                      :key="theme.id"
                      class="theme-menu-item"
                      :class="{ active: currentTheme === theme.id }"
                      @click="changeTheme(theme.id)"
                    >
                      <div class="theme-preview" :style="{ background: theme.preview }"></div>
                      <span class="theme-name">{{ theme.name }}</span>
                      <Check v-if="currentTheme === theme.id" :size="16" class="theme-check" />
                    </button>
                    <button
                      class="theme-menu-item"
                      :class="{ active: currentTheme === 'custom' }"
                      @click="showCustomThemePicker = true"
                    >
                      <div class="theme-preview custom-theme-preview"></div>
                      <span class="theme-name">{{ t('header.customTheme') }}</span>
                      <Check v-if="currentTheme === 'custom'" :size="16" class="theme-check" />
                    </button>
                  </div>
                </Transition>
                
                <!-- Color Picker -->
                <ColorPicker 
                  v-model:show="showCustomThemePicker"
                  @apply="handleCustomThemeApply"
                />
                
                <!-- Appearance Mode -->
                <button 
                  class="user-menu-item"
                  @click="showAppearanceMenu = !showAppearanceMenu"
                >
                  <component :is="currentAppearance === 'light' ? Sun : currentAppearance === 'dark' ? Moon : Monitor" :size="18" class="menu-item-icon" />
                  <span class="menu-item-label">{{ t('header.appearance') }}</span>
                  <ChevronRight :size="16" class="menu-item-arrow" />
                </button>
                
                <!-- Appearance Menu -->
                <Transition name="slide">
                  <div v-if="showAppearanceMenu" class="appearance-submenu">
                    <button
                      v-for="mode in appearanceModes"
                      :key="mode.id"
                      class="appearance-menu-item"
                      :class="{ active: currentAppearance === mode.id }"
                      @click="changeAppearance(mode.id)"
                    >
                      <component :is="mode.icon" :size="18" class="appearance-icon" />
                      <span class="appearance-name">{{ t(mode.labelKey) }}</span>
                      <Check v-if="currentAppearance === mode.id" :size="16" class="appearance-check" />
                    </button>
                  </div>
                </Transition>
                
                <div class="user-menu-divider"></div>
                
                <!-- Logout -->
                <button 
                  class="user-menu-item logout-item"
                  @click="handleLogout"
                >
                  <LogOut :size="18" class="menu-item-icon" />
                  <span class="menu-item-label">{{ t('header.logout') }}</span>
                </button>
              </div>
            </Transition>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="main-content" :class="{ 'pwa-mode': isPWA }">
      <!-- Mobile Menu Button (only visible on mobile) -->
      <button class="mobile-menu-btn-fixed" @click="toggleSidebar">
        <Menu :size="20" />
      </button>

      <!-- User Menu (Desktop) -->
      <div class="user-menu-container desktop-only">
        <button 
          class="user-menu-btn"
          @click="showUserMenu = !showUserMenu"
        >
          <UserIcon :size="24" />
        </button>
        
        <Transition name="dropdown">
          <div v-if="showUserMenu" class="user-menu glass-card" @click.stop>
            <!-- User Info -->
            <div class="user-info">
              <div class="user-avatar-large">
                <UserIcon :size="40" />
              </div>
              <div class="user-details">
                <div class="user-name">{{ appStore.userName || '管理员' }}</div>
              </div>
            </div>
            
            <div class="user-menu-divider"></div>
            
            <!-- Language Option -->
            <button 
              class="user-menu-item"
              @click.stop="showLanguageMenu = !showLanguageMenu"
            >
              <span class="menu-item-icon">{{ availableLocales.find(l => l.code === locale)?.flag || '' }}</span>
              <span class="menu-item-label">
                {{ availableLocales.find(l => l.code === locale)?.name || locale }}
              </span>
              <ChevronRight :size="16" class="menu-item-arrow" />
            </button>
            
            <!-- Language Submenu -->
            <Transition name="slide">
              <div v-if="showLanguageMenu" class="user-submenu">
                <button
                  v-for="lang in availableLocales"
                  :key="lang.code"
                  class="user-menu-item submenu-item"
                  :class="{ active: locale === lang.code }"
                  @click.stop="changeLanguage(lang.code)"
                >
                  <span class="menu-item-icon">{{ lang.flag }}</span>
                  <span class="menu-item-label">{{ lang.name }}</span>
                  <Check v-if="locale === lang.code" :size="14" class="menu-item-check" />
                </button>
              </div>
            </Transition>
            
            <div class="user-menu-divider"></div>
            
            <!-- Theme Option -->
            <button 
              class="user-menu-item"
              @click.stop="showThemeMenu = !showThemeMenu"
            >
              <Palette :size="18" class="menu-item-icon" />
              <span class="menu-item-label">{{ t('header.theme') }}</span>
              <ChevronRight :size="16" class="menu-item-arrow" />
            </button>
            
            <!-- Theme Submenu -->
            <Transition name="slide">
              <div v-if="showThemeMenu" class="user-submenu">
                <button
                  v-for="theme in availableThemes"
                  :key="theme.id"
                  class="user-menu-item submenu-item"
                  :class="{ active: currentTheme === theme.id }"
                  @click.stop="changeTheme(theme.id)"
                >
                  <div class="theme-preview" :style="{ background: theme.preview }"></div>
                  <span class="menu-item-label">{{ theme.name }}</span>
                  <Check v-if="currentTheme === theme.id" :size="14" class="menu-item-check" />
                </button>
                <button
                  class="user-menu-item submenu-item"
                  @click.stop="showCustomThemePicker = true"
                >
                  <Palette :size="14" class="menu-item-icon" />
                  <span class="menu-item-label">{{ t('header.customTheme') }}</span>
                </button>
              </div>
            </Transition>
            
            <!-- Color Picker -->
            <ColorPicker 
              v-model:show="showCustomThemePicker"
              @apply="handleCustomThemeApply"
            />
            
            <!-- Appearance Option -->
            <button 
              class="user-menu-item"
              @click.stop="showAppearanceMenu = !showAppearanceMenu"
            >
              <component :is="currentAppearance === 'light' ? Sun : currentAppearance === 'dark' ? Moon : Monitor" :size="18" class="menu-item-icon" />
              <span class="menu-item-label">{{ t('header.appearance') }}</span>
              <ChevronRight :size="16" class="menu-item-arrow" />
            </button>
            
            <!-- Appearance Submenu -->
            <Transition name="slide">
              <div v-if="showAppearanceMenu" class="user-submenu">
                <button
                  v-for="mode in appearanceModes"
                  :key="mode.id"
                  class="user-menu-item submenu-item"
                  :class="{ active: currentAppearance === mode.id }"
                  @click.stop="changeAppearance(mode.id)"
                >
                  <component :is="mode.icon" :size="14" class="menu-item-icon" />
                  <span class="menu-item-label">{{ t(mode.labelKey) }}</span>
                  <Check v-if="currentAppearance === mode.id" :size="14" class="menu-item-check" />
                </button>
              </div>
            </Transition>
            
            <div class="user-menu-divider"></div>
            
            <!-- Logout -->
            <button 
              class="user-menu-item logout-item"
              @click="handleLogout"
            >
              <LogOut :size="18" class="menu-item-icon" />
              <span class="menu-item-label">{{ t('header.logout') }}</span>
            </button>
          </div>
        </Transition>
      </div>


      <!-- Content Area -->
      <div class="content-area" :class="{ 'pwa-mode': isPWA }">
        <Transition name="fade-slide" mode="out-in">
          <component :is="currentPanel" :key="activeTab" />
        </Transition>
      </div>
    </main>

    <!-- Mobile Bottom TabBar -->
    <nav class="mobile-tabbar">
      <button
        v-for="item in mobileTabMainItems"
        :key="item.id"
        class="tabbar-item"
        :class="{ active: activeTab === item.id && !showMoreMenu }"
        @click="selectTab(item.id)"
      >
        <component :is="item.icon" :size="18" />
        <span class="tabbar-label">{{ item.label }}</span>
      </button>
      
      <!-- More Button -->
      <div class="tabbar-more-wrapper">
        <button
          class="tabbar-item"
          :class="{ active: isMoreTabActive || showMoreMenu }"
          @click="showMoreMenu = !showMoreMenu"
        >
          <MoreHorizontal :size="20" />
          <span class="tabbar-label">{{ t('nav.more') }}</span>
        </button>
        
        <!-- More Menu Popup -->
        <Transition name="slide-up">
          <div v-if="showMoreMenu" class="more-menu glass-card">
            <button
              v-for="item in mobileTabMoreItems"
              :key="item.id"
              class="more-menu-item"
              :class="{ active: activeTab === item.id }"
              @click="selectFromMoreMenu(item.id)"
            >
              <component :is="item.icon" :size="18" />
              <span>{{ item.label }}</span>
            </button>
          </div>
        </Transition>
      </div>
    </nav>
    
    <!-- More Menu Overlay -->
    <Transition name="fade">
      <div 
        v-if="showMoreMenu" 
        class="more-menu-overlay"
        @click="showMoreMenu = false"
      />
    </Transition>

    <!-- Custom Alert Dialog -->
    <Transition name="dialog">
      <div v-if="showAlertDialog" class="dialog-overlay" @click="handleAlertClose">
        <div class="dialog-container" @click.stop>
          <h3 class="dialog-title">{{ dialogTitle }}</h3>
          <p class="dialog-message">{{ dialogMessage }}</p>
          <div class="dialog-actions">
            <button class="btn btn-primary" @click="handleAlertClose">{{ t('common.confirm') }}</button>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Custom Confirm Dialog -->
    <Transition name="dialog">
      <div v-if="showConfirmDialogRef" class="dialog-overlay" @click="handleDialogCancel">
        <div class="dialog-container" @click.stop>
          <h3 class="dialog-title">{{ dialogTitle }}</h3>
          <p class="dialog-message">{{ dialogMessage }}</p>
          <div class="dialog-actions">
            <button class="btn btn-secondary" @click="handleDialogCancel">{{ t('common.cancel') }}</button>
            <button class="btn btn-primary" @click="handleDialogConfirm">{{ t('common.confirm') }}</button>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.dashboard {
  position: relative;
  display: flex;
  min-height: 100vh;
  --sidebar-width: 260px;
  --sidebar-collapsed-width: 80px;
  --header-height: 64px;
}

.dashboard.sidebar-collapsed {
  --sidebar-width: var(--sidebar-collapsed-width);
}

/* ========== Mobile Overlay ========== */
.mobile-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  z-index: 40;
  display: none;
}

@media (max-width: 1023px) {
  .mobile-overlay {
    display: block;
  }
}

/* ========== Sidebar ========== */
.sidebar {
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  width: var(--sidebar-width);
  display: flex;
  flex-direction: column;
  padding: var(--space-4);
  background: rgba(34, 30, 47, 0.85);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-right: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 0;
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.05) inset;
  z-index: 50;
  transition: all var(--duration-normal) var(--ease-spring);
}

.sidebar:hover {
  border-color: var(--color-glass-border);
  box-shadow: none;
}

@media (max-width: 1023px) {
  .dashboard {
    --sidebar-width: 220px;
    --sidebar-collapsed-width: 220px;
  }
  
  .sidebar {
    transform: translateX(-100%);
    width: 220px !important;
  }
  
  .sidebar.mobile-open {
    transform: translateX(0);
  }
  
  /* 移动端始终显示完整侧边栏，不收起 */
  .dashboard.sidebar-collapsed .sidebar {
    width: 220px !important;
  }
  
  .dashboard.sidebar-collapsed .sidebar-logo {
    justify-content: flex-start;
    padding: var(--space-2) var(--space-3);
    gap: var(--space-3);
  }
  
  .dashboard.sidebar-collapsed .logo-text {
    display: block !important;
    opacity: 1 !important;
  }
  
  .dashboard.sidebar-collapsed .nav-item span {
    display: inline !important;
    opacity: 1 !important;
  }
  
  .dashboard.sidebar-collapsed .nav-label {
    display: block !important;
  }
  
  /* 移动端隐藏展开/收起按钮 */
  .desktop-only {
    display: none !important;
  }
}

/* ========== Logo ========== */
.sidebar-logo {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-2) var(--space-3);
  margin-bottom: var(--space-6);
  position: relative;
  min-height: 48px;
}

.dashboard.sidebar-collapsed .sidebar-logo {
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: var(--space-3) var(--space-2);
  gap: var(--space-2);
}

.dashboard.sidebar-collapsed .logo-icon {
  margin-right: 0;
}

/* 统一的展开/收起按钮样式 */
.sidebar-toggle-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  color: var(--color-text-tertiary);
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
  opacity: 0.6;
  flex-shrink: 0;
}

.sidebar-toggle-btn:hover {
  color: var(--color-text-primary);
  background: var(--color-glass);
  opacity: 1;
}

/* 展开状态时，按钮在右侧 */
.sidebar-logo .sidebar-toggle-btn {
  margin-left: auto;
}

/* 收起状态时，按钮在icon下方 */
.dashboard.sidebar-collapsed .sidebar-logo .sidebar-toggle-btn {
  margin-left: 0;
  margin-top: var(--space-1);
}

.logo-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: var(--radius-lg);
  flex-shrink: 0;
  overflow: hidden;
}

.logo-icon svg {
  width: 100%;
  height: 100%;
}

.logo-text {
  font-size: var(--text-lg);
  font-weight: var(--font-bold);
  background: var(--gradient-primary);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  white-space: nowrap;
}

/* ========== Navigation ========== */
.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
}

.nav-section {
  margin-bottom: var(--space-6);
}

.nav-label {
  display: block;
  padding: 0 var(--space-3);
  margin-bottom: var(--space-2);
  font-size: var(--text-xs);
  font-weight: var(--font-semibold);
  color: var(--color-text-quaternary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.nav-items {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.nav-item {
  position: relative;
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
  white-space: nowrap;
  overflow: hidden;
}

.nav-item:hover {
  color: var(--color-text-primary);
  background: var(--color-glass);
}

.nav-item.active {
  color: white;
  background: var(--gradient-primary);
  box-shadow: var(--shadow-glow);
}

.nav-item svg {
  flex-shrink: 0;
}

.active-dot {
  position: absolute;
  right: 8px;
  top: 50%;
  transform: translateY(-50%);
  width: 6px;
  height: 6px;
  background: white;
  border-radius: 50%;
}

.dashboard.sidebar-collapsed .nav-item {
  justify-content: center;
  padding: var(--space-3);
}

/* ========== Sidebar Footer ========== */
.sidebar-footer {
  padding-top: var(--space-4);
  border-top: 1px solid var(--color-glass-border);
}

.collapse-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 40px;
  color: var(--color-text-tertiary);
  background: var(--color-glass);
  border: none;
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.collapse-btn:hover {
  color: var(--color-text-primary);
  background: var(--color-glass-hover);
}

.collapse-btn svg {
  transition: transform var(--duration-normal) var(--ease-spring);
}

.collapse-btn svg.rotated {
  transform: rotate(180deg);
}

@media (max-width: 1023px) {
  .sidebar-footer {
    display: none;
  }
}

/* ========== User Menu ========== */
.user-menu-container {
  position: fixed;
  top: var(--space-4);
  right: var(--space-4);
  z-index: 100;
}

/* ========== Mobile Top Navigation Bar ========== */
.mobile-top-nav {
  display: none; /* Hidden by default, shown on mobile via media query */
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  /* Height = safe area padding + 56px content area */
  height: calc(env(safe-area-inset-top) + 56px);
  background: rgba(15, 12, 24, 0.5);
  backdrop-filter: blur(40px) saturate(180%);
  -webkit-backdrop-filter: blur(40px) saturate(180%);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 
    0 2px 16px rgba(0, 0, 0, 0.15),
    0 1px 4px rgba(0, 0, 0, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.08),
    inset 0 -1px 0 rgba(0, 0, 0, 0.05);
  z-index: 90;
  overflow: visible; /* Allow user menu to overflow */
}

.mobile-top-nav::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.06) 0%,
    rgba(255, 255, 255, 0.02) 50%,
    rgba(255, 255, 255, 0.06) 100%
  );
  pointer-events: none;
}

/* Mobile navigation content - positioned at bottom of nav bar */
.mobile-nav-content {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 var(--space-4);
  box-sizing: border-box;
  /* 确保所有子元素都垂直居中 */
  line-height: 0;
  /* DEBUG: 临时调试边框 - 取消注释来查看容器边界 */
  /* border: 1px solid red; */
}

.nav-spacer {
  width: 40px;
  height: 40px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.mobile-nav-title {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
  height: 56px;
  box-sizing: border-box;
  /* 确保与按钮使用完全相同的对齐方式 */
  line-height: 0;
  position: relative;
}

.nav-title-text {
  font-size: 17px;
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
  white-space: nowrap;
  text-align: center;
  /* 关键：使用精确的 line-height，确保文字垂直居中 */
  /* 16px 字体，使用 line-height: 1.0 让文字更紧凑，与图标对齐 */
  line-height: 1.0;
  display: block;
  margin: 0;
  padding: 0;
  /* 使用 flexbox 对齐，与按钮保持一致 */
  align-self: center;
  /* 确保文字在导航栏中垂直居中 */
  position: relative;
  /* 文字保持在导航栏中心，不需要额外调整 */
  top: 0;
}

.mobile-nav-actions {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  position: relative;
  z-index: 205;
  box-sizing: border-box;
}

/* Mobile user menu container */
.user-menu-container.mobile-user-menu {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

/* Mobile user menu dropdown position */
.user-menu-container.mobile-user-menu .user-menu {
  z-index: 350;
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
}

@media (max-width: 768px) {
  .mobile-top-nav {
    display: block;
  }
  
  .user-menu-container.desktop-only {
    display: none;
  }
}

/* Mobile user menu button */
.mobile-user-menu-btn {
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  width: 40px !important;
  height: 40px !important;
  min-width: 40px !important;
  min-height: 40px !important;
  max-width: 40px !important;
  max-height: 40px !important;
  background: transparent;
  border: none;
  margin: 0 !important;
  padding: 0 !important;
  line-height: 0 !important;
  font-size: 0 !important;
  outline: none;
  -webkit-appearance: none;
  cursor: pointer;
  transition: transform 0.2s ease;
  color: var(--color-text-primary);
  position: relative;
  z-index: 210;
  pointer-events: auto;
  -webkit-tap-highlight-color: transparent;
  touch-action: manipulation;
  box-sizing: border-box !important;
}

.mobile-user-menu-btn:active {
  transform: scale(0.9);
}

/* 使用 :deep() 确保覆盖 UserIcon 组件的 scoped 样式 */
.mobile-user-menu-btn :deep(svg),
.mobile-user-menu-btn svg {
  display: block !important;
  width: 22px !important;
  height: 22px !important;
  min-width: 22px !important;
  min-height: 22px !important;
  max-width: 22px !important;
  max-height: 22px !important;
  flex-shrink: 0 !important;
  margin: 0 !important;
  padding: 0 !important;
  border: none !important;
  outline: none !important;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3));
  /* 确保 SVG 在 flexbox 中完全居中 */
  align-self: center !important;
  
  /* ========== 图标垂直对齐修复 ========== */
  /* 
   * 根据实际测试，图标需要向上移动 15px 才能与文字垂直居中对齐。
   * 
   * 原因分析：
   * 1. SVG 图标的视觉重心不在几何中心：
   *    - viewBox: "0 0 24 24"
   *    - circle 在 cy="7"（顶部）
   *    - path 在 y="21"（底部）
   *    - 视觉重心偏下，导致图标在 flexbox 中看起来偏下
   * 
   * 2. 文字的 line-height: 1.0 让文字在容器中垂直居中
   * 
   * 3. 通过 transform 调整图标的视觉位置，使其与文字的视觉中心对齐
   * 
   * 如果后续发现对齐仍有偏差，可以微调这个值：
   * - 图标偏上：减少负值（如 -14px, -13px）
   * - 图标偏下：增加负值（如 -16px, -17px）
   */
  transform: translateY(-15px);
}

.user-menu-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  background: rgba(34, 30, 47, 0.85);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
  color: var(--color-text-primary);
}

.user-menu-btn:hover {
  background: rgba(42, 37, 56, 0.9);
  border-color: rgba(255, 255, 255, 0.25);
  transform: translateY(-1px);
}

.user-menu-btn svg {
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.4));
}

.user-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: var(--gradient-primary);
  border-radius: var(--radius-md);
  color: white;
}

.user-menu {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 280px;
  padding: var(--space-4);
  background: rgba(34, 30, 47, 0.95);
  backdrop-filter: blur(24px) saturate(180%);
  -webkit-backdrop-filter: blur(24px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: var(--radius-xl);
  z-index: 300;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
}

.user-info {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding-bottom: var(--space-4);
}

.user-avatar-large {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  background: transparent;
  border-radius: var(--radius-lg);
  color: var(--color-text-primary);
  flex-shrink: 0;
}

.user-avatar-large svg {
  filter: drop-shadow(0 2px 6px rgba(0, 0, 0, 0.4));
}

.user-details {
  flex: 1;
  min-width: 0;
}

.user-name {
  font-size: var(--text-base);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
  margin-bottom: var(--space-1);
}

.user-role {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
}

.user-menu-divider {
  height: 1px;
  background: rgba(255, 255, 255, 0.1);
  margin: var(--space-3) 0;
}

.user-menu-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  width: 100%;
  padding: var(--space-3);
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
  text-align: left;
  font-size: var(--text-sm);
}

.user-menu-item:hover {
  color: var(--color-text-primary);
  background: rgba(255, 255, 255, 0.08);
}

.user-menu-item.active {
  color: var(--color-accent);
  background: rgba(59, 130, 246, 0.15);
}

.user-menu-item.logout-item {
  color: var(--color-error);
}

.user-menu-item.logout-item:hover {
  background: rgba(239, 68, 68, 0.15);
}

.menu-item-icon {
  flex-shrink: 0;
  font-size: var(--text-lg);
}

.menu-item-label {
  flex: 1;
}

.menu-item-arrow {
  flex-shrink: 0;
  color: var(--color-text-tertiary);
}

.menu-item-check {
  flex-shrink: 0;
  color: var(--color-accent);
}

.user-submenu {
  margin-left: var(--space-4);
  margin-top: var(--space-3);
  margin-bottom: var(--space-3);
  padding-left: var(--space-4);
  border-left: 1px solid rgba(255, 255, 255, 0.1);
}

.user-submenu .user-menu-item {
  padding: var(--space-3) var(--space-4);
  margin-bottom: var(--space-1);
}

.user-submenu .user-menu-item:last-child {
  margin-bottom: 0;
}

/* ========== Language, Theme & Appearance Submenus (Mobile) ========== */
.language-submenu,
.theme-submenu,
.appearance-submenu {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  margin-top: var(--space-2);
  margin-bottom: var(--space-2);
  margin-left: var(--space-4);
  padding: var(--space-3);
  background: rgba(0, 0, 0, 0.2);
  border-radius: var(--radius-lg);
}

.lang-menu-item,
.theme-menu-item,
.appearance-menu-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  width: 100%;
  padding: var(--space-3) var(--space-4);
  color: var(--color-text-secondary);
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
  text-align: left;
  font-size: var(--text-sm);
}

.lang-menu-item:hover,
.theme-menu-item:hover,
.appearance-menu-item:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--color-text-primary);
}

.lang-menu-item.active,
.theme-menu-item.active,
.appearance-menu-item.active {
  background: rgba(59, 130, 246, 0.2);
  color: var(--color-accent);
}

.appearance-icon {
  flex-shrink: 0;
  color: inherit;
}

.appearance-name {
  flex: 1;
}

.appearance-check {
  flex-shrink: 0;
  color: var(--color-accent);
}

.lang-flag {
  font-size: var(--text-base);
  flex-shrink: 0;
}

.lang-name,
.theme-name {
  flex: 1;
}

.lang-check,
.theme-check {
  flex-shrink: 0;
  color: var(--color-accent);
}

.theme-preview {
  width: 20px;
  height: 20px;
  border-radius: var(--radius-md);
  flex-shrink: 0;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.slide-enter-active,
.slide-leave-active {
  transition: all var(--duration-fast) var(--ease-default);
  max-height: 200px;
  overflow: hidden;
}

.slide-enter-from,
.slide-leave-to {
  max-height: 0;
  opacity: 0;
}

/* ========== Main Content ========== */
.main-content {
  flex: 1;
  margin-left: var(--sidebar-width);
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  padding: var(--space-6);
  padding-top: calc(var(--space-4) + 44px + var(--space-4));
  transition: margin-left var(--duration-normal) var(--ease-spring);
}

@media (max-width: 1023px) {
  .main-content {
    margin-left: 0;
    padding: var(--space-4);
    padding-top: calc(var(--space-4) + 44px + var(--space-4));
  }
}

/* PWA mobile: account for top navigation bar height */
@media (max-width: 768px) {
  .main-content {
    /* Top nav height (safe-area + 56px) + spacing (12px) */
    padding-top: calc(env(safe-area-inset-top) + 56px + 12px) !important;
  }
}

/* ========== Mobile Menu Button ========== */
.mobile-menu-btn-fixed {
  position: fixed;
  top: var(--space-4);
  left: var(--space-4);
  z-index: 50;
  display: none;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  color: var(--color-text-primary);
  background: var(--color-glass);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--color-glass-border);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.mobile-menu-btn-fixed:hover {
  background: var(--color-glass-hover);
  border-color: var(--color-glass-border-hover);
}

@media (max-width: 1023px) {
  .mobile-menu-btn-fixed {
    display: flex;
  }
}


/* ========== Transitions ========== */
.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--duration-fast) var(--ease-default);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-fade-enter-active {
  transition: all var(--duration-normal) var(--ease-spring);
}

.slide-fade-leave-active {
  transition: all var(--duration-fast) var(--ease-default);
}

.slide-fade-enter-from {
  opacity: 0;
  transform: translateY(12px);
}

.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* ========== Fade Slide Transition (Page Switch) ========== */
.fade-slide-enter-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-slide-leave-active {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateX(8px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateX(-8px);
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

/* ========== Custom Dialog ========== */
.dialog-overlay {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  z-index: 1000;
}

.dialog-container {
  min-width: 400px;
  max-width: 90vw;
  padding: var(--space-6);
  background: rgba(15, 22, 41, 0.95);
  backdrop-filter: blur(24px) saturate(180%);
  -webkit-backdrop-filter: blur(24px) saturate(180%);
  border: 1px solid var(--color-glass-border);
  border-radius: var(--radius-2xl);
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
}

@media (max-width: 768px) {
  .dialog-container {
    min-width: auto;
    max-width: 320px;
    width: calc(100vw - 32px);
    margin: 0 var(--space-4);
    padding: var(--space-5);
  }
  
  .dialog-title {
    font-size: var(--text-lg);
  }
  
  .dialog-message {
    font-size: var(--text-sm);
    margin-bottom: var(--space-5);
  }
  
  .dialog-actions {
    gap: var(--space-2);
  }
  
  .dialog-actions .btn {
    flex: 1;
    min-width: 0;
  }
}

.dialog-title {
  font-size: var(--text-xl);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
  margin: 0 0 var(--space-4);
}

.dialog-message {
  font-size: var(--text-base);
  color: var(--color-text-secondary);
  margin: 0 0 var(--space-6);
  line-height: var(--leading-relaxed);
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
}

.dialog-enter-active,
.dialog-leave-active {
  transition: all var(--duration-normal) var(--ease-default);
}

.dialog-enter-from,
.dialog-leave-to {
  opacity: 0;
}

.dialog-container {
  animation: dialogSlideIn var(--duration-normal) var(--ease-spring);
}

@keyframes dialogSlideIn {
  from {
    opacity: 0;
    transform: scale(0.95) translateY(-10px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

/* ========== Mobile Bottom TabBar - Premium Glass Design ========== */
.mobile-tabbar {
  display: none;
  position: fixed;
  bottom: calc(var(--space-4) + env(safe-area-inset-bottom));
  left: var(--space-4);
  right: var(--space-4);
  height: 60px;
  background: rgba(15, 12, 24, 0.55);
  backdrop-filter: blur(40px) saturate(180%);
  -webkit-backdrop-filter: blur(40px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.12);
  /* 完美胶囊形：border-radius = height / 2 */
  border-radius: 30px;
  box-shadow: 
    0 8px 32px rgba(0, 0, 0, 0.3),
    0 2px 8px rgba(0, 0, 0, 0.2),
    inset 0 1px 0 rgba(255, 255, 255, 0.1),
    inset 0 -1px 0 rgba(0, 0, 0, 0.1);
  z-index: 200;
  overflow: visible;
}

.mobile-tabbar::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.05) 0%,
    rgba(255, 255, 255, 0.02) 50%,
    rgba(255, 255, 255, 0.05) 100%
  );
  pointer-events: none;
  /* 完美胶囊形：border-radius = height / 2 */
  border-radius: 30px;
}

@media (max-width: 768px) {
  .mobile-tabbar {
    display: flex;
    align-items: center;
    justify-content: space-around;
    padding: 7px;
    gap: 4px;
  }
  
  .main-content {
    padding-bottom: calc(76px + env(safe-area-inset-bottom));
  }
  
  .content-area {
    padding-bottom: var(--space-4);
  }
  
  /* Ensure all components have consistent top spacing */
  .content-area > :first-child {
    margin-top: 0;
  }
  
  /* 隐藏移动端顶部菜单按钮 */
  .mobile-menu-btn-fixed {
    display: none !important;
  }
  
  /* 隐藏侧边栏在移动端 - 使用 tabbar 替代 */
  .sidebar {
    display: none !important;
  }
  
  .mobile-overlay {
    display: none !important;
  }
}

.tabbar-item {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 3px;
  flex: 1;
  height: 48px;
  padding: 5px 10px;
  color: rgba(255, 255, 255, 0.55);
  background: transparent;
  border: none;
  /* 完美胶囊形：border-radius = height / 2 */
  border-radius: 24px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 1;
}

.tabbar-item svg {
  position: relative;
  z-index: 2;
  transition: color 0.3s ease, transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.tabbar-item::before {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: 16px;
  opacity: 0;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  transform: scale(0.95);
  z-index: 0;
}

.tabbar-item:active {
  transform: scale(0.95);
}

.tabbar-item.active {
  color: #fff;
}

.tabbar-item.active::before {
  opacity: 1;
  /* 使用主题渐变，跟随当前主题颜色 */
  background: var(--gradient-primary);
  /* 完美胶囊形：border-radius = height / 2 */
  border-radius: 24px;
  transform: scale(1);
  /* 使用主题颜色创建阴影效果，跟随当前主题 */
  /* 注意：使用 CSS 变量配合透明度，确保兼容性 */
  box-shadow: 
    0 4px 16px rgba(var(--color-accent-rgb, 59, 130, 246), 0.5),
    0 2px 8px rgba(var(--color-secondary-rgb, 139, 92, 246), 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.3),
    inset 0 -1px 0 rgba(0, 0, 0, 0.1);
}


.tabbar-item.active svg {
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3));
  color: #fff !important;
}

@keyframes iconBounce {
  0% {
    transform: scale(0.9);
  }
  50% {
    transform: scale(1.2);
  }
  100% {
    transform: scale(1.1);
  }
}

.tabbar-item.active .tabbar-label {
  color: #fff;
  font-weight: var(--font-semibold);
  animation: labelBounce 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.tabbar-label {
  font-size: 11px;
  font-weight: var(--font-medium);
  white-space: nowrap;
  letter-spacing: 0.01em;
  transition: color 0.3s ease;
  position: relative;
  z-index: 2;
}

@keyframes labelBounce {
  0% {
    transform: scale(0.95);
  }
  50% {
    transform: scale(1.15);
  }
  100% {
    transform: scale(1);
  }
}

/* ========== More Menu ========== */
.tabbar-more-wrapper {
  position: relative;
  flex: 1;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 204;
  isolation: isolate; /* Create new stacking context */
}

.more-menu-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  z-index: 199;
  pointer-events: auto;
}

.more-menu {
  position: absolute;
  bottom: calc(100% + 12px);
  right: -8px;
  min-width: 200px;
  padding: 8px;
  background: rgba(15, 12, 24, 0.55);
  backdrop-filter: blur(40px) saturate(180%);
  -webkit-backdrop-filter: blur(40px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 20px;
  box-shadow: 
    0 -8px 40px rgba(0, 0, 0, 0.4),
    0 -2px 8px rgba(0, 0, 0, 0.2),
    inset 0 1px 0 rgba(255, 255, 255, 0.1),
    inset 0 -1px 0 rgba(0, 0, 0, 0.1);
  z-index: 204;
  overflow: hidden;
  pointer-events: auto;
  isolation: isolate; /* Create new stacking context */
}

.more-menu::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.05) 0%,
    rgba(255, 255, 255, 0.02) 50%,
    rgba(255, 255, 255, 0.05) 100%
  );
  pointer-events: none;
  border-radius: 20px;
}

.more-menu-item {
  position: relative;
  display: flex;
  align-items: center;
  gap: var(--space-3);
  width: 100%;
  padding: var(--space-3) var(--space-4);
  color: rgba(255, 255, 255, 0.65);
  background: transparent;
  border: none;
  border-radius: 14px;
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: var(--text-sm);
  z-index: 1;
}

.more-menu-item:hover {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.95);
}

.more-menu-item.active {
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.3) 0%, rgba(99, 102, 241, 0.3) 100%);
  color: #c4b5fd;
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.15),
    inset 0 -1px 0 rgba(0, 0, 0, 0.05);
}

/* Slide up animation */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

</style>

<!-- Light Theme Global Styles (non-scoped) -->
<style>
/* ========== Light Theme Overrides ========== */
[data-theme="light"] .dashboard .sidebar {
  background: rgba(255, 255, 255, 0.92) !important;
  border-color: rgba(0, 0, 0, 0.1) !important;
}

[data-theme="light"] .dashboard .sidebar:hover {
  border-color: rgba(0, 0, 0, 0.15) !important;
}

[data-theme="light"] .dashboard .logo-text {
  color: #1e293b !important;
}

[data-theme="light"] .dashboard .nav-label {
  color: #64748b !important;
}

[data-theme="light"] .dashboard .nav-item {
  color: #475569 !important;
}

[data-theme="light"] .dashboard .nav-item:hover {
  background: rgba(0, 0, 0, 0.06) !important;
  color: #1e293b !important;
}

[data-theme="light"] .dashboard .nav-item.active {
  background: var(--gradient-primary) !important;
  color: #fff !important;
}

[data-theme="light"] .dashboard .sidebar-footer {
  border-color: rgba(0, 0, 0, 0.1) !important;
}

[data-theme="light"] .dashboard .collapse-btn {
  color: #475569 !important;
}

[data-theme="light"] .dashboard .collapse-btn:hover {
  background: rgba(0, 0, 0, 0.06) !important;
  color: #1e293b !important;
}

[data-theme="light"] .dashboard .sidebar-toggle-btn {
  color: #64748b !important;
}

[data-theme="light"] .dashboard .sidebar-toggle-btn:hover {
  color: #1e293b !important;
  background: rgba(0, 0, 0, 0.06) !important;
}

/* Light Theme - Mobile Top Nav */
[data-theme="light"] .mobile-top-nav {
  background: rgba(255, 255, 255, 0.92) !important;
  border-color: rgba(0, 0, 0, 0.1) !important;
}

[data-theme="light"] .mobile-top-nav::before {
  background: linear-gradient(
    135deg,
    rgba(0, 0, 0, 0.03) 0%,
    rgba(0, 0, 0, 0.01) 50%,
    rgba(0, 0, 0, 0.03) 100%
  ) !important;
}

[data-theme="light"] .nav-title-text {
  color: #1e293b !important;
}

[data-theme="light"] .mobile-user-menu-btn {
  color: #1e293b !important;
}

/* Light Theme - Mobile TabBar */
[data-theme="light"] .mobile-tabbar {
  background: rgba(255, 255, 255, 0.92) !important;
  border-color: rgba(0, 0, 0, 0.1) !important;
}

[data-theme="light"] .mobile-tabbar::before {
  background: linear-gradient(
    135deg,
    rgba(0, 0, 0, 0.03) 0%,
    rgba(0, 0, 0, 0.01) 50%,
    rgba(0, 0, 0, 0.03) 100%
  ) !important;
}

[data-theme="light"] .tabbar-item {
  color: #64748b !important;
}

[data-theme="light"] .tabbar-item:hover {
  color: #475569 !important;
}

[data-theme="light"] .tabbar-item.active {
  color: #fff !important;
}

[data-theme="light"] .tabbar-label {
  color: inherit !important;
}

/* Light Theme - More Menu */
[data-theme="light"] .more-menu {
  background: rgba(255, 255, 255, 0.95) !important;
  border-color: rgba(0, 0, 0, 0.1) !important;
}

[data-theme="light"] .more-menu::before {
  background: linear-gradient(
    135deg,
    rgba(0, 0, 0, 0.03) 0%,
    rgba(0, 0, 0, 0.01) 50%,
    rgba(0, 0, 0, 0.03) 100%
  ) !important;
}

[data-theme="light"] .more-menu-item {
  color: #475569 !important;
}

[data-theme="light"] .more-menu-item:hover {
  background: rgba(0, 0, 0, 0.06) !important;
  color: #1e293b !important;
}

[data-theme="light"] .more-menu-item.active {
  background: var(--gradient-primary) !important;
  color: #fff !important;
}

/* Light Theme - User Menu */
[data-theme="light"] .user-menu {
  background: rgba(255, 255, 255, 0.98) !important;
  border-color: rgba(0, 0, 0, 0.1) !important;
}

[data-theme="light"] .user-info {
  border-color: rgba(0, 0, 0, 0.08) !important;
}

[data-theme="light"] .user-name {
  color: #1e293b !important;
}

[data-theme="light"] .user-menu-item {
  color: #475569 !important;
}

[data-theme="light"] .user-menu-item:hover {
  background: rgba(0, 0, 0, 0.06) !important;
  color: #1e293b !important;
}

[data-theme="light"] .user-menu-item .menu-item-icon {
  color: #64748b !important;
}

[data-theme="light"] .user-menu-item:hover .menu-item-icon {
  color: #1e293b !important;
}

[data-theme="light"] .user-menu-item .menu-item-arrow {
  color: #94a3b8 !important;
}

[data-theme="light"] .user-menu-divider {
  background: rgba(0, 0, 0, 0.08) !important;
}

[data-theme="light"] .user-submenu {
  background: rgba(0, 0, 0, 0.04) !important;
}

[data-theme="light"] .language-submenu,
[data-theme="light"] .theme-submenu,
[data-theme="light"] .appearance-submenu {
  background: rgba(0, 0, 0, 0.04) !important;
}

[data-theme="light"] .lang-menu-item,
[data-theme="light"] .theme-menu-item,
[data-theme="light"] .appearance-menu-item {
  color: #475569 !important;
}

[data-theme="light"] .lang-menu-item:hover,
[data-theme="light"] .theme-menu-item:hover,
[data-theme="light"] .appearance-menu-item:hover {
  background: rgba(0, 0, 0, 0.06) !important;
  color: #1e293b !important;
}

[data-theme="light"] .lang-menu-item.active,
[data-theme="light"] .theme-menu-item.active,
[data-theme="light"] .appearance-menu-item.active {
  background: rgba(59, 130, 246, 0.12) !important;
  color: var(--color-accent) !important;
}

/* Light Theme - Main Content Area */
[data-theme="light"] .main-content {
  background: linear-gradient(180deg, #f8fafc 0%, #f1f5f9 100%) !important;
}

[data-theme="light"] .content-area {
  background: transparent !important;
}

/* Light Theme - Glass Card */
[data-theme="light"] .glass-card {
  background: rgba(255, 255, 255, 0.9) !important;
  border-color: rgba(0, 0, 0, 0.08) !important;
}

/* Light Theme - User Menu Button (Desktop) */
[data-theme="light"] .user-menu-btn {
  background: rgba(255, 255, 255, 0.85) !important;
  border-color: rgba(0, 0, 0, 0.1) !important;
  color: #1e293b !important;
  backdrop-filter: blur(20px) saturate(180%) !important;
  -webkit-backdrop-filter: blur(20px) saturate(180%) !important;
}

[data-theme="light"] .user-menu-btn:hover {
  background: rgba(255, 255, 255, 0.95) !important;
  border-color: rgba(0, 0, 0, 0.15) !important;
}

[data-theme="light"] .user-menu-btn svg {
  filter: none !important;
  color: #475569 !important;
}

/* Light Theme - User Avatar */
[data-theme="light"] .user-avatar {
  color: #475569 !important;
}

/* Light Theme - Glassmorphism Effects */

/* Sidebar with Glass Effect */
[data-theme="light"] .dashboard .sidebar {
  background: rgba(255, 255, 255, 0.75) !important;
  backdrop-filter: blur(24px) saturate(180%) !important;
  -webkit-backdrop-filter: blur(24px) saturate(180%) !important;
  border-color: rgba(0, 0, 0, 0.08) !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08) !important;
}

/* Mobile Top Nav with Glass Effect */
[data-theme="light"] .mobile-top-nav {
  background: rgba(255, 255, 255, 0.75) !important;
  backdrop-filter: blur(24px) saturate(180%) !important;
  -webkit-backdrop-filter: blur(24px) saturate(180%) !important;
  border-color: rgba(0, 0, 0, 0.08) !important;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.06) !important;
}

/* Mobile TabBar with Glass Effect */
[data-theme="light"] .mobile-tabbar {
  background: rgba(255, 255, 255, 0.75) !important;
  backdrop-filter: blur(24px) saturate(180%) !important;
  -webkit-backdrop-filter: blur(24px) saturate(180%) !important;
  border-color: rgba(0, 0, 0, 0.08) !important;
  box-shadow: 0 -4px 24px rgba(0, 0, 0, 0.06) !important;
}

/* User Menu with Glass Effect */
[data-theme="light"] .user-menu {
  background: rgba(255, 255, 255, 0.85) !important;
  backdrop-filter: blur(24px) saturate(180%) !important;
  -webkit-backdrop-filter: blur(24px) saturate(180%) !important;
  border-color: rgba(0, 0, 0, 0.1) !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12) !important;
}

/* More Menu with Glass Effect */
[data-theme="light"] .more-menu {
  background: rgba(255, 255, 255, 0.85) !important;
  backdrop-filter: blur(24px) saturate(180%) !important;
  -webkit-backdrop-filter: blur(24px) saturate(180%) !important;
  border-color: rgba(0, 0, 0, 0.1) !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12) !important;
}

/* Mobile User Menu Button */
[data-theme="light"] .mobile-user-menu-btn {
  color: #1e293b !important;
}

[data-theme="light"] .mobile-user-menu-btn svg {
  color: #475569 !important;
}

/* Light Theme - Dialog */
[data-theme="light"] .dialog-overlay {
  background: rgba(0, 0, 0, 0.3) !important;
}

[data-theme="light"] .dialog-container {
  background: rgba(255, 255, 255, 0.95) !important;
  border: 1px solid rgba(0, 0, 0, 0.1) !important;
  backdrop-filter: blur(24px) saturate(180%) !important;
  -webkit-backdrop-filter: blur(24px) saturate(180%) !important;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.15) !important;
}

[data-theme="light"] .dialog-title {
  color: #1e293b !important;
}

[data-theme="light"] .dialog-message {
  color: #475569 !important;
}
</style>

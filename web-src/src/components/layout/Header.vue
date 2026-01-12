<script setup lang="ts">
/**
 * Header - Top Navigation Bar
 */

import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore, useConfigStore, type TabId } from '@/stores'
import { Menu, Save, Loader2, Globe } from 'lucide-vue-next'
import { availableLocales, setLocale } from '@/i18n'
import { ref } from 'vue'

const { t, locale } = useI18n()
const appStore = useAppStore()
const configStore = useConfigStore()

// 当前页面标题
const pageTitles: Record<TabId, string> = {
  logs: 'nav.logs',
  basic: 'nav.basic',
  integrations: 'nav.integrations',
  mappings: 'nav.mappings',
  ignore: 'nav.ignore',
  advanced: 'nav.advanced',
  oauth: 'nav.oauth'
}

const pageTitle = computed(() => t(pageTitles[appStore.activeTab] || 'nav.basic'))

// Language menu
const showLangMenu = ref(false)

const currentLocaleName = computed(() => {
  return availableLocales.find(l => l.code === locale.value)?.name || locale.value
})

function changeLocale(code: string) {
  setLocale(code)
  showLangMenu.value = false
}

// 保存配置
async function handleSave() {
  const success = await configStore.saveConfig()
  if (!success) {
    alert(t('common.error') + ': ' + (configStore.error || ''))
  }
}
</script>

<template>
  <header class="header glass-card">
    <div class="header-left">
      <!-- Mobile menu button -->
      <button 
        class="menu-btn btn-ghost btn-icon btn-sm"
        @click="appStore.toggleSidebar"
      >
        <Menu :size="20" />
      </button>

      <h1 class="page-title">{{ pageTitle }}</h1>
    </div>

    <div class="header-right">
      <!-- Status -->
      <div class="status-badge" :class="{ connected: appStore.isConnected }">
        <span class="status-dot" />
        <span class="status-text">
          {{ appStore.isConnected ? t('header.status.running') : t('header.status.disconnected') }}
        </span>
      </div>

      <!-- Language -->
      <div class="lang-dropdown">
        <button class="lang-btn btn-ghost btn-sm" @click="showLangMenu = !showLangMenu">
          <Globe :size="18" />
          <span class="lang-name">{{ currentLocaleName }}</span>
        </button>
        <Transition name="dropdown">
          <div v-if="showLangMenu" class="lang-menu">
            <button 
              v-for="lang in availableLocales" 
              :key="lang.code"
              class="lang-option"
              :class="{ active: locale === lang.code }"
              @click="changeLocale(lang.code)"
            >
              {{ lang.name }}
            </button>
          </div>
        </Transition>
      </div>

      <!-- Save -->
      <button 
        class="save-btn btn btn-primary btn-sm"
        @click="handleSave"
        :disabled="configStore.isSaving"
      >
        <Loader2 v-if="configStore.isSaving" :size="16" class="animate-spin" />
        <Save v-else :size="16" />
        <span class="save-text">{{ configStore.isSaving ? t('header.saving') : t('header.save') }}</span>
      </button>
    </div>
  </header>
</template>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: var(--header-height);
  padding: 0 var(--space-6);
  border-radius: 0;
  border-left: none;
  border-right: none;
  border-top: none;
  position: sticky;
  top: 0;
  z-index: var(--z-sticky);
  background: var(--color-glass);
  backdrop-filter: blur(20px);
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.menu-btn {
  display: flex;
}

@media (min-width: 1024px) {
  .menu-btn {
    display: none;
  }
}

.page-title {
  font-size: var(--text-lg);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
}

.header-right {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

/* ========== Status Badge ========== */
.status-badge {
  display: none;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-1) var(--space-3);
  background: var(--color-glass);
  border-radius: var(--radius-full);
}

@media (min-width: 640px) {
  .status-badge {
    display: flex;
  }
}

.status-dot {
  width: 8px;
  height: 8px;
  background: var(--color-error);
  border-radius: 50%;
}

.status-badge.connected .status-dot {
  background: var(--color-success);
  animation: pulse 2s ease-in-out infinite;
}

.status-text {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
}

/* ========== Language Dropdown ========== */
.lang-dropdown {
  position: relative;
}

.lang-btn {
  gap: var(--space-2);
}

.lang-name {
  display: none;
}

@media (min-width: 768px) {
  .lang-name {
    display: inline;
  }
}

.lang-menu {
  position: absolute;
  top: calc(100% + var(--space-2));
  right: 0;
  min-width: 140px;
  padding: var(--space-2);
  background: var(--color-bg-elevated);
  border: 1px solid var(--color-glass-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xl);
  z-index: var(--z-dropdown);
}

.lang-option {
  display: block;
  width: 100%;
  padding: var(--space-2) var(--space-3);
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
  text-align: left;
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.lang-option:hover {
  background: var(--color-glass);
  color: var(--color-text-primary);
}

.lang-option.active {
  background: var(--color-accent-light);
  color: var(--color-accent);
}

/* ========== Save Button ========== */
.save-btn {
  gap: var(--space-2);
}

.save-text {
  display: none;
}

@media (min-width: 640px) {
  .save-text {
    display: inline;
  }
}

/* ========== Transitions ========== */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all var(--duration-fast) var(--ease-default);
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

@keyframes pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.7; transform: scale(1.1); }
}
</style>

<script setup lang="ts">
/**
 * Sidebar - Navigation Sidebar with Lucide Icons
 */

import { computed, inject } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore, type TabId } from '@/stores'
import { 
  Settings, 
  Link2, 
  FolderTree, 
  Ban, 
  Wrench, 
  KeyRound,
  FileText,
  ChevronLeft,
  LogOut
} from 'lucide-vue-next'

const { t } = useI18n()
const appStore = useAppStore()
const logout = inject<() => void>('logout')

// 导航项
const navItems = computed(() => [
  { id: 'basic' as TabId, label: t('nav.basic'), icon: Settings },
  { id: 'integrations' as TabId, label: t('nav.integrations'), icon: Link2 },
  { id: 'mappings' as TabId, label: t('nav.mappings'), icon: FolderTree },
  { id: 'ignore' as TabId, label: t('nav.ignore'), icon: Ban },
  { id: 'advanced' as TabId, label: t('nav.advanced'), icon: Wrench },
  { id: 'oauth' as TabId, label: t('nav.oauth'), icon: KeyRound },
])

// 移动端显示日志 Tab
const mobileNavItems = computed(() => [
  { id: 'logs' as TabId, label: t('nav.logs'), icon: FileText },
  ...navItems.value
])

const displayItems = computed(() => {
  return appStore.isMobile ? mobileNavItems.value : navItems.value
})

function handleLogout() {
  if (confirm(t('user.logoutConfirm'))) {
    // Clear session cookie
    document.cookie = 'gd_session=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;'
    logout?.()
  }
}
</script>

<template>
  <aside 
    class="sidebar"
    :class="{ 
      collapsed: appStore.sidebarCollapsed,
      'sidebar-mobile': appStore.isMobile
    }"
  >
    <!-- Header -->
    <div class="sidebar-header">
      <div class="brand">
        <div class="brand-icon">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <Transition name="fade">
          <span v-if="!appStore.sidebarCollapsed" class="brand-name">
            {{ t('app.name') }}
          </span>
        </Transition>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="sidebar-nav">
      <div class="nav-section">
        <Transition name="fade">
          <span v-if="!appStore.sidebarCollapsed" class="nav-label">
            {{ t('nav.menu') }}
          </span>
        </Transition>
        
        <div class="nav-items">
          <button
            v-for="item in displayItems"
            :key="item.id"
            class="nav-item"
            :class="{ active: appStore.activeTab === item.id }"
            :title="appStore.sidebarCollapsed ? item.label : undefined"
            @click="appStore.setActiveTab(item.id)"
          >
            <component :is="item.icon" :size="20" class="nav-icon" />
            <Transition name="fade">
              <span v-if="!appStore.sidebarCollapsed" class="nav-text">
                {{ item.label }}
              </span>
            </Transition>
            <span 
              v-if="appStore.sidebarCollapsed && appStore.activeTab === item.id" 
              class="active-indicator"
            />
          </button>
        </div>
      </div>
    </nav>

    <!-- Footer -->
    <div class="sidebar-footer">
      <!-- Collapse Toggle -->
      <button 
        v-if="!appStore.isMobile"
        class="collapse-btn"
        @click="appStore.toggleSidebarCollapse"
      >
        <ChevronLeft 
          :size="18" 
          :class="{ rotated: appStore.sidebarCollapsed }"
        />
      </button>

      <!-- User -->
      <div class="user-section" @click="handleLogout">
        <div class="user-avatar">
          <span>A</span>
          <span class="user-status" />
        </div>
        <Transition name="fade">
          <div v-if="!appStore.sidebarCollapsed" class="user-info">
            <span class="user-name">Admin</span>
            <span class="user-role">{{ t('user.admin') }}</span>
          </div>
        </Transition>
        <Transition name="fade">
          <LogOut 
            v-if="!appStore.sidebarCollapsed" 
            :size="16" 
            class="logout-icon"
          />
        </Transition>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  width: var(--sidebar-width);
  height: 100vh;
  background: var(--color-bg-secondary);
  border-right: 1px solid var(--color-glass-border);
  transition: width var(--duration-slow) var(--ease-default);
  flex-shrink: 0;
  position: relative;
  z-index: var(--z-sticky);
}

.sidebar.collapsed {
  width: var(--sidebar-collapsed);
}

.sidebar-mobile {
  position: fixed;
  left: 0;
  top: 0;
  z-index: var(--z-fixed);
  box-shadow: var(--shadow-xl);
}

/* ========== Header ========== */
.sidebar-header {
  display: flex;
  align-items: center;
  height: var(--header-height);
  padding: 0 var(--space-5);
  border-bottom: 1px solid var(--color-glass-border);
}

.brand {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.brand-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-accent);
  color: white;
  border-radius: var(--radius-xl);
  flex-shrink: 0;
}

.brand-name {
  font-size: var(--text-lg);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
  white-space: nowrap;
}

/* ========== Navigation ========== */
.sidebar-nav {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-6) var(--space-3);
}

.nav-section {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.nav-label {
  padding: 0 var(--space-3);
  margin-bottom: var(--space-2);
  font-size: var(--text-xs);
  font-weight: var(--font-semibold);
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--color-text-quaternary);
}

.nav-items {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  width: 100%;
  padding: var(--space-3);
  background: transparent;
  border: none;
  border-radius: var(--radius-lg);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
  position: relative;
}

.collapsed .nav-item {
  justify-content: center;
}

.nav-item:hover {
  background: var(--color-glass);
  color: var(--color-text-primary);
}

.nav-item.active {
  background: var(--color-accent);
  color: white;
  box-shadow: 0 4px 12px var(--color-accent-glow);
}

.nav-icon {
  flex-shrink: 0;
  transition: transform var(--duration-fast) var(--ease-default);
}

.nav-item:not(.active):hover .nav-icon {
  transform: scale(1.1);
}

.nav-text {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  white-space: nowrap;
}

.active-indicator {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 6px;
  height: 6px;
  background: white;
  border-radius: 50%;
  animation: pulse 2s ease-in-out infinite;
}

/* ========== Footer ========== */
.sidebar-footer {
  padding: var(--space-3);
  border-top: 1px solid var(--color-glass-border);
}

.collapse-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 36px;
  margin-bottom: var(--space-3);
  background: transparent;
  border: 1px solid var(--color-glass-border);
  border-radius: var(--radius-lg);
  color: var(--color-text-tertiary);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.collapse-btn:hover {
  background: var(--color-glass);
  color: var(--color-text-primary);
}

.collapse-btn svg {
  transition: transform var(--duration-normal) var(--ease-default);
}

.collapse-btn svg.rotated {
  transform: rotate(180deg);
}

.user-section {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.collapsed .user-section {
  justify-content: center;
}

.user-section:hover {
  background: var(--color-glass);
}

.user-avatar {
  position: relative;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-glass-active);
  color: var(--color-text-primary);
  font-weight: var(--font-semibold);
  border-radius: var(--radius-full);
  flex-shrink: 0;
}

.user-status {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 12px;
  height: 12px;
  background: var(--color-success);
  border: 2px solid var(--color-bg-secondary);
  border-radius: 50%;
}

.user-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
  flex: 1;
}

.user-name {
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
}

.user-role {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
}

.logout-icon {
  color: var(--color-text-tertiary);
  transition: color var(--duration-fast) var(--ease-default);
}

.user-section:hover .logout-icon {
  color: var(--color-error);
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

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}
</style>

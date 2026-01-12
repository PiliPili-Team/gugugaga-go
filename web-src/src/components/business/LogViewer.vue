<script setup lang="ts">
/**
 * LogViewer - Real-time Log Viewer
 */

import { ref, watch, nextTick, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useLogsStore } from '@/stores'
import { FileText, Trash2, ArrowUpDown, CheckSquare, ArrowDownUp } from 'lucide-vue-next'

const { t } = useI18n()
const logsStore = useLogsStore()
const logContainer = ref<HTMLElement | null>(null)
const isReversed = ref(true) // 默认逆序（最新在前）
const isPWA = ref(false)

onMounted(() => {
  // Check if running in standalone mode (PWA)
  if (window.matchMedia('(display-mode: standalone)').matches || 
      (window.navigator as any).standalone === true) {
    isPWA.value = true
  }
})

// 逆序后的日志列表
const displayLogs = computed(() => {
  const logs = logsStore.parsedLogs
  return isReversed.value ? [...logs].reverse() : logs
})

// 自动滚动 - 逆序时滚动到顶部，正序时滚动到底部
watch(
  () => logsStore.parsedLogs.length,
  async () => {
    if (logsStore.autoScroll && logContainer.value) {
      await nextTick()
      if (isReversed.value) {
        logContainer.value.scrollTop = 0
      } else {
        logContainer.value.scrollTop = logContainer.value.scrollHeight
      }
    }
  }
)

// 切换排序
function toggleOrder() {
  isReversed.value = !isReversed.value
}

// 清空日志
async function handleClearLogs() {
  await logsStore.clearMemLogs()
}
</script>

<template>
  <div class="panel">
    <div class="panel-content">
      <div class="log-viewer glass-card">
        <!-- Header -->
        <div class="log-header">
          <div class="log-title">
            <FileText :size="18" />
            <span>{{ t('logs.title') }}</span>
            <span class="log-count">{{ logsStore.count }}</span>
          </div>
          <div class="log-actions" :class="{ 'pwa-mode': isPWA }">
            <button 
              class="order-btn" 
              @click="toggleOrder" 
              :title="isReversed ? t('logs.newestFirst') : t('logs.oldestFirst')"
              :class="{ active: isReversed }"
            >
              <ArrowDownUp :size="16" />
              <span class="action-label">{{ isReversed ? t('logs.newestFirst') : t('logs.oldestFirst') }}</span>
            </button>
            <label class="auto-scroll">
              <input type="checkbox" v-model="logsStore.autoScroll" />
              <CheckSquare :size="16" class="auto-scroll-icon" />
              <span class="action-label">{{ t('logs.autoScroll') }}</span>
            </label>
            <button class="clear-btn" @click="handleClearLogs" :title="t('logs.clear')">
              <Trash2 :size="16" />
            </button>
          </div>
        </div>

        <!-- Content -->
        <div ref="logContainer" class="log-content">
          <TransitionGroup name="log-list">
            <div 
              v-for="log in displayLogs" 
              :key="log.id"
              class="log-line"
              :class="`log-${log.type}`"
            >
              <span class="log-message">{{ log.content }}</span>
            </div>
          </TransitionGroup>

          <!-- Empty state -->
          <div v-if="logsStore.count === 0" class="log-empty">
            <FileText :size="32" />
            <p>{{ t('logs.empty') }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Panel wrapper styles */
.panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
}

.panel-content {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-6);
}

@media (max-width: 768px) {
  .panel-content {
    padding: var(--space-3);
  }
}

.log-viewer {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-height: 300px;
  overflow: hidden;
}

/* ========== Header ========== */
.log-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--color-glass-border);
  flex-shrink: 0;
}

.log-title {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
}

.log-count {
  padding: 2px 8px;
  background: var(--color-glass-active);
  color: var(--color-text-secondary);
  font-size: var(--text-xs);
  font-weight: var(--font-bold);
  border-radius: var(--radius-full);
}

.log-actions {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.log-actions.pwa-mode .action-label {
  display: none;
}

.log-actions.pwa-mode .auto-scroll-icon {
  display: block;
}

.log-actions.pwa-mode .auto-scroll input[type="checkbox"] {
  display: none;
}

@media (max-width: 768px) {
  .log-actions .action-label {
    display: none;
  }
  
  .log-actions .auto-scroll-icon {
    display: block;
  }
  
  .log-actions .auto-scroll input[type="checkbox"] {
    display: none;
  }
}

.auto-scroll {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  cursor: pointer;
  user-select: none;
}

.auto-scroll-icon {
  display: none;
  width: 16px;
  height: 16px;
  color: var(--color-text-tertiary);
  transition: color var(--duration-fast) var(--ease-default);
}

.auto-scroll:hover .auto-scroll-icon {
  color: var(--color-text-secondary);
}

.auto-scroll input[type="checkbox"]:checked ~ .auto-scroll-icon {
  color: var(--color-accent);
}

.auto-scroll input[type="checkbox"] {
  position: relative;
  width: 18px;
  height: 18px;
  margin: 0;
  cursor: pointer;
  appearance: none;
  -webkit-appearance: none;
  background: var(--color-bg-secondary);
  border: 2px solid var(--color-glass-border);
  border-radius: 4px;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  flex-shrink: 0;
}

.auto-scroll input[type="checkbox"]:hover {
  border-color: var(--color-accent);
  background: var(--color-bg-tertiary);
  transform: scale(1.08);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.auto-scroll input[type="checkbox"]:checked {
  background: var(--color-accent);
  border-color: var(--color-accent);
  animation: checkboxPop 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.auto-scroll input[type="checkbox"]:checked::after {
  content: '';
  position: absolute;
  left: 5px;
  top: 2px;
  width: 5px;
  height: 9px;
  border: solid white;
  border-width: 0 2px 2px 0;
  transform: rotate(45deg);
  animation: checkmarkDraw 0.2s ease-out 0.1s both;
}

@keyframes checkboxPop {
  0% { transform: scale(1); }
  50% { transform: scale(1.15); }
  100% { transform: scale(1); }
}

@keyframes checkmarkDraw {
  0% { opacity: 0; transform: rotate(45deg) scale(0); }
  100% { opacity: 1; transform: rotate(45deg) scale(1); }
}

.order-btn {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-md);
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  background: transparent;
  border: 1px solid var(--color-glass-border);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.action-label {
  display: inline;
}

.order-btn:hover {
  color: var(--color-text-secondary);
  background: var(--color-glass);
  border-color: var(--color-glass-border-hover);
}

.order-btn.active {
  color: var(--color-accent);
  border-color: var(--color-accent);
  background: var(--color-accent-light);
}

.clear-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: var(--radius-md);
  color: var(--color-text-tertiary);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.order-btn.pwa-mode,
.log-actions.pwa-mode .order-btn {
  width: 32px;
  height: 32px;
  padding: 0;
  justify-content: center;
}

.log-actions.pwa-mode .auto-scroll {
  width: 32px;
  height: 32px;
  padding: 0;
  justify-content: center;
  border: 1px solid var(--color-glass-border);
  border-radius: var(--radius-md);
  background: transparent;
  transition: all var(--duration-fast) var(--ease-default);
}

.log-actions.pwa-mode .auto-scroll:hover {
  background: var(--color-glass);
  border-color: var(--color-glass-border-hover);
}

.log-actions.pwa-mode .auto-scroll input[type="checkbox"]:checked ~ .auto-scroll-icon {
  color: var(--color-accent);
}

@media (max-width: 768px) {
  .order-btn {
    width: 32px;
    height: 32px;
    padding: 0;
    justify-content: center;
  }
  
  .auto-scroll {
    width: 32px;
    height: 32px;
    padding: 0;
    justify-content: center;
    border: 1px solid var(--color-glass-border);
    border-radius: var(--radius-md);
    background: transparent;
    transition: all var(--duration-fast) var(--ease-default);
  }
  
  .auto-scroll:hover {
    background: var(--color-glass);
    border-color: var(--color-glass-border-hover);
  }
  
  .auto-scroll input[type="checkbox"]:checked ~ .auto-scroll-icon {
    color: var(--color-accent);
  }
}

.clear-btn:hover {
  color: var(--color-error);
  background: var(--color-error-light);
}

/* ========== Content ========== */
.log-content {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-3);
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  line-height: var(--leading-relaxed);
}

/* ========== Log Lines ========== */
.log-line {
  display: flex;
  gap: var(--space-2);
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-sm);
  transition: background var(--duration-fast) var(--ease-default);
}

.log-line:hover {
  background: var(--color-glass);
}

.log-message {
  flex: 1;
  word-break: break-all;
  color: var(--color-text-secondary);
}

/* Log types */
.log-info .log-message {
  color: var(--color-text-secondary);
}

.log-warn .log-message {
  color: var(--color-warning);
}

.log-error .log-message {
  color: var(--color-error);
}

.log-debug .log-message {
  color: var(--color-text-tertiary);
}

/* ========== Empty State ========== */
.log-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 200px;
  gap: var(--space-3);
  color: var(--color-text-quaternary);
}

.log-empty p {
  font-size: var(--text-sm);
  margin: 0;
}

/* ========== Transitions ========== */
.log-list-enter-active {
  transition: all var(--duration-fast) var(--ease-default);
}

.log-list-enter-from {
  opacity: 0;
  transform: translateX(-8px);
}
</style>

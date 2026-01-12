<script setup lang="ts">
/**
 * QuickActions - Quick Action Panel
 */

import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useLogsStore } from '@/stores'
import { api } from '@/services/api'
import { 
  Zap, 
  RefreshCw, 
  TreeDeciduous, 
  Eraser, 
  Trash2,
  Loader2,
  Send
} from 'lucide-vue-next'

const { t } = useI18n()
const logsStore = useLogsStore()

// 测试路径
const testPath = ref('/云端硬盘/Test/1.mkv')

// Loading states
const isTestingWebhook = ref(false)
const isRefreshingRclone = ref(false)
const isRefreshingTree = ref(false)
const isTriggeringSync = ref(false)

// 触发同步
async function handleTriggerSync() {
  isTriggeringSync.value = true
  try {
    await api.triggerSync()
  } catch (e) {
    // Silently handle error
  } finally {
    isTriggeringSync.value = false
  }
}

// 测试 Symedia 通知
async function handleTestWebhook() {
  isTestingWebhook.value = true
  try {
    await api.testSymedia({ path: testPath.value })
  } catch (e) {
    // Silently handle error
  } finally {
    isTestingWebhook.value = false
  }
}

// 强刷 Rclone
async function handleRcloneFull() {
  if (!confirm(t('actions.confirm.rcloneFull'))) return
  
  isRefreshingRclone.value = true
  try {
    await api.triggerRcloneFull()
  } catch (e) {
    // Silently handle error
  } finally {
    isRefreshingRclone.value = false
  }
}

// 刷新文件树
async function handleRefreshTree() {
  if (!confirm(t('actions.confirm.rebuildTree'))) return
  
  isRefreshingTree.value = true
  try {
    await api.refreshTree()
  } catch (e) {
    // Silently handle error
  } finally {
    isRefreshingTree.value = false
  }
}

// 清空日志
async function handleClearLogs() {
  await logsStore.clearMemLogs()
}

// 清空日志文件
async function handleClearFiles() {
  if (!confirm(t('actions.confirm.clearFiles'))) return
  
  try {
    await logsStore.clearFileLogs()
  } catch (e) {
    // Silently handle error
  }
}
</script>

<template>
  <div class="panel">
    <div class="panel-content">
      <div class="quick-actions glass-card">
        <div class="actions-header">
          <Zap :size="18" />
          <span>{{ t('actions.title') }}</span>
        </div>

        <!-- Test Path Input -->
        <div class="test-section">
          <label>{{ t('actions.testPath') }}</label>
          <div class="test-row">
            <input 
              type="text"
              v-model="testPath"
              :placeholder="t('actions.testPathPlaceholder')"
              class="input mono"
            />
            <button 
              class="btn btn-secondary btn-sm" 
              @click="handleTestWebhook"
              :disabled="isTestingWebhook"
            >
              <Loader2 v-if="isTestingWebhook" :size="14" class="animate-spin" />
              <Send v-else :size="14" />
              <span>{{ t('actions.test') }}</span>
            </button>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="actions-grid">
          <button 
            class="btn btn-success" 
            @click="handleTriggerSync"
            :disabled="isTriggeringSync"
          >
            <Loader2 v-if="isTriggeringSync" :size="16" class="animate-spin" />
            <Zap v-else :size="16" />
            <span>{{ isTriggeringSync ? t('actions.syncing') : t('actions.sync') }}</span>
          </button>

          <button 
            class="btn btn-warning" 
            @click="handleRcloneFull"
            :disabled="isRefreshingRclone"
          >
            <Loader2 v-if="isRefreshingRclone" :size="16" class="animate-spin" />
            <RefreshCw v-else :size="16" />
            <span>{{ isRefreshingRclone ? t('actions.rcloneFulling') : t('actions.rcloneFull') }}</span>
          </button>

          <button 
            class="btn btn-secondary" 
            @click="handleRefreshTree"
            :disabled="isRefreshingTree"
          >
            <Loader2 v-if="isRefreshingTree" :size="16" class="animate-spin" />
            <TreeDeciduous v-else :size="16" />
            <span>{{ isRefreshingTree ? t('actions.rebuildingTree') : t('actions.rebuildTree') }}</span>
          </button>

          <button class="btn btn-secondary" @click="handleClearLogs">
            <Eraser :size="16" />
            <span>{{ t('actions.clearPanel') }}</span>
          </button>

          <button class="btn btn-danger" @click="handleClearFiles">
            <Trash2 :size="16" />
            <span>{{ t('actions.clearFiles') }}</span>
          </button>
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

.quick-actions {
  padding: var(--space-4);
}

.actions-header {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
  margin-bottom: var(--space-4);
  padding-bottom: var(--space-3);
  border-bottom: 1px solid var(--color-glass-border);
}

/* ========== Test Section ========== */
.test-section {
  margin-bottom: var(--space-4);
}

.test-section label {
  display: block;
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  color: var(--color-text-tertiary);
  margin-bottom: var(--space-2);
}

.test-row {
  display: flex;
  gap: var(--space-2);
}

.test-row .input {
  flex: 1;
  height: 36px;
  font-size: var(--text-sm);
}

/* ========== Actions Grid ========== */
.actions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--space-2);
}

.actions-grid .btn {
  height: 40px;
  font-size: var(--text-xs);
  padding: 0 var(--space-3);
}

/* ========== Warning Button ========== */
.btn-warning {
  background: var(--color-warning-light);
  color: var(--color-warning);
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.btn-warning:hover:not(:disabled) {
  background: var(--color-warning);
  color: #fff;
  border-color: var(--color-warning);
}
</style>

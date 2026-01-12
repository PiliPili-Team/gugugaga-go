<script setup lang="ts">
/**
 * IntegrationsPanel - Service Integrations Configuration
 */

import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfigStore } from '@/stores'
import { HardDrive, Server, Bell, Plus, Trash2, Save, Loader2, ExternalLink } from 'lucide-vue-next'
import type { RcloneInstance } from '@/types'

const { t } = useI18n()
const configStore = useConfigStore()
const isSaving = ref(false)

async function handleSave() {
  isSaving.value = true
  try {
    await configStore.saveConfig()
  } finally {
    isSaving.value = false
  }
}

function updateConfig(path: string, value: any) {
  configStore.updateNested(path, value)
}

// Rclone instance management
function addRcloneInstance() {
  const current = configStore.config?.rclone?.instances || []
  updateConfig('rclone.instances', [
    ...current,
    { host: 'http://localhost:5572', endpoint: '/rc/vfs/refresh', wait_for_data: true }
  ])
}

function updateRcloneInstance(index: number, field: keyof RcloneInstance, value: any) {
  const current = configStore.config?.rclone?.instances || []
  const updated = [...current]
  updated[index] = { ...updated[index], [field]: value }
  updateConfig('rclone.instances', updated)
}

function removeRcloneInstance(index: number) {
  const current = configStore.config?.rclone?.instances || []
  updateConfig('rclone.instances', current.filter((_, i) => i !== index))
}

// Symedia headers management
function updateSymediaHeader(key: string, value: string) {
  const currentHeaders = configStore.config?.symedia?.headers || {}
  const updated = { ...currentHeaders, [key]: value }
  updateConfig('symedia.headers', updated)
}

function removeSymediaHeader(key: string) {
  const currentHeaders = configStore.config?.symedia?.headers || {}
  const updated = { ...currentHeaders }
  delete updated[key]
  updateConfig('symedia.headers', updated)
}

function addSymediaHeader() {
  const currentHeaders = configStore.config?.symedia?.headers || {}
  const newKey = `header-${Date.now()}`
  updateConfig('symedia.headers', { ...currentHeaders, [newKey]: '' })
}

function updateSymediaHeaderKey(oldKey: string, newKey: string) {
  const currentHeaders = configStore.config?.symedia?.headers || {}
  const updated = { ...currentHeaders }
  const oldValue = updated[oldKey]
  delete updated[oldKey]
  updated[newKey] = oldValue
  updateConfig('symedia.headers', updated)
}
</script>

<template>
  <div class="panel">
    <div class="panel-content">
      <!-- Google Drive Section -->
      <section class="config-section">
        <h3>
          <HardDrive :size="16" />
          {{ t('panels.integrations.googleDrive') }}
        </h3>

        <div class="form-grid">
          <div class="form-group">
            <label>{{ t('panels.integrations.qps') }}</label>
            <input
              type="number"
              class="input"
              :value="configStore.config?.google?.qps || 2"
              @input="updateConfig('google.qps', Number(($event.target as HTMLInputElement).value))"
              min="1"
              max="100"
            />
            <span class="hint">
              {{ t('panels.integrations.qpsHint') }}
              <a href="https://developers.google.com/drive/api/v3/performance" target="_blank" rel="noopener noreferrer">
                {{ t('panels.integrations.learnMore') }}
                <ExternalLink :size="12" />
              </a>
            </span>
          </div>

          <div class="form-group">
            <label>{{ t('panels.integrations.personalDriveName') }}</label>
            <input
              type="text"
              class="input"
              :value="configStore.config?.google?.personal_drive_name || ''"
              @input="updateConfig('google.personal_drive_name', ($event.target as HTMLInputElement).value)"
              :placeholder="t('panels.integrations.personalDriveNamePlaceholder')"
            />
          </div>
        </div>

        <!-- Save Button -->
        <div class="section-footer">
          <button 
            class="btn btn-primary"
            @click="handleSave"
            :disabled="isSaving"
          >
            <Loader2 v-if="isSaving" :size="18" class="animate-spin" />
            <Save v-else :size="18" />
            <span>{{ isSaving ? t('common.saving') : t('common.save') }}</span>
          </button>
        </div>
      </section>

      <!-- Rclone Section -->
      <section class="config-section">
        <div class="section-header">
          <h3>
            <Server :size="16" />
            {{ t('panels.integrations.rclone') }}
          </h3>
          <button class="add-btn" @click="addRcloneInstance">
            <Plus :size="14" />
            <span>{{ t('common.add') }}</span>
          </button>
        </div>

        <div class="instances-list">
          <div
            v-for="(instance, index) in configStore.config?.rclone?.instances || []"
            :key="index"
            class="instance-card"
          >
            <div class="instance-header">
              <span class="instance-index">#{{ index + 1 }}</span>
              <button class="remove-btn" @click="removeRcloneInstance(index)">
                <Trash2 :size="14" />
              </button>
            </div>

            <div class="form-grid compact">
              <div class="form-group">
                <label>{{ t('panels.integrations.rcloneHost') }}</label>
                <input
                  type="text"
                  class="input mono"
                  :value="instance.host"
                  @input="updateRcloneInstance(index, 'host', ($event.target as HTMLInputElement).value)"
                  placeholder="http://localhost:5572"
                />
              </div>

              <div class="form-group">
                <label>{{ t('panels.integrations.rcloneEndpoint') }}</label>
                <input
                  type="text"
                  class="input mono"
                  :value="instance.endpoint"
                  @input="updateRcloneInstance(index, 'endpoint', ($event.target as HTMLInputElement).value)"
                  placeholder="/rc/vfs/refresh"
                />
              </div>

              <div class="form-group full-width">
                <label class="checkbox-label">
                  <input
                    type="checkbox"
                    :checked="instance.wait_for_data !== false"
                    @change="updateRcloneInstance(index, 'wait_for_data', ($event.target as HTMLInputElement).checked)"
                  />
                  <span>{{ t('panels.integrations.waitForData') }}</span>
                </label>
              </div>
            </div>
          </div>

          <div v-if="!configStore.config?.rclone?.instances?.length" class="empty-state">
            <Server :size="32" />
            <p>{{ t('panels.integrations.noRcloneInstances') }}</p>
            <button class="btn btn-secondary btn-sm" @click="addRcloneInstance">
              <Plus :size="14" />
              {{ t('panels.integrations.addRcloneInstance') }}
            </button>
          </div>
        </div>

        <!-- Save Button -->
        <div class="section-footer">
          <button 
            class="btn btn-primary"
            @click="handleSave"
            :disabled="isSaving"
          >
            <Loader2 v-if="isSaving" :size="18" class="animate-spin" />
            <Save v-else :size="18" />
            <span>{{ isSaving ? t('common.saving') : t('common.save') }}</span>
          </button>
        </div>
      </section>

      <!-- Symedia Section -->
      <section class="config-section">
        <h3>
          <Bell :size="16" />
          {{ t('panels.integrations.symedia') }}
        </h3>

        <div class="form-grid">
          <div class="form-group">
            <label>{{ t('panels.integrations.symediaHost') }}</label>
            <input
              type="text"
              class="input mono"
              :value="configStore.config?.symedia?.host || ''"
              @input="updateConfig('symedia.host', ($event.target as HTMLInputElement).value)"
              placeholder="http://localhost:8095"
            />
          </div>

          <div class="form-group">
            <label>{{ t('panels.integrations.symediaEndpoint') }}</label>
            <input
              type="text"
              class="input mono"
              :value="configStore.config?.symedia?.endpoint || ''"
              @input="updateConfig('symedia.endpoint', ($event.target as HTMLInputElement).value)"
              placeholder="/api/v1/library/match"
            />
          </div>

          <div class="form-group full-width">
            <label>{{ t('panels.integrations.symediaTemplate') }}</label>
            <textarea
              class="input mono"
              rows="3"
              :value="configStore.config?.symedia?.body_template || ''"
              @input="updateConfig('symedia.body_template', ($event.target as HTMLInputElement).value)"
              :placeholder='`{"path": "{{path}}", "action": "{{action}}"}`'
            ></textarea>
            <span class="hint">
              {{ t('panels.integrations.symediaTemplateHint') }}
              <code v-pre>{{path}}</code>, <code v-pre>{{action}}</code>, <code v-pre>{{name}}</code>
            </span>
          </div>

          <div class="form-group full-width">
            <label class="checkbox-label">
              <input
                type="checkbox"
                :checked="configStore.config?.symedia?.notify_unmatched !== false"
                @change="updateConfig('symedia.notify_unmatched', ($event.target as HTMLInputElement).checked)"
              />
              <span>{{ t('panels.integrations.notifyUnmatched') }}</span>
            </label>
          </div>

          <!-- Headers Section -->
          <div class="form-group full-width">
            <div class="headers-section">
              <div class="headers-header">
                <label>{{ t('panels.integrations.headers') }}</label>
                <button type="button" class="add-btn btn-sm" @click="addSymediaHeader">
                  <Plus :size="12" />
                  <span>{{ t('common.add') }}</span>
                </button>
              </div>
              
              <div class="headers-list">
                <div
                  v-for="(value, key) in configStore.config?.symedia?.headers || {}"
                  :key="key"
                  class="header-item"
                >
                  <input
                    type="text"
                    class="input mono header-key"
                    :value="key"
                    @input="updateSymediaHeaderKey(key, ($event.target as HTMLInputElement).value)"
                    placeholder="Header Name"
                  />
                  <span class="header-separator">:</span>
                  <input
                    type="text"
                    class="input mono header-value"
                    :value="value"
                    @input="updateSymediaHeader(key, ($event.target as HTMLInputElement).value)"
                    placeholder="Header Value"
                  />
                  <button
                    type="button"
                    class="remove-btn btn-sm"
                    @click="removeSymediaHeader(key)"
                    :title="t('common.remove') || 'Remove'"
                  >
                    <Trash2 :size="12" />
                  </button>
                </div>
                
                <div v-if="!configStore.config?.symedia?.headers || Object.keys(configStore.config.symedia.headers).length === 0" class="empty-headers">
                  <span class="hint">{{ t('panels.integrations.noHeaders') }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Save Button -->
        <div class="section-footer">
          <button 
            class="btn btn-primary"
            @click="handleSave"
            :disabled="isSaving"
          >
            <Loader2 v-if="isSaving" :size="18" class="animate-spin" />
            <Save v-else :size="18" />
            <span>{{ isSaving ? t('common.saving') : t('common.save') }}</span>
          </button>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
@import './panel.css';

/* ========== Section Header ========== */
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-4);
}

.section-header h3 {
  margin-bottom: 0;
}

.add-btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: white;
  background: var(--color-accent);
  border: none;
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.add-btn:hover {
  background: var(--color-accent-hover);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
  color: white;
}

.add-btn:active {
  transform: translateY(0);
}

/* ========== Instances List ========== */
.instances-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.instance-card {
  padding: var(--space-4);
  background: var(--color-glass);
  border: 1px solid var(--color-glass-border);
  border-radius: var(--radius-lg);
}

.instance-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-3);
  padding-bottom: var(--space-2);
  border-bottom: 1px solid var(--color-glass-border);
}

.instance-index {
  font-size: var(--text-xs);
  font-weight: var(--font-bold);
  color: var(--color-accent);
  background: var(--color-accent-light);
  padding: 2px 8px;
  border-radius: var(--radius-full);
}

.remove-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: var(--radius-md);
  color: var(--color-text-tertiary);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.remove-btn:hover {
  color: var(--color-error);
  background: var(--color-error-light);
}

/* ========== Compact Form Grid ========== */
.form-grid.compact {
  gap: var(--space-3);
}

.form-grid.compact .form-group label {
  font-size: var(--text-xs);
}

.form-grid.compact .input {
  height: 36px;
  font-size: var(--text-sm);
}

/* ========== Empty State ========== */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-8);
  text-align: center;
  color: var(--color-text-quaternary);
}

.empty-state p {
  margin: 0;
  font-size: var(--text-sm);
}

/* ========== Headers Section ========== */
.headers-section {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.headers-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.headers-header label {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
  margin: 0;
}

.add-btn.btn-sm {
  padding: var(--space-1) var(--space-3);
  font-size: var(--text-xs);
  height: 28px;
}

.headers-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.header-item {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.header-key {
  flex: 0 0 150px;
  font-size: var(--text-sm);
  height: 36px;
}

.header-separator {
  color: var(--color-text-tertiary);
  font-weight: var(--font-semibold);
  flex-shrink: 0;
}

.header-value {
  flex: 1;
  font-size: var(--text-sm);
  height: 36px;
}

.remove-btn.btn-sm {
  width: 36px;
  height: 36px;
  padding: 0;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-headers {
  padding: var(--space-2);
  text-align: center;
  color: var(--color-text-quaternary);
}
</style>

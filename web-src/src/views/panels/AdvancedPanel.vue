<script setup lang="ts">
/**
 * AdvancedPanel - Advanced Settings Configuration
 */

import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfigStore } from '@/stores'
import { Settings, Clock, FolderOpen, Database, Save, Loader2 } from 'lucide-vue-next'
import CronEditor from '@/components/business/CronEditor.vue'

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
</script>

<template>
  <div class="panel">
    <div class="panel-content">
      <!-- Sync Settings -->
      <section class="config-section">
        <h3>
          <Clock :size="16" />
          {{ t('panels.advanced.syncSettings') }}
        </h3>

        <div class="form-grid">
          <div class="form-group">
            <label>{{ t('panels.advanced.debounce') }}</label>
            <div class="input-with-suffix">
              <input
                type="number"
                class="input"
                :value="configStore.config?.advanced?.debounce_seconds || 5"
                @input="updateConfig('advanced.debounce_seconds', Number(($event.target as HTMLInputElement).value))"
                min="1"
                max="300"
              />
              <span class="suffix">{{ t('common.seconds') }}</span>
            </div>
            <span class="hint">
              {{ t('panels.advanced.debounceHint') }}
            </span>
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

      <!-- Logging Settings -->
      <section class="config-section">
        <h3>
          <FolderOpen :size="16" />
          {{ t('panels.advanced.logging') }}
        </h3>

        <div class="form-grid">
          <div class="form-group">
            <label>{{ t('panels.advanced.logDir') }}</label>
            <input
              type="text"
              class="input mono"
              :value="configStore.config?.advanced?.log_dir || './logs'"
              @input="updateConfig('advanced.log_dir', ($event.target as HTMLInputElement).value)"
            />
          </div>

          <div class="form-group">
            <label>{{ t('panels.advanced.logLevel') }}</label>
            <select
              class="input"
              :value="configStore.config?.advanced?.log_level || 1"
              @change="updateConfig('advanced.log_level', Number(($event.target as HTMLSelectElement).value))"
            >
              <option :value="0">{{ t('panels.advanced.logLevels.quiet') }}</option>
              <option :value="1">{{ t('panels.advanced.logLevels.info') }}</option>
              <option :value="2">{{ t('panels.advanced.logLevels.debug') }}</option>
            </select>
          </div>

          <div class="form-group full-width">
            <label class="checkbox-label">
              <input
                type="checkbox"
                :checked="configStore.config?.advanced?.log_save_enabled !== false"
                @change="updateConfig('advanced.log_save_enabled', ($event.target as HTMLInputElement).checked)"
              />
              <span>{{ t('panels.advanced.enableLogSave') }}</span>
            </label>
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

      <!-- Log Cleanup -->
      <section class="config-section">
        <h3>
          <Database :size="16" />
          {{ t('panels.advanced.logCleanup') }}
        </h3>

        <div class="form-grid">
          <div class="form-group full-width">
            <label class="checkbox-label">
              <input
                type="checkbox"
                :checked="configStore.config?.advanced?.log_cleanup?.enabled || false"
                @change="updateConfig('advanced.log_cleanup.enabled', ($event.target as HTMLInputElement).checked)"
              />
              <span>{{ t('panels.advanced.enableLogCleanup') }}</span>
            </label>
          </div>

          <template v-if="configStore.config?.advanced?.log_cleanup?.enabled">
            <div class="form-group">
              <label>{{ t('panels.advanced.retentionDays') }}</label>
              <div class="input-with-suffix">
                <input
                  type="number"
                  class="input"
                  :value="configStore.config?.advanced?.log_cleanup?.retention_days || 30"
                  @input="updateConfig('advanced.log_cleanup.retention_days', Number(($event.target as HTMLInputElement).value))"
                  min="1"
                  max="365"
                />
                <span class="suffix">{{ t('common.days') }}</span>
              </div>
            </div>

            <div class="form-group full-width">
              <label>{{ t('panels.advanced.cleanupCron') }}</label>
              <CronEditor
                :model-value="configStore.config?.advanced?.log_cleanup?.cron || '0 0 * * *'"
                @update:model-value="updateConfig('advanced.log_cleanup.cron', $event)"
              />
              <span class="hint">
                <a href="https://crontab.guru/" target="_blank" rel="noopener noreferrer">
                  {{ t('panels.advanced.cronLearnMore') }}
                  <ExternalLink :size="12" />
                </a>
              </span>
            </div>
          </template>
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

/* ========== Input with Suffix ========== */
.input-with-suffix {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.input-with-suffix .input {
  flex: 1;
  max-width: 120px;
}

.suffix {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
}
</style>

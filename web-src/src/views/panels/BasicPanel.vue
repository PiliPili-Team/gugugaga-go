<script setup lang="ts">
/**
 * BasicPanel - Basic Settings Configuration Panel
 */

import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfigStore } from '@/stores'
import { Server, Lock, Globe, Shield, Link, Save, Loader2, ExternalLink, Eye, EyeOff } from 'lucide-vue-next'

const { t } = useI18n()
const configStore = useConfigStore()
const isSaving = ref(false)
const showPassword = ref(false)

// Helper to update nested config
function updateConfig(path: string, value: any) {
  configStore.updateNested(path, value)
}

// Save this panel's config
async function handleSave() {
  isSaving.value = true
  try {
    await configStore.saveConfig()
  } finally {
    isSaving.value = false
  }
}
</script>

<template>
  <div class="panel">
    <div class="panel-content">
      <!-- Server Section -->
      <section class="config-section">
        <h3>
          <Server :size="16" />
          {{ t('panels.basic.server') }}
        </h3>

        <div class="form-grid">
          <div class="form-group">
            <label>{{ t('panels.basic.port') }}</label>
            <input
              type="number"
              class="input"
              :value="configStore.config?.server?.port || 8448"
              @input="updateConfig('server.port', Number(($event.target as HTMLInputElement).value))"
              min="1"
              max="65535"
            />
            <span class="hint">{{ t('panels.basic.portHint') }}</span>
          </div>

          <div class="form-group">
            <label>{{ t('panels.basic.publicUrl') }}</label>
            <div class="input-with-icon">
              <Globe :size="16" />
              <input
                type="text"
                class="input"
                :value="configStore.config?.server?.public_url || ''"
                @input="updateConfig('server.public_url', ($event.target as HTMLInputElement).value)"
                :placeholder="t('panels.basic.publicUrlPlaceholder')"
              />
            </div>
            <span class="hint">
              {{ t('panels.basic.publicUrlHint') }}
              <a href="https://developers.google.com/drive/api/v3/push" target="_blank" rel="noopener noreferrer">
                {{ t('panels.basic.learnMore') }}
                <ExternalLink :size="12" />
              </a>
            </span>
          </div>

          <div class="form-group">
            <label>{{ t('panels.basic.webhookPath') }}</label>
            <div class="input-with-icon">
              <Link :size="16" />
              <input
                type="text"
                class="input mono"
                :value="configStore.config?.server?.webhook_path || '/webhook'"
                @input="updateConfig('server.webhook_path', ($event.target as HTMLInputElement).value)"
                placeholder="/webhook"
              />
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
            <Loader2 v-if="isSaving" :size="16" class="animate-spin" />
            <Save v-else :size="16" />
            <span>{{ isSaving ? t('common.saving') : t('common.save') }}</span>
          </button>
        </div>
      </section>

      <!-- Authentication Section -->
      <section class="config-section">
        <h3>
          <Lock :size="16" />
          {{ t('panels.basic.auth') }}
        </h3>

        <div class="form-grid">
          <div class="form-group">
            <label>{{ t('panels.basic.username') }}</label>
            <input
              type="text"
              class="input"
              :value="configStore.config?.auth?.username || ''"
              @input="updateConfig('auth.username', ($event.target as HTMLInputElement).value)"
              autocomplete="username"
            />
          </div>

          <div class="form-group">
            <label>{{ t('panels.basic.password') }}</label>
            <div class="input-with-icon password-input">
              <Lock :size="16" class="password-lock-icon" />
              <input
                :type="showPassword ? 'text' : 'password'"
                class="input password-input-field"
                :value="configStore.config?.auth?.password || ''"
                @input="updateConfig('auth.password', ($event.target as HTMLInputElement).value)"
                autocomplete="current-password"
                placeholder="••••••••"
              />
              <button
                type="button"
                class="password-toggle"
                @click="showPassword = !showPassword"
                :title="showPassword ? t('panels.basic.hidePassword') : t('panels.basic.showPassword')"
              >
                <Eye v-if="showPassword" :size="18" />
                <EyeOff v-else :size="18" />
              </button>
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
            <Loader2 v-if="isSaving" :size="16" class="animate-spin" />
            <Save v-else :size="16" />
            <span>{{ isSaving ? t('common.saving') : t('common.save') }}</span>
          </button>
        </div>
      </section>

      <!-- SSL Section -->
      <section class="config-section">
        <h3>
          <Shield :size="16" />
          {{ t('panels.basic.ssl') }}
        </h3>

        <div class="form-grid">
          <div class="form-group full-width">
            <label class="checkbox-label">
              <input
                type="checkbox"
                :checked="configStore.config?.server?.ssl_enabled || false"
                @change="updateConfig('server.ssl_enabled', ($event.target as HTMLInputElement).checked)"
              />
              <span>{{ t('panels.basic.enableSsl') }}</span>
            </label>
          </div>

          <template v-if="configStore.config?.server?.ssl_enabled">
            <div class="form-group">
              <label>{{ t('panels.basic.sslCert') }}</label>
              <input
                type="text"
                class="input mono"
                :value="configStore.config?.server?.ssl_cert || ''"
                @input="updateConfig('server.ssl_cert', ($event.target as HTMLInputElement).value)"
                :placeholder="t('panels.basic.sslCertPlaceholder')"
              />
            </div>

            <div class="form-group">
              <label>{{ t('panels.basic.sslKey') }}</label>
              <input
                type="text"
                class="input mono"
                :value="configStore.config?.server?.ssl_key || ''"
                @input="updateConfig('server.ssl_key', ($event.target as HTMLInputElement).value)"
                :placeholder="t('panels.basic.sslKeyPlaceholder')"
              />
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
            <Loader2 v-if="isSaving" :size="16" class="animate-spin" />
            <Save v-else :size="16" />
            <span>{{ isSaving ? t('common.saving') : t('common.save') }}</span>
          </button>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
@import './panel.css';
</style>

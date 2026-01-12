<script setup lang="ts">
/**
 * OAuthPanel - OAuth Configuration
 */

import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfigStore } from '@/stores'
import { api } from '@/services/api'
import { Key, ExternalLink, Loader2, CheckCircle, AlertCircle, Save, Eye, EyeOff } from 'lucide-vue-next'

const { t } = useI18n()
const configStore = useConfigStore()

const isGettingUrl = ref(false)
const oauthStatus = ref<'idle' | 'success' | 'error'>('idle')
const isSaving = ref(false)
const showSecret = ref(false)

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

async function handleGoToOAuth() {
  isGettingUrl.value = true
  oauthStatus.value = 'idle'
  
  try {
    const response = await api.getOAuthLoginURL()
    if (response.url) {
      oauthStatus.value = 'success'
      window.open(response.url, '_blank')
    } else {
      oauthStatus.value = 'error'
    }
  } catch (error) {
    oauthStatus.value = 'error'
  } finally {
    isGettingUrl.value = false
  }
}
</script>

<template>
  <div class="panel">
    <div class="panel-content">
      <!-- OAuth Credentials -->
      <section class="config-section">
        <h3>
          <Key :size="16" />
          {{ t('panels.oauth.credentials') }}
        </h3>

        <div class="form-grid">
          <div class="form-group">
            <label>Client ID</label>
            <input
              type="text"
              class="input mono"
              :value="configStore.config?.oauth?.client_id || ''"
              @input="updateConfig('oauth.client_id', ($event.target as HTMLInputElement).value)"
              :placeholder="t('panels.oauth.clientIdPlaceholder')"
            />
          </div>

          <div class="form-group">
            <label>Client Secret</label>
            <div class="password-input">
              <input
                :type="showSecret ? 'text' : 'password'"
                class="input mono"
                :value="configStore.config?.oauth?.client_secret || ''"
                @input="updateConfig('oauth.client_secret', ($event.target as HTMLInputElement).value)"
                :placeholder="t('panels.oauth.clientSecretPlaceholder')"
              />
              <button 
                type="button" 
                class="password-toggle"
                @click="showSecret = !showSecret"
                :title="showSecret ? t('common.hide') : t('common.show')"
              >
                <EyeOff v-if="showSecret" :size="18" />
                <Eye v-else :size="18" />
              </button>
            </div>
          </div>

          <div class="form-group full-width">
            <label>{{ t('panels.oauth.redirectUri') }}</label>
            <input
              type="text"
              class="input mono"
              :value="configStore.config?.oauth?.redirect_uri || ''"
              @input="updateConfig('oauth.redirect_uri', ($event.target as HTMLInputElement).value)"
              :placeholder="t('panels.oauth.redirectUriPlaceholder')"
            />
            <span class="hint">{{ t('panels.oauth.redirectUriHint') }}</span>
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

      <!-- OAuth Action -->
      <section class="config-section">
        <h3>
          <ExternalLink :size="16" />
          {{ t('panels.oauth.authorization') }}
        </h3>

        <div class="oauth-action">
          <p class="action-desc">{{ t('panels.oauth.actionDesc') }}</p>
          
          <button 
            class="btn btn-primary btn-lg oauth-action-btn"
            @click="handleGoToOAuth"
            :disabled="isGettingUrl || !configStore.config?.oauth?.client_id"
          >
            <Loader2 v-if="isGettingUrl" :size="20" class="animate-spin" />
            <ExternalLink v-else :size="20" />
            <span>{{ t('panels.oauth.goToGoogle') }}</span>
          </button>

          <div v-if="oauthStatus === 'success'" class="status-message success">
            <CheckCircle :size="16" />
            <span>{{ t('panels.oauth.urlOpened') }}</span>
          </div>

          <div v-if="oauthStatus === 'error'" class="status-message error">
            <AlertCircle :size="16" />
            <span>{{ t('panels.oauth.urlError') }}</span>
          </div>

          <div v-if="!configStore.config?.oauth?.client_id" class="status-message warning">
            <AlertCircle :size="16" />
            <span>{{ t('panels.oauth.noClientId') }}</span>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
@import './panel.css';

/* ========== Password Input ========== */
.password-input {
  position: relative;
  display: flex;
  align-items: center;
}

.password-input .input {
  width: 100%;
  padding-right: 44px;
}

.password-toggle {
  position: absolute;
  right: var(--space-2);
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  color: var(--color-text-tertiary);
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.password-toggle:hover {
  color: var(--color-text-primary);
  background: var(--color-glass);
}

/* ========== OAuth Action ========== */
.oauth-action {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.action-desc {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
  margin: 0;
}

.btn-lg {
  height: 48px;
  padding: 0 var(--space-6);
  font-size: var(--text-base);
  gap: var(--space-3);
}

.oauth-action-btn {
  font-weight: var(--font-semibold);
  font-size: var(--text-lg);
  height: 52px;
  padding: 0 var(--space-8);
  color: white;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.oauth-action-btn:hover:not(:disabled) {
  color: white;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(59, 130, 246, 0.4);
}

/* ========== Status Messages ========== */
.status-message {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-4);
  border-radius: var(--radius-lg);
  font-size: var(--text-sm);
}

.status-message.success {
  background: var(--color-success-light);
  color: var(--color-success);
  border: 1px solid rgba(34, 197, 94, 0.2);
}

.status-message.error {
  background: var(--color-error-light);
  color: var(--color-error);
  border: 1px solid rgba(239, 68, 68, 0.2);
}

.status-message.warning {
  background: var(--color-warning-light);
  color: var(--color-warning);
  border: 1px solid rgba(245, 158, 11, 0.2);
}
</style>

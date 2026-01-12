<script setup lang="ts">
/**
 * MappingsPanel - Path Mapping Rules Configuration
 */

import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfigStore } from '@/stores'
import { Route, ArrowLeftRight, Save, Loader2 } from 'lucide-vue-next'
import MappingList from '@/components/business/MappingList.vue'
import type { MappingRule } from '@/types'

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

function updateSaMappings(rules: MappingRule[]) {
  configStore.updateNested('symedia.path_mappings', rules)
}

function updateRcloneMappings(rules: MappingRule[]) {
  configStore.updateNested('rclone.path_mappings', rules)
}
</script>

<template>
  <div class="panel">
    <div class="panel-content">
      <!-- SA Mappings -->
      <section class="config-section">
        <div class="section-header">
          <h3>
            <ArrowLeftRight :size="16" />
            {{ t('panels.mappings.saMappings') }}
          </h3>
        </div>
        <p class="section-desc">{{ t('panels.mappings.saMappingsDesc') }}</p>

        <MappingList
          :model-value="configStore.config?.symedia?.path_mappings || []"
          @update:model-value="updateSaMappings"
          :title="t('panels.mappings.saMappings')"
        />

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

      <!-- Rclone Mappings -->
      <section class="config-section">
        <div class="section-header">
          <h3>
            <ArrowLeftRight :size="16" />
            {{ t('panels.mappings.rcloneMappings') }}
          </h3>
        </div>
        <p class="section-desc">{{ t('panels.mappings.rcloneMappingsDesc') }}</p>

        <MappingList
          :model-value="configStore.config?.rclone?.path_mappings || []"
          @update:model-value="updateRcloneMappings"
          :title="t('panels.mappings.rcloneMappings')"
        />

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

.section-desc {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
  margin-bottom: var(--space-4);
}
</style>

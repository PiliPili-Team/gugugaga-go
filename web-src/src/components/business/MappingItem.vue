<script setup lang="ts">
/**
 * MappingItem - Single Mapping Rule Item
 */

import { useI18n } from 'vue-i18n'
import { Trash2 } from 'lucide-vue-next'

const { t } = useI18n()

defineProps<{
  regex: string
  replacement: string
  index: number
}>()

const emit = defineEmits<{
  (e: 'update:regex', value: string): void
  (e: 'update:replacement', value: string): void
  (e: 'remove'): void
}>()
</script>

<template>
  <div class="mapping-item">
    <div class="mapping-fields">
      <div class="field-group">
        <label>{{ t('mappings.regex') }}</label>
        <input
          type="text"
          class="input mono"
          :value="regex"
          @input="emit('update:regex', ($event.target as HTMLInputElement).value)"
          :placeholder="t('mappings.regexPlaceholder')"
        />
      </div>
      <div class="field-group">
        <label>{{ t('mappings.replacement') }}</label>
        <input
          type="text"
          class="input mono"
          :value="replacement"
          @input="emit('update:replacement', ($event.target as HTMLInputElement).value)"
          :placeholder="t('mappings.replacementPlaceholder')"
        />
      </div>
    </div>
    <button class="remove-btn" @click="emit('remove')" :title="t('common.delete')">
      <Trash2 :size="14" />
    </button>
  </div>
</template>

<style scoped>
.mapping-item {
  display: flex;
  align-items: flex-start;
  gap: var(--space-3);
  padding: var(--space-3);
  background: var(--color-glass);
  border: 1px solid var(--color-glass-border);
  border-radius: var(--radius-lg);
  transition: all var(--duration-fast) var(--ease-default);
}

.mapping-item:hover {
  border-color: var(--color-glass-border-hover);
}

.mapping-fields {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-3);
}

.field-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.field-group label {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  font-weight: var(--font-medium);
}

.field-group .input {
  height: 36px;
  font-size: var(--text-sm);
}

.remove-btn {
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
  margin-top: 20px;
  transition: all var(--duration-fast) var(--ease-default);
}

.remove-btn:hover {
  color: var(--color-error);
  background: var(--color-error-light);
}

/* Responsive */
@media (max-width: 640px) {
  .mapping-fields {
    grid-template-columns: 1fr;
  }
}
</style>

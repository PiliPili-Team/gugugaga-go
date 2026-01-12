<script setup lang="ts">
/**
 * BaseSelect - 下拉选择组件
 */

import { computed } from 'vue'

export interface SelectOption {
  value: string | number
  label: string
  disabled?: boolean
}

export type SelectSize = 'sm' | 'md' | 'lg'

interface Props {
  modelValue: string | number
  options: SelectOption[]
  size?: SelectSize
  placeholder?: string
  disabled?: boolean
  error?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  placeholder: '请选择...',
  disabled: false,
  error: false
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | number): void
  (e: 'change', value: string | number): void
}>()

const selectValue = computed({
  get: () => props.modelValue,
  set: (value) => {
    emit('update:modelValue', value)
    emit('change', value)
  }
})
</script>

<template>
  <div 
    class="base-select"
    :class="[
      `select-${size}`,
      {
        'select-error': error,
        'select-disabled': disabled
      }
    ]"
  >
    <select
      v-model="selectValue"
      :disabled="disabled"
      class="select-element"
    >
      <option v-if="placeholder" value="" disabled>
        {{ placeholder }}
      </option>
      <option
        v-for="option in options"
        :key="option.value"
        :value="option.value"
        :disabled="option.disabled"
      >
        {{ option.label }}
      </option>
    </select>

    <!-- Arrow icon -->
    <span class="select-arrow">
      <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
        <path d="M3 4.5L6 7.5L9 4.5" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
    </span>
  </div>
</template>

<style scoped>
.base-select {
  position: relative;
  display: flex;
  align-items: center;
  background: var(--color-bg-base);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  transition: all var(--transition-fast);
}

.base-select:focus-within {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px var(--color-primary-light);
}

/* ========== Sizes ========== */
.select-sm {
  height: 32px;
}

.select-sm .select-element {
  padding: 0 var(--space-6) 0 var(--space-2);
  font-size: var(--font-size-xs);
}

.select-md {
  height: 40px;
}

.select-md .select-element {
  padding: 0 var(--space-8) 0 var(--space-3);
  font-size: var(--font-size-sm);
}

.select-lg {
  height: 48px;
}

.select-lg .select-element {
  padding: 0 var(--space-10) 0 var(--space-4);
  font-size: var(--font-size-base);
}

/* ========== Select Element ========== */
.select-element {
  flex: 1;
  min-width: 0;
  height: 100%;
  background: transparent;
  border: none;
  outline: none;
  color: var(--color-text-primary);
  font-family: inherit;
  cursor: pointer;
  appearance: none;
  -webkit-appearance: none;
}

.select-element option {
  background: var(--color-bg-surface);
  color: var(--color-text-primary);
}

/* ========== Arrow ========== */
.select-arrow {
  position: absolute;
  right: var(--space-3);
  pointer-events: none;
  color: var(--color-text-tertiary);
  display: flex;
  align-items: center;
  justify-content: center;
}

/* ========== States ========== */
.select-error {
  border-color: var(--color-danger);
}

.select-error:focus-within {
  box-shadow: 0 0 0 2px var(--color-danger-light);
}

.select-disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: var(--color-bg-surface);
}

.select-disabled .select-element {
  cursor: not-allowed;
}
</style>

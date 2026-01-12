<script setup lang="ts">
/**
 * BaseCheckbox - 复选框组件
 */

import { computed } from 'vue'

interface Props {
  modelValue: boolean
  label?: string
  disabled?: boolean
  size?: 'sm' | 'md'
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false,
  size: 'md'
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'change', value: boolean): void
}>()

const checked = computed({
  get: () => props.modelValue,
  set: (value) => {
    emit('update:modelValue', value)
    emit('change', value)
  }
})
</script>

<template>
  <label 
    class="base-checkbox"
    :class="[
      `checkbox-${size}`,
      {
        'checkbox-checked': modelValue,
        'checkbox-disabled': disabled
      }
    ]"
  >
    <input
      type="checkbox"
      v-model="checked"
      :disabled="disabled"
      class="checkbox-input"
    />

    <span class="checkbox-box">
      <svg class="checkbox-icon" viewBox="0 0 12 12" fill="none">
        <path 
          d="M2.5 6L5 8.5L9.5 3.5" 
          stroke="currentColor" 
          stroke-width="2" 
          stroke-linecap="round" 
          stroke-linejoin="round"
        />
      </svg>
    </span>

    <span v-if="label || $slots.default" class="checkbox-label">
      <slot>{{ label }}</slot>
    </span>
  </label>
</template>

<style scoped>
.base-checkbox {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  cursor: pointer;
  user-select: none;
}

.checkbox-input {
  position: absolute;
  opacity: 0;
  pointer-events: none;
}

/* ========== Checkbox Box ========== */
.checkbox-box {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: var(--color-bg-base);
  border: 2px solid var(--color-border);
  border-radius: var(--radius-sm);
  transition: all var(--transition-fast);
}

.checkbox-sm .checkbox-box {
  width: 16px;
  height: 16px;
}

.checkbox-md .checkbox-box {
  width: 20px;
  height: 20px;
}

/* ========== Checkbox Icon ========== */
.checkbox-icon {
  opacity: 0;
  transform: scale(0.5);
  transition: all var(--transition-fast);
  color: var(--color-bg-base);
}

.checkbox-sm .checkbox-icon {
  width: 10px;
  height: 10px;
}

.checkbox-md .checkbox-icon {
  width: 12px;
  height: 12px;
}

/* ========== States ========== */
.base-checkbox:hover:not(.checkbox-disabled) .checkbox-box {
  border-color: var(--color-primary);
}

.checkbox-checked .checkbox-box {
  background: var(--color-primary);
  border-color: var(--color-primary);
}

.checkbox-checked .checkbox-icon {
  opacity: 1;
  transform: scale(1);
}

.checkbox-disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* ========== Label ========== */
.checkbox-label {
  font-size: var(--font-size-sm);
  color: var(--color-text-primary);
  line-height: var(--line-height-normal);
}

.checkbox-sm .checkbox-label {
  font-size: var(--font-size-xs);
}

/* Focus state */
.checkbox-input:focus-visible + .checkbox-box {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
}
</style>

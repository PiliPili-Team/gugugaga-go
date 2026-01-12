<script setup lang="ts">
/**
 * BaseInput - 输入框组件
 */

import { computed } from 'vue'

export type InputType = 'text' | 'password' | 'number' | 'email' | 'url' | 'tel'
export type InputSize = 'sm' | 'md' | 'lg'

interface Props {
  modelValue: string | number
  type?: InputType
  size?: InputSize
  placeholder?: string
  disabled?: boolean
  readonly?: boolean
  error?: boolean
  mono?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  size: 'md',
  placeholder: '',
  disabled: false,
  readonly: false,
  error: false,
  mono: false
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | number): void
  (e: 'focus', event: FocusEvent): void
  (e: 'blur', event: FocusEvent): void
}>()

const inputValue = computed({
  get: () => props.modelValue,
  set: (value) => {
    if (props.type === 'number') {
      emit('update:modelValue', value === '' ? '' : Number(value))
    } else {
      emit('update:modelValue', value)
    }
  }
})

function handleFocus(event: FocusEvent) {
  emit('focus', event)
}

function handleBlur(event: FocusEvent) {
  emit('blur', event)
}
</script>

<template>
  <div 
    class="base-input"
    :class="[
      `input-${size}`,
      {
        'input-error': error,
        'input-disabled': disabled,
        'input-mono': mono
      }
    ]"
  >
    <!-- Prefix slot -->
    <span v-if="$slots.prefix" class="input-prefix">
      <slot name="prefix" />
    </span>

    <!-- Input element -->
    <input
      v-model="inputValue"
      :type="type"
      :placeholder="placeholder"
      :disabled="disabled"
      :readonly="readonly"
      class="input-element"
      @focus="handleFocus"
      @blur="handleBlur"
    />

    <!-- Suffix slot -->
    <span v-if="$slots.suffix" class="input-suffix">
      <slot name="suffix" />
    </span>
  </div>
</template>

<style scoped>
.base-input {
  display: flex;
  align-items: center;
  background: var(--color-bg-base);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  transition: all var(--transition-fast);
  overflow: hidden;
}

.base-input:focus-within {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 2px var(--color-primary-light);
}

/* ========== Sizes ========== */
.input-sm {
  height: 32px;
}

.input-sm .input-element {
  padding: 0 var(--space-2);
  font-size: var(--font-size-xs);
}

.input-md {
  height: 40px;
}

.input-md .input-element {
  padding: 0 var(--space-3);
  font-size: var(--font-size-sm);
}

.input-lg {
  height: 48px;
}

.input-lg .input-element {
  padding: 0 var(--space-4);
  font-size: var(--font-size-base);
}

/* ========== Input Element ========== */
.input-element {
  flex: 1;
  min-width: 0;
  height: 100%;
  background: transparent;
  border: none;
  outline: none;
  color: var(--color-text-primary);
  font-family: inherit;
}

.input-element::placeholder {
  color: var(--color-text-tertiary);
}

/* ========== States ========== */
.input-error {
  border-color: var(--color-danger);
}

.input-error:focus-within {
  box-shadow: 0 0 0 2px var(--color-danger-light);
}

.input-disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background: var(--color-bg-surface);
}

.input-disabled .input-element {
  cursor: not-allowed;
}

/* ========== Mono ========== */
.input-mono .input-element {
  font-family: var(--font-family-mono);
  font-weight: var(--font-weight-medium);
  letter-spacing: 0.5px;
}

/* ========== Prefix & Suffix ========== */
.input-prefix,
.input-suffix {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 var(--space-3);
  color: var(--color-text-tertiary);
  font-size: var(--font-size-sm);
}

.input-prefix {
  border-right: 1px solid var(--color-border);
  background: var(--color-bg-surface);
}

.input-suffix {
  border-left: 1px solid var(--color-border);
  background: var(--color-bg-surface);
}
</style>

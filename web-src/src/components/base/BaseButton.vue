<script setup lang="ts">
/**
 * BaseButton - 基础按钮组件
 */

export type ButtonVariant = 'primary' | 'secondary' | 'success' | 'warning' | 'danger' | 'ghost'
export type ButtonSize = 'sm' | 'md' | 'lg'

interface Props {
  variant?: ButtonVariant
  size?: ButtonSize
  loading?: boolean
  disabled?: boolean
  block?: boolean
  iconOnly?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  loading: false,
  disabled: false,
  block: false,
  iconOnly: false
})

const emit = defineEmits<{
  (e: 'click', event: MouseEvent): void
}>()

function handleClick(event: MouseEvent) {
  if (!props.loading && !props.disabled) {
    emit('click', event)
  }
}
</script>

<template>
  <button
    class="base-btn"
    :class="[
      `btn-${variant}`,
      `btn-${size}`,
      {
        'btn-block': block,
        'btn-loading': loading,
        'btn-icon-only': iconOnly
      }
    ]"
    :disabled="disabled || loading"
    @click="handleClick"
  >
    <!-- Loading spinner -->
    <span v-if="loading" class="btn-spinner">
      <svg class="spinner-icon" viewBox="0 0 24 24">
        <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3" fill="none" stroke-dasharray="60" stroke-linecap="round" />
      </svg>
    </span>

    <!-- Icon slot -->
    <span v-if="$slots.icon && !loading" class="btn-icon-slot">
      <slot name="icon" />
    </span>

    <!-- Default slot -->
    <span v-if="$slots.default && !iconOnly" class="btn-content">
      <slot />
    </span>
  </button>
</template>

<style scoped>
.base-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  border-radius: var(--radius-md);
  font-family: var(--font-family-ui);
  font-weight: var(--font-weight-medium);
  cursor: pointer;
  transition: all var(--transition-fast);
  border: 1px solid transparent;
  white-space: nowrap;
  user-select: none;
}

.base-btn:focus-visible {
  outline: 2px solid var(--color-primary);
  outline-offset: 2px;
}

.base-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* ========== Sizes ========== */
.btn-sm {
  height: 32px;
  padding: 0 var(--space-3);
  font-size: var(--font-size-xs);
}

.btn-md {
  height: 40px;
  padding: 0 var(--space-4);
  font-size: var(--font-size-sm);
}

.btn-lg {
  height: 48px;
  padding: 0 var(--space-6);
  font-size: var(--font-size-base);
}

.btn-icon-only.btn-sm {
  width: 32px;
  padding: 0;
}

.btn-icon-only.btn-md {
  width: 40px;
  padding: 0;
}

.btn-icon-only.btn-lg {
  width: 48px;
  padding: 0;
}

/* ========== Variants ========== */
.btn-primary {
  background: var(--color-primary);
  color: var(--color-bg-base);
  border-color: var(--color-primary);
}

.btn-primary:hover:not(:disabled) {
  background: var(--color-primary-hover);
  border-color: var(--color-primary-hover);
  box-shadow: var(--shadow-glow-sm);
}

.btn-secondary {
  background: var(--color-secondary);
  color: var(--color-bg-base);
  border-color: var(--color-secondary);
}

.btn-secondary:hover:not(:disabled) {
  background: var(--color-secondary-hover);
  border-color: var(--color-secondary-hover);
}

.btn-success {
  background: var(--color-success-light);
  color: var(--color-success);
  border-color: var(--color-success);
}

.btn-success:hover:not(:disabled) {
  background: var(--color-success);
  color: var(--color-bg-base);
}

.btn-warning {
  background: var(--color-warning-light);
  color: var(--color-warning);
  border-color: var(--color-warning);
}

.btn-warning:hover:not(:disabled) {
  background: var(--color-warning);
  color: var(--color-bg-base);
}

.btn-danger {
  background: var(--color-danger-light);
  color: var(--color-danger);
  border-color: var(--color-danger);
}

.btn-danger:hover:not(:disabled) {
  background: var(--color-danger);
  color: var(--color-bg-base);
}

.btn-ghost {
  background: transparent;
  color: var(--color-text-secondary);
  border-color: var(--color-border);
}

.btn-ghost:hover:not(:disabled) {
  background: var(--color-bg-hover);
  color: var(--color-text-primary);
  border-color: var(--color-text-tertiary);
}

/* ========== Block ========== */
.btn-block {
  width: 100%;
}

/* ========== Loading ========== */
.btn-loading {
  pointer-events: none;
}

.btn-spinner {
  display: flex;
  align-items: center;
  justify-content: center;
}

.spinner-icon {
  width: 16px;
  height: 16px;
  animation: spin 1s linear infinite;
}

.btn-lg .spinner-icon {
  width: 20px;
  height: 20px;
}

/* ========== Icon ========== */
.btn-icon-slot {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.1em;
}
</style>

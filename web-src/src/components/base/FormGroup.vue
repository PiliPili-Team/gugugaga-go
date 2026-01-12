<script setup lang="ts">
/**
 * FormGroup - 表单项包装组件
 */

interface Props {
  label?: string
  required?: boolean
  error?: string
  hint?: string
  labelWidth?: string
  inline?: boolean
}

withDefaults(defineProps<Props>(), {
  required: false,
  inline: false
})
</script>

<template>
  <div 
    class="form-group"
    :class="{ 'form-group-inline': inline }"
    :style="labelWidth ? { '--label-width': labelWidth } : {}"
  >
    <label v-if="label" class="form-label">
      {{ label }}
      <span v-if="required" class="label-required">*</span>
    </label>

    <div class="form-content">
      <slot />

      <p v-if="hint && !error" class="form-hint">
        {{ hint }}
      </p>

      <p v-if="error" class="form-error">
        {{ error }}
      </p>
    </div>
  </div>
</template>

<style scoped>
.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.form-group-inline {
  flex-direction: row;
  align-items: flex-start;
}

.form-group-inline .form-label {
  width: var(--label-width, 120px);
  flex-shrink: 0;
  padding-top: var(--space-2);
}

/* ========== Label ========== */
.form-label {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  font-size: var(--font-size-sm);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-secondary);
}

.label-required {
  color: var(--color-danger);
}

/* ========== Content ========== */
.form-content {
  flex: 1;
  min-width: 0;
}

/* ========== Hint & Error ========== */
.form-hint,
.form-error {
  margin-top: var(--space-1);
  font-size: var(--font-size-xs);
  line-height: var(--line-height-normal);
}

.form-hint {
  color: var(--color-text-tertiary);
}

.form-error {
  color: var(--color-danger);
}
</style>

<script setup lang="ts">
/**
 * BaseCard - 卡片容器组件
 */

interface Props {
  title?: string
  subtitle?: string
  padding?: 'none' | 'sm' | 'md' | 'lg'
  bordered?: boolean
  hoverable?: boolean
}

withDefaults(defineProps<Props>(), {
  padding: 'md',
  bordered: false,
  hoverable: false
})
</script>

<template>
  <div 
    class="base-card"
    :class="[
      `card-padding-${padding}`,
      {
        'card-bordered': bordered,
        'card-hoverable': hoverable
      }
    ]"
  >
    <!-- Header -->
    <div v-if="title || $slots.header || $slots.actions" class="card-header">
      <div class="card-header-content">
        <slot name="header">
          <div class="card-titles">
            <h3 v-if="title" class="card-title">{{ title }}</h3>
            <p v-if="subtitle" class="card-subtitle">{{ subtitle }}</p>
          </div>
        </slot>
      </div>
      <div v-if="$slots.actions" class="card-actions">
        <slot name="actions" />
      </div>
    </div>

    <!-- Body -->
    <div class="card-body">
      <slot />
    </div>

    <!-- Footer -->
    <div v-if="$slots.footer" class="card-footer">
      <slot name="footer" />
    </div>
  </div>
</template>

<style scoped>
.base-card {
  background: var(--color-bg-surface);
  border-radius: var(--radius-lg);
  transition: all var(--transition-fast);
}

/* ========== Padding ========== */
.card-padding-none .card-body {
  padding: 0;
}

.card-padding-sm .card-body {
  padding: var(--space-3);
}

.card-padding-md .card-body {
  padding: var(--space-4);
}

.card-padding-lg .card-body {
  padding: var(--space-6);
}

/* ========== Bordered ========== */
.card-bordered {
  border: 1px solid var(--color-border);
}

/* ========== Hoverable ========== */
.card-hoverable {
  cursor: pointer;
}

.card-hoverable:hover {
  background: var(--color-bg-hover);
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

/* ========== Header ========== */
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4);
  border-bottom: 1px solid var(--color-border);
}

.card-header-content {
  flex: 1;
  min-width: 0;
}

.card-titles {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.card-title {
  font-size: var(--font-size-md);
  font-weight: var(--font-weight-semibold);
  color: var(--color-text-primary);
  margin: 0;
}

.card-subtitle {
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  margin: 0;
}

.card-actions {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  margin-left: var(--space-4);
}

/* ========== Footer ========== */
.card-footer {
  padding: var(--space-4);
  border-top: 1px solid var(--color-border);
}

/* ========== Dark variant for log viewer ========== */
.base-card.card-dark {
  background: var(--color-bg-base);
}
</style>

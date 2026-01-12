<script setup lang="ts">
/**
 * CronEditor - Cron Expression Editor with Visual Feedback
 */

import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Clock, AlertCircle, CheckCircle } from 'lucide-vue-next'

const { t } = useI18n()

const props = defineProps<{
  modelValue: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

// Parse cron expression
const parsedCron = computed(() => {
  const parts = props.modelValue.trim().split(/\s+/)
  if (parts.length < 5) return null

  return {
    minute: parts[0],
    hour: parts[1],
    day: parts[2],
    month: parts[3],
    weekday: parts[4],
    command: parts.slice(5).join(' ')
  }
})

// Validate cron expression
const isValidCron = computed(() => {
  if (!parsedCron.value) return false

  const { minute, hour, day, month, weekday } = parsedCron.value

  const validateField = (value: string, min: number, max: number): boolean => {
    if (value === '*') return true
    if (value.includes('/')) {
      const [, step] = value.split('/')
      return !isNaN(Number(step))
    }
    if (value.includes('-')) {
      const [start, end] = value.split('-').map(Number)
      return start >= min && end <= max && start <= end
    }
    if (value.includes(',')) {
      return value.split(',').every(v => !isNaN(Number(v)) && Number(v) >= min && Number(v) <= max)
    }
    const num = Number(value)
    return !isNaN(num) && num >= min && num <= max
  }

  return (
    validateField(minute, 0, 59) &&
    validateField(hour, 0, 23) &&
    validateField(day, 1, 31) &&
    validateField(month, 1, 12) &&
    validateField(weekday, 0, 6)
  )
})

// Human readable description
const humanReadable = computed(() => {
  if (!isValidCron.value || !parsedCron.value) {
    return t('cron.invalid')
  }

  const { minute, hour, day, month, weekday } = parsedCron.value

  // Simple description logic
  let time = ''
  if (minute === '0' && hour === '*') {
    time = t('cron.everyHour')
  } else if (minute === '*' && hour === '*') {
    time = t('cron.everyMinute')
  } else if (minute !== '*' && hour !== '*') {
    time = t('cron.atTime', { hour, minute: minute.padStart(2, '0') })
  } else {
    time = t('cron.custom')
  }

  let frequency = ''
  if (day === '*' && month === '*' && weekday === '*') {
    frequency = t('cron.daily')
  } else if (weekday !== '*' && day === '*') {
    frequency = t('cron.weekly')
  } else if (day !== '*' && month === '*') {
    frequency = t('cron.monthly')
  }

  return `${time}${frequency ? ', ' + frequency : ''}`
})

// Cron field descriptions
const fields = computed(() => [
  { label: t('cron.fields.minute'), value: parsedCron.value?.minute || '-', range: '0-59' },
  { label: t('cron.fields.hour'), value: parsedCron.value?.hour || '-', range: '0-23' },
  { label: t('cron.fields.day'), value: parsedCron.value?.day || '-', range: '1-31' },
  { label: t('cron.fields.month'), value: parsedCron.value?.month || '-', range: '1-12' },
  { label: t('cron.fields.weekday'), value: parsedCron.value?.weekday || '-', range: '0-6' }
])
</script>

<template>
  <div class="cron-editor">
    <!-- Input -->
    <div class="cron-input-group">
      <Clock :size="16" class="input-icon" />
      <input
        type="text"
        class="cron-input mono"
        :value="modelValue"
        @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
        :placeholder="t('cron.placeholder')"
      />
      <span class="status-indicator" :class="{ valid: isValidCron, invalid: !isValidCron }">
        <CheckCircle v-if="isValidCron" :size="14" />
        <AlertCircle v-else :size="14" />
      </span>
    </div>

    <!-- Human readable -->
    <div class="cron-description" :class="{ invalid: !isValidCron }">
      {{ humanReadable }}
    </div>

    <!-- Fields breakdown -->
    <div class="cron-fields">
      <div v-for="field in fields" :key="field.label" class="cron-field">
        <span class="field-value mono">{{ field.value }}</span>
        <span class="field-label">{{ field.label }}</span>
        <span class="field-range">{{ field.range }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.cron-editor {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

/* ========== Input ========== */
.cron-input-group {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: var(--space-3);
  color: var(--color-text-tertiary);
  pointer-events: none;
}

.cron-input {
  width: 100%;
  height: 44px;
  padding: 0 var(--space-10) 0 var(--space-10);
  font-size: var(--text-base);
  letter-spacing: 0.02em;
  background: var(--color-glass);
  border: 1px solid var(--color-glass-border);
  border-radius: var(--radius-lg);
  color: var(--color-text-primary);
  transition: all var(--duration-fast) var(--ease-default);
}

.cron-input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px var(--color-accent-light);
}

.status-indicator {
  position: absolute;
  right: var(--space-3);
  display: flex;
  align-items: center;
}

.status-indicator.valid {
  color: var(--color-success);
}

.status-indicator.invalid {
  color: var(--color-error);
}

/* ========== Description ========== */
.cron-description {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  padding: var(--space-2) var(--space-3);
  background: var(--color-success-light);
  border-radius: var(--radius-md);
  border: 1px solid rgba(34, 197, 94, 0.2);
}

.cron-description.invalid {
  background: var(--color-error-light);
  color: var(--color-error);
  border-color: rgba(239, 68, 68, 0.2);
}

/* ========== Fields ========== */
.cron-fields {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: var(--space-2);
}

.cron-field {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: var(--space-2);
  background: var(--color-glass);
  border: 1px solid var(--color-glass-border);
  border-radius: var(--radius-md);
}

.field-value {
  font-size: var(--text-lg);
  font-weight: var(--font-bold);
  color: var(--color-accent);
}

.field-label {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  margin-top: var(--space-1);
}

.field-range {
  font-size: 10px;
  color: var(--color-text-quaternary);
}

/* Responsive */
@media (max-width: 640px) {
  .cron-fields {
    grid-template-columns: repeat(3, 1fr);
  }
}
</style>

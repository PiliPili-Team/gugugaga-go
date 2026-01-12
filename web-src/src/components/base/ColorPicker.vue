<template>
  <div class="color-picker-overlay" v-if="show" @click.self="close">
    <div class="color-picker" @click.stop>
      <div class="color-picker-header">
        <h3>{{ t('header.customTheme') }}</h3>
        <button class="close-btn" @click="close">
          <X :size="18" />
        </button>
      </div>
      
      <div class="color-picker-content">
        <div class="color-group">
          <label>{{ t('theme.primaryColor') }}</label>
          <div class="color-input-group">
            <input 
              type="color" 
              v-model="colors.primary"
              @input="updateColors"
            />
            <input 
              type="text" 
              v-model="colors.primary"
              class="color-text-input"
              @input="updateColors"
            />
          </div>
        </div>
        
        <div class="color-group">
          <label>{{ t('theme.secondaryColor') }}</label>
          <div class="color-input-group">
            <input 
              type="color" 
              v-model="colors.secondary"
              @input="updateColors"
            />
            <input 
              type="text" 
              v-model="colors.secondary"
              class="color-text-input"
              @input="updateColors"
            />
          </div>
        </div>
        
        <div class="color-group">
          <label>{{ t('theme.opacity') }}</label>
          <div class="slider-group">
            <input 
              type="range" 
              min="0" 
              max="100" 
              v-model.number="opacity"
              @input="updateOpacity"
              class="opacity-slider"
            />
            <span class="opacity-value">{{ opacity }}%</span>
          </div>
        </div>
      </div>
      
      <div class="color-picker-footer">
        <button class="btn btn-secondary" @click="close">{{ t('common.cancel') }}</button>
        <button class="btn btn-primary" @click="apply">{{ t('common.apply') }}</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { X } from 'lucide-vue-next'

const props = defineProps<{
  show: boolean
}>()

const emit = defineEmits<{
  (e: 'update:show', value: boolean): void
  (e: 'apply', colors: { primary: string; secondary: string; opacity: number }): void
}>()

const { t } = useI18n()

const colors = ref({
  primary: '#3b82f6',
  secondary: '#8b5cf6'
})

const opacity = ref(100)

function close() {
  emit('update:show', false)
}

function updateColors() {
  // Colors updated
}

function updateOpacity() {
  // Opacity updated
}

function apply() {
  emit('apply', {
    primary: colors.value.primary,
    secondary: colors.value.secondary,
    opacity: opacity.value
  })
  close()
}

watch(() => props.show, (newVal) => {
  if (newVal) {
    // Load current theme colors
    const root = document.documentElement
    const primary = getComputedStyle(root).getPropertyValue('--color-accent').trim() || '#3b82f6'
    const secondary = getComputedStyle(root).getPropertyValue('--color-secondary').trim() || '#8b5cf6'
    colors.value = { primary, secondary }
  }
})
</script>

<style scoped>
.color-picker-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-4);
}

.color-picker {
  background: var(--color-bg-tertiary);
  border: 1px solid var(--color-glass-border);
  border-radius: var(--radius-xl);
  width: 100%;
  max-width: 400px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4);
}

.color-picker-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4);
  border-bottom: 1px solid var(--color-glass-border);
}

.color-picker-header h3 {
  margin: 0;
  font-size: var(--text-lg);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: transparent;
  border: none;
  color: var(--color-text-tertiary);
  cursor: pointer;
  border-radius: var(--radius-md);
  transition: all var(--duration-fast) var(--ease-default);
}

.close-btn:hover {
  color: var(--color-text-primary);
  background: var(--color-glass);
}

.color-picker-content {
  padding: var(--space-4);
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.color-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.color-group label {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--color-text-secondary);
}

.color-input-group {
  display: flex;
  gap: var(--space-2);
  align-items: center;
}

.color-input-group input[type="color"] {
  width: 60px;
  height: 44px;
  border: 1px solid var(--color-glass-border);
  border-radius: var(--radius-lg);
  cursor: pointer;
  background: transparent;
}

.color-text-input {
  flex: 1;
  height: 44px;
  padding: 0 var(--space-3);
  font-size: var(--text-sm);
  font-family: var(--font-mono);
  color: var(--color-text-primary);
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-glass-border);
  border-radius: var(--radius-lg);
}

.slider-group {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.opacity-slider {
  flex: 1;
  height: 6px;
  border-radius: var(--radius-full);
  background: var(--color-bg-secondary);
  outline: none;
  -webkit-appearance: none;
}

.opacity-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: var(--color-accent);
  cursor: pointer;
  border: 2px solid var(--color-bg-primary);
}

.opacity-slider::-moz-range-thumb {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: var(--color-accent);
  cursor: pointer;
  border: 2px solid var(--color-bg-primary);
}

.opacity-value {
  min-width: 50px;
  text-align: right;
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  font-family: var(--font-mono);
}

.color-picker-footer {
  display: flex;
  gap: var(--space-2);
  padding: var(--space-4);
  border-top: 1px solid var(--color-glass-border);
  justify-content: flex-end;
}
</style>

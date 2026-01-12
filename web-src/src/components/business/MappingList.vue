<script setup lang="ts">
/**
 * MappingList - Dynamic Mapping Rules List
 */

import { useI18n } from 'vue-i18n'
import { Plus } from 'lucide-vue-next'
import MappingItem from './MappingItem.vue'
import type { MappingRule } from '@/types'

const { t } = useI18n()

const props = defineProps<{
  modelValue: MappingRule[]
  title?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: MappingRule[]): void
}>()

function addRule() {
  emit('update:modelValue', [...props.modelValue, { regex: '', replacement: '' }])
}

function updateRegex(index: number, value: string) {
  const newRules = [...props.modelValue]
  newRules[index] = { ...newRules[index], regex: value }
  emit('update:modelValue', newRules)
}

function updateReplacement(index: number, value: string) {
  const newRules = [...props.modelValue]
  newRules[index] = { ...newRules[index], replacement: value }
  emit('update:modelValue', newRules)
}

function removeRule(index: number) {
  const newRules = props.modelValue.filter((_, i) => i !== index)
  emit('update:modelValue', newRules)
}
</script>

<template>
  <div class="mapping-list">
    <!-- Header -->
    <div v-if="title" class="list-header">
      <h4>{{ title }}</h4>
      <button class="add-btn" @click="addRule">
        <Plus :size="14" />
        <span>{{ t('common.add') }}</span>
      </button>
    </div>

    <!-- Items -->
    <div class="list-items">
      <TransitionGroup name="list">
        <MappingItem
          v-for="(rule, index) in modelValue"
          :key="index"
          :regex="rule.regex"
          :replacement="rule.replacement"
          :index="index"
          @update:regex="updateRegex(index, $event)"
          @update:replacement="updateReplacement(index, $event)"
          @remove="removeRule(index)"
        />
      </TransitionGroup>

      <!-- Empty state -->
      <div v-if="modelValue.length === 0" class="empty-state">
        <p>{{ t('mappings.empty') }}</p>
        <button class="btn btn-secondary btn-sm" @click="addRule">
          <Plus :size="14" />
          <span>{{ t('mappings.addFirst') }}</span>
        </button>
      </div>
    </div>

    <!-- Add button at bottom when items exist -->
    <div v-if="modelValue.length > 0 && !title" class="list-footer">
      <button class="btn btn-ghost btn-sm" @click="addRule">
        <Plus :size="14" />
        <span>{{ t('common.add') }}</span>
      </button>
    </div>
  </div>
</template>

<style scoped>
.mapping-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-4);
}

.list-header .add-btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: white;
  background: var(--color-accent);
  border: none;
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.list-header .add-btn:hover {
  background: var(--color-accent-hover);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
  color: white;
}

.list-header .add-btn:active {
  transform: translateY(0);
}

.list-header h4 {
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--color-text-primary);
  margin: 0;
}

.add-btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-1) var(--space-2);
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  color: var(--color-accent);
  background: var(--color-accent-light);
  border: none;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.add-btn:hover {
  background: var(--color-accent);
  color: var(--color-bg-primary);
}

.list-items {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-6);
  text-align: center;
}

.empty-state p {
  font-size: var(--text-sm);
  color: var(--color-text-quaternary);
  margin: 0;
}

.list-footer {
  display: flex;
  justify-content: center;
  margin-top: var(--space-2);
}

/* Transitions */
.list-enter-active,
.list-leave-active {
  transition: all var(--duration-normal) var(--ease-default);
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(-12px);
}

.list-move {
  transition: transform var(--duration-normal) var(--ease-default);
}
</style>

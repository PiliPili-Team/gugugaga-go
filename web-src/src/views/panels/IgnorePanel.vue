<script setup lang="ts">
/**
 * IgnorePanel - Ignored Parent IDs Configuration
 */

import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfigStore } from '@/stores'
import { EyeOff, Plus, Trash2, Save, Loader2, Edit3, Check, X } from 'lucide-vue-next'

const { t } = useI18n()
const configStore = useConfigStore()
const newParentId = ref('')
const newNote = ref('')
const isSaving = ref(false)
const editingId = ref<string | null>(null)
const editingNote = ref('')

// 从 localStorage 获取备注
const NOTES_KEY = 'gd_ignored_parent_notes'

function getNotes(): Record<string, string> {
  try {
    return JSON.parse(localStorage.getItem(NOTES_KEY) || '{}')
  } catch {
    return {}
  }
}

function saveNotes(notes: Record<string, string>) {
  localStorage.setItem(NOTES_KEY, JSON.stringify(notes))
}

const notes = ref<Record<string, string>>(getNotes())

// 带备注的ID列表
const idsWithNotes = computed(() => {
  const ids = configStore.config?.google?.ignored_parent_ids || []
  return ids.map(id => ({
    id,
    note: notes.value[id] || ''
  }))
})

async function handleSave() {
  isSaving.value = true
  try {
    await configStore.saveConfig()
  } finally {
    isSaving.value = false
  }
}

function addParentId() {
  if (!newParentId.value.trim()) return
  
  const current = configStore.config?.google?.ignored_parent_ids || []
  const idToAdd = newParentId.value.trim()
  
  if (current.includes(idToAdd)) {
    return
  }
  
  configStore.updateNested('google.ignored_parent_ids', [...current, idToAdd])
  
  // 保存备注
  if (newNote.value.trim()) {
    notes.value[idToAdd] = newNote.value.trim()
    saveNotes(notes.value)
  }
  
  newParentId.value = ''
  newNote.value = ''
}

function removeParentId(id: string) {
  const current = configStore.config?.google?.ignored_parent_ids || []
  configStore.updateNested('google.ignored_parent_ids', current.filter(i => i !== id))
  
  // 删除备注
  delete notes.value[id]
  saveNotes(notes.value)
}

function startEditNote(id: string) {
  editingId.value = id
  editingNote.value = notes.value[id] || ''
}

function saveNote() {
  if (editingId.value) {
    if (editingNote.value.trim()) {
      notes.value[editingId.value] = editingNote.value.trim()
    } else {
      delete notes.value[editingId.value]
    }
    saveNotes(notes.value)
  }
  editingId.value = null
  editingNote.value = ''
}

function cancelEditNote() {
  editingId.value = null
  editingNote.value = ''
}

function handleKeyDown(e: KeyboardEvent) {
  if (e.key === 'Enter') {
    e.preventDefault()
    addParentId()
  }
}

function handleNoteKeyDown(e: KeyboardEvent) {
  if (e.key === 'Enter') {
    e.preventDefault()
    saveNote()
  } else if (e.key === 'Escape') {
    cancelEditNote()
  }
}

onMounted(() => {
  notes.value = getNotes()
})
</script>

<template>
  <div class="panel">
    <div class="panel-content">
      <section class="config-section">
        <!-- Add new ID -->
        <div class="add-form">
          <div class="add-row">
            <input
              type="text"
              v-model="newParentId"
              class="input mono"
              :placeholder="t('panels.ignore.placeholder')"
              @keydown="handleKeyDown"
            />
            <input
              type="text"
              v-model="newNote"
              class="input note-input"
              :placeholder="t('panels.ignore.notePlaceholder')"
              @keydown="handleKeyDown"
            />
            <button class="btn btn-primary" @click="addParentId" :disabled="!newParentId.trim()">
              <Plus :size="16" />
              <span>{{ t('common.add') }}</span>
            </button>
          </div>
        </div>

        <!-- IDs list -->
        <div class="ids-list">
          <TransitionGroup name="list">
            <div
              v-for="item in idsWithNotes"
              :key="item.id"
              class="id-item"
            >
              <div class="id-info">
                <code class="id-value">{{ item.id }}</code>
                <!-- 编辑备注模式 -->
                <div v-if="editingId === item.id" class="note-edit">
                  <input
                    type="text"
                    v-model="editingNote"
                    class="input note-edit-input"
                    :placeholder="t('panels.ignore.notePlaceholder')"
                    @keydown="handleNoteKeyDown"
                    autofocus
                  />
                  <button class="note-action-btn save" @click="saveNote">
                    <Check :size="14" />
                  </button>
                  <button class="note-action-btn cancel" @click="cancelEditNote">
                    <X :size="14" />
                  </button>
                </div>
                <!-- 显示备注 -->
                <div v-else class="note-display" @click="startEditNote(item.id)">
                  <span v-if="item.note" class="note-text">{{ item.note }}</span>
                  <span v-else class="note-placeholder">{{ t('panels.ignore.addNote') }}</span>
                  <Edit3 :size="12" class="edit-icon" />
                </div>
              </div>
              <button class="remove-btn" @click="removeParentId(item.id)">
                <Trash2 :size="14" />
              </button>
            </div>
          </TransitionGroup>

          <!-- Empty state -->
          <div v-if="!idsWithNotes.length" class="empty-state">
            <EyeOff :size="32" />
            <p>{{ t('panels.ignore.empty') }}</p>
            <span class="hint">{{ t('panels.ignore.emptyHint') }}</span>
          </div>
        </div>

        <!-- Save Button -->
        <div class="section-footer">
          <button 
            class="btn btn-primary"
            @click="handleSave"
            :disabled="isSaving"
          >
            <Loader2 v-if="isSaving" :size="16" class="animate-spin" />
            <Save v-else :size="16" />
            <span>{{ isSaving ? t('common.saving') : t('common.save') }}</span>
          </button>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
@import './panel.css';

/* ========== Add Form ========== */
.add-form {
  margin-bottom: var(--space-4);
}

.add-row {
  display: flex;
  gap: var(--space-3);
  flex-wrap: wrap;
}

.add-row .input.mono {
  flex: 2;
  min-width: 200px;
}

.add-row .note-input {
  flex: 1;
  min-width: 150px;
}

/* ========== IDs List ========== */
.ids-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.id-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-3) var(--space-4);
  background: var(--color-glass);
  border: 1px solid var(--color-glass-border);
  border-radius: var(--radius-lg);
  transition: all var(--duration-fast) var(--ease-default);
  gap: var(--space-3);
}

.id-item:hover {
  border-color: var(--color-glass-border-hover);
}

.id-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  flex: 1;
  min-width: 0;
}

.id-value {
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  color: var(--color-accent);
  background: var(--color-accent-light);
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-sm);
  word-break: break-all;
  display: inline-block;
  width: fit-content;
}

/* ========== Note Display ========== */
.note-display {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  cursor: pointer;
  padding: var(--space-1) 0;
  transition: all var(--duration-fast) var(--ease-default);
}

.note-display:hover {
  color: var(--color-accent);
}

.note-display:hover .edit-icon {
  opacity: 1;
}

.note-text {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.note-placeholder {
  font-size: var(--text-xs);
  color: var(--color-text-quaternary);
  font-style: italic;
}

.edit-icon {
  color: var(--color-text-quaternary);
  opacity: 0;
  transition: opacity var(--duration-fast) var(--ease-default);
}

/* ========== Note Edit ========== */
.note-edit {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.note-edit-input {
  flex: 1;
  height: 32px;
  font-size: var(--text-sm);
  padding: var(--space-1) var(--space-2);
}

.note-action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: var(--radius-md);
  border: none;
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease-default);
}

.note-action-btn.save {
  color: var(--color-success);
  background: var(--color-success-light);
}

.note-action-btn.save:hover {
  background: var(--color-success);
  color: white;
}

.note-action-btn.cancel {
  color: var(--color-text-tertiary);
  background: var(--color-glass);
}

.note-action-btn.cancel:hover {
  color: var(--color-error);
  background: var(--color-error-light);
}

/* ========== Remove Button ========== */
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
  flex-shrink: 0;
  transition: all var(--duration-fast) var(--ease-default);
}

.remove-btn:hover {
  color: var(--color-error);
  background: var(--color-error-light);
}

/* ========== Empty State ========== */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-8);
  text-align: center;
  color: var(--color-text-quaternary);
}

.empty-state p {
  margin: 0;
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
}

.empty-state .hint {
  font-size: var(--text-xs);
}

/* ========== Transitions ========== */
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

/* ========== Responsive ========== */
@media (max-width: 640px) {
  .add-row {
    flex-direction: column;
  }
  
  .add-row .input.mono,
  .add-row .note-input {
    flex: none;
    width: 100%;
  }
}
</style>

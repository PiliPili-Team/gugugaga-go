<script setup lang="ts">
/**
 * DashboardPanel - System Status Dashboard
 */

import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { api } from '@/services/api'
import { Activity, Clock, CheckCircle, History, Cpu, HardDrive } from 'lucide-vue-next'

const { t } = useI18n()

// System status
const systemStatus = ref<{
  status: string
  uptime_seconds: number
  uptime_display: string
  start_time: string
  today_completed_tasks?: number
  history_completed_tasks?: number
  cpu_usage?: number
  memory_usage?: number
  memory_alloc_mb?: number
  memory_sys_mb?: number
  goroutines?: number
} | null>(null)

// Task statistics
const todayCompletedTasks = computed(() => systemStatus.value?.today_completed_tasks ?? 0)
const historyCompletedTasks = computed(() => systemStatus.value?.history_completed_tasks ?? 0)
const cpuUsage = computed(() => systemStatus.value?.cpu_usage ?? 0)
const memoryUsage = computed(() => systemStatus.value?.memory_usage ?? 0)

const isLoading = ref(true)
const error = ref('')
let pollInterval: number | null = null

// Fetch system status
async function fetchStatus() {
  try {
    const response = await api.fetchSystemStatus()
    systemStatus.value = response
    error.value = ''
  } catch (e: any) {
    error.value = 'Failed to fetch status'
  } finally {
    isLoading.value = false
  }
}

// Format uptime for display
const formattedUptime = computed(() => {
  if (!systemStatus.value) return '--'
  return systemStatus.value.uptime_display
})

// Status color
const statusColor = computed(() => {
  if (!systemStatus.value) return 'offline'
  return systemStatus.value.status === 'online' ? 'online' : 'offline'
})

// Format numbers for display
function formatNumber(num: number): string {
  if (num >= 1000000) {
    return (num / 1000000).toFixed(1) + 'M'
  }
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num.toString()
}

// Format percentage
function formatPercent(num: number): string {
  return num.toFixed(1) + '%'
}

const formattedTodayTasks = computed(() => formatNumber(todayCompletedTasks.value))
const formattedHistoryTasks = computed(() => formatNumber(historyCompletedTasks.value))
const formattedCpuUsage = computed(() => formatPercent(cpuUsage.value))
const formattedMemoryUsage = computed(() => formatPercent(memoryUsage.value))

// Start polling
onMounted(() => {
  fetchStatus()
  pollInterval = window.setInterval(fetchStatus, 30000) // Update every 30 seconds
})

onUnmounted(() => {
  if (pollInterval) {
    clearInterval(pollInterval)
  }
})
</script>

<template>
  <div class="panel">
    <div class="panel-content">
      <div class="dashboard-grid">
        <!-- System Status Card -->
        <div class="status-card" :class="statusColor">
          <div class="card-icon">
            <Activity :size="24" />
          </div>
          <div class="card-content">
            <span class="card-label">{{ t('dashboard.systemStatus') }}</span>
            <span class="card-value status-value">
              {{ systemStatus?.status === 'online' ? t('dashboard.online') : t('dashboard.offline') }}
            </span>
          </div>
          <div class="card-glow"></div>
        </div>

        <!-- Uptime Card -->
        <div class="status-card uptime">
          <div class="card-icon">
            <Clock :size="24" />
          </div>
          <div class="card-content">
            <span class="card-label">{{ t('dashboard.uptime') }}</span>
            <span class="card-value">{{ formattedUptime }}</span>
          </div>
        </div>

        <!-- CPU Usage Card -->
        <div class="status-card cpu-usage">
          <div class="card-icon">
            <Cpu :size="24" />
          </div>
          <div class="card-content">
            <span class="card-label">{{ t('dashboard.cpuUsage') }}</span>
            <span class="card-value">{{ formattedCpuUsage }}</span>
          </div>
        </div>

        <!-- Memory Usage Card -->
        <div class="status-card memory-usage">
          <div class="card-icon">
            <HardDrive :size="24" />
          </div>
          <div class="card-content">
            <span class="card-label">{{ t('dashboard.memoryUsage') }}</span>
            <span class="card-value">{{ formattedMemoryUsage }}</span>
          </div>
        </div>

        <!-- History Completed Tasks Card -->
        <div class="status-card history-tasks">
          <div class="card-icon">
            <History :size="24" />
          </div>
          <div class="card-content">
            <span class="card-label">{{ t('dashboard.historyCompletedTasks') }}</span>
            <span class="card-value">{{ formattedHistoryTasks }}</span>
          </div>
        </div>

        <!-- Today Completed Tasks Card -->
        <div class="status-card today-tasks">
          <div class="card-icon">
            <CheckCircle :size="24" />
          </div>
          <div class="card-content">
            <span class="card-label">{{ t('dashboard.todayCompletedTasks') }}</span>
            <span class="card-value">{{ formattedTodayTasks }}</span>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<style scoped>
@import './panel.css';

.panel-content {
  padding: var(--space-4);
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-4);
  max-width: 1000px;
}

@media (max-width: 1024px) {
  .dashboard-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .panel-content {
    padding: var(--space-3);
  }
  
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}

/* ========== Status Cards ========== */
.status-card {
  position: relative;
  display: flex;
  align-items: flex-start;
  gap: var(--space-4);
  padding: var(--space-5);
  background: var(--color-glass);
  border: 1px solid var(--color-glass-border);
  border-radius: var(--radius-xl);
  overflow: hidden;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: default;
}

.status-card::before {
  content: '';
  position: absolute;
  inset: 0;
  background: transparent;
  opacity: 0;
  transition: opacity 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.status-card:hover {
  transform: translateY(-2px);
  border-color: var(--color-glass-border-hover);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
}

.status-card:hover::before {
  opacity: 1;
}

/* Online Status - Cyan/Teal theme */
.status-card.online {
  background: 
    linear-gradient(135deg, rgba(20, 184, 166, 0.12) 0%, rgba(6, 182, 212, 0.04) 50%, rgba(20, 184, 166, 0.08) 100%),
    radial-gradient(ellipse at 0% 0%, rgba(6, 182, 212, 0.15) 0%, transparent 50%);
  border-color: rgba(20, 184, 166, 0.25);
  box-shadow: 
    inset 0 1px 0 rgba(20, 184, 166, 0.1),
    0 0 20px rgba(6, 182, 212, 0.05);
}

.status-card.online::before {
  background: 
    linear-gradient(135deg, rgba(20, 184, 166, 0.18) 0%, rgba(6, 182, 212, 0.06) 50%, rgba(20, 184, 166, 0.1) 100%),
    radial-gradient(ellipse at 0% 0%, rgba(6, 182, 212, 0.2) 0%, transparent 60%);
}

.status-card.online:hover {
  border-color: rgba(20, 184, 166, 0.4);
  box-shadow: 
    inset 0 1px 0 rgba(20, 184, 166, 0.15),
    0 8px 32px rgba(0, 0, 0, 0.2),
    0 0 30px rgba(6, 182, 212, 0.1);
}

.status-card.online .card-icon {
  color: #14b8a6;
  background: rgba(20, 184, 166, 0.15);
}

.status-card.online .status-value {
  color: #14b8a6;
}

/* Offline Status */
.status-card.offline {
  background: 
    linear-gradient(135deg, rgba(239, 68, 68, 0.12) 0%, rgba(239, 68, 68, 0.04) 50%, rgba(239, 68, 68, 0.08) 100%),
    radial-gradient(ellipse at 0% 0%, rgba(239, 68, 68, 0.15) 0%, transparent 50%);
  border-color: rgba(239, 68, 68, 0.25);
  box-shadow: 
    inset 0 1px 0 rgba(239, 68, 68, 0.1),
    0 0 20px rgba(239, 68, 68, 0.05);
}

.status-card.offline::before {
  background: 
    linear-gradient(135deg, rgba(239, 68, 68, 0.18) 0%, rgba(239, 68, 68, 0.06) 50%, rgba(239, 68, 68, 0.1) 100%),
    radial-gradient(ellipse at 0% 0%, rgba(239, 68, 68, 0.2) 0%, transparent 60%);
}

.status-card.offline:hover {
  border-color: rgba(239, 68, 68, 0.4);
  box-shadow: 
    inset 0 1px 0 rgba(239, 68, 68, 0.15),
    0 8px 32px rgba(0, 0, 0, 0.2),
    0 0 30px rgba(239, 68, 68, 0.1);
}

.status-card.offline .card-icon {
  color: var(--color-error);
  background: rgba(239, 68, 68, 0.15);
}

.status-card.offline .status-value {
  color: var(--color-error);
}

/* Uptime Card */
.status-card.uptime {
  background: 
    linear-gradient(135deg, rgba(59, 130, 246, 0.08) 0%, rgba(59, 130, 246, 0.02) 50%, rgba(59, 130, 246, 0.05) 100%),
    radial-gradient(ellipse at 0% 0%, rgba(59, 130, 246, 0.1) 0%, transparent 50%);
  border-color: rgba(59, 130, 246, 0.2);
}

.status-card.uptime .card-icon {
  color: var(--color-accent);
  background: rgba(59, 130, 246, 0.15);
}

.status-card.uptime::before {
  background: 
    linear-gradient(135deg, rgba(59, 130, 246, 0.12) 0%, rgba(59, 130, 246, 0.04) 50%, rgba(59, 130, 246, 0.08) 100%),
    radial-gradient(ellipse at 0% 0%, rgba(59, 130, 246, 0.15) 0%, transparent 60%);
}

.status-card.uptime:hover {
  border-color: rgba(59, 130, 246, 0.35);
  box-shadow: 
    inset 0 1px 0 rgba(59, 130, 246, 0.1),
    0 8px 32px rgba(0, 0, 0, 0.2),
    0 0 25px rgba(59, 130, 246, 0.08);
}

/* History Completed Tasks Card - Purple theme */
.status-card.history-tasks {
  background: 
    linear-gradient(135deg, rgba(139, 92, 246, 0.12) 0%, rgba(167, 139, 250, 0.04) 50%, rgba(139, 92, 246, 0.08) 100%),
    radial-gradient(ellipse at 0% 0%, rgba(167, 139, 250, 0.1) 0%, transparent 50%);
  border-color: rgba(139, 92, 246, 0.25);
}

.status-card.history-tasks .card-icon {
  color: #8b5cf6;
  background: rgba(139, 92, 246, 0.15);
}

.status-card.history-tasks .card-value {
  color: #8b5cf6;
}

.status-card.history-tasks::before {
  background: 
    linear-gradient(135deg, rgba(139, 92, 246, 0.18) 0%, rgba(167, 139, 250, 0.06) 50%, rgba(139, 92, 246, 0.1) 100%),
    radial-gradient(ellipse at 0% 0%, rgba(167, 139, 250, 0.15) 0%, transparent 60%);
}

.status-card.history-tasks:hover {
  border-color: rgba(139, 92, 246, 0.4);
  box-shadow: 
    inset 0 1px 0 rgba(139, 92, 246, 0.1),
    0 8px 32px rgba(0, 0, 0, 0.2),
    0 0 25px rgba(167, 139, 250, 0.08);
}

/* Today Completed Tasks Card - Green theme */
.status-card.today-tasks {
  background: 
    linear-gradient(135deg, rgba(34, 197, 94, 0.12) 0%, rgba(16, 185, 129, 0.04) 50%, rgba(34, 197, 94, 0.08) 100%),
    radial-gradient(ellipse at 0% 0%, rgba(16, 185, 129, 0.1) 0%, transparent 50%);
  border-color: rgba(34, 197, 94, 0.25);
}

.status-card.today-tasks .card-icon {
  color: #22c55e;
  background: rgba(34, 197, 94, 0.15);
}

.status-card.today-tasks .card-value {
  color: #22c55e;
}

.status-card.today-tasks::before {
  background: 
    linear-gradient(135deg, rgba(34, 197, 94, 0.18) 0%, rgba(16, 185, 129, 0.06) 50%, rgba(34, 197, 94, 0.1) 100%),
    radial-gradient(ellipse at 0% 0%, rgba(16, 185, 129, 0.15) 0%, transparent 60%);
}

.status-card.today-tasks:hover {
  border-color: rgba(34, 197, 94, 0.4);
  box-shadow: 
    inset 0 1px 0 rgba(34, 197, 94, 0.1),
    0 8px 32px rgba(0, 0, 0, 0.2),
    0 0 25px rgba(16, 185, 129, 0.08);
}

/* CPU Usage Card - Orange/Yellow theme */
.status-card.cpu-usage {
  background: 
    linear-gradient(135deg, rgba(245, 158, 11, 0.12) 0%, rgba(251, 191, 36, 0.04) 50%, rgba(245, 158, 11, 0.08) 100%),
    radial-gradient(ellipse at 0% 0%, rgba(251, 191, 36, 0.1) 0%, transparent 50%);
  border-color: rgba(245, 158, 11, 0.25);
}

.status-card.cpu-usage .card-icon {
  color: #f59e0b;
  background: rgba(245, 158, 11, 0.15);
}

.status-card.cpu-usage .card-value {
  color: #f59e0b;
}

.status-card.cpu-usage::before {
  background: 
    linear-gradient(135deg, rgba(245, 158, 11, 0.18) 0%, rgba(251, 191, 36, 0.06) 50%, rgba(245, 158, 11, 0.1) 100%),
    radial-gradient(ellipse at 0% 0%, rgba(251, 191, 36, 0.15) 0%, transparent 60%);
}

.status-card.cpu-usage:hover {
  border-color: rgba(245, 158, 11, 0.4);
  box-shadow: 
    inset 0 1px 0 rgba(245, 158, 11, 0.1),
    0 8px 32px rgba(0, 0, 0, 0.2),
    0 0 25px rgba(251, 191, 36, 0.08);
}

/* Memory Usage Card - Pink theme */
.status-card.memory-usage {
  background: 
    linear-gradient(135deg, rgba(236, 72, 153, 0.12) 0%, rgba(244, 114, 182, 0.04) 50%, rgba(236, 72, 153, 0.08) 100%),
    radial-gradient(ellipse at 0% 0%, rgba(244, 114, 182, 0.1) 0%, transparent 50%);
  border-color: rgba(236, 72, 153, 0.25);
}

.status-card.memory-usage .card-icon {
  color: #ec4899;
  background: rgba(236, 72, 153, 0.15);
}

.status-card.memory-usage .card-value {
  color: #ec4899;
}

.status-card.memory-usage::before {
  background: 
    linear-gradient(135deg, rgba(236, 72, 153, 0.18) 0%, rgba(244, 114, 182, 0.06) 50%, rgba(236, 72, 153, 0.1) 100%),
    radial-gradient(ellipse at 0% 0%, rgba(244, 114, 182, 0.15) 0%, transparent 60%);
}

.status-card.memory-usage:hover {
  border-color: rgba(236, 72, 153, 0.4);
  box-shadow: 
    inset 0 1px 0 rgba(236, 72, 153, 0.1),
    0 8px 32px rgba(0, 0, 0, 0.2),
    0 0 25px rgba(244, 114, 182, 0.08);
}

/* Card Components */
.card-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: var(--radius-lg);
  flex-shrink: 0;
}

.card-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  flex: 1;
  min-width: 0;
}

.card-label {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
  font-weight: var(--font-medium);
}

.card-value {
  font-size: var(--text-2xl);
  font-weight: var(--font-bold);
  color: var(--color-accent);
  line-height: var(--leading-tight);
}

.card-glow {
  position: absolute;
  top: -50%;
  right: -50%;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, currentColor 0%, transparent 70%);
  opacity: 0.1;
  pointer-events: none;
}

.status-card.online .card-glow {
  color: #14b8a6;
}

.status-card.offline .card-glow {
  color: var(--color-error);
}

/* ========== Responsive ========== */
@media (max-width: 768px) {
  .status-card {
    padding: var(--space-4);
  }
  
  .card-icon {
    width: 40px;
    height: 40px;
  }
  
  .card-value {
    font-size: var(--text-xl);
  }
}
</style>

<!-- Light Theme Styles (non-scoped) -->
<style>
[data-theme="light"] .status-card {
  background: rgba(255, 255, 255, 0.8) !important;
  border-color: rgba(0, 0, 0, 0.1) !important;
}

[data-theme="light"] .status-card:hover {
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1) !important;
}

[data-theme="light"] .status-card.online {
  background: 
    linear-gradient(135deg, rgba(20, 184, 166, 0.15) 0%, rgba(6, 182, 212, 0.08) 50%, rgba(20, 184, 166, 0.1) 100%),
    rgba(255, 255, 255, 0.85) !important;
  border-color: rgba(20, 184, 166, 0.3) !important;
}

[data-theme="light"] .status-card.uptime {
  background: 
    linear-gradient(135deg, rgba(59, 130, 246, 0.1) 0%, rgba(59, 130, 246, 0.05) 50%, rgba(59, 130, 246, 0.08) 100%),
    rgba(255, 255, 255, 0.85) !important;
  border-color: rgba(59, 130, 246, 0.25) !important;
}

[data-theme="light"] .status-card.history-tasks {
  background: 
    linear-gradient(135deg, rgba(139, 92, 246, 0.15) 0%, rgba(167, 139, 250, 0.08) 50%, rgba(139, 92, 246, 0.1) 100%),
    rgba(255, 255, 255, 0.85) !important;
  border-color: rgba(139, 92, 246, 0.3) !important;
}

[data-theme="light"] .status-card.today-tasks {
  background: 
    linear-gradient(135deg, rgba(34, 197, 94, 0.15) 0%, rgba(16, 185, 129, 0.08) 50%, rgba(34, 197, 94, 0.1) 100%),
    rgba(255, 255, 255, 0.85) !important;
  border-color: rgba(34, 197, 94, 0.3) !important;
}

[data-theme="light"] .status-card.cpu-usage {
  background: 
    linear-gradient(135deg, rgba(245, 158, 11, 0.15) 0%, rgba(251, 191, 36, 0.08) 50%, rgba(245, 158, 11, 0.1) 100%),
    rgba(255, 255, 255, 0.85) !important;
  border-color: rgba(245, 158, 11, 0.3) !important;
}

[data-theme="light"] .status-card.memory-usage {
  background: 
    linear-gradient(135deg, rgba(236, 72, 153, 0.15) 0%, rgba(244, 114, 182, 0.08) 50%, rgba(236, 72, 153, 0.1) 100%),
    rgba(255, 255, 255, 0.85) !important;
  border-color: rgba(236, 72, 153, 0.3) !important;
}

[data-theme="light"] .card-label {
  color: #64748b !important;
}

[data-theme="light"] .card-value {
  color: #1e293b !important;
}
</style>

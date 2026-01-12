import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/services/api'

interface ParsedLog {
  id: string
  time: string
  type: 'info' | 'warn' | 'error' | 'debug'
  content: string
}

export const useLogsStore = defineStore('logs', () => {
  // State
  const rawLogs = ref<string[]>([])
  const autoScroll = ref(true)
  const isLoading = ref(false)
  
  // Computed
  const count = computed(() => rawLogs.value.length)
  
  const parsedLogs = computed<ParsedLog[]>(() => {
    return rawLogs.value.map((log, index) => {
      // 首先尝试匹配格式: [TIME] LEVEL: MESSAGE
      const bracketMatch = log.match(/^\[([^\]]+)\]\s*(\w+)?:?\s*(.*)$/)
      
      if (bracketMatch) {
        const [, time, level, content] = bracketMatch
        let type: ParsedLog['type'] = 'info'
        
        if (level) {
          const lowerLevel = level.toLowerCase()
          if (lowerLevel.includes('error') || lowerLevel.includes('err')) {
            type = 'error'
          } else if (lowerLevel.includes('warn')) {
            type = 'warn'
          } else if (lowerLevel.includes('debug') || lowerLevel.includes('dbg')) {
            type = 'debug'
          }
        }
        
        // Also check content for error indicators
        if (content.toLowerCase().includes('error') || content.toLowerCase().includes('failed')) {
          type = 'error'
        }
        
        return {
          id: `log-${index}-${time}`,
          time,
          type,
          content: content || log
        }
      }
      
      return {
        id: `log-${index}`,
        time: new Date().toLocaleTimeString(),
        type: 'info' as const,
        content: log
      }
    })
  })
  
  // Actions
  async function fetchLogs() {
    try {
      const response = await api.fetchLogs()
      if (response.logs) {
        rawLogs.value = response.logs
      }
    } catch (e) {
      // Silently handle error
    }
  }
  
  async function clearMemLogs() {
    try {
      await api.clearMemLogs()
      rawLogs.value = []
    } catch (e) {
      // Silently handle error
    }
  }
  
  async function clearFileLogs() {
    try {
      await api.clearFiles()
    } catch (e) {
      // Silently handle error
    }
  }
  
  return {
    rawLogs,
    autoScroll,
    isLoading,
    count,
    parsedLogs,
    fetchLogs,
    clearMemLogs,
    clearFileLogs
  }
})

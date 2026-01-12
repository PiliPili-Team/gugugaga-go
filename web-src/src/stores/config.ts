import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '@/services/api'
import type { Config } from '@/types'
import { adaptBackendConfig, adaptFrontendConfig } from '@/utils/configAdapter'

export const useConfigStore = defineStore('config', () => {
  // State
  const config = ref<Config | null>(null)
  const isLoading = ref(false)
  const isSaving = ref(false)
  const error = ref<string | null>(null)
  
  // Computed
  const isLoaded = computed(() => config.value !== null)
  
  // Actions
  async function fetchConfig() {
    isLoading.value = true
    error.value = null
    
    try {
      let backendData = await api.fetchConfig()
      
      // Handle string response (if backend returns JSON string)
      if (typeof backendData === 'string') {
        backendData = JSON.parse(backendData)
      }
      
      // Convert backend format to frontend format
      config.value = adaptBackendConfig(backendData as any)
    } catch (e: any) {
      error.value = e.message || 'Failed to fetch config'
    } finally {
      isLoading.value = false
    }
  }
  
  async function saveConfig() {
    if (!config.value) return
    
    isSaving.value = true
    error.value = null
    
    try {
      // Convert frontend format to backend format
      const backendData = adaptFrontendConfig(config.value)
      await api.updateConfig(backendData as any)
    } catch (e: any) {
      error.value = e.message || 'Failed to save config'
      throw e
    } finally {
      isSaving.value = false
    }
  }
  
  // Update a top-level config field
  function update<K extends keyof Config>(key: K, value: Config[K]) {
    if (config.value) {
      config.value = { ...config.value, [key]: value }
    }
  }
  
  // Update nested config field using dot notation
  function updateNested(path: string, value: any) {
    if (!config.value) return
    
    const keys = path.split('.')
    const newConfig = JSON.parse(JSON.stringify(config.value))
    
    let obj: any = newConfig
    for (let i = 0; i < keys.length - 1; i++) {
      const key = keys[i]
      if (obj[key] === undefined) {
        obj[key] = {}
      }
      obj = obj[key]
    }
    
    obj[keys[keys.length - 1]] = value
    config.value = newConfig
  }
  
  return {
    config,
    isLoading,
    isSaving,
    error,
    isLoaded,
    fetchConfig,
    saveConfig,
    update,
    updateNested
  }
})

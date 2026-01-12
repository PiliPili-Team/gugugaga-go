/**
 * API Service - Centralized API calls
 */

import { ofetch } from 'ofetch'
import type { 
  Config, 
  LogsResponse, 
  OAuthLoginURLResponse,
  LoginRequest,
  LoginResponse,
  BingWallpaperResponse,
  TestSymediaRequest
} from '@/types'

// Create fetch instance with defaults
const apiFetch = ofetch.create({
  baseURL: '/api',
  credentials: 'include',
  headers: {
    'Content-Type': 'application/json'
  },
  onResponseError({ response }) {
    // Handle 401 unauthorized
    if (response.status === 401) {
      localStorage.removeItem('auth_token')
      localStorage.removeItem('user_name')
      window.location.reload()
    }
  }
})

export const api = {
  // ========== Auth ==========
  
  async login(data: LoginRequest): Promise<LoginResponse> {
    try {
      const response = await apiFetch<{ status?: string; success?: boolean }>('/auth/login', {
        method: 'POST',
        body: data
      })
      
      // Handle different response formats from backend
      // Backend returns {"status":"ok"} on success
      if (response.status === 'ok' || response.success === true) {
        // Store basic auth credentials for subsequent requests
        const credentials = btoa(`${data.username}:${data.password}`)
        localStorage.setItem('basic_auth', credentials)
        return { success: true }
      }
      
      return { success: false, message: 'Authentication failed' }
    } catch (e: any) {
      // Fallback: Try with basic auth header directly
      const credentials = btoa(`${data.username}:${data.password}`)
      try {
        await apiFetch('/config', {
          headers: {
            'Authorization': `Basic ${credentials}`
          }
        })
        // If we get here, auth succeeded
        localStorage.setItem('basic_auth', credentials)
        return { success: true }
      } catch {
        return { success: false, message: 'Authentication failed' }
      }
    }
  },
  
  // ========== Wallpapers ==========
  
  async fetchBingWallpaper(): Promise<BingWallpaperResponse> {
    try {
      return await apiFetch<BingWallpaperResponse>('/bing/wallpaper')
    } catch {
      return { url: '', copyright: '', source: 'bing' }
    }
  },
  
  async fetchTMDBWallpaper(): Promise<BingWallpaperResponse> {
    try {
      return await apiFetch<BingWallpaperResponse>('/tmdb/wallpaper')
    } catch {
      return { url: '', copyright: '', source: 'tmdb' }
    }
  },
  
  // Randomly fetch wallpaper from Bing or TMDB
  async fetchRandomWallpaper(): Promise<BingWallpaperResponse> {
    const usesTMDB = Math.random() > 0.5
    try {
      if (usesTMDB) {
        const result = await this.fetchTMDBWallpaper()
        if (result.url) return result
        // Fallback to Bing if TMDB fails
        return await this.fetchBingWallpaper()
      } else {
        const result = await this.fetchBingWallpaper()
        if (result.url) return result
        // Fallback to TMDB if Bing fails
        return await this.fetchTMDBWallpaper()
      }
    } catch {
      return { url: '', copyright: '', source: '' }
    }
  },
  
  // ========== System Status ==========
  
  async fetchSystemStatus(): Promise<{
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
  }> {
    try {
      return await apiFetch('/status')
    } catch {
      return {
        status: 'offline',
        uptime_seconds: 0,
        uptime_display: '--',
        start_time: '',
        today_completed_tasks: 0,
        history_completed_tasks: 0,
        cpu_usage: 0,
        memory_usage: 0
      }
    }
  },
  
  // ========== Config ==========
  
  async fetchConfig(): Promise<Config> {
    const auth = localStorage.getItem('basic_auth')
    const headers: Record<string, string> = {}
    if (auth) {
      headers['Authorization'] = `Basic ${auth}`
    }
    
    return await apiFetch<Config>('/config/get', { headers })
  },
  
  async updateConfig(config: Config): Promise<void> {
    const auth = localStorage.getItem('basic_auth')
    const headers: Record<string, string> = {}
    if (auth) {
      headers['Authorization'] = `Basic ${auth}`
    }
    await apiFetch('/config/update', {
      method: 'POST',
      body: config,
      headers
    })
  },
  
  // ========== Logs ==========
  
  async fetchLogs(): Promise<LogsResponse> {
    const auth = localStorage.getItem('basic_auth')
    const headers: Record<string, string> = {}
    if (auth) {
      headers['Authorization'] = `Basic ${auth}`
    }
    return await apiFetch<LogsResponse>('/logs', { headers })
  },
  
  async clearMemLogs(): Promise<void> {
    const auth = localStorage.getItem('basic_auth')
    const headers: Record<string, string> = {}
    if (auth) {
      headers['Authorization'] = `Basic ${auth}`
    }
    await apiFetch('/logs/clear_mem', { method: 'POST', headers })
  },
  
  async clearFiles(): Promise<void> {
    const auth = localStorage.getItem('basic_auth')
    const headers: Record<string, string> = {}
    if (auth) {
      headers['Authorization'] = `Basic ${auth}`
    }
    await apiFetch('/logs/clear-files', { method: 'POST', headers })
  },
  
  // ========== OAuth ==========
  
  async getOAuthLoginURL(): Promise<OAuthLoginURLResponse> {
    const auth = localStorage.getItem('basic_auth')
    const headers: Record<string, string> = {}
    if (auth) {
      headers['Authorization'] = `Basic ${auth}`
    }
    return await apiFetch<OAuthLoginURLResponse>('/auth/login_url', { headers })
  },
  
  // ========== Actions ==========
  
  async triggerSync(): Promise<void> {
    const auth = localStorage.getItem('basic_auth')
    const headers: Record<string, string> = {}
    if (auth) {
      headers['Authorization'] = `Basic ${auth}`
    }
    await apiFetch('/trigger', { method: 'POST', headers })
  },
  
  async triggerRcloneFull(): Promise<void> {
    const auth = localStorage.getItem('basic_auth')
    const headers: Record<string, string> = {}
    if (auth) {
      headers['Authorization'] = `Basic ${auth}`
    }
    await apiFetch('/rclone_full', { method: 'POST', headers })
  },
  
  async testSymedia(data: TestSymediaRequest): Promise<void> {
    const auth = localStorage.getItem('basic_auth')
    const headers: Record<string, string> = {}
    if (auth) {
      headers['Authorization'] = `Basic ${auth}`
    }
    await apiFetch('/test_symedia', { 
      method: 'POST', 
      body: data,
      headers
    })
  },
  
  async refreshTree(): Promise<void> {
    const auth = localStorage.getItem('basic_auth')
    const headers: Record<string, string> = {}
    if (auth) {
      headers['Authorization'] = `Basic ${auth}`
    }
    await apiFetch('/tree/refresh', { method: 'POST', headers })
  }
}

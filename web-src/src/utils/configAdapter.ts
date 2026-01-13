/**
 * Config Adapter - Convert backend config format to frontend format
 */

import type { Config } from '@/types'

// Backend config format (from Go)
interface BackendConfig {
  auth: {
    username: string
    password: string
  }
  oauth_config?: {
    client_id: string
    client_secret: string
    redirect_uri: string
  }
  advanced: {
    log_level: number
    log_save_enabled: boolean
    log_dir: string
    log_max_size_mb?: number
    debounce_seconds: number
    rclone_wait_seconds?: number
    log_cleanup_enabled?: boolean
    log_retention_days?: number
    log_cleanup_cron?: string
  }
  server: {
    listen_port: number
    public_url: string
    webhook_path: string
    ssl?: {
      enabled: boolean
      cert_path: string
      key_path: string
      restrict_to_domain?: boolean
    }
  }
  google: {
    rate_limit_qps: number
    personal_drive_name?: string
    my_drive_name?: string
    target_drive_ids?: string[]
    list_delay?: number
    batch_sleep_interval?: number
    ignored_parents?: string[]
  }
  rclone?: Array<{
    name?: string
    host: string
    endpoint: string
    mapping?: Array<{
      regex: string
      replacement: string
    }>
  }>
  symedia?: {
    host: string
    endpoint: string
    notify_unmatched?: boolean
    headers?: Record<string, string>
    body_template?: any
  }
  path_mapping?: Array<{
    regex: string
    replacement: string
  }>
}

/**
 * Convert backend config to frontend config format
 */
export function adaptBackendConfig(backend: BackendConfig): Config {
  return {
    auth: {
      username: backend.auth?.username ?? '',
      password: backend.auth?.password ?? ''
    },
    oauth: {
      client_id: backend.oauth_config?.client_id ?? '',
      client_secret: backend.oauth_config?.client_secret ?? '',
      redirect_uri: backend.oauth_config?.redirect_uri ?? ''
    },
    advanced: {
      debounce_seconds: backend.advanced?.debounce_seconds ?? 5,
      log_dir: backend.advanced?.log_dir ?? './logs',
      log_level: backend.advanced?.log_level ?? 1,
      log_save_enabled: backend.advanced?.log_save_enabled !== false,
      log_cleanup: backend.advanced?.log_cleanup_enabled ? {
        enabled: true,
        retention_days: backend.advanced?.log_retention_days ?? 7,
        cron: backend.advanced?.log_cleanup_cron ?? '0 0 3 * * ?'
      } : {
        enabled: false,
        retention_days: 7,
        cron: '0 0 3 * * ?'
      }
    },
    server: {
      port: backend.server?.listen_port ?? 8448,
      public_url: backend.server?.public_url ?? '',
      webhook_path: backend.server?.webhook_path ?? '/gd-webhook',
      ssl_enabled: backend.server?.ssl?.enabled ?? false,
      ssl_cert: backend.server?.ssl?.cert_path ?? '',
      ssl_key: backend.server?.ssl?.key_path ?? ''
    },
    google: {
      qps: backend.google?.rate_limit_qps ?? 5,
      personal_drive_name: backend.google?.personal_drive_name || backend.google?.my_drive_name || '',
      target_drive_ids: backend.google?.target_drive_ids || [],
      list_delay: backend.google?.list_delay ?? 1000,
      batch_sleep_interval: backend.google?.batch_sleep_interval ?? 300
    },
    rclone: {
      instances: (backend.rclone || []).map(instance => ({
        host: instance.host || '',
        endpoint: instance.endpoint || '',
        wait_for_data: true
      })),
      path_mappings: (backend.rclone || [])
        .flatMap(instance => instance.mapping || [])
        .map(m => ({ regex: m.regex || '', replacement: m.replacement || '' }))
    },
    symedia: {
      host: backend.symedia?.host ?? '',
      endpoint: backend.symedia?.endpoint ?? '',
      body_template: backend.symedia?.body_template
        ? (typeof backend.symedia.body_template === 'string'
          ? backend.symedia.body_template
          : JSON.stringify(backend.symedia.body_template))
        : '',
      path_mappings: (backend.path_mapping || []).map(m => ({
        regex: m.regex ?? '',
        replacement: m.replacement ?? ''
      })),
      notify_unmatched: backend.symedia?.notify_unmatched ?? false,
      headers: backend.symedia?.headers || {}
    }
  }
}

/**
 * Convert frontend config to backend config format
 */
export function adaptFrontendConfig(frontend: Config): BackendConfig {
  return {
    auth: {
      username: frontend.auth.username,
      password: frontend.auth.password
    },
    oauth_config: {
      client_id: frontend.oauth.client_id,
      client_secret: frontend.oauth.client_secret,
      redirect_uri: frontend.oauth.redirect_uri
    },
    advanced: {
      log_level: frontend.advanced.log_level,
      log_save_enabled: frontend.advanced.log_save_enabled,
      log_dir: frontend.advanced.log_dir,
      debounce_seconds: frontend.advanced.debounce_seconds,
      log_cleanup_enabled: frontend.advanced.log_cleanup?.enabled || false,
      log_retention_days: frontend.advanced.log_cleanup?.retention_days || 7,
      log_cleanup_cron: frontend.advanced.log_cleanup?.cron || '0 0 3 * * ?'
    },
    server: {
      listen_port: frontend.server.port,
      public_url: frontend.server.public_url,
      webhook_path: frontend.server.webhook_path,
      ssl: {
        enabled: frontend.server.ssl_enabled || false,
        cert_path: frontend.server.ssl_cert || '',
        key_path: frontend.server.ssl_key || '',
        restrict_to_domain: false
      }
    },
    google: {
      rate_limit_qps: frontend.google.qps,
      personal_drive_name: frontend.google.personal_drive_name,
      target_drive_ids: frontend.google.target_drive_ids || [],
      list_delay: frontend.google.list_delay || 1000,
      batch_sleep_interval: frontend.google.batch_sleep_interval || 300
    },
    rclone: frontend.rclone.instances.map((instance, index) => ({
      name: `instance_${index}`,
      host: instance.host,
      endpoint: instance.endpoint,
      mapping: frontend.rclone.path_mappings
    })),
    symedia: {
      host: frontend.symedia.host,
      endpoint: frontend.symedia.endpoint,
      notify_unmatched: frontend.symedia.notify_unmatched,
      body_template: frontend.symedia.body_template
        ? JSON.parse(frontend.symedia.body_template)
        : {},
      headers: frontend.symedia.headers || {}
    },
    path_mapping: frontend.symedia.path_mappings
  }
}

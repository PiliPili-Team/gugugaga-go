/**
 * Config Types - Backend configuration interfaces
 */

export interface MappingRule {
  regex: string
  replacement: string
}

export interface RcloneInstance {
  host: string
  endpoint: string
  wait_for_data?: boolean
  timeout?: number
}

export interface AuthConfig {
  username: string
  password: string
}

export interface OAuthConfig {
  client_id: string
  client_secret: string
  redirect_uri: string
}

export interface AdvancedConfig {
  debounce_seconds: number
  log_dir: string
  log_level: number
  log_save_enabled: boolean
  log_cleanup?: {
    enabled: boolean
    cron: string
    retention_days: number
  }
}

export interface ServerConfig {
  port: number
  public_url: string
  webhook_path: string
  ssl_enabled?: boolean
  ssl_cert?: string
  ssl_key?: string
}

export interface GoogleConfig {
  qps: number
  personal_drive_name: string
  ignored_parent_ids: string[]
}

export interface RcloneConfig {
  instances: RcloneInstance[]
  path_mappings: MappingRule[]
}

export interface SymediaConfig {
  host: string
  endpoint: string
  body_template: string
  path_mappings: MappingRule[]
  notify_unmatched?: boolean
  headers?: Record<string, string>
  timeout?: number
}

export interface Config {
  auth: AuthConfig
  oauth: OAuthConfig
  advanced: AdvancedConfig
  server: ServerConfig
  google: GoogleConfig
  rclone: RcloneConfig
  symedia: SymediaConfig
}

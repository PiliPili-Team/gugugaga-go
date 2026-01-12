/**
 * API Types - Request and Response interfaces
 */

export interface LogsResponse {
  logs: string[]
  count: number
}

export interface OAuthLoginURLResponse {
  url: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  success: boolean
  message?: string
  token?: string
}

export interface BingWallpaperResponse {
  url: string
  copyright?: string
  title?: string
  source?: string
}

export interface TestSymediaRequest {
  path: string
}

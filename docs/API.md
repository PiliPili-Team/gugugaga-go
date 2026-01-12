# GD Watcher API Documentation

[中文版](API_CN.md)

This document describes the REST API endpoints provided by GD Watcher.

## Base URL

```
http://localhost:8448
```

## Authentication

Most API endpoints require authentication. After logging in, a session cookie `gd_session` is set.

### Login

```http
POST /api/login
Content-Type: application/json

{
  "username": "admin",
  "password": "your-password"
}
```

**Response:**
```json
{
  "status": "ok"
}
```

**Cookies Set:**
- `gd_session` - Session token (HttpOnly, 24h expiry)

### Logout

```http
POST /api/logout
```

**Response:**
```json
{
  "status": "ok"
}
```

---

## Configuration

### Get Configuration

Retrieve the current system configuration.

```http
GET /api/config
```

**Response:**
```json
{
  "auth": {
    "username": "admin",
    "password": "***"
  },
  "oauth_config": {
    "client_id": "xxx.apps.googleusercontent.com",
    "client_secret": "***",
    "redirect_uri": "https://example.com/oauth/callback"
  },
  "advanced": {
    "log_level": 1,
    "log_save_enabled": true,
    "log_dir": "userdata/logs",
    "log_max_size_mb": 10,
    "debounce_seconds": 5,
    "rclone_wait_seconds": 2,
    "log_cleanup_enabled": true,
    "log_retention_days": 7,
    "log_cleanup_cron": "0 0 3 * * ?"
  },
  "server": {
    "listen_port": 8448,
    "public_url": "https://example.com",
    "webhook_path": "/gd-webhook",
    "ssl": {
      "enabled": false,
      "cert_path": "",
      "key_path": "",
      "restrict_to_domain": false
    }
  },
  "google": {
    "rate_limit_qps": 5,
    "personal_drive_name": "My Drive",
    "ignored_parents": []
  },
  "rclone": [
    {
      "name": "MyRclone",
      "host": "http://localhost:5572",
      "endpoint": "/vfs/refresh",
      "mapping": [
        {
          "regex": "^/My Drive/(.*)$",
          "replacement": "/$1"
        }
      ]
    }
  ],
  "symedia": {
    "host": "http://localhost:8096",
    "endpoint": "/emby/Library/Media/Updated",
    "notify_unmatched": false,
    "headers": {
      "X-Emby-Token": "your-api-key"
    },
    "body_template": {
      "Updates": [
        {
          "Path": "{{path}}",
          "UpdateType": "{{action}}"
        }
      ]
    }
  },
  "path_mapping": [
    {
      "regex": "^/My Drive/(.*)$",
      "replacement": "/mnt/media/$1"
    }
  ]
}
```

### Update Configuration

```http
POST /api/config
Content-Type: application/json

{
  // Full or partial configuration object
  // Empty password field preserves existing password
}
```

**Response:**
```
ok
```

---

## System Status

### Get System Status

```http
GET /api/status
```

**Response:**
```json
{
  "status": "online",
  "uptime_seconds": 86400,
  "uptime_display": "1d 0h 0m",
  "start_time": "2024-01-01T00:00:00Z",
  "app_name": "GD Watcher",
  "app_version": "4.0",
  "today_completed_tasks": 5,
  "history_completed_tasks": 42,
  "cpu_usage": 15.5,
  "memory_usage": 32.8,
  "memory_alloc_mb": 24.5,
  "memory_sys_mb": 74.7,
  "goroutines": 12
}
```

| Field | Type | Description |
|-------|------|-------------|
| `status` | string | System status ("online" or "offline") |
| `uptime_seconds` | int | Server uptime in seconds |
| `uptime_display` | string | Human-readable uptime |
| `start_time` | string | Server start time (RFC3339) |
| `app_name` | string | Application name |
| `app_version` | string | Application version |
| `today_completed_tasks` | int | Tasks completed today |
| `history_completed_tasks` | int | Total tasks completed historically |
| `cpu_usage` | float | Estimated CPU usage percentage |
| `memory_usage` | float | Memory usage percentage |
| `memory_alloc_mb` | float | Allocated memory in MB |
| `memory_sys_mb` | float | System memory in MB |
| `goroutines` | int | Number of active goroutines |

---

## Logs

### Get Real-time Logs

Retrieve logs from memory buffer with optional pagination.

```http
GET /api/logs?since=0
```

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `since` | int | Log index to start from (for incremental updates) |

**Response:**
```json
{
  "logs": [
    "2024/01/01 12:00:00 main.go:26: 🚀 GD Watcher v4.0 启动中...",
    "2024/01/01 12:00:01 main.go:51: ⏰ 日志清理计划任务已设定: [0 0 3 * * ?]"
  ],
  "next_idx": 2
}
```

### Clear Memory Logs

```http
POST /api/logs/clear/mem
```

**Response:**
```
ok
```

### Clear Log Files

Delete all log files except the current one.

```http
POST /api/logs/clear/files
```

**Response:**
```
ok
```

---

## Operations

### Trigger Manual Sync

Manually trigger a synchronization cycle.

```http
POST /api/trigger
```

**Response:**
```
ok
```

### Force Rclone Full Refresh

Trigger a full VFS refresh for all Rclone instances.

```http
POST /api/rclone/full
```

**Response:**
```
ok
```

### Refresh File Tree

Force rebuild the file tree cache from Google Drive.

```http
POST /api/tree/refresh
```

**Response:**
```json
{
  "status": "ok",
  "message": "Tree refresh started"
}
```

### Test Symedia Webhook

Send a test webhook to the configured Symedia endpoint.

```http
POST /api/symedia/test
Content-Type: application/json

{
  "path": "/test/path/file.mkv"
}
```

**Response:**
```
ok
```

---

## OAuth

### Get OAuth Login URL

Generate Google OAuth authorization URL.

```http
GET /oauth/url
```

**Response:**
```json
{
  "url": "https://accounts.google.com/o/oauth2/auth?client_id=xxx&redirect_uri=xxx&..."
}
```

### OAuth Callback

Handle OAuth callback from Google. This endpoint is called by Google after user authorization.

```http
GET /oauth/callback?code=xxx&state=xxx
```

**Response:** Redirects to `/` on success.

---

## Webhook

### Google Drive Webhook Endpoint

This endpoint receives push notifications from Google Drive.

```http
POST /{webhook_path}
```

**Headers (from Google):**
| Header | Description |
|--------|-------------|
| `X-Goog-Resource-State` | Resource state (sync, add, remove, update, trash, untrash, change) |
| `X-Goog-Resource-Id` | Resource identifier |
| `X-Goog-Channel-Id` | Channel identifier |
| `X-Goog-Message-Number` | Message sequence number |

**Response:** HTTP 200 OK

---

## Wallpapers

### Get Bing Daily Wallpaper

```http
GET /api/wallpaper/bing
```

**Response:**
```json
{
  "url": "https://www.bing.com/th?id=OHR.xxx_1920x1080.jpg",
  "copyright": "Photo description © Photographer",
  "title": "Photo Title",
  "source": "bing"
}
```

### Get TMDB Trending Backdrop

```http
GET /api/wallpaper/tmdb
```

**Response:**
```json
{
  "url": "https://image.tmdb.org/t/p/original/xxx.jpg",
  "copyright": "Movie/TV Show Title",
  "title": "Movie/TV Show Title",
  "source": "tmdb"
}
```

---

## Error Responses

### HTTP Status Codes

| Code | Description |
|------|-------------|
| 200 | Success |
| 400 | Bad Request - Invalid input |
| 401 | Unauthorized - Not logged in |
| 405 | Method Not Allowed |
| 500 | Internal Server Error |

### Error Response Format

```json
{
  "error": "Error message description"
}
```

Or plain text:
```
Error message
```

---

## Integration Examples

### Rclone VFS Refresh

GD Watcher sends requests to Rclone RC API:

```http
POST http://localhost:5572/vfs/refresh?_async=true
Content-Type: application/json

{
  "dir": "/path/to/refresh",
  "recursive": "true"
}
```

### Symedia Webhook (Emby Example)

```http
POST http://localhost:8096/emby/Library/Media/Updated
Content-Type: application/json
X-Emby-Token: your-api-key

{
  "Updates": [
    {
      "Path": "/mnt/media/Movies/Movie.mkv",
      "UpdateType": "Created"
    }
  ]
}
```

---

## Rate Limiting

- Google Drive API calls are rate-limited based on `google.rate_limit_qps` config
- Rclone refresh requests are limited to 5 concurrent requests
- Symedia webhooks have no artificial rate limiting

---

## WebSocket (Future)

WebSocket support for real-time log streaming is planned for future versions.

---

## Frontend Development

For frontend developers who want to build a custom UI:

1. All API endpoints are prefixed with `/api/` (except OAuth and webhook)
2. Authentication is cookie-based (`gd_session`)
3. CORS is not enabled by default (same-origin only)
4. Content-Type for POST requests should be `application/json`

### Minimal Frontend Integration

```javascript
// Login
const login = async (username, password) => {
  const res = await fetch('/api/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
    credentials: 'include'
  });
  return res.ok;
};

// Get status
const getStatus = async () => {
  const res = await fetch('/api/status', { credentials: 'include' });
  return res.json();
};

// Get logs with polling
const pollLogs = async (since = 0) => {
  const res = await fetch(`/api/logs?since=${since}`, { credentials: 'include' });
  return res.json();
};

// Trigger sync
const triggerSync = async () => {
  await fetch('/api/trigger', { method: 'POST', credentials: 'include' });
};
```

# GD Watcher API 文档

[English](API.md)

本文档描述 GD Watcher 提供的 REST API 接口。

## 基础 URL

```
http://localhost:8448
```

## 认证

大多数 API 接口需要认证。登录后会设置会话 cookie `gd_session`。

### 登录

```http
POST /api/login
Content-Type: application/json

{
  "username": "admin",
  "password": "your-password"
}
```

**响应：**
```json
{
  "status": "ok"
}
```

**设置的 Cookie：**
- `gd_session` - 会话令牌（HttpOnly，24小时过期）

### 登出

```http
POST /api/logout
```

**响应：**
```json
{
  "status": "ok"
}
```

---

## 配置

### 获取配置

获取当前系统配置。

```http
GET /api/config
```

**响应：**
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
  "rclone": [...],
  "symedia": {...},
  "path_mapping": [...]
}
```

### 更新配置

```http
POST /api/config
Content-Type: application/json

{
  // 完整或部分配置对象
  // 密码字段为空则保留现有密码
}
```

**响应：**
```
ok
```

---

## 系统状态

### 获取系统状态

```http
GET /api/status
```

**响应：**
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

| 字段 | 类型 | 说明 |
|------|------|------|
| `status` | string | 系统状态 ("online" 或 "offline") |
| `uptime_seconds` | int | 服务器运行时间（秒） |
| `uptime_display` | string | 人类可读的运行时间 |
| `start_time` | string | 服务器启动时间 (RFC3339) |
| `app_name` | string | 应用名称 |
| `app_version` | string | 应用版本 |
| `today_completed_tasks` | int | 今日完成任务数 |
| `history_completed_tasks` | int | 历史完成任务总数 |
| `cpu_usage` | float | 估计 CPU 使用率（百分比） |
| `memory_usage` | float | 内存使用率（百分比） |
| `memory_alloc_mb` | float | 已分配内存（MB） |
| `memory_sys_mb` | float | 系统内存（MB） |
| `goroutines` | int | 活跃的 goroutine 数量 |

---

## 日志

### 获取实时日志

从内存缓冲区获取日志，支持分页。

```http
GET /api/logs?since=0
```

**查询参数：**
| 参数 | 类型 | 描述 |
|------|------|------|
| `since` | int | 开始的日志索引（用于增量更新） |

**响应：**
```json
{
  "logs": [
    "2024/01/01 12:00:00 main.go:26: 🚀 GD Watcher v4.0 starting...",
    "2024/01/01 12:00:01 main.go:51: ⏰ Log cleanup cron scheduled: [0 0 3 * * ?]"
  ],
  "next_idx": 2
}
```

### 清空内存日志

```http
POST /api/logs/clear/mem
```

### 清空日志文件

删除除当前日志外的所有日志文件。

```http
POST /api/logs/clear/files
```

---

## 操作

### 手动触发同步

手动触发同步周期。

```http
POST /api/trigger
```

### 强制 Rclone 全量刷新

为所有 Rclone 实例触发完整 VFS 刷新。

```http
POST /api/rclone/full
```

### 刷新文件树

强制从 Google Drive 重建文件树缓存。

```http
POST /api/tree/refresh
```

**响应：**
```json
{
  "status": "ok",
  "message": "Tree refresh started"
}
```

### 测试 Symedia Webhook

向配置的 Symedia 端点发送测试 webhook。

```http
POST /api/symedia/test
Content-Type: application/json

{
  "path": "/test/path/file.mkv"
}
```

---

## OAuth

### 获取 OAuth 登录 URL

生成 Google OAuth 授权 URL。

```http
GET /oauth/url
```

**响应：**
```json
{
  "url": "https://accounts.google.com/o/oauth2/auth?client_id=xxx&redirect_uri=xxx&..."
}
```

### OAuth 回调

处理来自 Google 的 OAuth 回调。此端点在用户授权后由 Google 调用。

```http
GET /oauth/callback?code=xxx&state=xxx
```

**响应：** 成功后重定向到 `/`

---

## Webhook

### Google Drive Webhook 端点

此端点接收来自 Google Drive 的推送通知。

```http
POST /{webhook_path}
```

**Headers（来自 Google）：**
| Header | 描述 |
|--------|------|
| `X-Goog-Resource-State` | 资源状态（sync, add, remove, update, trash, untrash, change） |
| `X-Goog-Resource-Id` | 资源标识符 |
| `X-Goog-Channel-Id` | 频道标识符 |
| `X-Goog-Message-Number` | 消息序列号 |

**响应：** HTTP 200 OK

---

## 壁纸

### 获取 Bing 每日壁纸

```http
GET /api/wallpaper/bing
```

### 获取 TMDB 热门背景图

```http
GET /api/wallpaper/tmdb
```

---

## 错误响应

### HTTP 状态码

| 状态码 | 描述 |
|--------|------|
| 200 | 成功 |
| 400 | 请求错误 - 无效输入 |
| 401 | 未授权 - 未登录 |
| 405 | 方法不允许 |
| 500 | 服务器内部错误 |

---

## 集成示例

### Rclone VFS 刷新

GD Watcher 向 Rclone RC API 发送请求：

```http
POST http://localhost:5572/vfs/refresh?_async=true
Content-Type: application/json

{
  "dir": "/path/to/refresh",
  "recursive": "true"
}
```

### Symedia Webhook（Emby 示例）

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

## 速率限制

- Google Drive API 调用根据 `google.rate_limit_qps` 配置进行限速
- Rclone 刷新请求限制为最多 5 个并发
- Symedia webhook 无人为速率限制

---

## 前端开发

对于想构建自定义 UI 的前端开发者：

1. 所有 API 端点以 `/api/` 为前缀（OAuth 和 webhook 除外）
2. 认证基于 cookie（`gd_session`）
3. 默认未启用 CORS（仅同源）
4. POST 请求的 Content-Type 应为 `application/json`

### 最小前端集成示例

```javascript
// 登录
const login = async (username, password) => {
  const res = await fetch('/api/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
    credentials: 'include'
  });
  return res.ok;
};

// 获取状态
const getStatus = async () => {
  const res = await fetch('/api/status', { credentials: 'include' });
  return res.json();
};

// 轮询日志
const pollLogs = async (since = 0) => {
  const res = await fetch(`/api/logs?since=${since}`, { credentials: 'include' });
  return res.json();
};

// 触发同步
const triggerSync = async () => {
  await fetch('/api/trigger', { method: 'POST', credentials: 'include' });
};
```

# 系统架构

[English](ARCHITECTURE.md)

## 概述

GD Watcher 是一个实时 Google Drive 文件变更监控与自动同步系统。它使用 Google Drive 推送通知（Webhook）来即时检测文件变更，并触发下游服务如 Rclone 和媒体服务器。

## 架构图

```mermaid
flowchart TB
    subgraph Google["Google Cloud"]
        GD[Google Drive]
        GD -->|推送通知| WH[Webhook]
    end

    subgraph Server["GD Watcher 服务器"]
        WH -->|HTTP POST| Handler[Webhook 处理器]
        Handler -->|触发| Sync[同步服务]
        Sync -->|获取变更| DriveAPI[Drive API 客户端]
        DriveAPI -->|列出变更| GD
        
        Sync -->|更新| FileTree[文件树缓存]
        Sync -->|通知| Rclone[Rclone 服务]
        Sync -->|通知| Symedia[Symedia 服务]
    end

    subgraph Integrations["外部服务"]
        Rclone -->|VFS 刷新| RcloneRC[Rclone RC API]
        Symedia -->|Webhook| MediaServer[Emby/Jellyfin/Plex]
    end

    subgraph UI["Web 界面"]
        WebUI[Vue 3 仪表盘]
        WebUI -->|REST API| Handler
    end
```

## 组件说明

### 核心组件

| 组件 | 描述 |
|------|------|
| **Webhook 处理器** | 接收 Google Drive 推送通知 |
| **同步服务** | 协调同步流程，管理防抖 |
| **Drive API 客户端** | 与 Google Drive API 交互获取文件变更 |
| **文件树缓存** | 文件/文件夹结构的内存缓存 |
| **Rclone 服务** | 触发 Rclone VFS 刷新已挂载的网盘 |
| **Symedia 服务** | 向媒体服务器（Emby/Jellyfin/Plex）发送 webhook |

### 数据流

1. **变更检测**：Google Drive 向 webhook 端点发送推送通知
2. **防抖处理**：同步服务等待可配置的防抖时间以批量处理变更
3. **获取变更**：调用 Changes API 获取已修改文件列表
4. **树更新**：使用新增/修改/删除的文件更新文件树缓存
5. **路径映射**：使用正则规则转换文件路径
6. **通知**：并行通知 Rclone 和 Symedia 服务

## 工作流程

```mermaid
sequenceDiagram
    participant GD as Google Drive
    participant WH as GD Watcher
    participant RC as Rclone
    participant MS as 媒体服务器

    Note over GD,MS: 初始化设置
    WH->>GD: 注册 Webhook 频道
    GD-->>WH: 频道 ID 和资源 ID

    Note over GD,MS: 文件变更检测
    GD->>WH: 推送通知 (X-Goog-Resource-State: change)
    WH->>GD: 列出变更 (changes.list API)
    GD-->>WH: 变更文件列表
    
    Note over GD,MS: 同步流程
    WH->>WH: 更新文件树缓存
    WH->>WH: 应用路径映射规则
    
    par 并行通知
        WH->>RC: VFS 刷新 (/vfs/refresh?_async=true)
        RC-->>WH: 任务 ID
    and
        WH->>MS: Symedia Webhook
        MS-->>WH: OK
    end
```

## 关键特性

### Webhook 自动续期
- Google Drive webhook 7 天后过期
- GD Watcher 每 6 天自动续期

### 速率限制
- Google Drive API 可配置 QPS 限制
- Rclone 并发请求限制（最大 5 个）

### 缓存
- 文件树缓存到磁盘以加速启动
- 增量更新减少 API 调用

### 路径映射
- 基于正则的路径转换
- Rclone 和 Symedia 独立的映射规则

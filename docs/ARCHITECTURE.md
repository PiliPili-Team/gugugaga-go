# System Architecture

[中文版](ARCHITECTURE_CN.md)

## Overview

GD Watcher is a real-time Google Drive file change monitoring and auto-sync system. It uses Google Drive Push Notifications (Webhook) to detect file changes instantly and triggers downstream services like Rclone and media servers.

## Architecture Diagram

```mermaid
flowchart TB
    subgraph Google["Google Cloud"]
        GD[Google Drive]
        GD -->|Push Notification| WH[Webhook]
    end

    subgraph Server["GD Watcher Server"]
        WH -->|HTTP POST| Handler[Webhook Handler]
        Handler -->|Trigger| Sync[Sync Service]
        Sync -->|Fetch Changes| DriveAPI[Drive API Client]
        DriveAPI -->|List Changes| GD
        
        Sync -->|Update| FileTree[File Tree Cache]
        Sync -->|Notify| Rclone[Rclone Service]
        Sync -->|Notify| Symedia[Symedia Service]
    end

    subgraph Integrations["External Services"]
        Rclone -->|VFS Refresh| RcloneRC[Rclone RC API]
        Symedia -->|Webhook| MediaServer[Emby]
    end

    subgraph UI["Web Interface"]
        WebUI[Vue 3 Dashboard]
        WebUI -->|REST API| Handler
    end
```

## Component Description

### Core Components

| Component | Description |
|-----------|-------------|
| **Webhook Handler** | Receives push notifications from Google Drive |
| **Sync Service** | Orchestrates the sync process, manages debouncing |
| **Drive API Client** | Interfaces with Google Drive API for file changes |
| **File Tree Cache** | In-memory cache of file/folder structure |
| **Rclone Service** | Triggers Rclone VFS refresh for mounted drives |
| **Symedia Service** | Sends webhooks to media servers (Emby) |

### Data Flow

1. **Change Detection**: Google Drive sends push notification to webhook endpoint
2. **Debouncing**: Sync service waits for configurable debounce period to batch changes
3. **Change Fetching**: Changes API is called to get list of modified files
4. **Tree Update**: File tree cache is updated with new/modified/deleted files
5. **Path Mapping**: File paths are transformed using regex rules
6. **Notification**: Rclone and Symedia services are notified in parallel

## System Flow

```mermaid
sequenceDiagram
    participant GD as Google Drive
    participant WH as GD Watcher
    participant RC as Rclone
    participant MS as Media Server

    Note over GD,MS: Initial Setup
    WH->>GD: Register Webhook Channel
    GD-->>WH: Channel ID & Resource ID

    Note over GD,MS: File Change Detection
    GD->>WH: Push Notification (X-Goog-Resource-State: change)
    WH->>GD: List Changes (changes.list API)
    GD-->>WH: Changed Files List
    
    Note over GD,MS: Sync Process
    WH->>WH: Update File Tree Cache
    WH->>WH: Apply Path Mappings
    
    par Parallel Notifications
        WH->>RC: VFS Refresh (/vfs/refresh?_async=true)
        RC-->>WH: Job ID
    and
        WH->>MS: Symedia Webhook
        MS-->>WH: OK
    end
```

## Key Features

### Webhook Auto-Renewal
- Google Drive webhooks expire after 7 days
- GD Watcher automatically renews every 6 days

### Rate Limiting
- Configurable QPS limit for Google Drive API
- Concurrent request limiting for Rclone (max 5)

### Caching
- File tree is cached to disk for fast startup
- Incremental updates reduce API calls

### Path Mapping
- Regex-based path transformation
- Separate mapping rules for Rclone and Symedia

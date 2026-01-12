# GD Watcher

<div align="center">

<img src="web-src/public/pwa-icon.svg" alt="GD Watcher Logo" width="120" height="120">

**Real-time Google Drive file change monitoring and auto-sync system**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![Vue Version](https://img.shields.io/badge/Vue-3.4+-4FC08D?style=flat-square&logo=vue.js)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)](LICENSE)

[English](#features) | [中文文档](README_CN.md)

</div>

---

## Features

- 🔔 **Real-time Monitoring** - Uses Google Drive Push Notifications (Webhook) for instant file change detection
- 🔄 **Rclone Integration** - Automatically triggers Rclone VFS refresh when files change
- 📺 **Media Server Support** - Sends notifications to Emby via Symedia webhook
- 🌳 **Smart File Tree** - Caches and incrementally updates the file tree structure
- 🎨 **Modern Web UI** - Beautiful glassmorphism design with Vue 3 + TypeScript
- 📱 **PWA Support** - Installable on mobile devices with native-like experience
- 🌍 **Multi-language** - Supports English, Simplified Chinese
- 🌓 **Theme Support** - Light/Dark/System appearance modes

## Screenshots

### Web Interface

<p align="center">
  <img src="docs/img/web/01.png" width="45%" alt="Web Screenshot 1">
  <img src="docs/img/web/02.png" width="45%" alt="Web Screenshot 2">
</p>
<p align="center">
  <img src="docs/img/web/03.png" width="45%" alt="Web Screenshot 3">
  <img src="docs/img/web/04.png" width="45%" alt="Web Screenshot 4">
</p>

### Mobile PWA

<p align="center">
  <img src="docs/img/mobile/01.jpg" width="22%" alt="Mobile Screenshot 1">
  <img src="docs/img/mobile/02.jpg" width="22%" alt="Mobile Screenshot 2">
  <img src="docs/img/mobile/03.jpg" width="22%" alt="Mobile Screenshot 3">
  <img src="docs/img/mobile/04.jpg" width="22%" alt="Mobile Screenshot 4">
</p>

## Documentation

| Document | Description |
|----------|-------------|
| [Architecture](docs/ARCHITECTURE.md) | System architecture and workflow diagrams |
| [API Reference](docs/API.md) | Complete REST API documentation |
| [Development Guide](docs/DEVELOPMENT.md) | Setup, building, and contribution guide |

## Quick Start

### Prerequisites

| Software | Version | Purpose |
|----------|---------|---------|
| **Go** | 1.21+ | Backend compilation |
| **Node.js** | 18+ | Frontend build |
| **npm** | 9+ | Package management |

### Google Cloud Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable **Google Drive API**
4. Create OAuth 2.0 credentials:
   - Application type: **Web application**
   - Authorized redirect URIs: `https://your-domain.com/oauth/callback`
5. Download the credentials JSON

> 📖 Detailed guide: [Google Drive API Quickstart](https://developers.google.com/drive/api/quickstart/go)

## Installation

### Option 1: Docker

> ⚠️ **Note**: Docker deployment has not been fully tested. Please use manual build or install script for production environments.

```bash
# Clone repository
git clone https://github.com/PiliPili-Team/gugugaga-go.git
cd gugugaga-go

# Start with Docker Compose
docker-compose up -d
```

### Option 2: Build from Source (Recommended)

```bash
# Clone repository
git clone https://github.com/PiliPili-Team/gugugaga-go.git
cd gugugaga-go

# Build frontend
cd web-src
npm install
npm run build
cd ..

# Build backend
go build -o gd-webhook-server ./src

# Run
./gd-webhook-server
```

### Option 3: Quick Install Script (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/PiliPili-Team/gugugaga-go/main/onekey-install.sh | bash
```

## Configuration

Access the Web UI at `http://localhost:8448` after starting the server.

### First-time Setup

1. **Login** with your configured credentials (set in `userdata/config/config.json`)
   - Default username: `admin`
   - Password: Set your own password in config file
   
2. **OAuth Configuration**: Enter your Google OAuth credentials
3. **Authorize**: Click "Authorize" to complete Google Drive authorization
4. **Configure Integrations**: Set up Rclone and/or Symedia as needed

### Configuration File

The configuration is stored in `userdata/config/config.json`，template file reference [`config.json.exmaple`](./config.json.exmaple):

```json
{
  "auth": {
    "username": "admin",
    "password": "password"
  },
  "oauth_config": {
    "client_id": "",
    "client_secret": "",
    "redirect_uri": ""
  },
  "advanced": {
    "log_level": 2,
    "log_save_enabled": true,
    "log_dir": "logs",
    "log_max_size_mb": 10,
    "debounce_seconds": 5,
    "rclone_wait_seconds": 5,
    "log_cleanup_enabled": false,
    "log_retention_days": 7,
    "log_cleanup_cron": "0 0 3 * * ?"
  },
  "server": {
    "listen_port": 8448,
    "public_url": "",
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
    "my_drive_name": "",
    "ignored_parents": [
    ]
  },
  "rclone": [
    {
      "name": "anime",
      "host": "http://127.0.0.1:5173",
      "endpoint": "/vfs/refresh",
      "mapping": [
        {
          "regex": "",
          "replacement": ""
        }
      ]
    }
  ],
  "symedia": {
    "host": "http://127.0.0.1:8095",
    "endpoint": "/api/v1/webhook/clouddrive2/file_notify",
    "notify_unmatched": false,
    "headers": {
      "content-type": "application/json",
      "user-agent": "clouddrive2/0.9.8",
      "authorization": "basic usernamepassword"
    },
    "body_template": {
      "data": [
        {
          "action": "{{ACTION}}",
          "destination_file": "",
          "is_dir": "{{IS_DIR}}",
          "source_file": "{{FILE_PATH}}"
        }
      ],
      "device_name": "Manual-Test",
      "event_category": "file",
      "event_name": "notify",
      "type": "notify",
      "user_name": "admin",
      "version": "0.9.8"
    }
  },
  "path_mapping": [
    {
      "regex": "",
      "replacement": ""
    }
  ]
}
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_NAME` | GD Watcher | Application name |
| `APP_VERSION` | 4.0 | Application version |
| `TZ` | UTC | Timezone |

## External References

### Google APIs

- [Google Drive API Overview](https://developers.google.com/drive/api/guides/about-sdk)
- [Push Notifications (Webhook)](https://developers.google.com/drive/api/guides/push)
- [Changes API](https://developers.google.com/drive/api/reference/rest/v3/changes)
- [OAuth 2.0 for Web Server Applications](https://developers.google.com/identity/protocols/oauth2/web-server)

### Rclone

- [Rclone Documentation](https://rclone.org/docs/)
- [Rclone RC API](https://rclone.org/rc/)
- [VFS Commands](https://rclone.org/commands/rclone_rc_vfs_refresh/)

### Media Servers

- [Emby API Documentation](https://github.com/MediaBrowser/Emby/wiki/Api-Documentation)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Google Drive API](https://developers.google.com/drive)
- [Rclone](https://rclone.org/)
- [Vue.js](https://vuejs.org/)
- [Lucide Icons](https://lucide.dev/)

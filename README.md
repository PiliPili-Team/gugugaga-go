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
- 📺 **Media Server Support** - Sends notifications to Emby/Jellyfin/Plex via Symedia webhook
- 🌳 **Smart File Tree** - Caches and incrementally updates the file tree structure
- 🎨 **Modern Web UI** - Beautiful glassmorphism design with Vue 3 + TypeScript
- 📱 **PWA Support** - Installable on mobile devices with native-like experience
- 🌍 **Multi-language** - Supports English, Simplified Chinese, Traditional Chinese
- 🌓 **Theme Support** - Light/Dark/System appearance modes

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
git clone https://github.com/yourusername/gd-webhook.git
cd gd-webhook

# Start with Docker Compose
docker-compose up -d
```

### Option 2: Build from Source (Recommended)

```bash
# Clone repository
git clone https://github.com/yourusername/gd-webhook.git
cd gd-webhook

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
curl -fsSL https://raw.githubusercontent.com/yourusername/gd-webhook/main/install.sh | bash
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

The configuration is stored in `userdata/config/config.json`:

```json
{
  "auth": {
    "username": "admin",
    "password": "your-secure-password"
  },
  "server": {
    "listen_port": 8448,
    "public_url": "https://your-domain.com",
    "webhook_path": "/gd-webhook"
  },
  "rclone": [
    {
      "name": "MyRclone",
      "host": "http://localhost:5572",
      "endpoint": "/vfs/refresh",
      "mapping": [...]
    }
  ],
  "symedia": {
    "host": "http://localhost:8096",
    "endpoint": "/emby/Library/Media/Updated",
    "headers": {
      "X-Emby-Token": "your-api-key"
    }
  }
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
- [Jellyfin API](https://api.jellyfin.org/)
- [Plex Webhooks](https://support.plex.tv/articles/115002267687-webhooks/)

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

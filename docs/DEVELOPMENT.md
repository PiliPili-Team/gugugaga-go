# Development Guide

[中文版](DEVELOPMENT_CN.md)

## Prerequisites

### Go Installation

**macOS (Homebrew):**
```bash
brew install go
```

**Ubuntu/Debian:**
```bash
sudo apt update
sudo apt install golang-go
```

**Windows:**
Download installer from: https://golang.org/dl/

**Verify installation:**
```bash
go version
# Output: go version go1.21.x ...
```

### Node.js Installation

**Recommended: Use nvm (Node Version Manager):**

```bash
# Install nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

# Reload shell
source ~/.bashrc  # or ~/.zshrc

# Install Node.js
nvm install 18
nvm use 18

# Verify
node --version  # v18.x.x
npm --version   # 9.x.x
```

**Or direct install:**

- macOS: `brew install node`
- Ubuntu: `sudo apt install nodejs npm`
- Windows: https://nodejs.org/

## Project Structure

```
gd-webhook/
├── src/                    # Go backend source
│   ├── config/            # Configuration management
│   │   └── config.go      # Config loading, saving, env vars
│   ├── logger/            # Logging utilities
│   │   └── logger.go      # File & memory logging
│   ├── model/             # Data models
│   │   └── model.go       # Config, FileNode, API types
│   ├── server/            # HTTP handlers & middleware
│   │   ├── handler.go     # API endpoint handlers
│   │   ├── middleware.go  # Auth, session management
│   │   └── server.go      # HTTP server setup
│   ├── service/           # Business logic
│   │   ├── drive_service.go   # Google Drive API wrapper
│   │   ├── file_tree.go       # In-memory file tree cache
│   │   ├── rclone_service.go  # Rclone VFS refresh
│   │   ├── symedia_service.go # Media server webhooks
│   │   └── sync_service.go    # Sync orchestration
│   ├── web/               # Static file embedding
│   │   └── assets.go      # embed.FS for static files
│   └── main.go            # Application entry point
├── web-src/               # Vue 3 frontend source
│   ├── src/
│   │   ├── components/    # Vue components
│   │   │   ├── base/      # Base UI components
│   │   │   ├── business/  # Business components
│   │   │   └── layout/    # Layout components
│   │   ├── views/         # Page views
│   │   │   └── panels/    # Settings panels
│   │   ├── stores/        # Pinia state management
│   │   ├── services/      # API client
│   │   ├── i18n/          # Internationalization
│   │   ├── types/         # TypeScript types
│   │   └── utils/         # Utility functions
│   ├── public/            # Static assets (icons, manifest)
│   ├── index.html         # HTML entry point
│   ├── vite.config.ts     # Vite configuration
│   └── package.json       # npm dependencies
├── userdata/              # Runtime data (gitignored)
│   ├── config/            # config.json, token.json
│   ├── data/              # tree_cache.json, start_token.txt
│   └── logs/              # Log files
├── docs/                  # Documentation
├── Dockerfile             # Multi-stage Docker build
├── docker-compose.yml     # Docker Compose config
├── build.sh              # Build script
└── install.sh            # Installation script
```

## Development Workflow

### Frontend Development

```bash
cd web-src

# Install dependencies
npm install

# Start development server (hot reload)
npm run dev
# Access at http://localhost:5173

# Type check
npm run type-check

# Build for production
npm run build
```

### Backend Development

```bash
# Run directly
go run ./src

# Build binary
go build -o gd-webhook-server ./src

# With hot reload (using air)
go install github.com/cosmtrek/air@latest
air
```

### Full Build

```bash
# Using build script
./build.sh

# Or manually:
cd web-src && npm run build && cd ..
go build -ldflags="-s -w" -o gd-webhook-server ./src
```

## Code Style

### Go

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable names
- Add comments for exported functions
- Handle errors explicitly

### TypeScript/Vue

- Use TypeScript for type safety
- Follow Vue 3 Composition API patterns
- Use Pinia for state management
- Keep components focused and reusable

## Adding New Features

### Adding a New API Endpoint

1. Add handler function in `src/server/handler.go`:
```go
func (h *Handler) HandleNewEndpoint(w http.ResponseWriter, r *http.Request) {
    // Implementation
}
```

2. Register route in `src/server/server.go`:
```go
mux.HandleFunc("/api/new-endpoint", s.Handler.HandleNewEndpoint)
```

3. Add frontend API call in `web-src/src/services/api.ts`:
```typescript
export async function callNewEndpoint() {
  return await $fetch('/api/new-endpoint')
}
```

### Adding a New Configuration Option

1. Add field to `Config` struct in `src/model/model.go`
2. Add default value in `src/config/config.go` LoadConfig()
3. Add TypeScript type in `web-src/src/types/config.ts`
4. Add UI in appropriate panel component
5. Add i18n keys in `web-src/src/i18n/locales/`

## Testing

### Manual Testing

1. Start the server: `go run ./src`
2. Open browser: `http://localhost:8448`
3. Login and configure OAuth
4. Test file changes in Google Drive

### Debug Mode

Set `log_level` to `2` in config for debug output:
```json
{
  "advanced": {
    "log_level": 2
  }
}
```

## Building for Production

```bash
# Optimized build with stripped symbols
go build -ldflags="-s -w" -o gd-webhook-server ./src

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o gd-webhook-server ./src
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `APP_NAME` | GD Watcher | Application name displayed in UI |
| `APP_VERSION` | 4.0 | Application version |
| `TZ` | UTC | Timezone for logs |

## Troubleshooting

### Common Issues

**OAuth Error:**
- Check credentials.json is valid
- Verify redirect URI matches config

**Webhook Not Receiving:**
- Ensure public URL is accessible
- Check webhook path matches config

**Rclone Refresh Failing:**
- Verify Rclone RC is enabled and accessible
- Check path mapping rules

### Debug Logging

Enable debug mode to see detailed logs:
- Request/response details for Rclone
- Google Drive API calls
- Path resolution steps

# 开发指南

[English](DEVELOPMENT.md)

## 环境要求

### Go 安装

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
下载安装包: https://golang.org/dl/

**验证安装:**
```bash
go version
# 输出: go version go1.21.x ...
```

### Node.js 安装

**推荐使用 nvm (Node Version Manager):**

```bash
# 安装 nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

# 重新加载 shell
source ~/.bashrc  # 或 ~/.zshrc

# 安装 Node.js
nvm install 18
nvm use 18

# 验证
node --version  # v18.x.x
npm --version   # 9.x.x
```

**或者直接安装:**

- macOS: `brew install node`
- Ubuntu: `sudo apt install nodejs npm`
- Windows: https://nodejs.org/

## 项目结构

```
gd-webhook/
├── src/                    # Go 后端源码
│   ├── config/            # 配置管理
│   │   └── config.go      # 配置加载、保存、环境变量
│   ├── logger/            # 日志工具
│   │   └── logger.go      # 文件和内存日志
│   ├── model/             # 数据模型
│   │   └── model.go       # Config、FileNode、API 类型
│   ├── server/            # HTTP 处理器和中间件
│   │   ├── handler.go     # API 端点处理器
│   │   ├── middleware.go  # 认证、会话管理
│   │   └── server.go      # HTTP 服务器设置
│   ├── service/           # 业务逻辑
│   │   ├── drive_service.go   # Google Drive API 封装
│   │   ├── file_tree.go       # 内存文件树缓存
│   │   ├── rclone_service.go  # Rclone VFS 刷新
│   │   ├── symedia_service.go # 媒体服务器 webhook
│   │   └── sync_service.go    # 同步编排
│   ├── web/               # 静态文件嵌入
│   │   └── assets.go      # embed.FS 静态文件
│   └── main.go            # 应用入口点
├── web-src/               # Vue 3 前端源码
│   ├── src/
│   │   ├── components/    # Vue 组件
│   │   │   ├── base/      # 基础 UI 组件
│   │   │   ├── business/  # 业务组件
│   │   │   └── layout/    # 布局组件
│   │   ├── views/         # 页面视图
│   │   │   └── panels/    # 设置面板
│   │   ├── stores/        # Pinia 状态管理
│   │   ├── services/      # API 客户端
│   │   ├── i18n/          # 国际化
│   │   ├── types/         # TypeScript 类型
│   │   └── utils/         # 工具函数
│   ├── public/            # 静态资源（图标、manifest）
│   ├── index.html         # HTML 入口
│   ├── vite.config.ts     # Vite 配置
│   └── package.json       # npm 依赖
├── userdata/              # 运行时数据（已忽略）
│   ├── config/            # config.json, token.json
│   ├── data/              # tree_cache.json, start_token.txt
│   └── logs/              # 日志文件
├── docs/                  # 文档
├── Dockerfile             # 多阶段 Docker 构建
├── docker-compose.yml     # Docker Compose 配置
├── build.sh              # 构建脚本
└── install.sh            # 安装脚本
```

## 开发流程

### 前端开发

```bash
cd web-src

# 安装依赖
npm install

# 启动开发服务器（热重载）
npm run dev
# 访问 http://localhost:5173

# 类型检查
npm run type-check

# 生产构建
npm run build
```

### 后端开发

```bash
# 直接运行
go run ./src

# 编译二进制
go build -o gd-webhook-server ./src

# 热重载（使用 air）
go install github.com/cosmtrek/air@latest
air
```

### 完整构建

```bash
# 使用构建脚本
./build.sh

# 或手动：
cd web-src && npm run build && cd ..
go build -ldflags="-s -w" -o gd-webhook-server ./src
```

## 代码风格

### Go

- 遵循标准 Go 格式化（`gofmt`）
- 使用有意义的变量名
- 为导出函数添加注释
- 显式处理错误

### TypeScript/Vue

- 使用 TypeScript 保证类型安全
- 遵循 Vue 3 Composition API 模式
- 使用 Pinia 进行状态管理
- 保持组件专注和可复用

## 添加新功能

### 添加新 API 端点

1. 在 `src/server/handler.go` 添加处理函数：
```go
func (h *Handler) HandleNewEndpoint(w http.ResponseWriter, r *http.Request) {
    // 实现
}
```

2. 在 `src/server/server.go` 注册路由：
```go
mux.HandleFunc("/api/new-endpoint", s.Handler.HandleNewEndpoint)
```

3. 在 `web-src/src/services/api.ts` 添加前端 API 调用：
```typescript
export async function callNewEndpoint() {
  return await $fetch('/api/new-endpoint')
}
```

### 添加新配置选项

1. 在 `src/model/model.go` 的 `Config` 结构体添加字段
2. 在 `src/config/config.go` 的 LoadConfig() 添加默认值
3. 在 `web-src/src/types/config.ts` 添加 TypeScript 类型
4. 在适当的面板组件添加 UI
5. 在 `web-src/src/i18n/locales/` 添加 i18n 键

## 测试

### 手动测试

1. 启动服务器：`go run ./src`
2. 打开浏览器：`http://localhost:8448`
3. 登录并配置 OAuth
4. 在 Google Drive 中测试文件变更

### 调试模式

在配置中设置 `log_level` 为 `2` 以输出调试信息：
```json
{
  "advanced": {
    "log_level": 2
  }
}
```

## 生产构建

```bash
# 优化构建，去除符号
go build -ldflags="-s -w" -o gd-webhook-server ./src

# 交叉编译 Linux 版本
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o gd-webhook-server ./src
```

## 环境变量

| 变量 | 默认值 | 描述 |
|------|--------|------|
| `APP_NAME` | GD Watcher | UI 中显示的应用名称 |
| `APP_VERSION` | 4.0 | 应用版本 |
| `TZ` | UTC | 日志时区 |

## 故障排除

### 常见问题

**OAuth 错误：**
- 检查 credentials.json 是否有效
- 验证重定向 URI 与配置匹配

**Webhook 未接收：**
- 确保公共 URL 可访问
- 检查 webhook 路径与配置匹配

**Rclone 刷新失败：**
- 验证 Rclone RC 已启用且可访问
- 检查路径映射规则

### 调试日志

启用调试模式查看详细日志：
- Rclone 请求/响应详情
- Google Drive API 调用
- 路径解析步骤

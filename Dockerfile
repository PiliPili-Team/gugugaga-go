# fileName: Dockerfile

# ========== Stage 1: Build Frontend ==========
FROM node:20-alpine AS frontend-builder
WORKDIR /app/web-src

# Copy frontend source
COPY web-src/package*.json ./
RUN npm ci --silent

COPY web-src/ ./

# Build to dist folder for Docker
ENV BUILD_OUTPUT=dist
RUN npm run build

# ========== Stage 2: Build Backend ==========
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app

# 设置国内代理 (按需开启)
# ENV GOPROXY=https://goproxy.cn,direct

# --- 优化点：利用 Docker 缓存层 ---
# 先只复制依赖描述文件
COPY go.mod go.sum ./
# 下载依赖 (如果 go.mod 没变，这一步会直接使用缓存)
RUN go mod download

# 复制源码
COPY src/ ./src/

# 复制前端构建产物
COPY --from=frontend-builder /app/web-src/dist ./src/web/static/

# 编译
# -ldflags="-s -w" 可以减小二进制体积
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app ./src

# ========== Stage 3: Runtime ==========
FROM alpine:latest
WORKDIR /app

# 安装基础证书和时区数据 (日志时间会更准)
RUN apk add --no-cache ca-certificates tzdata

# 从 builder 阶段复制
COPY --from=backend-builder /app/app .

# 挂载点
VOLUME /app/userdata

# 暴露端口
EXPOSE 8448

CMD ["./app"]

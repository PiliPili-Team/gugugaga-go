#!/bin/bash

# Vue 项目构建脚本
# 输出到 ../src/web/static 供 Go embed 使用

set -e

echo "🔨 Building Vue Frontend..."

cd "$(dirname "$0")"

# 检查 node 环境
if ! command -v node &> /dev/null; then
    echo "❌ Error: Node.js not found, please install Node.js first"
    exit 1
fi

echo "✅ Node version: $(node -v)"
echo "✅ NPM version: $(npm -v)"

# 安装依赖
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm install
fi

# 清理旧构建
echo "🧹 Cleaning old build..."
rm -rf ../src/web/static/*

# 构建
echo "🏗️  Building for production..."
npm run build

# 检查构建结果
if [ -f "../src/web/static/index.html" ]; then
    echo "✅ Build successful!"
    echo "📦 Output: src/web/static/"
    ls -la ../src/web/static/
else
    echo "❌ Build failed - index.html not found"
    exit 1
fi

echo ""
echo "🎉 Frontend build complete!"
echo "   Now you can build the Go binary with: ./build.sh"

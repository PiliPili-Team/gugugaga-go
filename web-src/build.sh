#!/bin/bash

set -e

echo "🔨 Building Vue Frontend..."

cd "$(dirname "$0")"

if ! command -v node &> /dev/null; then
    echo "❌ Error: Node.js not found, please install Node.js first"
    exit 1
fi

echo "✅ Node version: $(node -v)"
echo "✅ NPM version: $(npm -v)"

if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm install
fi

echo "🧹 Cleaning old build..."
rm -rf ../src/web/static/*

echo "🏗️  Building for production..."
npm run build

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

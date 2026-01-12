#!/bin/bash

set -e

echo "🔨 Building GD-Webhook Server..."

# Navigate to project directory
cd "$(dirname "$0")"

# Check Go environment
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go not found, please install Go first"
    exit 1
fi

echo "✅ Go version: $(go version)"

# Build frontend (optional, skip if --skip-frontend flag is passed)
if [[ "$1" != "--skip-frontend" ]]; then
    if command -v node &> /dev/null; then
        echo ""
        echo "🎨 Building Vue Frontend..."
        
        cd web-src
        
        # Install dependencies if needed
        if [ ! -d "node_modules" ]; then
            echo "📦 Installing frontend dependencies..."
            npm install
        fi
        
        # Build frontend
        echo "🏗️  Compiling frontend..."
        npm run build
        
        cd ..
        
        echo "✅ Frontend build complete!"
    else
        echo "⚠️  Node.js not found, skipping frontend build"
        echo "   Frontend will use existing files in src/web/static/"
    fi
else
    echo "⏭️  Skipping frontend build (--skip-frontend)"
fi

echo ""

# Clean old build artifacts
if [ -f "gd-webhook-server" ]; then
    echo "🧹 Cleaning old build artifacts..."
    rm -f gd-webhook-server
fi

# Build Go binary
echo "🏗️  Compiling Go backend..."
go build -o gd-webhook-server ./src

# Check build result
if [ -f "gd-webhook-server" ]; then
    echo "✅ Build successful!"
    echo "📦 Output: $(pwd)/gd-webhook-server"
    
    # Show file info
    ls -lh gd-webhook-server
    
    # Add execute permission
    chmod +x gd-webhook-server
    echo "✅ Execute permission added"
    
    echo ""
    echo "🚀 Run with:"
    echo "   ./gd-webhook-server"
else
    echo "❌ Build failed"
    exit 1
fi

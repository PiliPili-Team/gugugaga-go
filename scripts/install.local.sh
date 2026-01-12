#!/bin/bash

# Define colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${GREEN}🚀 Starting GD-Webhook One-Click Setup...${NC}"

# Navigate to project directory
cd "$(dirname "$0")" || exit

# Function to print usage guide
print_install_guide() {
    local tool=$1
    echo -e "${BLUE}💡 How to install ${tool}:${NC}"
    if [ "$tool" == "Go" ]; then
        echo -e "  ${YELLOW}MacOS:${NC} brew install go"
        echo -e "  ${YELLOW}Linux:${NC} Follow https://go.dev/doc/install or use package manager (e.g., apt install golang)"
    elif [ "$tool" == "PM2" ]; then
        echo -e "  ${YELLOW}MacOS/Linux:${NC} First install Node.js/NPM, then run: ${GREEN}npm install -g pm2${NC}"
        echo -e "  ${YELLOW}Node.js:${NC} https://nodejs.org/ (or use nvm)"
    elif [ "$tool" == "Node.js" ]; then
        echo -e "  ${YELLOW}MacOS:${NC} brew install node"
        echo -e "  ${YELLOW}Linux:${NC} https://nodejs.org/ or use nvm"
    fi
}

# 1. Check Go Environment
echo -e "\n${YELLOW}🔧 Step 1: Checking Go Environment...${NC}"
if ! command -v go >/dev/null 2>&1; then
    echo -e "${RED}❌ Go compiler not found.${NC}"
    print_install_guide "Go"
    exit 1
fi
echo -e "${GREEN}✅ Go version: $(go version | awk '{print $3}')${NC}"

# 2. Build Frontend (optional)
echo -e "\n${YELLOW}🎨 Step 2: Building Frontend...${NC}"
if command -v node >/dev/null 2>&1; then
    echo -e "   Node.js version: $(node -v)"
    
    if [ -d "web-src" ]; then
        cd web-src
        
        # Install dependencies if needed
        if [ ! -d "node_modules" ]; then
            echo -e "   Installing frontend dependencies..."
            npm install --silent
        fi
        
        # Build frontend
        echo -e "   Compiling Vue frontend..."
        npm run build --silent
        
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}✅ Frontend build complete!${NC}"
        else
            echo -e "${YELLOW}⚠️ Frontend build failed, using existing files${NC}"
        fi
        
        cd ..
    else
        echo -e "${YELLOW}⚠️ web-src directory not found, skipping frontend build${NC}"
    fi
else
    echo -e "${YELLOW}⚠️ Node.js not found, skipping frontend build${NC}"
    echo -e "   Frontend will use existing files in src/web/static/"
    print_install_guide "Node.js"
fi

# 3. Build Go Application
echo -e "\n${YELLOW}🏗️ Step 3: Building Go Application...${NC}"

# Clean old binary
if [ -f "gd-webhook-server" ]; then
    rm -f gd-webhook-server
fi

go build -o gd-webhook-server ./src

if [ $? -eq 0 ]; then
    chmod +x gd-webhook-server
    echo -e "${GREEN}✅ Build successful! Binary: ./gd-webhook-server${NC}"
    ls -lh gd-webhook-server
else
    echo -e "${RED}❌ Build failed. Please check your code or Go version.${NC}"
    exit 1
fi

# 4. PM2 Setup
echo -e "\n${YELLOW}🚀 Step 4: Configuring PM2 Task...${NC}"
if ! command -v pm2 >/dev/null 2>&1; then
    echo -e "${RED}❌ PM2 not found.${NC}"
    print_install_guide "PM2"
    echo -e "\n${YELLOW}💡 You can run the server manually:${NC}"
    echo -e "   ./gd-webhook-server"
    exit 0
fi

# PM2 Management
echo -e "   Updating PM2 process list..."
pm2 stop gd-webhook >/dev/null 2>&1
pm2 delete gd-webhook >/dev/null 2>&1

# Start with ecosystem config
pm2 start ecosystem.config.js
if [ $? -eq 0 ]; then
    pm2 save
    echo -e "${GREEN}✅ PM2 task 'gd-webhook' started and config saved.${NC}"
    echo -e ""
    echo -e "📝 ${BLUE}Useful commands:${NC}"
    echo -e "   ${YELLOW}pm2 log gd-webhook${NC}     - View logs"
    echo -e "   ${YELLOW}pm2 restart gd-webhook${NC} - Restart service"
    echo -e "   ${YELLOW}pm2 stop gd-webhook${NC}    - Stop service"
    echo -e "   ${YELLOW}pm2 startup${NC}            - Enable auto-start on boot"
    echo -e ""
    echo -e "🌐 ${GREEN}WebUI available at: http://localhost:8448${NC}"
else
    echo -e "${RED}❌ Failed to start PM2 task.${NC}"
    exit 1
fi

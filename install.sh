#!/usr/bin/env bash
set -e

REPO="https://github.com/PiliPili-Team/gugugaga-go.git"
APP_NAME="gugugaga-go"
INSTALL_DIR="/opt/${APP_NAME}"

echo "🚀 Gugugaga-Go One-Click Installer"
echo "📦 Repository: $REPO"
echo "📂 Install Dir: $INSTALL_DIR"

if ! command -v git >/dev/null 2>&1; then
  echo "❌ git not found, please install git first"
  exit 1
fi

if [ -d "$INSTALL_DIR/.git" ]; then
  echo "🔄 Existing installation found, updating..."
  cd "$INSTALL_DIR"
  git pull
else
  echo "⬇️  Cloning repository..."
  sudo mkdir -p "$INSTALL_DIR"
  sudo chown -R "$(whoami)" "$INSTALL_DIR"
  git clone "$REPO" "$INSTALL_DIR"
  cd "$INSTALL_DIR"
fi

chmod +x scripts/install.local.sh
exec bash scripts/install.local.sh

#!/usr/bin/env bash
# One-time VPS bootstrap — run as user 'maker' on 103.172.205.183
# Usage: bash setup_vps.sh
set -euo pipefail

APP_DIR=/var/www/sespimma-api-go
NGINX_CONF=/etc/nginx/sites-available/sespimma

echo "=== Creating application directory ==="
sudo mkdir -p "$APP_DIR/uploads/profile"
sudo chown -R maker:maker "$APP_DIR"

echo "=== Copying .env to app directory ==="
# Place your .env file in the same directory as this script, then run:
if [ -f ".env" ]; then
  cp .env "$APP_DIR/.env"
  chmod 600 "$APP_DIR/.env"
  echo ".env copied"
else
  echo "WARNING: .env not found — copy .env.example to .env, fill in credentials, then rerun."
fi

echo "=== Installing Nginx config ==="
sudo cp "$(dirname "$0")/nginx.conf" "$NGINX_CONF"
sudo ln -sf "$NGINX_CONF" /etc/nginx/sites-enabled/sespimma
sudo nginx -t && sudo systemctl reload nginx
echo "Nginx configured"

echo "=== Ensuring PM2 is available ==="
export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && . "$NVM_DIR/nvm.sh"
if ! command -v pm2 &>/dev/null; then
  npm install -g pm2
  pm2 startup
fi

echo "=== Required GitHub Actions secrets ==="
echo "  VPS_HOST       = 103.172.205.183"
echo "  VPS_USERNAME   = maker"
echo "  VPS_PASSWORD   = (VPS password)"
echo "  DB_PASSWORD    = (DB password from .env)"
echo ""
echo "Setup complete. Push to main to trigger first deploy."

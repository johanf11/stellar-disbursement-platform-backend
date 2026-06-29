#!/usr/bin/env bash
# Deploy the SDP frontend to the DigitalOcean droplet.
# Usage: ./dev/deploy-frontend.sh
# Run from the repo root (stellar-disbursement-platform-backend).
# Expects the frontend repo at ../stellar-disbursement-platform-frontend

set -euo pipefail

DROPLET_IP="64.225.19.99"
REMOTE_USER="root"
SSH_KEY="${HOME}/.ssh/id_do_sdp"
FRONTEND_DIR="$(realpath "$(dirname "$(realpath "$0")")/../../sdp-frontend-theo-demo")"
export PATH="${HOME}/.npm-global/bin:/usr/local/bin:$PATH"

if [ ! -d "$FRONTEND_DIR" ]; then
  echo "ERROR: Frontend repo not found at $FRONTEND_DIR"
  exit 1
fi

echo "==> Building frontend..."
cd "$FRONTEND_DIR"
npm run build

echo "==> Uploading build to droplet..."
ssh -i "$SSH_KEY" -o IdentitiesOnly=yes "${REMOTE_USER}@${DROPLET_IP}" "mkdir -p /opt/sdp-frontend"
scp -i "$SSH_KEY" -o IdentitiesOnly=yes -r build/. "${REMOTE_USER}@${DROPLET_IP}:/opt/sdp-frontend/"

echo "==> Reloading nginx..."
ssh -i "$SSH_KEY" -o IdentitiesOnly=yes "${REMOTE_USER}@${DROPLET_IP}" "systemctl reload nginx"

echo ""
echo "==> Done!"
echo "    Frontend: http://${DROPLET_IP}:3000"

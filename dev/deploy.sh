#!/usr/bin/env bash
# Deploy the SDP backend to the DigitalOcean droplet.
# Usage: ./dev/deploy.sh
# Run from the repo root.

set -euo pipefail

DROPLET_IP="64.225.19.99"
REMOTE_USER="root"
APP_DIR="/opt/sdp"
BINARY="stellar-disbursement-platform"
SSH_KEY="${HOME}/.ssh/id_do_sdp"
SSH_OPTS="-i ${SSH_KEY} -o IdentitiesOnly=yes -o StrictHostKeyChecking=accept-new"

echo "==> Building linux/amd64 binary..."
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "$BINARY" .

echo "==> Creating app directory on droplet..."
ssh ${SSH_OPTS} "${REMOTE_USER}@${DROPLET_IP}" "mkdir -p ${APP_DIR}"

echo "==> Copying binary..."
scp ${SSH_OPTS} "$BINARY" "${REMOTE_USER}@${DROPLET_IP}:${APP_DIR}/${BINARY}"
ssh ${SSH_OPTS} "${REMOTE_USER}@${DROPLET_IP}" "chmod +x ${APP_DIR}/${BINARY}"

echo "==> Copying systemd unit files..."
scp ${SSH_OPTS} dev/sdp-api.service dev/sdp-tss.service "${REMOTE_USER}@${DROPLET_IP}:/etc/systemd/system/"

echo "==> Uploading env file (skipped if /opt/sdp/.env already exists)..."
ssh ${SSH_OPTS} "${REMOTE_USER}@${DROPLET_IP}" "test -f ${APP_DIR}/.env && echo '  .env already exists, skipping'" \
  || scp ${SSH_OPTS} dev/.env.production "${REMOTE_USER}@${DROPLET_IP}:${APP_DIR}/.env"

echo "==> Opening firewall ports..."
ssh ${SSH_OPTS} "${REMOTE_USER}@${DROPLET_IP}" "ufw allow ssh; ufw allow 8000/tcp; ufw allow 8003/tcp; ufw --force enable" 2>/dev/null || true

echo "==> Reloading systemd and restarting services..."
ssh ${SSH_OPTS} "${REMOTE_USER}@${DROPLET_IP}" bash <<'ENDSSH'
  systemctl daemon-reload
  systemctl enable sdp-api sdp-tss
  systemctl restart sdp-api sdp-tss
  sleep 3
  echo ""
  echo "--- sdp-api status ---"
  systemctl status sdp-api --no-pager -l
  echo ""
  echo "--- sdp-tss status ---"
  systemctl status sdp-tss --no-pager -l
ENDSSH

echo "==> Cleaning up local binary..."
rm -f "$BINARY"

echo ""
echo "==> Done!"
echo "    API:   http://${DROPLET_IP}:8000"
echo "    Admin: http://${DROPLET_IP}:8003"
echo "    Logs:  ssh ${REMOTE_USER}@${DROPLET_IP} journalctl -u sdp-api -f"

#!/bin/bash
set -euo pipefail

echo "=>> Deploying Server SFP SLA service"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BINARY_PATH="${SFP_SLA_BINARY:-$PROJECT_ROOT/build/Server_SFP_SLA}"
SERVICE_PATH="$PROJECT_ROOT/deploy/systemd/Server_SFP_SLA.service"
ENV_PATH="/etc/default/server-sfp-sla"

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

systemctl stop Server_SFP_SLA 2>/dev/null || true
sleep 1

if [ ! -x "$BINARY_PATH" ]; then
  if ! command -v go >/dev/null 2>&1; then
    echo "Binary not found and Go is unavailable: $BINARY_PATH"
    exit 1
  fi

  mkdir -p "$(dirname "$BINARY_PATH")"
  (cd "$PROJECT_ROOT" && go build -o "$BINARY_PATH" ./cmd/server-sfp-sla)
fi

cp -f "$BINARY_PATH" /usr/local/bin/Server_SFP_SLA
cp -f "$SERVICE_PATH" /lib/systemd/system/Server_SFP_SLA.service

if [ ! -f "$ENV_PATH" ]; then
  install -m 0640 /dev/null "$ENV_PATH"
  cat > "$ENV_PATH" <<EOF
SFP_SLA_DB_USER=${SFP_SLA_DB_USER:-sfp_user}
SFP_SLA_DB_PASSWORD=${SFP_SLA_DB_PASSWORD:-}
SFP_SLA_DB_NAME=${SFP_SLA_DB_NAME:-server_sfp_sla}
SFP_SLA_DB_ADDR=${SFP_SLA_DB_ADDR:-}
SFP_SLA_TIMEZONE=${SFP_SLA_TIMEZONE:-UTC}
SFP_SLA_VERBOSE=${SFP_SLA_VERBOSE:-0}
EOF
  chown root:root "$ENV_PATH"
fi

chmod 644 /lib/systemd/system/Server_SFP_SLA.service
chmod u+x /usr/local/bin/Server_SFP_SLA

systemctl daemon-reload
ulimit -n 16384
sysctl -w net.ipv4.ping_group_range="0 2147483647"
systemctl enable Server_SFP_SLA
systemctl start Server_SFP_SLA

echo "=>> Service deployed"

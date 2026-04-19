#!/bin/bash
set -euo pipefail

echo "=>> Removing Server SFP SLA"

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

: "${MYSQL_ROOT_PASSWORD:?Set MYSQL_ROOT_PASSWORD before running this script}"
SFP_SLA_DB_USER="${SFP_SLA_DB_USER:-sfp_user}"
SFP_SLA_DB_NAME="${SFP_SLA_DB_NAME:-server_sfp_sla}"

validate_sql_name() {
  local name="$1"
  local value="$2"
  if [[ ! "$value" =~ ^[A-Za-z0-9_]+$ ]]; then
    echo "$name must contain only letters, digits, and underscores"
    exit 1
  fi
}

validate_sql_name "SFP_SLA_DB_USER" "$SFP_SLA_DB_USER"
validate_sql_name "SFP_SLA_DB_NAME" "$SFP_SLA_DB_NAME"

systemctl stop Server_SFP_SLA 2>/dev/null || true
systemctl disable Server_SFP_SLA 2>/dev/null || true
rm -f /usr/local/bin/Server_SFP_SLA /lib/systemd/system/Server_SFP_SLA.service /etc/default/server-sfp-sla
systemctl daemon-reload

mysql -u root -p"$MYSQL_ROOT_PASSWORD" <<MY_QUERY
DROP DATABASE IF EXISTS \`$SFP_SLA_DB_NAME\`;
DROP USER IF EXISTS '$SFP_SLA_DB_USER'@'localhost';
FLUSH PRIVILEGES;
MY_QUERY

rm -rf /var/www/html/*
apt -y remove php php-cli php-common php-mysql php-snmp php-curl php-cgi libapache2-mod-php || true
apt -y remove mysql mysql-server apache2 net-tools pv || true

echo "=>> Removal completed"

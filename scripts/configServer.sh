#!/bin/bash
set -euo pipefail

echo "=>> Configuring MySQL"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

: "${MYSQL_ROOT_PASSWORD:?Set MYSQL_ROOT_PASSWORD before running this script}"
SFP_SLA_DB_USER="${SFP_SLA_DB_USER:-sfp_user}"
: "${SFP_SLA_DB_PASSWORD:?Set SFP_SLA_DB_PASSWORD before running this script}"
SFP_SLA_DB_NAME="${SFP_SLA_DB_NAME:-server_sfp_sla}"

validate_sql_name() {
  local name="$1"
  local value="$2"
  if [[ ! "$value" =~ ^[A-Za-z0-9_]+$ ]]; then
    echo "$name must contain only letters, digits, and underscores"
    exit 1
  fi
}

sql_escape() {
  printf "%s" "$1" | sed "s/'/''/g"
}

validate_sql_name "SFP_SLA_DB_USER" "$SFP_SLA_DB_USER"
validate_sql_name "SFP_SLA_DB_NAME" "$SFP_SLA_DB_NAME"
SFP_SLA_DB_PASSWORD_SQL="$(sql_escape "$SFP_SLA_DB_PASSWORD")"

# Keep compatibility with MySQL installations that still require native password authentication.
if [ -f /etc/mysql/mysql.conf.d/mysqld.cnf ] && ! grep -Eq '^[^#].*mysql_native_password' /etc/mysql/mysql.conf.d/mysqld.cnf; then
  echo "default-authentication-plugin=mysql_native_password" >> /etc/mysql/mysql.conf.d/mysqld.cnf
  systemctl restart mysql
  sleep 1
fi

mysql -u root -p"$MYSQL_ROOT_PASSWORD" <<MY_QUERY
CREATE DATABASE IF NOT EXISTS \`$SFP_SLA_DB_NAME\`;
USE mysql;
FLUSH PRIVILEGES;
CREATE USER IF NOT EXISTS '$SFP_SLA_DB_USER'@'localhost' IDENTIFIED BY '$SFP_SLA_DB_PASSWORD_SQL';
GRANT ALL PRIVILEGES ON \`$SFP_SLA_DB_NAME\`.* TO '$SFP_SLA_DB_USER'@'localhost';
MY_QUERY

mysql -u root -p"$MYSQL_ROOT_PASSWORD" --database="$SFP_SLA_DB_NAME" < "$PROJECT_ROOT/db/server_sfp_sla.sql"

echo "=>> MySQL configured"
"$SCRIPT_DIR/deploy.sh"

echo "=>> Publishing web console"
rm -rf /var/www/html/*
cp -a "$PROJECT_ROOT/web/htdocs/." /var/www/html/
chown -R www-data:www-data /var/www/html
systemctl enable apache2
systemctl restart apache2

echo "=>> Web console published"
echo "=>> Create a web administrator with: sudo -E $SCRIPT_DIR/add_user.sh"

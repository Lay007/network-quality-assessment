#!/bin/bash
set -euo pipefail

echo "== Add web user =="
read -r -p "Login: " login
read -r -s -p "Password: " password
echo

if [ -z "$login" ] || [ -z "$password" ]; then
  echo "Login and password must not be empty"
  exit 1
fi

if ! command -v php >/dev/null 2>&1; then
  echo "PHP is required to hash the web user password"
  exit 1
fi

: "${MYSQL_ROOT_PASSWORD:?Set MYSQL_ROOT_PASSWORD before running this script}"
SFP_SLA_DB_NAME="${SFP_SLA_DB_NAME:-server_sfp_sla}"

if [[ ! "$SFP_SLA_DB_NAME" =~ ^[A-Za-z0-9_]+$ ]]; then
  echo "SFP_SLA_DB_NAME must contain only letters, digits, and underscores"
  exit 1
fi

sql_escape() {
  printf "%s" "$1" | sed "s/'/''/g"
}

login_sql="$(sql_escape "$login")"
password_hash="$(php -r 'echo password_hash($argv[1], PASSWORD_DEFAULT);' "$password")"
password_sql="$(sql_escape "$password_hash")"

mysql -u root -p"$MYSQL_ROOT_PASSWORD" <<MY_QUERY
USE \`$SFP_SLA_DB_NAME\`;
INSERT INTO users (login, password, type)
VALUES ('$login_sql', '$password_sql', 'admin')
ON DUPLICATE KEY UPDATE password = VALUES(password), type = VALUES(type);
MY_QUERY

echo "=>> User added"

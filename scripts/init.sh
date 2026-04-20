#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

echo "=> Installing system dependencies"

apt update
apt -y install build-essential gcc
apt -y install apache2
apt -y install mysql mysql-server
apt -y install php php-cli php-common php-mysql php-snmp php-curl php-cgi libapache2-mod-php
apt -y install net-tools pv

chmod u+x "$SCRIPT_DIR/deploy.sh"
chmod u+x "$SCRIPT_DIR/configServer.sh"
chmod u+x "$SCRIPT_DIR/clean.sh"
chmod u+x "$SCRIPT_DIR/add_user.sh"
chmod u+x "$SCRIPT_DIR/check_web_links.py"

echo "=> System dependencies installed"

#!/bin/bash
set -euo pipefail

echo "=>> Installing development dependencies"

if [ "$EUID" -ne 0 ]; then
  echo "Please run as root"
  exit 1
fi

apt update
apt -y install build-essential gcc libpcap0.8-dev php-cli php-mysql php-snmp php-curl

echo "=>> Development dependencies installed"

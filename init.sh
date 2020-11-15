#!/bin/bash

echo "=> Server init"

apt update && apt upgrade
apt-get install libpcap0.8-dev
apt install apache2
apt install mysql
apt install mysql-server
apt install php php-cli php-common php-mysql php-snmp

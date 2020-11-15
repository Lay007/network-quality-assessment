#!/bin/bash

echo "=> Zabbix setup"

wget https://repo.zabbix.com/zabbix/4.0/ubuntu/pool/main/z/zabbix-release/zabbix-release_4.0-3+xenial_all.deb
dpkg -i zabbix-release_4.0-3+xenial_all.deb
apt update && apt upgrade
apt install zabbix-server-mysql zabbix-frontend-php zabbix-agent

echo "=> Zabbix setup succes"
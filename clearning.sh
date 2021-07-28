#!/bin/bash
echo " =>> Clean MySQL"
# Удаление таблиц server_sfp_sla и пользователя sfp_user 
mysql -u root -psecret <<MY_QUERY

CREATE DATABASE IF NOT EXISTS server_sfp_sla;
DROP DATABASE server_sfp_sla;

USE mysql;
drop user 'sfp_user'@'localhost';
flush privileges;

MY_QUERY

# Удаление установленных пакетов
apt -y purge libpcap0.8-dev
apt -y purge apache2
apt -y purge mysql
apt -y purge mysql-server
apt -y purge php php-cli php-common php-mysql php-snmp php-curl php-cgi libapache2-mod-php
apt -y purge php7.0  php7.0-cli php7.0-common php-mysql php-snmp php7.0-curl php7.0-cgi  libapache2-mod-php7.0
apt -y purge net-tools pv
apt -y purge zabbix-server-mysql zabbix-frontend-php zabbix-agent

echo " =>> Clean web Server_SFP_SLA"

echo " =>> Clean Zabbix"

rm ./zabbix-release_4.0*
mysql -u root -psecret <<MY_QUERY

CREATE DATABASE IF NOT EXISTS zabbix;
DROP DATABASE zabbix;

USE mysql;
drop user 'zabbix'@'localhost';
flush privileges;

MY_QUERY

apt -y autoremove

echo " =>> Clean end"
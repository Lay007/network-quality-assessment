#!/bin/bash
echo "=> Zabbix setup"
# Проверка на запуск от имени администратора
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

# Определение версии Ububtu Linux
source /etc/os-release
DNAME=$UBUNTU_CODENAME
if [[ -z "$DNAME" ]]; then
	source /etc/lsb-release
	DNAME=$DISTRIB_CODENAME
fi	

# Скачивание дистрибутива Zabbix
wget https://repo.zabbix.com/zabbix/4.0/ubuntu/pool/main/z/zabbix-release/zabbix-release_4.0-3+"$DNAME"_all.deb
# Установка дистрибутива Zabbix
dpkg -i zabbix-release_4.0-3+"$DNAME"_all.deb
apt update && apt upgrade
apt -y install zabbix-server-mysql zabbix-frontend-php zabbix-agent

# Начальная настройка базы данных для Zabbix
mysql -u root -psecret <<MY_QUERY
   create database zabbix character set utf8 collate utf8_bin;
   create user zabbix@localhost identified by 'zabbix';
   grant all privileges on zabbix.* to zabbix@localhost;
MY_QUERY

echo "Please, wait some minutes while we deploy Zabbix database"
zcat /usr/share/doc/zabbix-server-mysql*/create.sql.gz | pv |  mysql -u zabbix --password=zabbix zabbix

# Установка пароля для базы данных Zabbix
echo "DBPassword=zabbix" >> /etc/zabbix/zabbix_server.conf 
# Установка временной зоны Zabbix
sed -i 's%.*# php_value date.timezone Europe/Riga%        php_value date.timezone Europe/Moscow%g' /etc/apache2/conf-enabled/zabbix.conf

# Включение и перезапуск системы мониторинга Zabbix
systemctl restart zabbix-server zabbix-agent apache2
systemctl enable zabbix-server zabbix-agent apache2

echo "=> Zabbix setup success"

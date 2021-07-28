#!/bin/bash

# Проверка на запуск от имени администратора
echo "=> Server init"
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

# Проверка и установка текущих обновлений
apt update && apt upgrade
# Установка библиотеки libpcap для работы с сетевыми пакетами
apt -y install libpcap0.8-dev
# Установка Web-сервера Apache
apt -y install apache2
# Установка базы данных MySQL
apt -y install mysql
apt -y install mysql-server
# Установка модулей языка PHP
apt -y install php php-cli php-common php-mysql php-snmp php-curl php-cgi libapache2-mod-php
apt -y install php7.0  php7.0-cli php7.0-common php-mysql php-snmp php7.0-curl php7.0-cgi  libapache2-mod-php7.0
# Установка модулей для работы с сетевыми интерфейсами
apt -y install net-tools pv

chmod u+x ./deploy.sh
chmod u+x ./setupZabbix.sh
chmod u+x ./configServer.sh
chmod u+x ./clearning.sh
chmod u+x ./add_user.sh
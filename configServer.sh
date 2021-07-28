#!/bin/bash
echo " =>> MySQL config"

# Проверка на запуск от имени администратора
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi

# Добавление функции mysql_native_password для обратной совместимости
if egrep  ^[^#].*mysql_native_password /etc/mysql/mysql.conf.d/mysqld.cnf
 then echo "native password found"
 else
   echo "not found"
   echo "add native password to auth"
   echo "default-authentication-plugin=mysql_native_password" >> /etc/mysql/mysql.conf.d/mysqld.cnf
   systemctl restart mysql
   sleep 1
fi
# Начальная настройка базы данных
mysql -u root -psecret <<MY_QUERY

CREATE DATABASE server_sfp_sla;
USE server_sfp_sla;
SOURCE ./server_sfp_sla.sql;

USE mysql;
flush privileges;
create user 'sfp_user'@'localhost' identified by 'rootsfp';

grant all privileges on server_sfp_sla . * to 'sfp_user'@'localhost';

MY_QUERY

echo " =>> MySQL config success"
# Запуск скрипта
./deploy.sh

echo " =>> Deploy web-config"
# Размещение Web-конфигуратора
rm -rf /var/www/html/*
tar -C / -xvf httpServerSLA.tar.gz
# Включение и перезапуск Web-сервера Apache
systemctl enable apache2
systemctl restart apache2

echo " =>> Deploy web-config success"

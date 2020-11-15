#!/bin/bash
echo " =>> MySQL config"

mysql -u root -psecret <<MY_QUERY

CREATE DATABASE server_sfp_sla;
USE server_sfp_sla;
SOURCE ./server_sfp_sla.sql;

USE mysql;
flush privileges;
create user 'sfp_user'@'localhost' identified by 'rootsfp';

grant all privileges on server_sfp_sla . * to 'sfp_user'@'localhost';

MY_QUERY

echo " =>> MySQL config succes"

./deploy.sh

echo " =>> Deploy web-config"

rm -rf /var/www/html/*
tar -C / -xvf httpServerSLA.tar.gz
systemctl enable apache2
systemctl restart apache2

echo " =>> Deploy web-config succes"

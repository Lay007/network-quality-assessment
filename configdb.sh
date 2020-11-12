#!/bin/bash
echo " =>> MySQL config"

mysql -u root -psecret <<MY_QUERY

CREATE DATABASE IF NOT EXISTS server_sfp_sla;
TRUNC Table server_sfp_sla
USE server_sfp_sla;
SOURCE ./server_sfp_sla.sql;

GRANT USAGE ON *.* TO 'sfp_user'@'localhost';
DROP user 'sfp_user'@'localhost';

flush privileges;

USE mysql;
create user 'sfp_user'@'localhost' identified by 'rootsfp';
grant all privileges on server_sfp_sla . * to 'sfp_user'@'localhost';

MY_QUERY

echo " =>> MySQL config succes"

./deploy.sh

echo " =>> Deploy web-config"

tar -C / -xvf httpServerSLA.tar.gz

echo " =>> Deploy web-config succes"
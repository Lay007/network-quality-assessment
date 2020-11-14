#!/bin/bash
echo " =>> MySQL config"

mysql -u root -psecret <<MY_QUERY

CREATE DATABASE IF NOT EXISTS server_sfp_sla;
DROP DATABASE server_sfp_sla;
CREATE DATABASE server_sfp_sla;
USE server_sfp_sla;
SOURCE ./server_sfp_sla.sql;

USE mysql;
drop user 'sfp_user'@'localhost';
flush privileges;
create user 'sfp_user'@'localhost' identified by 'rootsfp';

grant all privileges on server_sfp_sla . * to 'sfp_user'@'localhost';

MY_QUERY

echo " =>> MySQL config succes"

./deploy.sh

echo " =>> Deploy web-config"

#tar -C / -xvf httpServerSLA.tar.gz

echo " =>> Deploy web-config succes"

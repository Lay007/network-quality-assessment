#!/bin/bash
echo " =>> Clean MySQL"
mysql -u root -psecret <<MY_QUERY

CREATE DATABASE IF NOT EXISTS server_sfp_sla;
DROP DATABASE server_sfp_sla;

USE mysql;
drop user 'sfp_user'@'localhost';
flush privileges;

MY_QUERY
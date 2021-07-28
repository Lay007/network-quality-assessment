#!/bin/bash
echo " == Добавление нового пользователя =="
read -p "Введите имя пользователя: " user
read -sp "Введите пароль: " pass

# Удаление таблиц server_sfp_sla и пользователя sfp_user 
mysql -u root -psecret <<MY_QUERY

USE server_sfp_sla;
INSERT INTO users (login, password, type) VALUES
('$user', '$pass', 'admin');

MY_QUERY

echo " =>> Добавление пользователя завершено"
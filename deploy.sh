#!/bin/bash
echo " =>> Deploy generator"
# Проверка на запуск от имени администратора
if [ "$EUID" -ne 0 ]
  then echo "Please run as root"
  exit
fi
# Остановка генератора и анализатора пакетов
systemctl stop Server_SFP_SLA
sleep 1

# Копироание новых файлов генератора и анализатора пакетов  
cp -f ./Server_SFP_SLA /usr/local/bin
cp -f ./Server_SFP_SLA.service /lib/systemd/system

# Присвоение прав доступа для файлов генератора и анализатора пакетов
chmod 644 /lib/systemd/system/Server_SFP_SLA.service
chmod u+x /usr/local/bin/Server_SFP_SLA

# Включение и перезапуск генератора и анализатора пакетов
systemctl daemon-reload
# исправление ошибки socket “Too many open files”
ulimit -n 16384
# Операция для реализации функции пинга
sysctl -w net.ipv4.ping_group_range="0 2147483647"
systemctl enable Server_SFP_SLA
systemctl start Server_SFP_SLA

echo " =>> Deploy generator success"

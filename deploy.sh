#!/bin/bash
echo " =>> Server SFP SLA"

#./go build

systemctl stop Server_SFP_SLA
sleep 1
cp -f ./Server_SFP_SLA /bin
cp -f ./Server_SFP_SLA.service /lib/systemd/system
chmod 644 /lib/systemd/system/Server_SFP_SLA.service

systemctl daemon-reload
systemctl enable Server_SFP_SLA
systemctl start Server_SFP_SLA


#git pull

#sudo ./Server_SFP_SLA
#!/bin/bash
echo " =>> Server SFP SLA"

systemctl stop Server_SFP_SLA
cp -f ./Server_SFP_SLA bin
systemctl daemon-reload
systemctl start Server_SFP_SLA


#git pull
#go build
#sudo ./Server_SFP_SLA
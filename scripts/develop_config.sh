#!/bin/bash
echo " =>> Develop configuration"

sudo curl -O https://dl.google.com/go/go1.14.linux-amd64.tar.gz
sudo tar -xvf go1.14.linux-amd64.tar.gz
sudo mv go /usr/local
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go/packages
source /etc/profile
go version

sudo apt install mc
sudo apt install gcc

go get github.com/mdlayher/ethernet
go get github.com/go-sql-driver/mysql
go get github.com/google/gopacket/pcap 
go get github.com/mdlayher/raw
go get -u github.com/go-ping/ping

go get github.com/newtools/zsocket
go get github.com/soniah/gosnmp
go get github.com/tatsushid/go-fastping

sysctl -w net.ipv4.ping_group_range="0 2147483647"


echo " =>>  Develop configuration succes"
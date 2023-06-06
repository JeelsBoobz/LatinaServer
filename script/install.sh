#!/usr/bin/env bash

DIR=$(dirname "$0")
PROJECT=$DIR/..

go build -o $PROJECT/latinaserver $PROJECT/cmd/latinaserver/main.go

sudo mkdir -p /usr/local/etc/latinaserver
sudo cp $PROJECT/resources/openresty/stream/reverse_default.conf /usr/local/etc/latinaserver/
sudo cp $PROJECT/config.json /usr/local/etc/latinaserver/

if [ ! -f /etc/systemd/system/latinaserver.service ]; then
    sudo cp ./latinaserver.service /etc/systemd/system/
else 
    sudo systemctl stop latinaserver
fi
sudo cp $PROJECT/latinaserver /usr/local/bin/
sudo systemctl daemon-reload

sudo systemctl start latinaserver
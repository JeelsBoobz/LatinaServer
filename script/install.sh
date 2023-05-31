#!/usr/bin/env bash

DIR=$(dirname "$0")
PROJECT=$DIR/..

go build -o $PROJECT/latinaserver $PROJECT/cmd/latinaserver/main.go

sudo mkdir -p /usr/local/etc/latinaserver
sudo cp $PROJECT/config.json /usr/local/etc/latinaserver/
sudo cp $PROJECT/latinaserver /usr/local/bin/ 
sudo cp ./latinaserver.service /etc/systemd/system/
sudo systemctl daemon-reload
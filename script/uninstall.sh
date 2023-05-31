#!/usr/bin/env bash

sudo systemctl stop latinaserver
sudo rm -rf /usr/local/etc/latinaserver
sudo rm -rf /usr/local/bin/latinaserver
sudo rm -rf /etc/systemd/system/latinaserver.service
sudo systemctl daemon-reload
#!/bin/bash
sudo sh -c "dockerd >/var/log/dockerd.log 2>&1 &"

zsh /setup/setup.sh

sudo service ssh start

touch ~/.setup_complete

sleep infinity

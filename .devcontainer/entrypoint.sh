#!/bin/bash
sudo sh -c "dockerd >/var/log/dockerd.log 2>&1 &"

fish /setup/setup.fish

sudo service ssh start

touch ~/.setup_complete

sleep infinity

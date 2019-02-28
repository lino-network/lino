#!/usr/bin/env bash

if [ "$#" -ne 3 ]; then
    echo "require 3 parameters: linux username, persistent_peers, private_peer_ids"
    exit 1
fi

USER=$1
CONFIG=/home/$USER/.lino/config/config.toml

sed -i 's/^db_backend = "leveldb"$/db_backend = "cleveldb"/g' $CONFIG

sed -i 's/^persistent_peers = ""$/persistent_peers = "'"$2"'"/g' $CONFIG

sed -i 's/^send_rate = 5120000$/send_rate = 128000000/g' $CONFIG

sed -i 's/^recv_rate = 5120000$/recv_rate = 128000000/g' $CONFIG

sed -i 's/^private_peer_ids = ""$/private_peer_ids = "'"$3"'"/g' $CONFIG

sed -i 's/^timeout_commit = "5s"$/timeout_commit = "3s"/g' $CONFIG

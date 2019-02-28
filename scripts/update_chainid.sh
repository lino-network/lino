#!/usr/bin/env bash

if [ "$#" -ne 2 ]; then
    echo "require 2 parameters: linux username, chainID"
    exit 1
fi

USER=$1
GENESIS=/home/$USER/.lino/config/genesis.json

sed -i 's/^.*"chain_id": .*$/  "chain_id": "'"$2"'",/g' $GENESIS

#!/usr/bin/env bash

if [ "$#" -ne 3 ]; then
    echo "require 3 parameters: linux username, chainID, genesistime"
    exit 1
fi

USER=$1
GENESIS=/home/$USER/.lino/config/genesis.json

sed -i 's/^.*"chain_id": .*$/  "chain_id": "'"$2"'",/g' $GENESIS
sed -i 's/^.*"genesis_time": .*$/  "genesis_time": "'"$3"'",/g' $GENESIS

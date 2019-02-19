#!/usr/bin/env bash

echo "info: MUST run dep ensure before this script"
diff -u ../../vendor/github.com/cosmos/cosmos-sdk/baseapp/baseapp.go ./baseapp > ./baseapp.patch
diff -u ../../vendor/github.com/cosmos/cosmos-sdk/store/iavlstore.go ./iavlstore > ./iavlstore.patch
diff -u ../../vendor/github.com/tendermint/tendermint/rpc/lib/server/http_server.go ./http_server > ./http_server.patch
echo "patches generated."


#!/usr/bin/env bash

echo "info: MUST run dep ensure before this script"
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
diff -u $DIR/../../vendor/github.com/cosmos/cosmos-sdk/server/constructors.go $DIR/constructors > cosmos-server-constructor.patch
echo "patches generated."


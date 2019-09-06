#!/usr/bin/env bash

### update golang to latest.
sudo rm -rf /usr/local/go
GOSOURCE=https://dl.google.com/go/go1.13.linux-amd64.tar.gz
GOTARGET=/usr/local
GOPATH=\$HOME/go
PROFILE=/home/ubuntu/.profile

curl -sSL $GOSOURCE -o /tmp/go.tar.gz
sudo tar -xzf /tmp/go.tar.gz -C $GOTARGET
sudo rm /tmp/go.tar.gz

source $PROFILE
sudo ldconfig
go version

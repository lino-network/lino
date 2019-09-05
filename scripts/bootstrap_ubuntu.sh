#!/usr/bin/env bash

### install golang

GOSOURCE=https://dl.google.com/go/go1.13.linux-amd64.tar.gz
GOTARGET=/usr/local
GOPATH=\$HOME/go
PROFILE=/home/ubuntu/.profile

curl -sSL $GOSOURCE -o /tmp/go.tar.gz
sudo tar -xzf /tmp/go.tar.gz -C $GOTARGET
sudo rm /tmp/go.tar.gz

# apply environment configuration to the user's .profile
touch $PROFILE
printf "\n" >> $PROFILE
printf "# golang configuration\n" >> $PROFILE
printf "export GOROOT=$GOTARGET/go\n" >> $PROFILE
printf "export GOPATH=$GOPATH\n" >> $PROFILE
printf "export PATH=\$PATH:$GOTARGET/go/bin:$GOPATH/bin\n" >> $PROFILE

source $PROFILE
sudo ldconfig
go version

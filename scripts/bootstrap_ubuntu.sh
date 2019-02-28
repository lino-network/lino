#!/usr/bin/env bash

### install golang

GOSOURCE=https://dl.google.com/go/go1.11.5.linux-amd64.tar.gz
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

# install snappy & cleveldb
sudo apt-get update
sudo apt-get install -y build-essential gcc g++ make

sudo apt-get install -y libsnappy-dev
homedir="$PWD"

wget https://github.com/google/leveldb/archive/v1.20.tar.gz && \
  tar -zxvf v1.20.tar.gz && \
  cd leveldb-1.20/ && \
  make && \
  sudo cp -r out-static/lib* out-shared/lib* /usr/local/lib/ && \
  cd include/ && \
  sudo cp -r leveldb /usr/local/include/ && \
  sudo ldconfig && \
  rm -f v1.20.tar.gz

cd "$homedir"

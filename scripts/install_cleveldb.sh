#!/usr/bin/env bash

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
  sudo ldconfig

cd "$homedir"

rm -f v1.20.tar.gz
rm -rf leveldb-1.20

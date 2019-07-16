#!/usr/bin/env bash

sudo apt-get update
sudo apt-get install -y build-essential gcc g++ make

sudo apt-get install -y libgflags-dev libzstd-dev
homedir="$PWD"

wget -O rocks.tar.gz https://github.com/facebook/rocksdb/archive/v6.1.2.tar.gz && \
  tar -zxvf rocks.tar.gz && mv rocksdb-6.1.2 rocksdb \
  cd rocksdb/ && \
  make -j4 && \
  sudo cp librocksdb.a /usr/local/lib/ && \
  cd include/ && \
  sudo cp -r rocksdb /usr/local/include/ && \
  sudo ldconfig

cd "$homedir"

rm -f rocksdb.tar.gz
rm -rf rocksdb

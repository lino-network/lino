#!/usr/bin/env bash

sudo apt-get update
sudo apt-get install -y build-essential gcc g++ make

sudo apt-get install -y libgflags-dev libzstd-dev libsnappy-dev zlib1g-dev libbz2-dev liblz4-dev
homedir="$PWD"

wget -O rocks.tar.gz https://github.com/facebook/rocksdb/archive/v6.1.2.tar.gz && \
  tar -zxvf rocks.tar.gz && \
  cd rocksdb-6.1.2/ && \
  make static_lib -j4 && \
  sudo cp librocksdb.a /usr/local/lib/ && \
  cd include/ && \
  sudo cp -r rocksdb /usr/local/include/ && \
  sudo ldconfig

cd "$homedir"

rm rocks.tar.gz
rm -rf rocksdb-6.1.2

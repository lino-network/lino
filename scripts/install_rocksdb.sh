#!/usr/bin/env bash

sudo apt-get update
sudo apt-get install -y build-essential gcc g++ make

sudo apt-get install -y libgflags-dev

homedir="$PWD"
wget -O zstd.tar.gz https://github.com/facebook/zstd/archive/v1.4.3.tar.gz && \
    tar -zxvf zstd.tar.gz && \
    cd zstd-1.4.3/ && \
    sudo make install -j4 && \
    sudo ldconfig
cd "$homedir"

wget -O rocks.tar.gz https://github.com/facebook/rocksdb/archive/v6.2.2.tar.gz && \
  tar -zxvf rocks.tar.gz && \
  cd rocksdb-6.2.2/ && \
  make static_lib -j4 && \
  sudo cp librocksdb.a /usr/local/lib/ && \
  cd include/ && \
  sudo rm -rf /usr/local/include/rocksdb && \
  sudo cp -r rocksdb /usr/local/include/ && \
  sudo ldconfig

cd "$homedir"

rm rocks.tar.gz
rm -rf rocksdb-6.2.2
rm zstd.tar.gz
rm -rf zstd-1.4.3

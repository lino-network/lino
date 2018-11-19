#!/bin/sh
wget https://s3.us-east-2.amazonaws.com/lino-blockchain-data/data.tar.gz
tar -xzvf data.tar.gz -C data/
rm -rf data.tar.gz
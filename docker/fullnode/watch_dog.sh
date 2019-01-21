#!/bin/bash
if [ -d /root/.lino/config ]
then
  echo "already init"
else
  ./lino init
  cp genesis.json /root/.lino/config/genesis.json
  cp config.toml /root/.lino/config/config.toml
  sed -i "11s/.*/moniker=\"$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 12 | head -n 1)\"/" ~/.lino/config/config.toml
  ./lino unsafe_reset_all
  wget https://s3-eu-west-1.amazonaws.com/lino-blockchain-data-ireland/data.tar.gz
  mkdir -p /root/.lino/data/
  tar -xzvf data.tar.gz -C /root/.lino/data/
  rm -rf data.tar.gz
fi

./lino start --log_level=error &
pid=$!
last_height=0
while true
 do
    sleep 30s
    status=$(curl --max-time 10 -s -o /dev/null -w "%{http_code}" http://localhost:26657/health)
    if [ $status -eq 200 ]
    then
      echo node is running
    else
      echo node is down
      kill -9 $pid
      sleep 10s
      ./lino start --log_level=error &
      pid=$!
      healthy=false
    fi
done
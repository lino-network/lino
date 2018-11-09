#!/bin/bash
/root/.local/bin/aws s3 cp s3://lino-blockchain-data/data.tar.gz data.tar.gz
RUN rm -rf /root/.lino/data
mkdir -p /root/.lino/data/
tar -xzvf data.tar.gz -C /root/.lino/data/
rm -rf data.tar.gz
./lino start --log_level=error &
pid=$!
last_height=0
while true
 do
    sleep 30s
    status=$(curl --max-time 10 -s -o /dev/null -w "%{http_code}" http://localhost:26657)
    height=$(curl --max-time 10 http://localhost:26657/status | jq '. | .result.sync_info.latest_block_height')
    echo "running at height $height"
    if [ $status -eq 200 ]
    then
      echo node is running
      if [ "$height" = "$last_height" ]
      then
        echo node is at the same height about 30s
        kill -9 $pid
        sleep 10s
        ./lino start --log_level=error &
        pid=$!
        healthy=false
      else
        echo node is healthy
        last_height=$height
      fi
    else
      echo node is down
      kill -9 $pid
      sleep 10s
      ./lino start --log_level=error &
      pid=$!
      healthy=false
    fi
done
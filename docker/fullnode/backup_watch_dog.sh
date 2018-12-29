#!/bin/bash
./lino start --log_level=error &>> node.log &
healthy=true
pid=$!
last_height=0
echo pid is $pid
counter=0
while true
 do
    sleep 30s
    catching_up=$(curl --max-time 10 http://localhost:26657/status | jq '. | .result.sync_info.catching_up')
    if [ "$catching_up" = true ] ; then
      echo 'still catching up!'
      continue
    fi
    status=$(curl --max-time 10 -s -o /dev/null -w "%{http_code}" http://localhost:26657)
    height=$(curl --max-time 10 http://localhost:26657/status | jq '. | .result.sync_info.latest_block_height')
    if [ $status -eq 200 ]
    then
      echo "node is running at height $height"
      if [ "$height" = "$last_height" ]
      then
        echo node is at the same height about 30s
        kill -INT $pid
        sleep 10s
        ./lino start --log_level=error &>> node.log &
        pid=$!
        healthy=false
        counter=0
      else
        echo node is healthy
        last_height=$height
        counter=$((counter+1))
        echo "counter is $counter"
        if [ "$counter" = 6 ]
        then
          echo "counter reach 6! $counter"
          counter=0
          kill -INT $pid
          tar -czvf data.tar.gz -C ~/.lino/data .
          mv data.tar.gz /backup/data_$(date +%F-%H:%M).tar.gz
          numOfFile=$(ls /backup | wc -l)
          if [ "$numOfFile" -gt 3 ]
          then
            rm /backup/$(ls -1 /backup | head -n 1)
          fi
        fi
        if [ "$healthy" = false ]
        then
          healthy=true
        fi
      fi
    else
      echo node is down
      kill -INT $pid
      sleep 10s
      ./lino start --log_level=error &>> node.log &
      pid=$!
      healthy=false
      counter=0
    fi
done
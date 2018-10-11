#!/bin/bash
node=$1
sed -i "11s/.*/moniker=\"$node\"/" config.toml
mv config.toml ~/.lino/config/
./lino start --log_level=error &
pid=$!
echo pid is $pid
while true
 do
    sleep 30s
    status=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:26657)
    if [ $status -eq 200 ]
    then
      echo node is running
    else
      echo node is down
      curl -s -d "to=%2B13105629182" -d "body=node: $node is down!" https://utils.lib.id/sms@1.0.4/
      curl -s -d "to=%2B18588595067" -d "body=node: $node is down!" https://utils.lib.id/sms@1.0.4/
      curl -s -d "to=%2B13106941760" -d "body=node: $node is down!" https://utils.lib.id/sms@1.0.4/
      curl -s -d "to=%2B19178737414" -d "body=node: $node is down!" https://utils.lib.id/sms@1.0.4/
      kill -9 $pid
      ./lino start &
      let pid=$!
    fi
done
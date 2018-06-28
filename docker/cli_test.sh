#/bin/bash
while ! nc -z lino 26657;
    do
        echo sleeping;
        sleep 1;
    done;
    echo Connected!

echo "show username"
./linocli username lino --node=lino:26657


echo "register user"
./linocli register --user=newuser --referrer=lino --amount=10 --priv-key=E1B0F79A20AECFA4549861801551DB876C3D54A1A729A030CC07BDEEB8935294CD51D6ADE2 --chain-id=lino-test --node=lino:26657


echo "transfer to new user"
./linocli transfer --sender=lino --receiver=newuser --amount=10 --memo=memo --priv-key=E1B0F79A20AECFA4549861801551DB876C3D54A1A729A030CC07BDEEB8935294CD51D6ADE2 --chain-id=lino-test --sequence=1 --node=lino:26657
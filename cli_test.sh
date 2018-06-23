#/bin/bash
while ! nc -z lino 46657;
    do
        echo sleeping;
        sleep 1;
    done;
    echo Connected!

echo "test show username"

./linocli username lino --node=lino:46657
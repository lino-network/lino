# Basic run through of using linocli....

To keep things clear, let's have two shells...

`$` is for linocoin (server), `%` is for linocli (client)

## Set up your linocli with a new key

```
% export BCHOME=~/.democli
% linocli keys new demo
% linocli keys get demo -o json
```

And set up a few more keys for fun...

```
% linocli keys new buddy
% linocli keys list
% ME=$(linocli keys get demo | awk '{print $2}')
% YOU=$(linocli keys get buddy | awk '{print $2}')
```

## Set up a clean linocoin, initialized with your account

```
$ export BCHOME=~/.demoserve
$ linocoin init $ME
$ linocoin start
```

## Connect your linocli the first time

```
% linocli init --chain-id test_chain_id --node tcp://localhost:46657
```

## Check your balances...

```
% linocli query account $ME
% linocli query account $YOU
```

## Send the money

```
% linocli tx send --name demo --amount 1000mycoin --sequence 1 --to $YOU
-> copy hash to HASH
% linocli query tx $HASH
% linocli query account $YOU
```

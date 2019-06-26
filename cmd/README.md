# Lino Blockchain Command

This cmd directory contains two command line tool: lino and linocli. _lino_ is used to luanch the Lino Blockchain node. _linocli_ can be used to interact with Lino Blockchain.

# Launch Blockchain
## Generate genesis file
```
$ ./lino init
```
## Start generate block as a validator
```
$ ./lino start
```

# Launch Client
## Transfer coin to a user
```
$ ./linocli transfer --sender=<username>  --receiver=<receiver> --amount=1 --chain-id=<chain id> --sequence=<sender's sequence number>
```

## Register an account
```
$ ./linocli register --referrer=<username> --user=<new user> --amount=1 --chain-id=<chain id> --sequence=<sender's sequence number>
```

## Follow & Unfollow
Follow
```
$ ./linocli follow --follower=<me> --followee=<other> --is-follow=true --sequence= --chain-id=<chain id> --sequence=<sender's sequence number>
```
Unfollow
```
$ ./linocli follow --follower=<me> --followee=<other> --is-follow=false --sequence= --chain-id=<chain id> --sequence=<sender's sequence number>
```
## Query Account
Check Bank
```
$ ./linocli username XXXXXXXX
```


## Others
List all keys 
```
$ ./linocli keys list
```



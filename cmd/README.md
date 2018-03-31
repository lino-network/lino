# Lino Blockchain Command

This cmd directory contains two command line tool: lino and linocli. _lino_ is used to luanch the Lino Blockchain. _linocli_ can be used to interact with Lino Blockchain.

# Luanch Blockchain
## Generate genesis file
```
$ ./lino init
```
## Start generate block as a validator
```
$ ./lino start
```

# Luanch Client
## Generate key pair

```
$ ./linocli keys add USERNAME
```
Enter a passphrase for this key.
Keep the generated address(for transfer) and seed phrase(for recover)

## Transfer coin to your address
```
$ ./linocli transfer --name=  --receiver_addr= --amount= --chain-id= --sequence=
```
User other account to transfer some coin to your address

## Register an account
```
$ ./linocli register --name= --chain-id=
```
## Follow & Unfollow
Follow
```
$ ./linocli follow --followee= --name= --sequence= --chain-id=
```
Unfollow
```
./linocli follow --is_follow=false --followee= --name= --sequence= --chain-id=
```
## Query Info
Check Bank
```
$ ./linocli address XXXXXXXX
```

Check Account
```
$ ./linocli username XXXXXXXX
```

## Others
List all keys 
```
$ ./linocli keys list
```



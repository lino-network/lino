# Authorization

## asymmetric cryptography

Like Bitcoin, account authorization will be verified by asymmetric cryptography. The account owner keeps their private key in secret, then register their public key on blockchain with register msg. The public key will be stored on the blockchain with their account. Next time when account owner send msg in the name of an account. The account owner sign the msg using private key and blockchain verify the signature by using public key.

## sequence number

To prevent replay attack, when user sign a msg, it should always contain the chain id and sequence number. The sequence number is increment only and starts from 0. If msg signed with an invalid sequence, it is rejected. If the sequence is valid, the msg pass the authentication phase and the sequence number is increased by 1. Next time the msg should be signed with increased sequence number. For more details see #160. If msg is signed by granted app, the app should use user’s sequence number.

## permission

The asymmetric cryptography is the only way to check the authentication of a msg, which means the user has to authenticate all msgs by providing their private key. It is neither safe nor convenient. To help user with their key management, on Lino Blockchain, each user has three keys with different permission scopes: reset key, transaction key and app key. For all user activities and their permission scope see User Activity.

### reset key

Reset key is used to recover account only (reset all key pairs). The reset key should be kept absolutely secret.

### transaction key

Transaction key has highest permission between left two keys. It can sign for all msg but recover msg. Balance related msg (except preauth and claim msg) can only be signed by transaction key. User should be careful when using their transaction key.

### app key

The app key can’t access account balance. The post permission can be granted to the app so app can use their app key sign msg for user (like publish post for user). User should be aware of their grant permission and set limitation to them (like number of signing times and expired date).

## grant permission

User can grant his pre-authorization or post permission to app developer so app can sign the msg for user using app’s transaction or app key. User can set the expire time and amount of preauthorization. User can revoke permission at will.

* [JSON RPC API](#JSON-RPC-API)  
* [JavaScript API](#JavaScript-API)
* [JSON-RPC Endpoint](#JSON-RPC-Endpoint)
* [JSON-RPC API Reference](#JSON-RPC-API-Reference)
    * [Blockchain Status](#Blockchain-Status)
    * [Block Info](#Block-Info)
    * [Tx Info](#Tx-Info)
    * [Account Info](#Account-Info)
    * [Account Bank](#Account-Bank)
    * [Post Info](#Post-Info)
    * [Stake Info](#Stake-Info)
    * [Validator Info](#Validator-Info)

# JSON RPC API
## JavaScript API
To talk to an Lino node from inside a JavaScript application use the [lino-js](https://github.com/lino-network/lino-js) library, which gives a convenient interface for the RPC methods. See the  [lino-js](https://github.com/lino-network/lino-js) library for more.

## JSON-RPC Endpoint
Default JSON-RPC endpoints: https://fullnode.lino.network/

## JSON-RPC API Reference
### Blockchain Status
Returns current blockchain status. Including fullnode information and sync info.

#### Parameters
none

#### Returns
String - Current blockchain status.

#### Example
```
// Request
$ curl -X POST --data-binary '{"jsonrpc":"2.0","method":"status"}' "https://fullnode.lino.network"

// Result
{
  "jsonrpc": "2.0",
  "id": "",
  "result": {
    "node_info": {
      "protocol_version": {
        "p2p": "7",
        "block": "10",
        "app": "0"
      },
      "id": "58f7c3e342647155a3c2b3635807f6890be33af2",
      "listen_addr": "tcp://0.0.0.0:26656",
      "network": "lino-testnet-upgrade2",
      "version": "0.32.2",
      "channels": "40202122233038",
      "moniker": "73pEAeSPBCAD",
      "other": {
        "tx_index": "on",
        "rpc_address": "tcp://0.0.0.0:26657"
      }
    },
    "sync_info": {
      "latest_block_hash": "C4B541D9019E87CCAB7A643BC9802BA249404776D234874FC33A3CD6862234B2",
      "latest_app_hash": "562CC637AF432F540C5215AEFCC12E9294CDF99E3988D708E14CA0601321302C",
      "latest_block_height": "237594",
      "latest_block_time": "2019-09-29T23:53:57.207992926Z",
      "catching_up": false
    },
    "validator_info": {
      "address": "852FB8F81F013BC350C8EDF34AFAAFFE8F6F77BC",
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "E5sXbFq7VgoECX99ZecDETkBeTlDzGyAmAXyl8yb50M="
      },
      "voting_power": "0"
    }
  }
}
```

### Block Info
Returns block information for a specific block. Including block meta and all transactions in the block.

#### Parameters
1. height - the height of the block.

#### Returns
String - Block information.

#### Example
```
// Request
$ curl -X POST --data-binary '{"jsonrpc":"2.0","method":"block", "params":{"height":"1"}}' "https://fullnode.lino.network" -s 'https://fullnode.lino.network/status'

// Result
{
  "jsonrpc": "2.0",
  "id": "jsonrpc-client",
  "result": {
    "block_meta": {
      "block_id": {
        "hash": "AA1D14CA41F7DE01BC47D9E8443779F212B5CCE959CEB930F5DED26AC12223D9",
        "parts": {
          "total": "1",
          "hash": "9DDF48A11FBAFDD2FE61C65C4A10D10D40F1B13A1FFE8604CBC34D5035E0AD9E"
        }
      },
      "header": {
        "version": {
          "block": "10",
          "app": "0"
        },
        "chain_id": "lino-testnet-upgrade2",
        "height": "1",
        "time": "2019-09-19T19:12:05.141699565Z",
        "num_txs": "0",
        "total_txs": "0",
        "last_block_id": {
          "hash": "",
          "parts": {
            "total": "0",
            "hash": ""
          }
        },
        "last_commit_hash": "",
        "data_hash": "",
        "validators_hash": "F29730F7417C82A26FE3FA55D07E2D5F771FC2B0EA2E387E2E2A3CF10B1B571F",
        "next_validators_hash": "F29730F7417C82A26FE3FA55D07E2D5F771FC2B0EA2E387E2E2A3CF10B1B571F",
        "consensus_hash": "048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
        "app_hash": "",
        "last_results_hash": "",
        "evidence_hash": "",
        "proposer_address": "1E1B60F12C837BB35218E1F370B935FEBC17B8A0"
      }
    },
    "block": {
      "header": {
        "version": {
          "block": "10",
          "app": "0"
        },
        "chain_id": "lino-testnet-upgrade2",
        "height": "1",
        "time": "2019-09-19T19:12:05.141699565Z",
        "num_txs": "0",
        "total_txs": "0",
        "last_block_id": {
          "hash": "",
          "parts": {
            "total": "0",
            "hash": ""
          }
        },
        "last_commit_hash": "",
        "data_hash": "",
        "validators_hash": "F29730F7417C82A26FE3FA55D07E2D5F771FC2B0EA2E387E2E2A3CF10B1B571F",
        "next_validators_hash": "F29730F7417C82A26FE3FA55D07E2D5F771FC2B0EA2E387E2E2A3CF10B1B571F",
        "consensus_hash": "048091BC7DDC283F77BFBF91D73C44DA58C3DF8A9CBC867405D8B7F3DAADA22F",
        "app_hash": "",
        "last_results_hash": "",
        "evidence_hash": "",
        "proposer_address": "1E1B60F12C837BB35218E1F370B935FEBC17B8A0"
      },
      "data": {
        "txs": null
      },
      "evidence": {
        "evidence": null
      },
      "last_commit": {
        "block_id": {
          "hash": "",
          "parts": {
            "total": "0",
            "hash": ""
          }
        },
        "precommits": null
      }
    }
  }
}
```
To parse a transaction in the block:
```
// Request
$ curl -X POST --data-binary '{"jsonrpc":"2.0","method":"block", "params":{"height":"20000"}}' "https://fullnode.lino.network" | jq -r .result.block.data.txs[0] | base64 -d

// Result
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "lino/register",
        "value": {
          "referrer": "dlivetv-50",
          "register_fee": "5.05",
          "new_username": "rizqienb1",
          "new_reset_public_key": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "A4hcZmY4TTMMLWu1Vbid6atcgZQmra9xx/cSfsL1wiqw"
          },
          "new_transaction_public_key": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "A/xXr5+3aKxqRMkf8c+2VsXgrdMvc3hbE4lmz49YwBne"
          },
          "new_app_public_key": {
            "type": "tendermint/PubKeySecp256k1",
            "value": "AoeI2TmhEpoCeySUQDkCzcFSqzUEmBhvRwMJg3U91REW"
          }
        }
      }
    ],
    "fee": {
      "amount": [
        {
          "denom": "linocoin",
          "amount": "100000"
        }
      ],
      "gas": "0"
    },
    "signatures": [
      {
        "pub_key": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AilEeb1MEK2EIObumdHHvVl6Qtert1OxRYVuZvRFMnj2"
        },
        "signature": "zFJymbFOFZZ8GpzQHQhU/vAwyDYv7kwMuxzWlJvXv3VxKcJzgxdYH1Riw9fLF64kx0m9mPthnJdabx7S5cM2mQ=="
      }
    ],
    "memo": ""
  }
}
```

### Tx Info
Returns a specific transaction status and execution result.

#### Parameters
1. hash - the hash of a transaction.

#### Returns
String - Tx information.

#### Example
```
// Request
$  curl -X POST --data-binary '{"jsonrpc":"2.0","method":"tx", "params":{"hash":"3fTjLeArr8uLbRYxL6zJiRRhvpp+NX9FSqLAgdAY+9A="}}' "https://fullnode.lino.network"

// Result
{
  "jsonrpc": "2.0",
  "id": "jsonrpc-client",
  "result": {
    "hash": "DDF4E32DE02BAFCB8B6D16312FACC9891461BE9A7E357F454AA2C081D018FBD0",
    "height": "237788",
    "index": 0,
    "tx_result": {
      "log": "[{\"msg_index\":0,\"success\":true,\"log\":\"\"}]",
      "gasUsed": "215100",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "YWN0aW9u",
              "value": "VHJhbnNmZXJNc2c="
            }
          ]
        }
      ]
    },
    "tx": "eyJ0eXBlIjoiYXV0aC9TdGRUeCIsInZhbHVlIjp7Im1zZyI6W3sidHlwZSI6Imxpbm8vdHJhbnNmZXIiLCJ2YWx1ZSI6eyJzZW5kZXIiOiJkbGl2ZXR2IiwicmVjZWl2ZXIiOiJkbGl2ZXR2LTEzIiwiYW1vdW50IjoiMy40NSIsIm1lbW8iOiIifX1dLCJmZWUiOnsiYW1vdW50IjpbeyJkZW5vbSI6Imxpbm9jb2luIiwiYW1vdW50IjoiMTAwMDAwIn1dLCJnYXMiOiIwIn0sInNpZ25hdHVyZXMiOlt7InB1Yl9rZXkiOnsidHlwZSI6InRlbmRlcm1pbnQvUHViS2V5U2VjcDI1NmsxIiwidmFsdWUiOiJBM2hzYkErUGVHVUZMazNmcFVvSDdVVXFvU1Iwc2pDVUNOU2IvR3B0WjVzRSJ9LCJzaWduYXR1cmUiOiJhOTVRU1lnN1l4QzNOamdnWEVCSitiQ2VTM3EzUEhEWGRqeFlRT2wwZDZ3Q0pMTk1yQjhhUGJ1MDJBRmtSdXpreHFLSnlTUFU4Z2I3V2NmeEw4OTI2QT09In1dLCJtZW1vIjoiIn19"
  }
}
```

Transaction hash can be parsed by following steps
```
// Request
$ curl -X POST --data-binary '{"jsonrpc":"2.0","id":"jsonrpc-client","method":"block", "params":{"height":"237788"}}' "https://fullnode.lino.network" | jq -r .result.block.data.txs[0] | base64 -d | sha256sum | xxd -r -p | base64

// Result
3fTjLeArr8uLbRYxL6zJiRRhvpp+NX9FSqLAgdAY+9A=
```


### Account Info
Returns a specific user's account information.

#### Parameters
1. username - the username of a Lino Blockchain user.

#### Returns
String - Account information, which includes username, create time in unix, public keys and address. Address can be derived from transaction public key.

#### Example
```
// Request, username is `ytu`
$ curl -X POST --data-binary '{"jsonrpc":"2.0","method":"abci_query","params":{"height":"0","trusted":false,"path":"/custom/account/info/ytu","data":""}}' "https://fullnode.lino.network" | jq -r .result.response.value | base64 -d | jq .

// Result
{
  "username": "ytu",
  "created_at": "1537817595",
  "signing_key": {
    "type": "tendermint/PubKeySecp256k1",
    "value": "AoxfvcigEx+LtU2t0aAanloux5CA5kjORVvBgKVt/Hip"
  },
  "transaction_key": {
    "type": "tendermint/PubKeySecp256k1",
    "value": "Awa7WFs9Oeyl5skmqmlV+eaN95ajWkQNbL8wzRdkx9+j"
  },
  "address": "lino1722lj3a89nnmt8teadp98h5rkvrcsc4e2ulm9s"
}
```

### Account Bank
Returns a specific user's bank information.

#### Parameters
1. username - the username of a Lino Blockchain user.

#### Returns
String - Bank information, which includes bank balance (in Lino Coin, 1 LINO = 100000 Lino Coin), frozen money list (pending Lino), public key (same as transaction public key above), sequence number and username.

#### Example
```
// Request, username is `ytu`
$ curl -X POST --data-binary '{"jsonrpc":"2.0","method":"abci_query","params":{"height":"0","trusted":false,"path":"/custom/account/bank/ytu","data":""}}' "https://fullnode.lino.network" | jq -r .result.response.value | base64 -d | jq .

// Result
{
  "saving": {
    "amount": "117057339"
  },
  "frozen_money_list": [
    {
      "amount": {
        "amount": "1000000"
      },
      "start_at": "1539034248",
      "times": "12",
      "interval": "604800"
    },
    
  ],
  "public_key": {
    "type": "tendermint/PubKeySecp256k1",
    "value": "Awa7WFs9Oeyl5skmqmlV+eaN95ajWkQNbL8wzRdkx9+j"
  },
  "sequence": "1865",
  "username": "ytu"
}
```


### Post Info
Returns a specific post's information.

#### Parameters
1. permlink - the permlink of the post. Permlink = username + "#" + postID.

#### Returns
String - Post information, which includes author, post id, title, content, create time, etc.

#### Example
```
// Request, permlink is `pika35#VxWqSm2Wg`
$  curl -X POST --data-binary '{"jsonrpc":"2.0","method":"abci_query","params":{"height":"0","trusted":false,"path":"/custom/post/info/pika35#VxWqSm2Wg","data":""}}' "https://fullnode.lino.network" | jq -r .result.response.value | base64 -d | jq .

// Result
{
  "post_id": "VxWqSm2Wg",
  "title": "💛 ⚡UYKUSUZ VE DENGESİZ",
  "content": "",
  "author": "pika35",
  "created_by": "dlivetv",
  "created_at": "1569793433",
  "updated_at": "1569793433",
  "is_deleted": false
}
```


### Stake Info
Returns a user's stake information.

#### Parameters
1. username - the username of a Lino Blockchain user.

#### Returns
String - Stake information, which includes total Lino stake (in Lino Coin, 1 LINO = 100000 Lino Coin), delegation info (ytu 09/29/2019: deprecated in next update), duty and frozen amount (ytu 09/29/2019: enable in next update).

#### Example
```
// Request, username is `dlivetv`
$  curl -X POST --data-binary '{"jsonrpc":"2.0","method":"abci_query","params":{"height":"0","trusted":false,"path":"/custom/vote/voter/dlivetv","data":""}}' "https://fullnode.lino.network" | jq -r .result.response.value | base64 -d | jq .

// Result
{
  "username": "dlivetv",
  "lino_stake": {
    "amount": "100000000000"
  },
  "delegated_power": {
    "amount": "0"
  },
  "delegate_to_others": {
    "amount": "0"
  },
  "last_power_change_at": "1568322539",
  "interest": {
    "amount": "0"
  },
  "duty": "0",
  "frozen_amount": {
    "amount": "0"
  }
}
```

### Stake Info
Returns a user's stake information.

#### Parameters
1. username - the username of a Lino Blockchain user.

#### Returns
String - Stake information, which includes total Lino stake (in Lino Coin, 1 LINO = 100000 Lino Coin), delegation info (ytu 09/29/2019: deprecated in next update), duty and frozen amount (ytu 09/29/2019: enable in next update).

#### Example
```
// Request, username is `dlivetv`
$  curl -X POST --data-binary '{"jsonrpc":"2.0","method":"abci_query","params":{"height":"0","trusted":false,"path":"/custom/vote/voter/dlivetv","data":""}}' "https://fullnode.lino.network" | jq -r .result.response.value | base64 -d | jq .

// Result
{
  "username": "dlivetv",
  "lino_stake": {
    "amount": "100000000000"
  },
  "delegated_power": {
    "amount": "0"
  },
  "delegate_to_others": {
    "amount": "0"
  },
  "last_power_change_at": "1568322539",
  "interest": {
    "amount": "0"
  },
  "duty": "0",
  "frozen_amount": {
    "amount": "0"
  }
}
```

### Validator Info
Returns a validator's information.

#### Parameters
1. username - the username of a Lino Blockchain user.

#### Returns
String - Validator information, which number of poduced blocks, deposit (ytu 09/29/2019: deprecated in next update), public key, commit power, etc.

#### Example
```
// Request, username is `validator1`
$ curl -X POST --data-binary '{"jsonrpc":"2.0","method":"abci_query","params":{"height":"0","trusted":false,"path":"/custom/validator/validator/validator1","data":""}}' "https://fullnode.lino.network" | jq -r .result.response.value | base64 -d | jq .

// Result
{
  "ABCIValidator": {
    "address": "Hhtg8SyDe7NSGOHzcLk1/rwXuKA=",
    "power": "1000"
  },
  "pubkey": {
    "type": "tendermint/PubKeyEd25519",
    "value": "4rgh/IbevTzo/2s3YJip1F/ih0gm153mLrIkKZTrNhI="
  },
  "username": "validator1",
  "deposit": {
    "amount": "29340000000"
  },
  "absent_commit": "0",
  "byzantine_commit": "0",
  "produced_blocks": "9267485",
  "link": ""
}
```

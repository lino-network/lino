# Consensus

## Tendermint

Currently the Lino Blockchain is built on the consensus engine: Tendermint. Tendermint is software for securely and consistently replicating an application on many machines. It implements the Byzantine fault tolerance (BFT) which tolerate up to ⅓ of validators failed in arbitrary ways. To achieve the better performance we predefine the size of validator set is 22. The validator can be changed by the locked Lino Stake.

For more information about tendermint, please refer to [here](https://tendermint.com/docs/).

## ABCI and Cosmos SDK

The tendermint consensus engine communicates with the Lino Blockchain via a socket protocol that satisfies the ABCI. Lino Blockchain implement ABCI through Cosmos SDK, which is a platform for building multi-asset Proof-of-Stake cryptocurrencies. The goals of the SDK are to abstract away the complexities of building a Tendermint ABCI application in Golang and to provide a framework for building interoperable blockchain applications in the Cosmos Network.

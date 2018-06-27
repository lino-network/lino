![banner](docs/graphics/banner.png)

# Lino Blockchain

[![API Reference](https://godoc.org/github.com/cosmos/cosmos-sdk?status.svg
)](https://docs.google.com/document/d/1Ytd57axPfJ13TSGVU_Yykv8ijW_VuWtx1s79ny6i5M8)
[![LoC](https://tokei.rs/b1/github/lino-network/lino)](https://github.com/lino-network/lino).

Welcome to Lino Blockchain. Lino aims to create a decentralized autonomous content economy, where content value can be recognized efficiently and all contributors can be incentivized directly and effectively to promote long-term economic growth. For more information about Lino refer to our [website](https://lino.network/).


## Get Source Code

    mkdir -p $GOPATH/src/github.com/lino-network
    cd $GOPATH/src/github.com/lino-network
    git clone https://github.com/lino-network/lino.git
    cd lino

## Get Tools & Dependencies

    dep ensure

## Compile Lino Blockchain Node

    cd cmd/lino
    go build


## Compile Lino Blockchain Client

    cd cmd/linocli
    go build

## run test image

    docker-compose up --build



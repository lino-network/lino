FROM golang:1.10

RUN apt-get update && \
    apt-get -y install wget curl sudo make && \
    apt-get install -y git

# get lino core code
RUN mkdir -p src/github.com/lino-network
WORKDIR src/github.com/lino-network
RUN git clone https://github.com/lino-network/lino.git
WORKDIR lino

# golang dep
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

# replace customize file
COPY docker/fullnode/http_server.go ./vendor/github.com/tendermint/tendermint/rpc/lib/server/http_server.go
COPY docker/fullnode/iavlstore.go ./vendor/github.com/cosmos/cosmos-sdk/store/iavlstore.go
WORKDIR cmd/lino
RUN go build
RUN ./lino init

COPY docker/fullnode/genesis.json /root/.lino/config/genesis.json
COPY docker/fullnode/config.toml ./config.toml

EXPOSE 26656
EXPOSE 26657

COPY docker/fullnode/watch_dog.sh watch_dog.sh
RUN chmod +x watch_dog.sh

ENTRYPOINT ["./watch_dog.sh"]
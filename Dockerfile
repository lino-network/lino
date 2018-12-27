FROM golang:1.10

RUN apt-get update && \
    apt-get -y install wget curl sudo make && \
    apt-get install -y git && \
    apt-get install -y jq && \
    apt install -y python3 \
        python3-pip \
        python3-setuptools && \
        pip3 install --upgrade pip && \
        apt-get clean

# get lino core code
RUN pip install awscli --upgrade --user
RUN mkdir -p src/github.com/lino-network
WORKDIR src/github.com/lino-network
RUN git clone https://github.com/lino-network/lino.git
WORKDIR lino
RUN git fetch
RUN git checkout staging

# golang dep
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

# replace customize file
COPY docker/fullnode/http_server ./vendor/github.com/tendermint/tendermint/rpc/lib/server/http_server.go
COPY docker/fullnode/iavlstore ./vendor/github.com/cosmos/cosmos-sdk/store/iavlstore.go
COPY docker/fullnode/baseapp ./vendor/github.com/cosmos/cosmos-sdk/baseapp/baseapp.go
WORKDIR cmd/lino
RUN go build

COPY docker/fullnode/genesis_staging.json genesis.json
COPY docker/fullnode/config_staging.toml config.toml

EXPOSE 26656
EXPOSE 26657

COPY docker/fullnode/watch_dog.sh watch_dog.sh
RUN chmod +x watch_dog.sh

ENTRYPOINT ["./lino", "start", "--log_level=error"]
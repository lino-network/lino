FROM golang:1.12

RUN apt-get update && \
    apt-get install -y make tar sudo wget curl

RUN mkdir -p src/github.com/lino-network/lino
WORKDIR src/github.com/lino-network/lino

COPY . .
RUN make get_tools
RUN make install

RUN lino init
COPY genesis/upgrade5/config.toml  /root/.lino/config/config.toml
COPY genesis/upgrade5/genesis.json /root/.lino/config/genesis.json
RUN cd /root/.lino && wget https://lino-blockchain-opendata.s3.amazonaws.com/prd/prevstates.tar.gz
RUN cd /root/.lino && tar -xf prevstates.tar.gz

# prometheus if enabled
EXPOSE 26660
# p2p
EXPOSE 26656
# tendermint rpc
EXPOSE 26657
# abci app
EXPOSE 26658

CMD ["lino", "unsafe-reset-all"]
CMD ["lino", "start"]

FROM golang:1.9


RUN apt-get update && \
    apt-get -y install wget curl sudo && \
    apt-get install -y git


RUN mkdir -p src/github.com/lino-network/lino
WORKDIR src/github.com/lino-network/lino
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep

RUN dep ensure

WORKDIR cmd/lino
RUN go build

RUN rm -rf /root/.lino
RUN mkdir -p /root/.lino/config/
COPY docker/genesis.json /root/.lino/config/genesis.json
COPY docker/priv_validator.json /root/.lino/config/priv_validator.json

EXPOSE 26657

CMD ["./lino", "unsafe_reset_all"]
CMD ["./lino", "start", "--log_level=error"]
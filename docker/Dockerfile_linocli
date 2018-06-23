FROM golang:1.9


RUN apt-get update && \
    apt-get -y install wget curl netcat sudo && \
    apt-get install -y git


RUN mkdir -p src/github.com/lino-network/lino
WORKDIR src/github.com/lino-network/lino
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep

RUN dep ensure

WORKDIR cmd/linocli
RUN go build
COPY docker/cli_test.sh .
RUN chmod +x cli_test.sh
CMD ./cli_test.sh
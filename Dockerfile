FROM golang:1.9

ENV GOBIN /go/bin

RUN mkdir /go/src/github.com/lino-network/lino
ADD . /go/src/github.com/lino-network/lino
WORKDIR /go/src/github.com/lino-network/lino

RUN go get -u github.com/golang/dep/...
RUN dep ensure

EXPOSE port 46656
EXPOSE port 46657

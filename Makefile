PACKAGES=$(shell go list ./... | grep -v '/vendor/')
COMMIT_HASH := $(shell git rev-parse --short HEAD)

all: get_tools get_vendor_deps install build test

get_tools:
	go get github.com/golang/dep/cmd/dep
build:
	go build -o bin/linocli cmd/linocli/main.go && go build -o bin/linod cmd/lino/main.go
build_prd:
	CGO_LDFLAGS="-lsnappy" CGO_ENABLED=1 go build -ldflags "-X github.com/tendermint/tendermint/version.GitCommit=`git rev-parse --short=8 HEAD`" -tags "tendermint gcc cgo"  -o bin/linod cmd/lino/main.go
	CGO_LDFLAGS="-lsnappy" CGO_ENABLED=1 go build -ldflags "-X github.com/tendermint/tendermint/version.GitCommit=`git rev-parse --short=8 HEAD`" -tags "tendermint gcc cgo"  -o bin/linocli cmd/linocli/main.go
install:
	go install ./cmd/lino
	go install ./cmd/linocli
get_vendor_deps:
	@rm -rf vendor/
	@dep ensure
test:
	@go test $(PACKAGES)
benchmark:
	@go test -bench=. $(PACKAGES)
.PHONY: all get_tools get_vendor_deps install build test
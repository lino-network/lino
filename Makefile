PACKAGES=$(shell go list ./... | grep -v '/vendor/')
COMMIT_HASH := $(shell git rev-parse --short HEAD)
LD_FLAGS := "-X github.com/tendermint/tendermint/version.GitCommit=`git rev-parse --short=8 HEAD`"
GO_TAGS := "tendermint gcc cgo"
CGO_LDFLAGS := "-lsnappy"

all: get_tools get_vendor_deps install build test

get_tools:
	go get github.com/golang/dep/cmd/dep
	cd scripts && ./install_cleveldb.sh

apply_patch:
	cp ./patches/general/constructors ./vendor/github.com/cosmos/cosmos-sdk/server/constructors.go

build: get_vendor_deps apply_patch
	CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=1 go build -ldflags $(LD_FLAGS) -tags $(GO_TAGS) -o bin/linod   cmd/lino/main.go
	CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=1 go build -ldflags $(LD_FLAGS) -tags $(GO_TAGS) -o bin/linocli cmd/linocli/main.go

install: get_vendor_deps apply_patch
	cd cmd/lino    && CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=1 go install -ldflags $(LD_FLAGS) -tags $(GO_TAGS)
	cd cmd/linocli && CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=1 go install -ldflags $(LD_FLAGS) -tags $(GO_TAGS)

get_vendor_deps:
	@rm -rf vendor/
	@dep ensure

test:
	CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=1 go test -ldflags $(LD_FLAGS) -tags $(GO_TAGS) ./...

benchmark:
	@go test -bench=. $(PACKAGES)
.PHONY: all get_tools get_vendor_deps install build test
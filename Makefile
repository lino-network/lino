GOPATH ?= $(shell $(GO) env GOPATH)
COMMIT := $(shell git log -1 --format='%H')
PACKAGES=$(shell go list ./... | grep -v '/vendor/')
LD_FLAGS := "-X github.com/tendermint/tendermint/version.GitCommit=$(COMMIT) -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb"
GO_TAGS := "tendermint cgo cleveldb"
CGO_LDFLAGS := "-lsnappy"
GO111MODULE = on

all: get_tools install build test

get_tools:
	cd scripts && ./install_cleveldb.sh

# apply_patch:
# 	(cd vendor/github.com/tendermint/tendermint && patch -p1 -t < ../../../../patches/fullnode/tendermint-cached-txindexer.patch); exit 0
# 	(cd vendor/github.com/cosmos/cosmos-sdk     && patch -p1 -t < ../../../../patches/fixes/cosmos-cleveldb-close-batch.patch); exit 0

update_mocks:
	GO111MODULE=$(GO111MODULE) go generate ./...

_raw_build_cmd:
	GO111MODULE=$(GO111MODULE) CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=1 go build -ldflags $(LD_FLAGS) -tags $(GO_TAGS) -o bin/linod   cmd/lino/main.go
	GO111MODULE=$(GO111MODULE) CGO_ENABLED=0 go build -ldflags $(LD_FLAGS) -o bin/linocli cmd/linocli/main.go

_raw_install_cmd:
	cd cmd/lino    && GO111MODULE=$(GO111MODULE) CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=1 go install -ldflags $(LD_FLAGS) -tags $(GO_TAGS)
	cd cmd/linocli && GO111MODULE=$(GO111MODULE) CGO_ENABLED=0 go install -ldflags $(LD_FLAGS)

build:
	make _raw_build_cmd

install:
	make _raw_install_cmd

install_cli:
	cd cmd/linocli && GO111MODULE=$(GO111MODULE) CGO_ENABLED=0 go install -ldflags $(LD_FLAGS)

build_cli:
	GO111MODULE=$(GO111MODULE) CGO_ENABLED=0 go build -ldflags $(LD_FLAGS) -o bin/linocli cmd/linocli/main.go

test:
	GO111MODULE=$(GO111MODULE) CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=1 go test -ldflags $(LD_FLAGS) -tags $(GO_TAGS) ./... -timeout 600s

benchmark:
	@go test -bench=. $(PACKAGES)

# lint
GOLANGCI_LINT_VERSION := v1.17.1
GOLANGCI_LINT_HASHSUM := f5fa647a12f658924d9f7d6b9628d505ab118e8e049e43272de6526053ebe08d

get_golangci_lint:
	cd scripts && bash install-golangci-lint.sh $(GOPATH)/bin $(GOLANGCI_LINT_VERSION) $(GOLANGCI_LINT_HASHSUM)

lint:
	@golangci-lint run
	go mod verify

lint-fix:
	@echo "--> Running linter auto fix"
	@golangci-lint run --fix
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

.PHONY: lint lint-fix


.PHONY: all get_tools install build test
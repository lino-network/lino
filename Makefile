NAME=lino
GOPATH ?= $(shell $(GO) env GOPATH)
COMMIT := $(shell git --no-pager describe --tags --always --dirty)
PACKAGES=$(shell go list ./... | grep -v '/vendor/')
LD_FLAGS := "-X github.com/lino-network/lino/app.Version=$(COMMIT) -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb"
GO_TAGS := "tendermint cgo cleveldb"
CGO_LDFLAGS := "-lsnappy"
GO111MODULE = on

all: get_tools install build test

get_tools:
	cd scripts && ./install_cleveldb.sh

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

docker-build:
	docker build -t $(NAME) .

docker-build-nc:
	docker build --no-cache -t $(NAME) .

docker-run:
	docker run --name=$(NAME) -it $(NAME)

docker-up: docker-build docker-run

docker-clean:
	docker stop $(NAME)
	docker rm $(NAME)

# lint
GOLANGCI_LINT_VERSION := v1.17.1
GOLANGCI_LINT_HASHSUM := f5fa647a12f658924d9f7d6b9628d505ab118e8e049e43272de6526053ebe08d

get_golangci_lint:
	cd scripts && bash install-golangci-lint.sh $(GOPATH)/bin $(GOLANGCI_LINT_VERSION) $(GOLANGCI_LINT_HASHSUM)

lint:
	GO111MODULE=$(GO111MODULE) golangci-lint run
	GO111MODULE=$(GO111MODULE) go mod verify
	GO111MODULE=$(GO111MODULE) go mod tidy

lint-fix:
	@echo "--> Running linter auto fix"
	GO111MODULE=$(GO111MODULE) golangci-lint run --fix
	GO111MODULE=$(GO111MODULE) find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	GO111MODULE=$(GO111MODULE) go mod verify
	GO111MODULE=$(GO111MODULE) go mod tidy

.PHONY: lint lint-fix


.PHONY: all get_tools install build test
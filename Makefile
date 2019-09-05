COMMIT := $(shell git log -1 --format='%H')
PACKAGES=$(shell go list ./... | grep -v '/vendor/')
VENDOR_PATH=github.com/lino-network/lino/vendor
LD_FLAGS := "-X $(VENDOR_PATH)/github.com/tendermint/tendermint/version.GitCommit=$(COMMIT) -X $(VENDOR_PATH)/github.com/cosmos/cosmos-sdk/types.DBBackend=rocksdb"
GO_TAGS := "tendermint gcc cgo rocksdb"
CGO_LDFLAGS := "-lrocksdb -lstdc++ -lm -lzstd -lsnappy"

all: get_tools get_vendor_deps install build test

get_tools:
	go get github.com/golang/dep/cmd/dep
	cd scripts && ./install_rocksdb.sh

apply_patch:
	dep ensure
	(cd vendor/github.com/cosmos/cosmos-sdk     && patch -p1 -t < ../../../../patches/fixes/cosmos-cleveldb-close-batch.patch); exit 0
	(cd vendor/github.com/cosmos/cosmos-sdk     && patch -p1 -t < ../../../../patches/general/cosmos-db-hack.patch); exit 0
	(cd vendor/github.com/cosmos/cosmos-sdk     && patch -p1 -t < ../../../../patches/fixes/cosmos-export-hack.patch); exit 0

_raw_build_cmd:
	CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=1 go build -ldflags $(LD_FLAGS) -tags $(GO_TAGS) -o bin/linod   cmd/lino/main.go
	CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=1 go build -ldflags $(LD_FLAGS) -tags $(GO_TAGS) -o bin/linocli cmd/linocli/main.go

_raw_install_cmd:
	cd cmd/lino    && CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=1 go install -ldflags $(LD_FLAGS) -tags $(GO_TAGS)
	cd cmd/linocli && CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=1 go install -ldflags $(LD_FLAGS) -tags $(GO_TAGS)

build: get_vendor_deps apply_patch
	make _raw_build_cmd

install: get_vendor_deps apply_patch
	make _raw_install_cmd

get_vendor_deps:
	@rm -rf vendor/
	@dep ensure

test:get_vendor_deps apply_patch
	CGO_LDFLAGS=$(CGO_LDFLAGS) CGO_ENABLED=1 go test -ldflags $(LD_FLAGS) -tags $(GO_TAGS) ./... -timeout 600s

benchmark:
	@go test -bench=. $(PACKAGES)

.PHONY: all get_tools get_vendor_deps install build test

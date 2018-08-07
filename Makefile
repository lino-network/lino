PACKAGES=$(shell go list ./... | grep -v '/vendor/')
COMMIT_HASH := $(shell git rev-parse --short HEAD)

all: get_tools get_vendor_deps install build test

get_tools:
	go get github.com/golang/dep/cmd/dep
build:
	go build -o bin/linocli cmd/linocli/main.go && go build -o bin/linod cmd/lino/main.go
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
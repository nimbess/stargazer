BINDIR=bin
DOCKER=docker
GO=go
BINARY=stargazer

TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)
export TAG

## Build the binary and image
default: image

## Build the binary and image
all: image

## Make the bin directory
bindir:
	mkdir -p ${BINDIR}

## Run tests
test:
	${GO} test -v ./...

## Show code test coverage
coverage:
	${GO} test -v -cover ./...

## Build the test code coverage profile and html output
coverage-html: bindir
	${GO} test -v -cover -coverprofile=${BINDIR}/${BINARY}.profile.out ./...
	${GO} tool cover -html=${BINDIR}/${BINARY}.profile.out -o ${BINDIR}/${BINARY}.profile.html

## Build the application binary
build: bindir
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GO} build -o ${BINDIR}/${BINARY} -ldflags "-X main.VERSION=$(TAG)" ./cmd/${BINARY}/

## Build a container image
image:
	${DOCKER} build -t nimbess/${BINARY}:${TAG} .

## Get the protobuf generator plugin
get-generators:
	go get -u github.com/golang/protobuf/protoc-gen-go

## Compile the protobuf files
proto:
	protoc --go_out=plugins=grpc:. ./pkg/model/node/*.proto

## Clean the build dirs
clean:
	go clean
	rm -rf ${BINDIR}/

.PHONY: help
## Display this help text.
help: # Some kind of magic from https://gist.github.com/rcmachado/af3db315e31383502660
	$(info Available targets)
	@awk '/^[a-zA-Z\-\_0-9\/]+:/ {                                      \
		nb = sub( /^## /, "", helpMsg );                                \
		if(nb == 0) {                                                   \
			helpMsg = $$0;                                              \
			nb = sub( /^[^:]*:.* ## /, "", helpMsg );                   \
		}                                                               \
		if (nb)                                                         \
			printf "\033[1;31m%-" width "s\033[0m %s\n", $$1, helpMsg;  \
	}                                                                   \
	{ helpMsg = $$0 }'                                                  \
	width=20                                                            \
	$(MAKEFILE_LIST)

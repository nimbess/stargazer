BINDIR=bin
DOCKER=docker
GO=go
BINARY=stargazer

TAG?=$(shell git rev-list HEAD --max-count=1 --abbrev-commit)
export TAG

default: image

all: image

bindir:
	mkdir -p ${BINDIR}

test:
	${GO} test -v ./...

coverage:
	${GO} test -v -cover ./...

coverage-html: bindir
	${GO} test -v -cover -coverprofile=${BINDIR}/${BINARY}.profile.out ./...
	${GO} tool cover -html=${BINDIR}/${BINARY}.profile.out -o ${BINDIR}/${BINARY}.profile.html

build: bindir
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GO} build -v -o ${BINDIR}/${BINARY} -ldflags "-X main.VERSION=$(TAG)" ./cmd/${BINARY}/

image: build
	${DOCKER} build -t nimbess/${BINARY}:${TAG} .

clean:
	go clean
	rm -rf ${BINDIR}/

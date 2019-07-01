FROM golang:1.11 AS builder-base

# Get dependencies
WORKDIR /go/src/github.com/nimbess/stargazer
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN go mod verify

# Build the binary
FROM builder-base AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -v -o /go/bin/stargazer -ldflags "-X main.VERSION=0.0.1" -a -installsuffix cgo ./cmd/stargazer/

# Copy the binary
FROM alpine
COPY --from=builder /go/bin/stargazer /go/bin/stargazer
COPY stargazer.yaml /etc/stargazer.yaml

# Run the binary
CMD ["/go/bin/stargazer", "-config-path", "/etc"]

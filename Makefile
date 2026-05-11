# M3TAL Core Makefile

BINARY_NAME=m3tal
API_NAME=m3tal-api

.PHONY: all build build-linux clean fmt tidy vendor

all: build

build:
	go build -o $(BINARY_NAME) ./cmd/m3tal
	go build -o $(API_NAME) ./cmd/api

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME)-linux ./cmd/m3tal
	GOOS=linux GOARCH=amd64 go build -o $(API_NAME)-linux ./cmd/api

clean:
	rm -f $(BINARY_NAME) $(API_NAME) $(BINARY_NAME)-linux $(API_NAME)-linux
	rm -f *.exe

fmt:
	gofmt -s -w .

tidy:
	go mod tidy

vendor:
	go mod vendor

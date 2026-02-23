.PHONY: build clean install test run

BINARY_NAME=floodguard
BUILD_DIR=build
INSTALL_PATH=/usr/bin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/floodguard

build-linux:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/floodguard

clean:
	rm -rf $(BUILD_DIR)
	go clean

install: build
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)/
	sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)

test:
	go test -v ./...

run:
	go run cmd/floodguard/main.go

fmt:
	go fmt ./...

lint:
	golangci-lint run

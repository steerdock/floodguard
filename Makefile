.PHONY: build clean install test run

BINARY_NAME=floodguard
BUILD_DIR=build
INSTALL_PATH=/usr/local/bin

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/floodguard/main.go

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 cmd/floodguard/main.go

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

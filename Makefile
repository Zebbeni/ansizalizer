APP_NAME := ansizalizer
DIST_DIR := dist
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

.PHONY: all clean build test dist

all: build

build:
	go build -o $(APP_NAME) .

test:
	go test ./...

clean:
	rm -rf $(DIST_DIR)
	rm -f $(APP_NAME) $(APP_NAME).exe

dist: clean
	@mkdir -p $(DIST_DIR)
	@echo "Building $(VERSION) for all platforms..."

	GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/$(APP_NAME)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 go build -o $(DIST_DIR)/$(APP_NAME)-windows-arm64.exe .
	GOOS=darwin  GOARCH=amd64 go build -o $(DIST_DIR)/$(APP_NAME)-macos-amd64 .
	GOOS=darwin  GOARCH=arm64 go build -o $(DIST_DIR)/$(APP_NAME)-macos-arm64 .
	GOOS=linux   GOARCH=amd64 go build -o $(DIST_DIR)/$(APP_NAME)-linux-amd64 .
	GOOS=linux   GOARCH=arm64 go build -o $(DIST_DIR)/$(APP_NAME)-linux-arm64 .

	@echo ""
	@echo "Built:"
	@ls -lh $(DIST_DIR)/

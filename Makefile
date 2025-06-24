APP_NAME = gokeeper
VERSION = 1.0.0
DATE    = $(shell date)

BUILD_FLAGS = -ldflags "-X main.Version=$(VERSION) -X 'main.BuildDate=$(DATE)'"

# Папка для релизов
BIN_DIR = build

all: clean release

release: darwin_amd64 darwin_arm64 linux_amd64 windows_amd64

darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BIN_DIR)/$(APP_NAME)-darwin-amd64 ./cmd/client

darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BIN_DIR)/$(APP_NAME)-darwin-arm64 ./cmd/client

linux_amd64:
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BIN_DIR)/$(APP_NAME)-linux-amd64 ./cmd/client

windows_amd64:
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe ./cmd/client

clean:
	rm -rf $(BIN_DIR)
	mkdir -p $(BIN_DIR)


BIN_DIR = bin
APP_NAME = efp

BUILD_CMD = go build -ldflags '-w -s -buildid='
# The -w and -s flags reduce binary sizes by excluding unnecessary symbols and debug info
# The -buildid= flag makes builds reproducible

all: windows-amd64

windows-386:
	GOARCH=386 GOOS=windows $(BUILD_CMD) -o $(BIN_DIR)/$(APP_NAME)-$@.exe

windows-amd64:
	GOARCH=amd64 GOOS=windows $(BUILD_CMD) -o $(BIN_DIR)/$(APP_NAME)-$@.exe


clean:
	rm $(BIN_DIR)/*

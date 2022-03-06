
BIN_DIR = bin
APP_NAME = efp

BUILD_CMD = go build -ldflags '-w -s -buildid='
# The -w and -s flags reduce binary sizes by excluding unnecessary symbols and debug info
# The -buildid= flag makes builds reproducible

all: windows-amd64 linux-armv8

windows-amd64:
	GOARCH=amd64 GOOS=windows $(BUILD_CMD) -o $(BIN_DIR)/$(APP_NAME)-$@.exe

linux-armv8:
	GOARCH=arm64 GOOS=linux $(BUILD_CMD) -o $(BIN_DIR)/$(APP_NAME)-$@

clean:
	rm $(BIN_DIR)/*

# ONLY FOR TESTING!

HOST = raspberrypi
PORT = 8488
AUTH = "root@$(HOST)"
CIPHER = aes-128-gcm
KEY = 1234567890abcdef1234567890abcdef

RUN_AS_SERVER = efp -s :${PORT} -cipher ${CIPHER} -key ${KEY} -v

deploy-server:
	ssh $(AUTH) "mkdir -p ~/bin"
	-(ssh $(AUTH) "pkill efp; rm -f ~/bin/efp")
	scp bin/efp-linux-armv8 $(AUTH):~/bin/efp
	ssh $(AUTH) "export PATH=\$$PATH:~/bin; chmod +x ~/bin/efp; setsid $(RUN_AS_SERVER) >/dev/null 2>&1 &"
	ssh $(AUTH) "ps ax | grep efp"

RUN_AS_CLIENT = efp -c $(HOST):${PORT} -cipher ${CIPHER} -key ${KEY} -socks :1080 -v
COPY_SRC_TO_DST = bin/efp-windows-amd64.exe c:/developer/bin/efp.exe

deploy-client:
	cp $(COPY_SRC_TO_DST)
	$(RUN_AS_CLIENT)

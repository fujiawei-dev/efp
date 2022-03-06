.PHONY: ;
#.SILENT: ;			   # no need for @
.ONESHELL: ;			 # recipes execute in same shell
.NOTPARALLEL: ;		  # wait for target to finish
.EXPORT_ALL_VARIABLES: ; # send all vars to shell
.IGNORE: clean;			# ignore all errors, keep going

APP_NAME = efp
APP_VERSION = 0.0.1

BIN_DIR = bin

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

RUN_AS_SERVER = ${APP_NAME} -s :${PORT} -cipher ${CIPHER} -key ${KEY} -v

deploy-server:
	ssh $(AUTH) "mkdir -p ~/bin"
	-(ssh $(AUTH) "pkill $(APP_NAME); rm -f ~/bin/$(APP_NAME)")
	scp bin/$(APP_NAME)-linux-armv8 $(AUTH):~/bin/$(APP_NAME)
	ssh $(AUTH) "export PATH=\$$PATH:~/bin; chmod +x ~/bin/$(APP_NAME); setsid $(RUN_AS_SERVER) >/dev/null 2>&1 &"
	ssh $(AUTH) "ps ax | grep $(APP_NAME)"

RUN_AS_CLIENT = $(APP_NAME) -c $(HOST):${PORT} -cipher ${CIPHER} -key ${KEY} -socks :1080 -v
COPY_SRC_TO_DST = bin/$(APP_NAME)-windows-amd64.exe c:/developer/bin/$(APP_NAME).exe

deploy-client:
	cp $(COPY_SRC_TO_DST)
	$(RUN_AS_CLIENT)

tag:
	git tag v$(APP_VERSION)
	git push origin v$(APP_VERSION)

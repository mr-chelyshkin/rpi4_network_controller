THIS_FILE := $(lastword $(MAKEFILE_LIST))
SHELL=/bin/sh

all:
	db_generate

##################    VARIABLES    ##################
ASM_FLAGS:=-trimpath=OPATH
GC_FLAGS:=-trimpath=OPATH
LD_FLAGS:=-w -s

MAIN_FILE:=./cmd
BIN_FILE:=network

CMD_BUILD=CGO_ENABLED=1 CFLAGS="-march=armv8-a" CC=aarch64-linux-gnu-gcc GOOS=$(1) GOARCH=$(2) go build\
	-asmflags="$(ASM_FLAGS)"\
	-ldflags="$(LD_FLAGS)"\
	-gcflags="$(GC_FLAGS)"\
	-o "$(BIN_FILE)"\
	$(MAIN_FILE)

##################    TARGETS    ##################
rpi4:
	$(call CMD_BUILD,linux,arm64)

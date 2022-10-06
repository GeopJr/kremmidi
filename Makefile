.PHONY: all fmt build onion-rings

KREMMIDI_SRC ?= ./src/
PNPM_VERSION := $(shell command -v pnpm 2>/dev/null)
NPM_VERSION := $(shell command -v npm 2>/dev/null)

all: fmt build

fmt:
	go fmt $(KREMMIDI_SRC)

build:
	go build -o kremmidi  $(KREMMIDI_SRC)*.go

onion-rings:
ifdef PNPM_VERSION
	cd onion-rings/ && pnpm i
	cd onion-rings/ && pnpm build
else ifdef NPM_VERSION
	cd onion-rings/ && npm i
	cd onion-rings/ && npm run build
else
	$(error "Please install pnpm or npm.")
endif
	

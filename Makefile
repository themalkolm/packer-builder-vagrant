OS             ?= $(shell go env GOOS)
ARCH           ?= $(shell go env GOARCH)

VERSION        ?= 0.1.1
GOPKG          ?= github.com/themalkolm/packer-builder-vagrant
DESTDIR        ?= dist
PACKER_VERSION ?= 0.12.2

BINARY         ?= packer-$(PACKER_VERSION)_packer-builder-vagrant_$(OS)_$(ARCH)
ifeq ($(OS),windows)
	BINARY := $(BINARY).exe
endif

all: build

.PHONY: fmt
fmt:
	find . -name \*.go -not -path "./vendor/*" | xargs gofmt -w

build/$(BINARY):
	mkdir -p build
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 \
		go build -v \
	    -ldflags="-s -w -X main.Version=$(VERSION)" \
	    -o $@ $(GOPKG)

.PHONY: build
build: build/$(BINARY)

.PHONY: dist
dist:
	OS=linux   ARCH=amd64 make build
	OS=darwin  ARCH=amd64 make build
	OS=windows ARCH=amd64 make build

.PHONY: install
install: build
	install -d                      $(DESTDIR)/bin
	install -m 0755 ./build/brew-me $(DESTDIR)/bin

.PHONY: clean
clean:
	rm -rf build dist

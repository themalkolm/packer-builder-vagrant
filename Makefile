OS             ?= $(shell go env GOOS)
ARCH           ?= $(shell go env GOARCH)

VERSION        ?= 2018.10.15
GOPKG          ?= github.com/themalkolm/packer-builder-vagrant
PACKER_VERSION ?= 1.2.5

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
build: vendor/ build/$(BINARY)

.PHONY: dist
dist:
	OS=linux   ARCH=amd64 make build
	OS=darwin  ARCH=amd64 make build
	OS=windows ARCH=amd64 make build

.PHONY: clean
clean:
	rm -rf build dist vendor/

vendor/:
	mkdir -p $(CURDIR)/vendor/github.com/hashicorp
	ln -s    $(CURDIR)/vendor/github.com/hashicorp $(CURDIR)/vendor/github.com/mitchellh
	git clone -b v$(PACKER_VERSION) --single-branch --depth 1 https://github.com/hashicorp/packer.git $(CURDIR)/vendor/github.com/hashicorp/packer
	rsync -azK $(CURDIR)/vendor/github.com/hashicorp/packer/vendor/ $(CURDIR)/vendor/
	rm     -rf $(CURDIR)/vendor/github.com/hashicorp/packer/vendor/
	git clone           --single-branch --depth 1 https://github.com/koding/vagrantutil   $(CURDIR)/vendor/github.com/koding/vagrantutil && \
	    cd $(CURDIR)/vendor/github.com/koding/vagrantutil && git checkout 70827343f1169931bbe84e8051fcbcf06da90eb7
	git clone           --single-branch           https://github.com/koding/logging.git $(CURDIR)/vendor/github.com/koding/logging && \
		cd $(CURDIR)/vendor/github.com/koding/logging && git checkout 8b5a689ed69b1c7cd1e3595276fc2a352d7818e0
	find $(CURDIR)/vendor -name .git | xargs rm -rf

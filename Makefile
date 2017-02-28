VERSION        ?= 0.1.1
GOPKG           = github.com/themalkolm/packer-builder-vagrant
DESTDIR        ?= dist
PACKER_VERSION  = 0.12.2

all: build

.PHONY: fmt
fmt:
	find . -name \*.go -not -path "./vendor/*" | xargs gofmt -w

.PHONY: build
build:
	mkdir -p build
	go build -v \
	    -ldflags="-s -w -X main.Version=$(VERSION)" \
	    -o ./build/packer-$(PACKER_VERSION)_packer-builder-vagrant_$(shell go env GOOS)_$(shell go env GOARCH) $(GOPKG)

.PHONY: dist
dist:
	GOOS=linux   GOARCH=amd64 make build
	GOOS=darwin  GOARCH=amd64 make build
	GOOS=windows GOARCH=amd64 make build

.PHONY: install
install: build
	install -d                      $(DESTDIR)/bin
	install -m 0755 ./build/brew-me $(DESTDIR)/bin

.PHONY: clean
clean:
	rm -rf build dist

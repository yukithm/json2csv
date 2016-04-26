PREFIX = /usr/local
BINDIR = $(PREFIX)/bin
DISTDIR = releases
GO_FLAGS = -ldflags="-s -w"
DEVTOOL_DIR = $(CURDIR)/devtool
GOX = $(DEVTOOL_DIR)/bin/gox
OSARCH = linux/amd64 linux/arm darwin/amd64 windows/386 windows/amd64
DIST_FORMAT = $(DISTDIR)/{{.Dir}}-{{.OS}}-{{.Arch}}

.PHONY: all test build install clean dist dist-clean

all: build

build: json2csv

json2csv: *.go jsonpointer/*.go cmd/json2csv/*.go
	go build $(GO_FLAGS) ./cmd/json2csv

test:
	go test $(shell go list ./... | grep -v "/vendor/")

install: all
	install -d $(BINDIR)
	install json2csv $(BINDIR)

clean:
	rm -f json2csv

dist: $(DEVTOOL_DIR)/bin/gox
	$(GOX) -osarch="$(OSARCH)" $(GO_FLAGS) -output="$(DIST_FORMAT)" ./cmd/json2csv

$(DEVTOOL_DIR)/bin/gox:
	mkdir -p $(DEVTOOL_DIR)/{bin,pkg,src}
	GOPATH=$(DEVTOOL_DIR) go get github.com/mitchellh/gox

dist-clean:
	rm -rf json2csv $(DISTDIR) $(DEVTOOL_DIR)

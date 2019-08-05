PREFIX = /usr/local
BINDIR = $(PREFIX)/bin
DISTDIR = releases
GO_FLAGS = -ldflags="-s -w"
DEVTOOLS_DIR = $(CURDIR)/devtools
DEVTOOLS_BIN = $(DEVTOOLS_DIR)/bin
GOX = $(DEVTOOLS_BIN)/gox
OSARCH = linux/amd64 linux/arm darwin/amd64 windows/386 windows/amd64
DIST_FORMAT = $(DISTDIR)/{{.Dir}}-{{.OS}}-{{.Arch}}

.PHONY: all test build install clean dist dist-clean $(GOX) devtools

all: build

build: json2csv

json2csv: *.go jsonpointer/*.go cmd/json2csv/*.go
	go build $(GO_FLAGS) ./cmd/json2csv

test:
	go test

install: all
	install -d $(BINDIR)
	install json2csv $(BINDIR)

clean:
	rm -f json2csv

dist: $(GOX)
	$(GOX) -osarch="$(OSARCH)" $(GO_FLAGS) -output="$(DIST_FORMAT)" ./cmd/json2csv

dist-clean:
	rm -rf json2csv json2csv.exe $(DISTDIR) $(DEVTOOLS_BIN)

$(GOX): devtools

devtools:
	go generate $(DEVTOOLS_DIR)/devtools.go

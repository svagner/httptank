PROGRAM=httptank
BINNAME=$(PROGRAM)
SRCPATH=$(shell pwd)
BUILDDIR=$(SRCPATH)/build
MAINTEINER=Stan Putrya <root.vagner@gmail.com>
VERSION?=$(shell git describe --abbrev=0 --tags)

.PHONY: all clean prepare server client

all: build

prepare:
	@mkdir -p $(BUILDDIR)

build: prepare
	go build -ldflags "-X main.version=$(VERSION)" -o $(BUILDDIR)/$(BINNAME) cmd/$(PROGRAM)/main.go 

clean:
	@rm -rf $(BUILDDIR)

PROGRAM=httptank
BINNAME=$(PROGRAM)
SRCPATH=$(shell pwd)
BUILDDIR=$(SRCPATH)/build
MAINTEINER=Stan Putrya <root.vagner@gmail.com>
VERSION?=$(shell git describe --abbrev=0 --tags)

.PHONY: all clean prepare server client

all: deps build

tools:
	go get -u github.com/kardianos/govendor

prepare:
	@mkdir -p $(BUILDDIR)

deps: tools
	cd $(SRCPATH)/src/$(PROGRAM) ;	GOPATH=$(SRCPATH) govendor init 
	cd $(SRCPATH)/src/$(PROGRAM) ;	GOPATH=$(SRCPATH) govendor fetch +missing 

build: prepare
	cd $(SRCPATH)/src/$(PROGRAM) ; GOPATH=$(SRCPATH) go build -ldflags "-X main.version=$(VERSION)" -o $(BUILDDIR)/$(BINNAME)

clean:
	@rm -rf $(BUILDDIR)

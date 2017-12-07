# Helper makefile for ryankurte/utils

VER=$(shell git describe --dirty)
ARGS= -ldflags "-X main.version=$(VER)"

build:
	go build $(ARGS) ./cmd/...

build-all:
	gox -output=build/{{.OS}}-{{.Arch}}/{{.Dir}} $(ARGS) ./cmd/...

deps:
	dep ensure

package:
	cd build; for i in *; do if [ -d "$$i" ]; then tar -czf "$$i-$(VER).tgz"  "$$i"; fi; done

install:
	go get -u github.com/golang/dep/cmd/dep
	go get github.com/mitchellh/gox
	dep ensure

clean:
	rm -rf build/*

.PHONY: build build-all package install


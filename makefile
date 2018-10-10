# Helper makefile for ryankurte/utils

VER=$(shell git describe --dirty)
ARGS= -ldflags "-X main.version=$(VER)"

build:
	go build $(ARGS) ./cmd/...

build-all:
	gox -output=build/{{.OS}}-{{.Arch}}/{{.Dir}} $(ARGS) ./cmd/...

package:
	cd build; for i in *; do if [ -d "$$i" ]; then tar -czf "$$i-$(VER).tgz"  "$$i"; fi; done

install:
	go get github.com/mitchellh/gox
	go get ./cmd/...

clean:
	rm -rf build/*

.PHONY: build build-all package install


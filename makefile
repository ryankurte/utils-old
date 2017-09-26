# Helper makefile for ryankurte/utils

VER=$(shell git describe --dirty)
ARGS= -ldflags "-X main.version=$(VER)"

build:
	go build $(ARGS) ./cmd/...

build-all:
	gox -output=build/{{.OS}}-{{.Arch}}/{{.Dir}} $(ARGS) ./cmd/...

package:
	mkdir -p packages
	for i in build/*; do if [ -d "$$i" ]; then tar -czf "$$i-$(VER).tgz" "$$i"; fi; done
	mv build/*.tgz packages/

install:
	go get github.com/mitchellh/gox
	go get -u ./cmd/...

clean:
	rm -rf build/*

.PHONY: build build-all package install

SHELL = /bin/bash

version := 0.1.0
sha := $(shell git rev-parse --short HEAD)

.PHONY: build
build:
	go build \
		-ldflags "-X kloudy/gitter/cmd.Version=$(version) -X kloudy/gitter/cmd.SHA=$(sha)" \
		-o bin/gitter

.PHONY: install
install: build
	cp bin/gitter "$$(go env GOPATH)/bin/gitter"

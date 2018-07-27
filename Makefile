PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
GOFMT ?= gofmt "-s"
BUILD ?= go build -o ./register cmd/register/main.go
PACK ?= gzip ./register

fmt:
	$(GOFMT) -w $(GOFILES)

vet:
	go vet $(PACKAGES)

.PHONY: build
build:
	$(BUILD)

.PHONY: pack
pack:
	$(PACK)

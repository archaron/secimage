CONTAINER = secimage
REPO = $(shell go list -m)
VERSION ?= $(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)
IMAGE_NAME = docker.pkg.github.com/archaron/secimage/$(CONTAINER):$(VERSION)

.PHONY: image deps

deps:
	@echo "=> Fetch dependencies:"
	@go mod tidy -v
	@go mod vendor

image: deps
	@docker build \
		--build-arg REPO=$(REPO) \
		--build-arg VERSION=$(VERSION) \
		-t $(IMAGE_NAME) .

ARCHS := amd64 386 arm arm64
OSES := linux windows
OUTPUT_DIR := dist
VERSION ?= $(shell git describe --tags --abbrev=0)

default: testacc


# Run acceptance tests
.PHONY: testacc build
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

build:
	echo "Building version $(VERSION)"
	for arch in $(ARCHS); do \
		for system in $(OSES); do \
			@echo "Building $${system}_$${arch}"; \
			GOARCH=$${arch} GOOS=$${system} go build -o $(OUTPUT_DIR)/$${system}_$${arch}/terraform-provider-planningcenter -ldflags="-X 'main.version=$(VERSION)'"; \
		done; \
	done

ARCHS := amd64 386 arm arm64
OSES := linux windows darwin
VERSION ?= $(shell git describe --tags --abbrev=0)
DIR_PATH := .terraform.d/plugins/github.com/HubbardHarvey3/planningcenter/$(VERSION)
default: testacc


# Run acceptance tests
.PHONY: testacc build
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

build:
	echo "Building version $(VERSION)"
	for arch in $(ARCHS); do \
		for system in $(OSES); do \
			echo "Building $${system}_$${arch}"; \
			GOARCH=$${arch} GOOS=$${system} go build -o $(DIR_PATH)/$${system}_$${arch}/terraform-provider-planningcenter -ldflags="-X 'main.version=$(VERSION)'"; \
			zip -r planningcenter_$(VERSION)_$${system}_$${arch} $(DIR_PATH)/$${system}_$${arch}; \
		done; \
	done

PROJECT ?= combi
PROJECT_VERSION ?= $(shell cat version)
PROJECT_GOVER ?= $(shell grep '^go ' go.mod | cut -d ' ' -f 2)
PROJECT_COMMIT ?= $(shell git rev-parse --short HEAD)

BINARY ?= $(PROJECT)

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

##@ Build

VERSION_PACKAGE_PATH ?= $(PROJECT)/internal/cmd/cmdver
BUILD_LDFLAGS_VERSION ?= -X $(VERSION_PACKAGE_PATH).version=$(PROJECT_VERSION)
BUILD_LDFLAGS_GOVER   ?= -X $(VERSION_PACKAGE_PATH).golang=$(PROJECT_GOVER)
BUILD_LDFLAGS_COMMIT  ?= -X $(VERSION_PACKAGE_PATH).commit=$(PROJECT_COMMIT)
BUILD_LDFLAGS_VALUE ?= "$(BUILD_LDFLAGS_VERSION) $(BUILD_LDFLAGS_GOVER) $(BUILD_LDFLAGS_COMMIT)"
BUILD_LDFLAGS ?= -ldflags $(BUILD_LDFLAGS_VALUE)

BINPATH ?= ./bin/$(BINARY)

$(BINPATH): fmt vet
	go build $(BUILD_LDFLAGS) -o bin/$(BINARY) cmd/$(PROJECT)/main.go

.PHONY: build
build: $(BINPATH) ## Build manager binary.

RUN_ARGS ?= version

.PHONY: run
run: build ## Run a command from your host (define RUN_ARGS to custom run).
	./bin/$(BINARY) $(RUN_ARGS)

.PHONY: vrun
vrun: fmt vet ## Run a command from your host (define RUN_ARGS to custom run).
	go run $(BUILD_LDFLAGS) cmd/$(PROJECT)/main.go $(RUN_ARGS)

##@ Container

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

# Image URL to use all building/pushing image targets
IMG_REGISTRY ?= ghcr.io/freepik-company
IMG_NAME ?= $(PROJECT)
IMG_TAG ?= latest
IMG ?= $(IMG_REGISTRY)/$(IMG_NAME):$(IMG_TAG)

# PLATFORMS defines the target platforms for the manager image be built to provide support to multiple
# architectures. (i.e. make docker-buildx IMG=myregistry/mypoperator:0.0.1). To use this option you need to:
# - be able to use docker buildx. More info: https://docs.docker.com/build/buildx/
# - have enabled BuildKit. More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image to your registry (i.e. if you do not set a valid value via IMG=<myregistry/image:<tag>> then the export will fail)
# To adequately provide solutions that are compatible with multiple platforms, you should consider using this option.
PLATFORMS ?= linux/arm64,linux/amd64,linux/s390x,linux/ppc64le
.PHONY: docker-buildx
docker-buildx: ## Build and push docker image for the manager for cross-platform support
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	# sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	sed -e 's/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/g' Dockerfile > Dockerfile.cross
	- $(CONTAINER_TOOL) buildx create --name project-builder
	$(CONTAINER_TOOL) buildx use project-builder
	- $(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --build-arg LDFLAGS_VALUE=$(BUILD_LDFLAGS_VALUE) --tag $(IMG) --file Dockerfile.cross .
	- $(CONTAINER_TOOL) buildx rm project-builder
	rm Dockerfile.cross

.PHONY: container-build
container-build: ## Build the container image
	$(CONTAINER_TOOL) build --build-arg LDFLAGS_VALUE=$(BUILD_LDFLAGS_VALUE) --no-cache --tag $(IMG) --file Dockerfile .

.PHONY: container-push
container-push: ## Push the container image
	$(CONTAINER_TOOL) push $(IMG)

CONTAINER_ARGS ?= version

.PHONY: container-run
container-run: ## Run the container image
	$(CONTAINER_TOOL) run $(IMG) $(CONTAINER_ARGS)

.PHONY: kind-load
kind-load: ## Loads the container image in kind cluster
	kind load docker-image $(IMG)

# Common
###############################################################################

#http_proxy := "http://www-proxy.us.oracle.com:80"
#https_proxy := "http://www-proxy.us.oracle.com:80"

GO ?= go
GOOS ?= linux
LD_FLAGS ?= -extldflags "-static"

BUILD_DIR := dist
VERSION := $(shell date +%y%m%d%H%M)

BIN_DIR := ${BUILD_DIR}/bin
MANIFEST_DIR := ${BUILD_DIR}/manifest

GO_SRC := $(shell find . -name "*.go")

DOCKER_REGISTRY ?= gcr.io
DOCKER_USER ?= fmwplt-gcp

# weblogic-operator
###############################################################################

OPERATOR_NAME := weblogic-operator

OPERATOR_BIN_NAME := ${OPERATOR_NAME}

OPERATOR_DOCKER_IMAGE_NAME ?= weblogic-operator
export OPERATOR_DOCKER_IMAGE_TAG ?= ${VERSION}

.PHONY: all
all: build

.PHONY: build
build: ${BIN_DIR}/${OPERATOR_BIN_NAME}

${BIN_DIR}/${OPERATOR_BIN_NAME}: ${GO_SRC}
	@mkdir -p ${BIN_DIR}
	GOOS=$(GOOS) CGO_ENABLED=0 $(GO) build -v -ldflags '${LD_FLAGS}' -o $@ ./cmd/weblogic-operator
	@cp -r scripts ${BUILD_DIR}/

.PHONY: docker-stage
docker-stage: ${BIN_DIR}/${OPERATOR_BIN_NAME}
	@sed "s/{{VERSION}}/$(OPERATOR_DOCKER_IMAGE_TAG)/g" manifests/weblogic-operator-template.yaml > manifests/weblogic-operator.yaml

	@mkdir -p docker-stage
	@cp -r ${BUILD_DIR}/scripts/weblogic/ docker-stage/scripts/
	@cp -r ${BIN_DIR}/weblogic-operator docker-stage/

.PHONY: image
image: ${BIN_DIR}/${OPERATOR_BIN_NAME}
	sed "s/{{VERSION}}/$(OPERATOR_DOCKER_IMAGE_TAG)/g" manifests/weblogic-operator-template.yaml > manifests/weblogic-operator.yaml
	@docker build \
		--build-arg=http_proxy \
		--build-arg=https_proxy \
		-t ${DOCKER_REGISTRY}/${DOCKER_USER}/${OPERATOR_DOCKER_IMAGE_NAME}:${OPERATOR_DOCKER_IMAGE_TAG} \
		-f Dockerfile \
		.
	#@docker tag ${DOCKER_REGISTRY}/${DOCKER_USER}/${OPERATOR_DOCKER_IMAGE_NAME}:${OPERATOR_DOCKER_IMAGE_TAG}

.PHONY: push
push: image
	#@docker login -u '$(DOCKER_REGISTRY_USERNAME)' -p '$(DOCKER_REGISTRY_PASSWORD)' $(DOCKER_REGISTRY)
	@docker push ${DOCKER_REGISTRY}/${DOCKER_USER}/${OPERATOR_DOCKER_IMAGE_NAME}:${OPERATOR_DOCKER_IMAGE_TAG}

.PHONY: fmt
fmt:
	@gofmt -s -e -d $(shell find . -name "*.go" | grep -v /vendor/)

.PHONY: vet
vet: ${GO_SRC}
	@go vet $(shell go list ./... | grep -v /vendor/)

.PHONY: vendor
vendor:
	@$(GO) get -u github.com/golang/dep/cmd/dep
	#@dep init -v
	@$(GOPATH)/bin/dep ensure -v

.PHONY: clean
clean:
	rm -rf ${BUILD_DIR}



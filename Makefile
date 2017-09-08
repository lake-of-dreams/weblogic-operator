# Common
###############################################################################

GO ?= go
GOOS ?= linux
LD_FLAGS ?= -extldflags "-static"

ifdef GITLAB_CI
BUILD_DIR := ${CI_PROJECT_DIR}/dist

ifdef CI_COMMIT_TAG
VERSIONi := ${CI_COMMIT_TAG}
else
VERSION = ${CI_COMMIT_SHA}
endif

else
# localbuild user timestamp it
BUILD_DIR := dist
VERSION := ${USER}-$(shell date +%Y%m%d%H%M%S)
endif

BIN_DIR := ${BUILD_DIR}/bin
MANIFEST_DIR := ${BUILD_DIR}/manifest

GO_SRC := $(shell find . -name "*.go")

DOCKER_REGISTRY ?= registry.oracledx.com
DOCKER_USER ?= skeppare

# mysql-operator
###############################################################################

OPERATOR_NAME := mysql-operator

OPERATOR_BIN_NAME := ${OPERATOR_NAME}

OPERATOR_DOCKER_IMAGE_NAME ?= mysql-operator
OPERATOR_DOCKER_IMAGE_TAG ?= ${VERSION}

.PHONY: all
all: build

.PHONY: build
build: ${BIN_DIR}/${OPERATOR_BIN_NAME}

${BIN_DIR}/${OPERATOR_BIN_NAME}: ${GO_SRC}
	@mkdir -p ${BIN_DIR}
	GOOS=$(GOOS) CGO_ENABLED=0 $(GO) build -v -ldflags '${LD_FLAGS}' -o $@ ./cmd/mysql-operator

.PHONY: image
image: ${BIN_DIR}/${OPERATOR_BIN_NAME}
	sed "s/{{VERSION}}/$(OPERATOR_DOCKER_IMAGE_TAG)/g" manifests/mysql-operator.yaml > \
		$(BUILD_DIR)/mysql-operator.yaml
	docker build \
		--build-arg=http_proxy \
		--build-arg=https_proxy \
		-t ${OPERATOR_DOCKER_IMAGE_NAME}:${OPERATOR_DOCKER_IMAGE_TAG} \
		-f Dockerfile \
		.
	docker tag ${OPERATOR_DOCKER_IMAGE_NAME}:${OPERATOR_DOCKER_IMAGE_TAG} ${DOCKER_REGISTRY}/${DOCKER_USER}/${OPERATOR_DOCKER_IMAGE_NAME}:${OPERATOR_DOCKER_IMAGE_TAG}

.PHONY: push
push: image
	@docker login -u '$(DOCKER_REGISTRY_USERNAME)' -p '$(DOCKER_REGISTRY_PASSWORD)' $(DOCKER_REGISTRY)
	@docker push ${DOCKER_REGISTRY}/${DOCKER_USER}/${OPERATOR_DOCKER_IMAGE_NAME}:${OPERATOR_DOCKER_IMAGE_TAG}

.PHONY: fmt
fmt:
	@gofmt -s -e -d $(shell find . -name "*.go" | grep -v /vendor/)

.PHONY: test
test: ${GO_SRC}
	go test -v ./pkg/... ./cmd/...

.PHONY: vet
vet: ${GO_SRC}
	@go vet $(shell go list ./... | grep -v /vendor/)

.PHONY: vendor
vendor:
	glide install -v

.PHONY: clean
clean:
	rm -rf ${BUILD_DIR}

.PHONY: deploy
deploy: push
	kubectl apply -f ${BUILD_DIR}/${OPERATOR_NAME}.yaml

.PHONY: start
start:
	kubectl apply -f ${BUILD_DIR}/${OPERATOR_NAME}.yaml

.PHONY: stop
stop:
	kubectl delete -f ${BUILD_DIR}/${OPERATOR_NAME}.yaml

.PHONY: run-dev
run-dev:
	@$(GO) run cmd/weblogic-operator/main.go --kubeconfig=${KUBECONFIG} --v=4

.PHONY: e2e
e2e:
	@go test -v ./test/e2e/ --kubeconfig=${KUBECONFIG}

# mysql-agent
###############################################################################

AGENT_NAME := mysql-agent

AGENT_BIN_NAME := ${AGENT_NAME}

AGENT_DOCKER_IMAGE_NAME ?= mysql-agent
AGENT_DOCKER_IMAGE_TAG ?= ${VERSION}

.PHONY: agent-build
agent-build: ${BIN_DIR}/${AGENT_BIN_NAME}

${BIN_DIR}/${AGENT_BIN_NAME}: ${GO_SRC}
	@mkdir -p ${BIN_DIR}
	GOOS=$(GOOS) CGO_ENABLED=0 $(GO) build -v -ldflags "${LD_FLAGS}" -o $@ ./cmd/mysql-agent

.PHONY: agent-image
agent-image: ${BIN_DIR}/${AGENT_BIN_NAME}
	docker build \
		--build-arg=http_proxy \
		--build-arg=https_proxy \
		-t ${AGENT_DOCKER_IMAGE_NAME}:${AGENT_DOCKER_IMAGE_TAG} \
		-f ./docker/mysql-agent/Dockerfile \
		.
	docker tag ${AGENT_DOCKER_IMAGE_NAME}:${AGENT_DOCKER_IMAGE_TAG} ${DOCKER_REGISTRY}/${DOCKER_USER}/${AGENT_DOCKER_IMAGE_NAME}:${AGENT_DOCKER_IMAGE_TAG}

.PHONY: agent-push
agent-push: agent-image
	@docker login -u '$(DOCKER_REGISTRY_USERNAME)' -p '$(DOCKER_REGISTRY_PASSWORD)' $(DOCKER_REGISTRY)
	@docker push ${DOCKER_REGISTRY}/${DOCKER_USER}/${AGENT_DOCKER_IMAGE_NAME}:${AGENT_DOCKER_IMAGE_TAG}

.PHONY: agent-test-container
agent-test-container: agent-build
	@docker build -t agent-test:latest . -f docker/mysql-agent/Dockerfile.dev
	@docker run --name agent-test -e MYSQL_ROOT_PASSWORD=root -v `pwd`/dist/bin/mysql-agent:/mysql-agent -itd agent-test
	@docker exec -it agent-test /bin/bash

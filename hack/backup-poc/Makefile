NAME:=mysql-operator-backup

ifdef GITLAB_CI
BUILD_DIR:=${CI_PROJECT_DIR}/dist

ifdef CI_COMMIT_TAG
VERSION:=${CI_COMMIT_TAG}
else
VERSION=${CI_COMMIT_SHA}
endif
else
# localbuild user timestamp it
BUILD_DIR:=dist
VERSION ?= ${USER}-$(shell date +%Y%m%d%H%M%S)
endif

DOCKER_REPO ?= registry.oracledx.com
DOCKER_USER ?= skeppare
DOCKER_IMAGE_NAME ?= mysql-operator-backup
DOCKER_IMAGE_TAG ?= ${VERSION}

# -------------------------------------------------------------------------------------------------
# docker dev build - for manually developing/testing bash scripts using oracle-linux-7.2 bash

DOCKER_DEV_IMAGE_NAME ?= ${DOCKER_IMAGE_NAME}-dev
DOCKER_DEV_IMAGE_TAG ?= ${USER}

# build the dev docker image - use a derivatvie mysql-enterprise-server:5.7.19-1.1.0 with mysql and adds some dev support tools.
# 
.PHONY: dev-docker-image
dev-docker-image:
	sed 's/FROM oraclelinux:7.3/FROM registry.oracledx.com\/skeppare\/mysql-enterprise-server:5.7.19-1.1.0/' Dockerfile > Dockerfile.dev
	sed -i "" 's/# ENV http_proxy/ENV http_proxy/' Dockerfile.dev
	sed -i "" 's/# ENV https_proxy/ENV https_proxy/' Dockerfile.dev
	sed -i "" 's/# RUN yum install -y vim/RUN yum install -y vim/' Dockerfile.dev

	# TODO: get these to work - then we can remove the above and Dockerfile entries
	# shell sed -i "" -e $'2 a\\\n' "ENV http_proxy http:\/\/www-proxy.us.oracle.com:80" Dockerfile.dev
	# shell sed -i "" -e $'3 a\\\n' "ENV https_proxy http:\/\/www-proxy.us.oracle.com:80" Dockerfile.dev
	# shell sed -i "" -e $'4 a\\\n' "RUN yum install -y vim" Dockerfile.dev

	docker build \
		--build-arg=http_proxy \
		--build-arg=https_proxy \
		-t ${DOCKER_DEV_IMAGE_NAME}:${DOCKER_DEV_IMAGE_TAG} \
		-f Dockerfile.dev \
		.
	rm -f Dockerfile.dev

# run the dev image - externally mount required script files to allow inline editing from host machine.
# 
dev-docker-start: dev-docker-image
	docker run \
		--name=${DOCKER_DEV_IMAGE_NAME} \
		--hostname=${DOCKER_DEV_IMAGE_NAME} \
		-v `pwd`/backup-poc.sh:/backup-poc.sh \
		-v `pwd`/backup-poc-create-user.sh:/backup-poc-create-user.sh \
		-e MYSQL_ROOT_PASSWORD=root \
		-itd ${DOCKER_DEV_IMAGE_NAME}:${DOCKER_DEV_IMAGE_TAG}

# shell into dev image
# 
dev-docker-shell:
	docker exec -it ${DOCKER_DEV_IMAGE_NAME} /bin/bash

# stop dev image
# 
dev-docker-stop:
	docker rm -f ${DOCKER_DEV_IMAGE_NAME}

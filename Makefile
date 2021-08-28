APP=asscRegistration

SHELL=bash
GOOS?=linux
GOARCH?=amd64
#REVISION=$(shell ./version.sh)
RELEASE?=$(shell git rev-parse --abbrev-ref HEAD)

DB_IMAGE=postgres:11
DB_USER?=assc
DB_PASSWORD?=blablabla
DB_HOST?=
DB_NAME?=${APP}
SCHEMA_VERSION?=

HAS_GOLANGCI_LINT := $(shell command -v golangci-lint;)
HAS_GINKGO := $(shell command -v ginkgo;)

mkfile_path := $(abspath $(MAKEFILE_LIST))
current_dir := $(patsubst %/,%,$(dir $(mkfile_path)))

export CGO_ENABLED=0
export GO111MODULE=on
export GOPROXY=https://proxy.golang.org

.PHONY: build clean vendor tests lint
.ONESHELL:

clean:
	@echo "+ $@"
	@rm -rf build
	@rm -rf tests/*/*.test
	@rm -f checkstyle.xml
	@rm -f coverage.*
	@rm -f *.coverprofile
	@rm -f ./**/junit.xml

vendor:
	@echo "+ $@"
	@go mod vendor

init: clean vendor
	@echo "+ $@"
ifndef HAS_GOLANGCI_LINT
	@GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
endif

lint:
	@echo "+ $@"
	@golangci-lint run

build:
	@echo "+ $@"
	@go build -trimpath -mod=readonly \
	-ldflags "-s -w -X ${PROJECT}/config.RELEASE=${RELEASE} -X ${PROJECT}/config.REPO=${REPO_INFO}" \
	-o build/${APP} ${PROJECT}/cmd/${APP}

db: db-stop
	@echo "+ $@"
	@docker run -p 5432:5432 -d \
		-e POSTGRES_USER=${DB_USER} \
		-e POSTGRES_PASSWORD=${DB_PASSWORD} \
		-e POSTGRES_DB=${DB_NAME} \
		--name=${APP}_db ${DB_IMAGE}

db-stop:
	@echo "+ $@"
	@docker stop ${APP}_db || true \
		&& docker rm -v ${APP}_db || true

ifeq ($(SCHEMA_VERSION),)
SCHEMA_VERSION:=up
else
SCHEMA_VERSION:=goto ${SCHEMA_VERSION}
endif

ifeq ($(DB_HOST),)
MIGRATE_NETWORK:=--net="container:${APP}_db"
else
MIGRATE_NETWORK:=
endif

migrate:
	@echo "+ $@"
	@docker run -i --rm \
		${MIGRATE_NETWORK} \
		-v "$(current_dir)/migrations:/migrations" \
		migrate/migrate:latest \
		-path /migrations/ \
		-database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:5432/${DB_NAME}?sslmode=disable" \
		${SCHEMA_VERSION}

new-migration:
	@echo "+ $@"
	@migrate create -ext sql -dir deploy/migrations -seq ${NAME}

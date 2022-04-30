BUILDDIR = bin
VERSION := $(shell git describe --tags --always --dirty)
BUILD := $(shell date +%Y-%m-%d\ %H:%M)
LDFLAGS=-ldflags="-w -s -X 'libcommon.Version=${VERSION}' -X 'libcommon.Build=${BUILD}'"

build:
	go build ${LDFLAGS} -o $(BUILDDIR)/ cmd/main.go

docker.dev:
	docker-compose -f docker/docker-compose.dev.yaml up

docker.dev.build:
	docker-compose -f docker/docker-compose.dev.yaml up --build
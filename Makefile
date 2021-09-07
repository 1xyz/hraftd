GO=go
DELETE=rm
BINARY=hraftd
BUILD_BINARY=bin/$(BINARY)
# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
SAFE_BRANCH = $(subst /,-,$(BRANCH))
# current git version short-hash
VER = $(shell git rev-parse --short HEAD)
GIT_RELEASE_TAG=$(shell git describe --tags)

DOCKER=docker
DOCKER_REPO=1xyz/hraftd
DOCKER_TAG = "$(SAFE_BRANCH)-$(VER)"



all: build

build: clean
	$(GO) build -o $(BUILD_BINARY) -v main.go

.PHONY: clean
clean:
	$(DELETE) -rf bin/
	$(GO) clean -cache

release/%: clean
	$(GO) test ./...
	@echo "build GOOS: $(subst release/,,$@) & GOARCH: amd64"
	GOOS=$(subst release/,,$@) GOARCH=amd64 $(GO) build -o bin/$(subst release/,,$@)/$(BINARY) -v main.go


# test w/ race detector on always
# https://golang.org/doc/articles/race_detector.html#Typical_Data_Races
.PHONY: test
test: build
	$(GO) test -v -race ./...


docker-build:
	$(DOCKER) build -t $(DOCKER_REPO):$(DOCKER_TAG) -f Dockerfile .

docker-scan:
	$(DOCKER) scan --accept-license  $(DOCKER_REPO):$(DOCKER_TAG)

docker-push: docker-build
	$(DOCKER) push $(DOCKER_REPO):$(DOCKER_TAG)

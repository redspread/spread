BASE := rsprd.com/spread
DIR := $(GOPATH)/src/$(BASE)
CMD_NAME := spread

EXEC_PKG := $(BASE)/cmd/$(CMD_NAME)
PKGS := ./pkg/... ./cli/... ./cmd/...

# set spread version if not provided by environment
ifndef SPREAD_VERSION
	SPREAD_VERSION := v0.0.0
endif

LIBGIT2_VERSION ?= v0.23.4

GOX_OS ?= linux darwin windows
GOX_ARCH ?= amd64

GO ?= go
GOX ?= gox
GOFMT ?= gofmt "-s"
GOLINT ?= golint
DOCKER ?= docker

GOFILES := find . -name '*.go' -not -path "./vendor/*"

VERSION_LDFLAG := -X main.Version=$(SPREAD_VERSION)

LIBGIT2_DIR := $(DIR)/vendor/libgit2
LIBGIT2_BUILD := $(LIBGIT2_DIR)/build
LIBGIT2_PKGCONFIG := $(LIBGIT2_BUILD)/libgit2.pc

LIBGIT2_CGOFLAG := -I$(LIBGIT2_DIR)/include

# Do not use sudo on OS X
ifneq ($(OS),Windows_NT)
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Darwin)
		SUDO :=
	else
		SUDO := sudo
	endif
endif

GOBUILD_LDFLAGS ?= $(VERSION_LDFLAG)
GOBUILD_FLAGS ?= -i -v
GOTEST_FLAGS ?= -v
GOX_FLAGS ?= -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}" -os="${GOX_OS}" -arch="${GOX_ARCH}"
CGO_ENV ?= CGO_CFLAGS="$(LIBGIT2_CGOFLAG)"

STATIC_LDFLAGS ?= -extldflags "-static" --s -w

GITLAB_CONTEXT ?= ./build/gitlab
LIBGIT2_URL ?= https://github.com/libgit2/libgit2/archive/$(LIBGIT2_VERSION).tar.gz

# image data
ORG ?= redspreadapps
NAME ?= gitlabci
TAG ?= latest

GITLAB_IMAGE_NAME = "$(ORG)/$(NAME):$(TAG)"

GOX_OS ?= linux darwin windows
GOX_ARCH ?= amd64

.PHONY: all
all: clean validate test

.PHONY: release
release: validate test crossbuild

.PHONY: test
test: unit integration

.PHONY: unit
unit: build
	$(GO) test $(GOTEST_FLAGS) $(PKGS)

.PHONY: integration
integration: build
	mkdir -p ./build
	./test/mattermost-demo.sh

.PHONY: validate
validate: lint checkgofmt vet

.PHONY: build
build: build/spread

build/spread:
	$(CGO_ENV) $(GO) build $(GOBUILD_FLAGS) -ldflags "$(GOBUILD_LDFLAGS)" -o $@ $(EXEC_PKG)

build/spread-linux-static:
	GOOS=linux $(GO) build -o $@ $(GOBUILD_FLAGS) -ldflags "$(GOBUILD_LDFLAGS) $(STATIC_LDFLAGS)" $(EXEC_PKG)
	chmod +x $@

.PHONY: crossbuild
crossbuild: deps gox-setup
	$(GOX) $(GOX_FLAGS) -gcflags="$(GOBUILD_FLAGS)" -ldflags="$(GOBUILD_LDFLAGS) $(STATIC_LDFLAGS)" $(EXEC_PKG)

.PHONY: build-gitlab
build-gitlab: build/spread-linux-static
	rm -rf $(GITLAB_CONTEXT)
	cp -r ./images/gitlabci $(GITLAB_CONTEXT)
	cp ./build/spread-linux-static $(GITLAB_CONTEXT)
	$(DOCKER) build $(DOCKER_OPTS) -t $(GITLAB_IMAGE_NAME) $(GITLAB_CONTEXT)

install-libgit2: build-libgit2
	cd $(LIBGIT2_BUILD) && $(SUDO) make install

build-libgit2: $(LIBGIT2_PKGCONFIG)
$(LIBGIT2_PKGCONFIG): vendor/libgit2
	./hack/build-libgit2.sh

vendor/libgit2: vendor/libgit2.tar.gz
	mkdir -p $@
	tar -zxf $< -C $@ --strip-components=1
vendor/libgit2.tar.gz:
	curl -L -o $@ $(LIBGIT2_URL)

.PHONY: vet
vet:
	$(GO) vet $(PKGS)

lint: .golint-install
	for pkg in $(PKGS); do \
		echo "Running golint on $$pkg:"; \
		golint $$pkg; \
	done;

.PHONY: checkgofmt
checkgofmt:
	# get all go files and run go fmt on them
	files=$$($(GOFILES) | xargs $(GOFMT) -l); if [ -n "$$files" ]; then \
		  echo "Error: '$(GOFMT)' needs to be run on:"; \
		  echo "$${files}"; \
		  exit 1; \
		  fi;

.PHONY: deps
deps: .golint-install .gox-install

.golint-install:
	$(GO) get -x github.com/golang/lint/golint > $@

PHONY: gox-setup
gox-setup: .gox-install

.gox-install:
	$(GO) get -x github.com/mitchellh/gox > $@

.PHONY: clean
clean:
	rm -vf .gox-* .golint-*
	rm -rfv ./build
	$(GO) clean $(PKGS) || true

.PHONY: godep
godep:
	go get -u -v github.com/tools/godep
	@echo "Recalculating godeps, removing Godeps and vendor if not canceled in 5 seconds"
	@sleep 5
	rm -rf Godeps vendor
	GO15VENDOREXPERIMENT="1" godep save -v ./pkg/... ./cli/... ./cmd/...

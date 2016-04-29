BASE := rsprd.com/spread
CMD_NAME := spread

EXEC_PKG := $(BASE)/cmd/$(CMD_NAME)
PKGS := ./pkg/... ./cli/... ./cmd/...

# set spread version if not provided by environment
ifndef SPREAD_VERSION
	SPREAD_VERSION := v0.0.0
endif

GOX_OS ?= linux darwin windows
GOX_ARCH ?= amd64

GO ?= go
GOX ?= gox
GOFMT ?= gofmt # eventually should add "-s"
GOLINT ?= golint
DOCKER ?= docker

VERSION_LDFLAG := -X main.Version=$(SPREAD_VERSION)

GOFILES := find . -name '*.go' -not -path "./vendor/*"

GOBUILD_LDFLAGS ?= $(VERSION_LDFLAG)
GOBUILD_FLAGS ?= -i -v
GOTEST_FLAGS ?= -v
GOX_FLAGS ?= -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}" -os="${GOX_OS}" -arch="${GOX_ARCH}"

STATIC_LDFLAGS ?= -extldflags "-static" --s -w

GITLAB_CONTEXT ?= ./build/gitlab

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
	$(GO) build $(GOBUILD_FLAGS) -ldflags "$(GOBUILD_LDFLAGS)" -o $@ $(EXEC_PKG)

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

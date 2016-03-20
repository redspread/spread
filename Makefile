BASE := rsprd.com/spread
CMD_NAME := spread

EXEC_PKG := $(BASE)/cmd/$(CMD_NAME)
PKGS := ./pkg/... ./cli/... ./cmd/...

GOX_OS ?= linux darwin windows
GOX_ARCH ?= amd64

GO ?= go
GOX ?= gox
GOFMT ?= gofmt # eventually should add "-s"
GOLINT ?= golint

GOFILES := find . -name '*.go' -not -path "./vendor/*"

GOBUILD_LDFLAGS ?=
GOBUILD_FLAGS ?=
GOTEST_FLAGS ?= -v
GOX_FLAGS ?= -output="build/{{.Dir}}_{{.OS}}_{{.Arch}}" -os="${GOX_OS}" -arch="${GOX_ARCH}"

GOX_OS ?= linux darwin windows
GOX_ARCH ?= amd64

.PHONY: all
all: clean validate test

.PHONY: release
release: validate test crossbuild

.PHONY: test
test: build
	$(GO) test $(GOTEST_FLAGS) $(PKGS)

.PHONY: validate
validate: lint checkgofmt vet

.PHONY: build
build:
	$(GO) install $(GOBUILD_FLAGS) $(GOBUILD_LDFLAGS) $(EXEC_PKG)

.PHONY: crossbuild
crossbuild: deps gox-setup
	$(GOX) $(GOX_FLAGS) -gcflags="$(GOBUILD_FLAGS)" -ldflags="$(GOBUILD_LDFLAGS)" $(EXEC_PKG)

.PHONY: vet
vet:
	$(GO) vet $(PKGS) 

lint: .golint-install
	for pkg in $(PKGS); do \
		echo "Running golint on $$i:"; \
		golint $$i; \
	done;

.PHONY: checkgofmt
checkgofmt:
	# get all go files and run go fmt on them
	files=$$($(GOFILES) | xargs $(GOFMT) -l); if [[ -n "$$files" ]]; then \
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

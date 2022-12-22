XC_OS="linux darwin"
XC_ARCH="amd64"
XC_PARALLEL="2"
BIN="bin"
SRC=$(shell find . -name "*.go")
GIT_SHA?=$(shell git rev-parse HEAD)
VERSION?=$(shell git describe --always --tags | cut -d "v" -f 2)
LINKER_FLAGS="-s -w -X github.com/spf13/cobra-cli/cmd.Version=${VERSION} -X github.com/spf13/cobra-cli/cmd.GitCommit=${GIT_SHA}"
ifeq (, $(shell which golangci-lint))
$(warning "could not find golangci-lint in $(PATH), run: curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh")
endif

ifeq (, $(shell which richgo))
$(warning "could not find richgo in $(PATH), run: go get github.com/kyoh86/richgo")
endif

ifeq (, $(shell which gox))
$(warning "could not find gox in $(PATH), run: go get github.com/mitchellh/gox")
endif

.PHONY: all build fmt lint test install_deps clean

default: all

all: fmt test build

build: install_deps
	gox \
		-os=$(XC_OS) \
		-arch=$(XC_ARCH) \
		-parallel=$(XC_PARALLEL) \
		-ldflags=$(LINKER_FLAGS) \
		-output=$(BIN)/{{.Dir}}_{{.OS}}_{{.Arch}} \
		;

fmt:
	$(info ******************** checking formatting ********************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint:
	$(info ******************** running lint tools ********************)
	golangci-lint run -v

test: install_deps
	$(info ******************** running tests ********************)
	richgo test -v ./...

install_deps:
	$(info ******************** downloading dependencies ********************)
	go get -v ./...

clean:
	rm -rf $(BIN)

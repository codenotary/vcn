SHELL=/bin/bash -o pipefail

VERSION=0.8.2
TARGETS=linux/amd64 windows/amd64 darwin/amd64 linux/s390x linux/ppc64le linux/arm-7 linux/arm64

GO ?= go
DOCKER ?= docker

GIT_REV := $(shell git rev-parse HEAD 2> /dev/null || true)
GIT_COMMIT := $(if $(shell git status --porcelain --untracked-files=no),${GIT_REV}-dirty,${GIT_REV})
GIT_BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)

export GO111MODULE=on
PWD=$(shell pwd)
LDFLAGS := -s -X github.com/vchain-us/vcn/pkg/meta.version=v${VERSION} \
			  -X github.com/vchain-us/vcn/pkg/meta.gitCommit=${GIT_COMMIT} \
			  -X github.com/vchain-us/vcn/pkg/meta.gitBranch=${GIT_BRANCH}
LDFLAGS_STATIC := ${LDFLAGS} \
				  -X github.com/vchain-us/vcn/pkg/meta.static=static \
				  -extldflags "-static"
TEST_FLAGS ?= -v -race
VCNEXE=vcn-v${VERSION}-windows-amd64.exe
SETUPEXE=codenotary_vcn_v${VERSION}_setup.exe

.PHONY: vcn
vcn:
	$(GO) build -ldflags '${LDFLAGS} -X github.com/vchain-us/vcn/pkg/meta.version=v${VERSION}-dev' ./cmd/vcn

.PHONY: vendor
vendor:
	$(GO) mod vendor

.PHONY: test
test:
	$(GO) vet ./...
	$(GO) test ${TEST_FLAGS} ./...

.PHONY: install
install: TEST_FLAGS=-v
install: vendor test
	$(GO) install -ldflags '${LDFLAGS}' ./cmd/vcn

.PHONY: static
static:
	$(GO) build -a -tags netgo -ldflags '${LDFLAGS_STATIC}' ./cmd/vcn

.PHONY: docs/cmd
docs/cmd:
	rm -rf docs/cmd/*.md
	$(GO) run docs/cmd/main.go

.PHONY: build/xgo
build/xgo:
	$(DOCKER) build \
			-f ./build/xgo/Dockerfile \
			-t vcn-xgo \
			--pull=true \
			./build/xgo

.PHONY: build/makensis
build/makensis:
	$(DOCKER) build \
		-f ./build/makensis/Dockerfile \
		-t vcn-makensis \
		./build/makensis

.PHONY: clean/dist
clean/dist: 
	rm -Rf ./dist

.PHONY: clean
clean: clean/dist
	rm -f ./vcn

.PHONY: CHANGELOG.md
CHANGELOG.md:
	git-chglog -o CHANGELOG.md

.PHONY: CHANGELOG.md.next-tag
CHANGELOG.md.next-tag:
	git-chglog -o CHANGELOG.md --next-tag v${VERSION}

.PHONY: dist
dist: clean/dist build/xgo
	mkdir -p dist
	$(GO) build -a -tags netgo -ldflags '${LDFLAGS_STATIC}' \
			-o ./dist/vcn-v${VERSION}-linux-amd64-static \
			./cmd/vcn 
	$(DOCKER) run --rm \
			-v ${PWD}/dist:/dist \
			-v ${PWD}:/source:ro \
			-e GO111MODULE=on \
			-e FLAG_LDFLAGS="${LDFLAGS}" \
			-e TARGETS="${TARGETS}" \
			-e PACK=cmd/vcn \
			-e OUT=vcn-v${VERSION} \
			vcn-xgo .
	mv ./dist/vcn-v${VERSION}-linux-arm-7 ./dist/vcn-v${VERSION}-linux-arm
	mv ./dist/vcn-v${VERSION}-windows-4.0-amd64.exe ./dist/${VCNEXE}
	mv ./dist/vcn-v${VERSION}-darwin-10.6-amd64 ./dist/vcn-v${VERSION}-darwin-amd64

.PHONY: dist/${VCNEXE} dist/${SETUPEXE}
dist/${VCNEXE} dist/${SETUPEXE}:
	echo ${SIGNCODE_PVK_PASSWORD} | $(DOCKER) run --rm -i \
		-v ${PWD}/dist:/dist \
		-v ${SIGNCODE_SPC}:/certs/f.spc:ro \
		-v ${SIGNCODE_PVK}:/certs/f.pvk:ro \
		mono:5.20 signcode \
		-spc /certs/f.spc -v /certs/f.pvk \
		-a sha1 -$ commercial \
		-n "CodeNotary vcn" \
		-i https://codenotary.io/ \
		-t http://timestamp.comodoca.com -tr 10 \
		$@
	rm -Rf $@.bak

.PHONY: dist/NSIS
dist/NSIS: build/makensis
	mkdir -p dist/NSIS
	cp -f ./dist/${VCNEXE} ./dist/NSIS/vcn.exe
	cp -f ./build/NSIS/* ./dist/NSIS/
	sed -e "s/{VCN_VERSION}/v${VERSION}/g" ./build/NSIS/setup.nsi > ./dist/NSIS/setup.nsi
	$(DOCKER) run --rm \
			-v ${PWD}/dist/NSIS/:/app \
			vcn-makensis /app/setup.nsi
	cp ./dist/NSIS/*_setup.exe ./dist/
	rm -Rf ./dist/NSIS

.PHONY: dist/sign
dist/sign: vendor vcn
	for f in ./dist/*; do ./vcn sign -p $$f; printf "\n\n"; done

.PHONY: dist/all
dist/all: dist dist/${VCNEXE} dist/NSIS dist/${SETUPEXE}

.PHONY: dist/binary.md
dist/binary.md:
	@for f in ./dist/*; do \
		ff=$$(basename $$f); \
		shm_id=$$(sha256sum $$f | awk '{print $$1}'); \
		printf "[$$ff](https://github.com/vchain-us/vcn/releases/download/v${VERSION}/$$ff) | $$shm_id \n" ; \
	done

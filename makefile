export GO111MODULE=on

SOLOFS_LDFLAGS += -X "soloos/solofs/version.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
SOLOFS_LDFLAGS += -X "soloos/solofs/version.GitHash=$(shell git rev-parse HEAD)"
# SOLOFS_PREFIX += GOTMPDIR=./go.build/tmp GOCACHE=./go.build/cache

SOLOOS_SOLOFS_PROTOS = $(shell find ./ -name '*.fbs'|grep -v vendor)
GENERATED_PROTOS = $(shell find ./ -name "*.fbs"|grep -v vendor| sed 's/\.fbs/\.fbs\.go/g')
SOURCES = $(shell find . -name '*.go') $(GENERATED_PROTOS)

GOBUILD = $(SOLOFS_PREFIX) go build -i -ldflags '$(SOLOFS_LDFLAGS)' 

clean-test-cache:
	go clean -testcache

%.fbs.go: $(SOLOOS_SOLOFS_PROTOS)
	flatc -o ./ -g $(SOLOOS_SOLOFS_PROTOS)

fbs: $(GENERATED_PROTOS)

all:solofsd solofsd-mock solofssdk

libsolofs:
	$(GOBUILD) -o ./bin/libsolofs.so -buildmode=c-shared ./apps/libsolofs

solofsd:
	$(GOBUILD) -o ./bin/solofsd ./apps/solofsd

solofsd-fuse:
	rm -f bin/solofsd-fuse
	$(GOBUILD) -o ./bin/solofsd-fuse ./apps/solofsd-fuse

include ./make/test
include ./make/bench

.PHONY:all solofsd solofsd-fuse libsolofs test

export GO111MODULE=on

SDFS_LDFLAGS += -X "soloos/sdfs/version.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
SDFS_LDFLAGS += -X "soloos/sdfs/version.GitHash=$(shell git rev-parse HEAD)"
# SDFS_PREFIX += GOTMPDIR=./go.build/tmp GOCACHE=./go.build/cache

SOLOOS_SDFS_PROTOS = $(shell find ./ -name '*.fbs'|grep -v vendor)
GENERATED_PROTOS = $(shell find ./ -name "*.fbs"|grep -v vendor| sed 's/\.fbs/\.fbs\.go/g')
SOURCES = $(shell find . -name '*.go') $(GENERATED_PROTOS)

GOBUILD = $(SDFS_PREFIX) go build -i -ldflags '$(SDFS_LDFLAGS)' 

clean-test-cache:
	go clean -testcache

%.fbs.go: $(SOLOOS_SDFS_PROTOS)
	flatc -o ./ -g $(SOLOOS_SDFS_PROTOS)

fbs: $(GENERATED_PROTOS)

all:sdfsd sdfsd-mock sdfssdk

libsdfs:
	$(GOBUILD) -o ./bin/libsdfs.so -buildmode=c-shared ./apps/libsdfs

sdfsd:
	$(GOBUILD) -o ./bin/sdfsd ./apps/sdfsd

sdfsd-fuse:
	rm -f bin/sdfsd-fuse
	$(GOBUILD) -o ./bin/sdfsd-fuse ./apps/sdfsd-fuse

include ./make/test
include ./make/bench

.PHONY:all sdfsd sdfsd-fuse libsdfs test

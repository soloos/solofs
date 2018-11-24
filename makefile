SDFS_LDFLAGS += -X "soloos/sdfs/version.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
SDFS_LDFLAGS += -X "soloos/sdfs/version.GitHash=$(shell git rev-parse HEAD)"
# SDFS_PREFIX += GOTMPDIR=./go.build/tmp GOCACHE=./go.build/cache

SOLOOS_SDFS_PROTOS = $(shell find lib/soloos/sdfs -name '*.fbs')
GENERATED_PROTOS = $(shell find lib/soloos/sdfs -name "*.fbs"| sed 's/\.fbs/\.fbs\.go/g')
SOURCES = $(shell find . -name '*.go') $(GENERATED_PROTOS)

%.fbs.go: $(SOLOOS_SDFS_PROTOS)
	flatc -o ./lib/soloos/sdfs -g $(SOLOOS_SDFS_PROTOS)

fbs: $(GENERATED_PROTOS)

all:sdfsd sdfsd-mock

sdfsd:
	$(SDFS_PREFIX) go build -i -ldflags '$(SDFS_LDFLAGS)' -o ./bin/sdfsd sdfsd

sdfsd-mock:
	$(SDFS_PREFIX) go build -i -ldflags '$(SDFS_LDFLAGS)' -o ./bin/sdfsd-mock sdfsd-mock

include ./make/test
include ./make/bench

.PHONY:all soloos-server test

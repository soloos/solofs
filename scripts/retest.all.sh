#!/bin/bash
make sdfsd-fuse
make sdfsd
./scripts/stop.cluster.sh
./scripts/stop.sdfsd-fuse.sh
./scripts/reformat.mysql.sh
./scripts/clean.cluster.sh
./scripts/start.cluster.sh
./scripts/start.sdfsd-fuse.sh

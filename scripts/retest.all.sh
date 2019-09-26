#!/bin/bash
make solofsd-fuse
make solofsd
./scripts/stop.cluster.sh
./scripts/stop.solofsd-fuse.sh
./scripts/reformat.mysql.sh
./scripts/clean.cluster.sh
./scripts/start.cluster.sh
./scripts/start.solofsd-fuse.sh

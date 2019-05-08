#!/bin/bash
./scripts/stop.cluster.sh
./scripts/reformat.mysql.sh
rm -rf /tmp/sdfs*

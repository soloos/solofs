#!/bin/bash
rm -rf /tmp/sdfs*
./bin/sdfsd ./scripts/namenode.1.json > ./logs/namenode.1.log 2>&1 &
sleep 1
./bin/sdfsd ./scripts/datanode.1.json > ./logs/datanode.1.log 2>&1 &
./bin/sdfsd ./scripts/datanode.2.json > ./logs/datanode.2.log 2>&1 &
./bin/sdfsd ./scripts/datanode.3.json > ./logs/datanode.3.log 2>&1 &
sleep 1

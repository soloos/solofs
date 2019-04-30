#!/bin/bash
rm -rf /tmp/sdfs*
./bin/sdfsd ./scripts/conf/namenode.1.json > ./logs/namenode.1.log 2>&1 &
sleep 1

for i in {1..3}
do
./bin/sdfsd ./scripts/conf/datanode.${i}.json > ./logs/datanode.1.log 2>&1 &
done

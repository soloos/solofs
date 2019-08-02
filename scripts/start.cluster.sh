#!/bin/bash
rm -rf /tmp/sdfs*
bash scripts/front.namenode.sh > ./logs/namenode.1.log 2>&1 &
sleep 1

for i in {0..3}
do
        ./bin/sdfsd ./scripts/conf/datanode.${i}.json > ./logs/datanode.${i}.log 2>&1 &
done

#!/bin/bash
rm -rf /tmp/solofs*
bash scripts/front.solonn.sh > ./logs/solonn.1.log 2>&1 &
sleep 1

for i in {0..3}
do
        ./bin/solofsd ./scripts/conf/solodn.${i}.json > ./logs/solodn.${i}.log 2>&1 &
done

#!/bin/bash
cd ./scripts
jq -s add ./conf_template/namenode.1.json ./conf_template/db.json > ./conf/namenode.1.json
for i in {1..3}
do
        jq -s add ./conf_template/datanode.${i}.json ./conf_template/db.json > ./conf/datanode.${i}.json
done
jq -s add ./conf_template/sdfsd-fuse.json ./conf_template/db.json > ./conf/sdfsd-fuse.json

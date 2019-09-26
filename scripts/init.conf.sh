#!/bin/bash
scripts/init.conf.template.py
cd ./scripts
jq -s add ./conf_template/solonn.1.json ./conf_template/solonn.common.json ./conf_template/common.json > ./conf/solonn.1.json
for i in {0..3}
do
        jq -s add ./conf_template/solodn.${i}.json ./conf_template/solodn.common.json ./conf_template/common.json > ./conf/solodn.${i}.json
done
jq -s add ./conf_template/solofsd-fuse.json ./conf_template/common.json > ./conf/solofsd-fuse.json

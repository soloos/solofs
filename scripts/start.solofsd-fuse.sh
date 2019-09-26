#!/bin/bash
nohup ./bin/solofsd-fuse ./scripts/conf/solofsd-fuse.json > logs/solofsd-fuse.log 2>&1 &
sleep 1
rm -rf /tmp/rocksdbtest-1002
ln -s /opt/soloos/solofsd-fuse-mnt/tmp/rocksdbtest-1002 /tmp/rocksdbtest-1002

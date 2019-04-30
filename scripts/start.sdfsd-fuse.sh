#!/bin/bash
nohup ./bin/sdfsd-fuse ./scripts/conf/sdfsd-fuse.json > logs/sdfsd-fuse.log 2>&1 &
sleep 1
rm -rf /tmp/rocksdbtest-1002
ln -s /opt/soloos/sdfsd-fuse-mnt/tmp/rocksdbtest-1002 /tmp/rocksdbtest-1002

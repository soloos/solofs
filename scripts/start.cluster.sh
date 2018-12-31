#!/bin/bash
./bin/sdfsd namenode 127.0.0.1:10601 > ./logs/namenode.10000.log 2>&1 &
sleep 1
./bin/sdfsd datanode 00001 127.0.0.1:10701 /tmp/sdfs_test.data.01 10000 127.0.0.1:10601 > ./logs/datanode.10000.log 2>&1 &

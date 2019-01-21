#!/bin/bash
#ps aux|grep sdfsd-fuse|grep bin|grep -v grep |awk '{print $2}'|xargs kill -SIGABRT
sudo bash -c 'umount /opt/soloos/sdfsd-fuse-mnt || true'

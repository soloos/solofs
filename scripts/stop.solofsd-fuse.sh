#!/bin/bash
#ps aux|grep solofsd-fuse|grep bin|grep -v grep |awk '{print $2}'|xargs kill -SIGABRT
sudo bash -c 'umount /opt/soloos/solofsd-fuse-mnt || true'

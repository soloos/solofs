#!/bin/bash
#ps aux|grep sdfsd-fuse|grep bin|grep -v grep |awk '{print $2}'|xargs kill -SIGABRT
sudo umount /opt/soloos/sdfsd-fuse-mnt

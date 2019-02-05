#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import os
import stat

filename = '/opt/soloos/sdfsd-fuse-mnt/testhi'
mode = 0600|stat.S_IRUSR

# 文件系统节点指定不同模式
os.mknod(filename, mode)

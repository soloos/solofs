#!/bin/bash
setfattr -n user.hello -v "hello" /opt/soloos/solofsd-fuse-mnt/tmp
setfattr -n user.hello -v "fuck" /opt/soloos/solofsd-fuse-mnt/tmp
setfattr -n user.hi -v "shit" /opt/soloos/solofsd-fuse-mnt/tmp
getfattr -d -m -  /opt/soloos/solofsd-fuse-mnt/tmp
setfattr -x user.hello /opt/soloos/solofsd-fuse-mnt/tmp
getfattr -d -m -  /opt/soloos/solofsd-fuse-mnt/tmp

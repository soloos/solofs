#!/bin/bash
setfattr -n user.hello -v "hello" /opt/soloos/sdfsd-fuse-mnt/tmp
setfattr -n user.hello -v "fuck" /opt/soloos/sdfsd-fuse-mnt/tmp
setfattr -n user.hi -v "shit" /opt/soloos/sdfsd-fuse-mnt/tmp
getfattr -d -m -  /opt/soloos/sdfsd-fuse-mnt/tmp
setfattr -x user.hello /opt/soloos/sdfsd-fuse-mnt/tmp
getfattr -d -m -  /opt/soloos/sdfsd-fuse-mnt/tmp

#!/bin/bash
ps aux|grep solofsd|grep fuse|grep -v cpu|grep -v pprof|grep bin|grep -v grep |awk '{print $2}'|xargs kill -SIGABRT 2>/dev/null

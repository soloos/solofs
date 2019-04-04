#!/bin/bash
ps aux|grep memstg|grep tmp|awk '{print $2}'|xargs kill -SIGABRT

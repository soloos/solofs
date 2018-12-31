#!/bin/bash
ps aux|grep sdfsd|grep bin|grep -v grep |awk '{print $2}'|xargs kill -9

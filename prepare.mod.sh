#!/bin/bash
export GOFLAGS=-mod=vendor
go mod vendor
rm -rf vendor/soloos
mkdir -p vendor/soloos
ln -s /soloos/sdbone ./vendor/soloos/
ln -s /soloos/common ./vendor/soloos/


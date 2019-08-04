#!/bin/bash
rm -rf vendor/soloos
mkdir vendor/soloos
ln -s /opt/soloos/common ./vendor/soloos/common
ln -s /opt/soloos/sdbone ./vendor/soloos/sdbone
ln -s /opt/soloos/sdfs ./vendor/soloos/sdfs
ln -s /opt/soloos/soloboat ./vendor/soloos/soloboat
ln -s /opt/soloos/swal ./vendor/soloos/swal
export GOFLAGS=-mod=vendor
go mod vendor

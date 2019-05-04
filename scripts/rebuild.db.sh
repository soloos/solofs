#!/bin/bash
mysql -uroot -phello mysql -e "drop database sdfs_test;"
mysql -uroot -phello mysql -e "CREATE DATABASE sdfs_test /*!40100 DEFAULT CHARACTER SET utf8 */;"

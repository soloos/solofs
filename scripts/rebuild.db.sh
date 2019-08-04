#!/bin/bash
mysql -uroot -phello mysql -e "drop database soloos_test;"
mysql -uroot -phello mysql -e "CREATE DATABASE soloos_test /*!40100 DEFAULT CHARACTER SET utf8 */;"

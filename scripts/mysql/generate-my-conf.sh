#!/bin/sh

DIR=$(cd $(dirname $0) && pwd)

cat << EOS > ${DIR}/my.conf
[client]
user = root
password = root
host = 127.0.0.1
port = ${GQB_MYSQL_PORT}
EOS

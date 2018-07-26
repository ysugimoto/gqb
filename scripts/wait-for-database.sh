#!/bin/bash

DIR=$(cd $(dirname $0) && pwd)
MYSQLCMD="mysql --defaults-extra-file=${DIR}/my.conf"

for i in $(seq 10); do
  $MYSQLCMD -e 'show databases' || (sleep 5; false) && break
done
echo "MySQL has been started"

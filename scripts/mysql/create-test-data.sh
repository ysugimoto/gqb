#!/bin/bash

DIR=$(cd $(dirname $0) && pwd)
MYSQLCMD="mysql --defaults-extra-file=${DIR}/my.conf"

echo "Creating database for testing..."
echo "CREATE DATABASE IF NOT EXISTS gqb_test;" | $MYSQLCMD

echo "Creating table for testing..."
$MYSQLCMD --database gqb_test << EOS
CREATE TABLE IF NOT EXISTS companies (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(255)  NOT NULL,
  PRIMARY KEY (id)
) DEFAULT CHARSET=utf8;
EOS

echo "Reset and insert 10000 records..."
echo "TRUNCATE TABLE companies;" | $MYSQLCMD --database gqb_test

COUNT=$1
i=0
j=0
SQL="INSERT INTO companies (name) VALUES "
VALUES=""
while [ $i -lt $COUNT ]; do
  VALUES="${VALUES}, ('company_$i')"
  j=$((++j))
  i=$((++i))
  if [ $j -eq 100 ]; then
    echo "${SQL}${VALUES:2}" | $MYSQLCMD --database gqb_test
    VALUES=""
    j=0
  fi
done
if [ "$VALUES" != "" ]; then
  echo "${SQL}${VALUES:2}" | $MYSQLCMD --database gqb_test
fi
echo "Done."

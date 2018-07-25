#!/bin/bash

echo "CREATE DATABASE IF NOT EXISTS gqb_test;" | mysql -uroot -p -h127.0.0.1 -proot

mysql -uroot -p -h127.0.0.1 -proot gqb_test << EOS
CREATE TABLE IF NOT EXISTS companies (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(255)  NOT NULL,
  PRIMARY KEY (id)
) DEFAULT CHARSET=utf8;
EOS

echo "TRUNCATE TABLE companies;" | mysql -uroot -p -h127.0.0.1 -proot gqb_test

COUNT=10000
i=0
j=0
SQL="INSERT INTO companies (name) VALUES "
VALUES=""
while [ $i -lt $COUNT ]; do
  VALUES="${VALUES}, ('company_$i')"
  j=$((++j))
  i=$((++i))
  if [ $j -eq 100 ]; then
    echo "RUN INSERT"
    echo "${SQL}${VALUES:2}" | mysql -uroot -p -h127.0.0.1 -proot gqb_test
    VALUES=""
    j=0
  fi
done
if [ "$VALUES" != "" ]; then
  echo "${SQL}${VALUES:2}" | mysql -uroot -p -h127.0.0.1 -proot gqb_test
fi

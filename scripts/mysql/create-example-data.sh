#!/bin/bash

DIR=$(cd $(dirname $0) && pwd)
MYSQLCMD="mysql --defaults-extra-file=${DIR}/my.conf"

echo "Creating database for example..."
echo "CREATE DATABASE IF NOT EXISTS example;" | $MYSQLCMD

echo "Creating table for testing..."
$MYSQLCMD --database example << EOS
DROP TABLE IF EXISTS companies;

CREATE TABLE companies (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(255) NOT NULL,
  created_at datetime NOT NULL,
  PRIMARY KEY (id)
) DEFAULT CHARSET=utf8;

INSERT INTO companies (name, created_at) VALUES ('Google', '2018-07-30 00:00:00'), ('Apple', '2018-07-30 00:00:00'), ('Microsoft', '2018-07-30 00:00:00');

DROP TABLE IF EXISTS company_attributes;

CREATE TABLE company_attributes (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  company_id int(11) unsigned NOT NULL,
  url varchar(255) NOT NULL,
  PRIMARY KEY (id)
) DEFAULT CHARSET=utf8;

INSERT INTO company_attributes (company_id, url) VALUES (1, 'https://google.com'), (2, 'https://apple.com'), (3, 'https://microsoft.com');
EOS

echo "Done."

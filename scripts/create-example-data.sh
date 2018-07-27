#!/bin/bash

DIR=$(cd $(dirname $0) && pwd)
MYSQLCMD="mysql --defaults-extra-file=${DIR}/my.conf"

echo "Creating database for example..."
echo "CREATE DATABASE IF NOT EXISTS example;" | $MYSQLCMD

echo "Creating table for testing..."
$MYSQLCMD --database gqb_test << EOS
CREATE TABLE IF NOT EXISTS companies (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(255) NOT NULL,
  PRIMARY KEY (id)
) DEFAULT CHARSET=utf8;

TRUNCATE TABLE companies;
INSERT INTO companies (name) VALUES ('Google'), ('Apple'), ('Microsoft');

CREATE TABLE IF NOT EXISTS company_attributes (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  company_id int(11) unsigned NOT NULL,
  url varchar(255) NOT NULL,
  PRIMARY KEY (id)
) DEFAULT CHARSET=utf8;

TRUNCATE TABLE company_attributes;
INSERT INTO company_attributes (company_id, url) VALUES (1, 'https://google.com'), (2, 'https://apple.com'), (3, 'https://microsoft.com');
EOS

echo "Done."

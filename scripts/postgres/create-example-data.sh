#!/bin/bash

DIR=$(cd $(dirname $0) && pwd)
PSQLCMD="psql -U postgres -h 127.0.0.1"

echo "Creating table for testing..."
$PSQLCMD -d postgres << EOS
DROP TABLE IF EXISTS companies;
CREATE TABLE companies (
  id SERIAL,
  name varchar(255) NOT NULL,
  created_at timestamp
);

INSERT INTO companies (id, name, created_at) VALUES (1, 'Google', '2018-07-30 00:00:00'), (2, 'Apple', '2018-07-30 00:00:00'), (3, 'Microsoft', '2018-07-30 00:00:00');

DROP TABLE IF EXISTS company_attributes;
CREATE TABLE company_attributes (
  id SERIAL,
  company_id int NOT NULL,
  url varchar(255) NOT NULL
);

TRUNCATE TABLE company_attributes;
INSERT INTO company_attributes (id, company_id, url) VALUES (1, 1, 'https://google.com'), (2, 2, 'https://apple.com'), (3, 3, 'https://microsoft.com');
EOS

echo "Done."

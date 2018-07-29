#!/bin/bash

DIR=$(cd $(dirname $0) && pwd)
PSQLCMD="psql -U postgres -h 127.0.0.1"

echo "Creating table for testing..."
$PSQLCMD -d postgres << EOS
DROP TABLE IF EXISTS companies;
CREATE TABLE companies (
  id SERIAL,
  name varchar(255) NOT NULL
);

INSERT INTO companies (id, name) VALUES (1, 'Google'), (2, 'Apple'), (3, 'Microsoft');

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

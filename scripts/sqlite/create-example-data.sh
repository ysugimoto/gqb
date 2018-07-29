#!/bin/bash

DIR=$(cd $(dirname $0) && pwd)
SQLITECMD="sqlite3 /tmp/gqb_test.sqlite "

echo "Creating table for testing..."
$SQLITECMD << EOS
DROP TABLE If EXISTS companies;
CREATE TABLE companies (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT
);

INSERT INTO companies (name) VALUES ('Google'), ('Apple'), ('Microsoft');

DROP TABLE IF EXISTS company_attributes;
CREATE TABLE company_attributes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  company_id INT,
  url TEXT
);

INSERT INTO company_attributes (company_id, url) VALUES (1, 'https://google.com'), (2, 'https://apple.com'), (3, 'https://microsoft.com');
EOS

echo "Done."

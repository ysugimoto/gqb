#!/bin/bash

DIR=$(cd $(dirname $0) && pwd)
SQLITECMD="sqlite3 /tmp/gqb_test.sqlite "

echo "Creating table for testing..."
$SQLITECMD << EOS
DROP TABLE If EXISTS companies;
CREATE TABLE companies (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT,
  created_at TEXT
);

INSERT INTO companies (id, name, created_at) VALUES (1, 'Google', '2018-07-30 00:00:00'), (2, 'Apple', '2018-07-30 00:00:00'), (3, 'Microsoft', '2018-07-30 00:00:00');

DROP TABLE IF EXISTS company_attributes;
CREATE TABLE company_attributes (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  company_id INT,
  url TEXT
);

INSERT INTO company_attributes (company_id, url) VALUES (1, 'https://google.com'), (2, 'https://apple.com'), (3, 'https://microsoft.com');
EOS

echo "Done."

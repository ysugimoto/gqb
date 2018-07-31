#!/bin/bash

DIR=$(cd $(dirname $0) && pwd)
SQLITECMD="sqlite3 /tmp/gqb_test.sqlite "

echo "Creating table for testing..."
$SQLITECMD < "${DIR}/../../examples/sqlite/schema.sql"

echo "Done."

#!/bin/bash

DIR=$(cd $(dirname $0) && pwd)
MYSQLCMD="mysql --defaults-extra-file=${DIR}/my.conf"

echo "Creating database for example..."
echo "CREATE DATABASE IF NOT EXISTS example;" | $MYSQLCMD

echo "Creating table for testing..."
$MYSQLCMD --database example < "${DIR}/../../examples/mysql/schema.sql"

echo "Done."

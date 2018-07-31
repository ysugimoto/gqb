#!/bin/bash

DIR=$(cd $(dirname $0) && pwd)
PSQLCMD="psql -U postgres -h 127.0.0.1"

echo "Creating table for testing..."
$PSQLCMD -d postgres < "${DIR}/../../examples/postgres/schema.sql"

echo "Done."

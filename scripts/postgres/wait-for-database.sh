#!/bin/bash

DIR=$(cd $(dirname $0) && pwd)
PSQLCMD="psql -U postgres -h 127.0.0.1"

for i in $(seq 10); do
  $PSQLCMD -c '\l' || (sleep 5; false) && break
done
echo "PostgreSQL has been started"

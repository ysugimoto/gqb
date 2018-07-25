#!/bin/bash

for i in $(seq 10); do
  mysql -h 127.0.0.1 -u root -proot -e 'show databases' || (sleep 5; false) && break
done
echo "MySQL has been started"

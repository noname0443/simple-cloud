#!/bin/bash

( echo "cat <<EOF" ; cat /etc/configs/mysync.yaml ) | sh > /etc/mysync.yaml

if [ ! -f /tmp/usedspace ]; then
    echo 10 > /tmp/usedspace
fi
if [ ! -f /tmp/readonly ]; then
    echo "false" > /tmp/readonly
fi

if [ ! -f /tmp/cluster_initialized ]; then
    rm -rf /var/lib/mysql/*
    mysqld --initialize --init-file=/etc/configs/init.sql
    echo 'true' > /tmp/cluster_initialized
fi

chown mysql:mysql -R /var/lib/mysql

service mysql start

mysql -proot -e "SELECT 1;"
while [ $? != 0 ]
do
mysql -proot -e "SELECT 1;"
done

while :
do
if ! pgrep -x "mysync" > /dev/null
then
	service mysync-init start
fi

if ! pgrep -x "mysqld" > /dev/null
then
	service mysql start
fi
done
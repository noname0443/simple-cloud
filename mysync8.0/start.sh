#!/bin/bash

export SERVER_ID=$(($(echo $POD_NAME | grep -Eo "[0-9]*") + 1000))

( echo "cat <<EOF" ; cat /etc/configs/mysync.yaml ) | sh > /etc/mysync.yaml
( echo "cat <<EOF" ; cat /etc/configs/server.cnf ) | sh > /etc/mysql/conf.d/server.cnf

if [ ! -f /tmp/usedspace ]; then
    echo 10 > /tmp/usedspace
fi

if [ ! -f /tmp/readonly ]; then
    echo "false" > /tmp/readonly
fi

if [ -z "$(ls -A /var/lib/mysql)" ]; then
    mysqld --initialize --init-file=/etc/configs/init.sql
    chown mysql:mysql -R /var/lib/mysql
fi

service mysql start

while :
do
if ! pgrep -x "mysqld" > /dev/null
then
	service mysql start
fi

if ! pgrep -x "mysync" > /dev/null
then
	mysync >> /var/log/mysync.log 2>&1 &
fi
done
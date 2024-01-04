#!/bin/bash

chown mysql:mysql -R /var/lib/mysql

service mysql start

mysql -e "SELECT 1;"
while [ $? != 0 ]
do
mysql -e "SELECT 1;"
done

mysql -e "CREATE USER 'root'@'10.%.%.%' IDENTIFIED BY '$MYSQL_ROOT_PASSWORD';GRANT ALL PRIVILEGES ON *.* TO 'root'@'10.%.%.%' WITH GRANT OPTION;FLUSH PRIVILEGES;"
mysql -e "ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY '$MYSQL_ROOT_PASSWORD';"
mysql -proot -e "CREATE USER 'user'@'%' IDENTIFIED WITH mysql_native_password BY '$MYSQL_ROOT_PASSWORD';GRANT ALL PRIVILEGES ON *.* TO 'user'@'%';FLUSH PRIVILEGES;"

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
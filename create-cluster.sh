#!/bin/bash
if [ -z "$1" ]
  then
    echo "ERROR: enter name of cluster"
fi

if [ -z "$2" ]
  then
    echo "ERROR: enter replica count"
fi

if [ -z "$2" ]
  then
    echo "ERROR: enter mysql password"
fi

helm install $1 mysql-chart --atomic --set replicaCount=$2 --set mysql.password=$3
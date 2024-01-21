#!/bin/zsh

name=$1
count=$2
pass=$3

if [ -z "$name" ]; then
	echo "ERROR: enter name of cluster"
  exit 1;
fi

if [ -z "$count" ]; then
	echo "ERROR: enter replica count"
  exit 2;
fi

helm upgrade $name mysql-chart --atomic --set replicaCount=$count --timeout 10m0s
